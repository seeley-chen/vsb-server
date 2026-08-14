package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Vanselyn/vsb-server/pkg/idgen"
	"go.uber.org/zap"
)

// 请求级可变信息载体：由 Logger 注入 context，内层中间件（如 Auth）写入字段，Logger 在请求结束后读取。
// 这样外层 Logger 能拿到内层 Auth 设置的 user_id（context.WithValue 本身不可回传到外层 request）。
type reqInfo struct {
	requestID string
	userID    string
}

type contextKey string

const (
	// UserIDKey 用户 ID 的 context key（保留供 handler 直接取用）
	UserIDKey contextKey = "user_id"
	// RequestIDKey 请求 ID 的 context key
	RequestIDKey contextKey = "request_id"
	// reqInfoKey 内部传递 reqInfo 指针的 key
	reqInfoKey contextKey = "req_info"
)

// body 截断阈值，超过此长度的 body 在日志中截断显示
const maxBodyLogLen = 1024

type statusRecorder struct {
	http.ResponseWriter
	status   int
	response bytes.Buffer
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	// 记录响应体（截断到合理长度，避免大响应吃内存）
	if r.response.Len() < maxBodyLogLen {
		remain := maxBodyLogLen - r.response.Len()
		if len(b) <= remain {
			r.response.Write(b)
		} else {
			r.response.Write(b[:remain])
		}
	}
	return r.ResponseWriter.Write(b)
}

// truncateBody 截断 body 用于日志展示
func truncateBody(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	s := string(b)
	if len(s) > maxBodyLogLen {
		return s[:maxBodyLogLen] + "...(truncated)"
	}
	return s
}

// formatBody 根据 bodyMode 处理 body 用于日志展示：
//   - full:   完整内容（截断到 maxBodyLogLen）
//   - masked: 敏感字段（password/secret/token 等）值替换为 *** 后展示
//   - off:    只展示长度，不展示内容
func formatBody(b []byte, mode string) string {
	if len(b) == 0 {
		return ""
	}
	switch mode {
	case "off":
		return fmt.Sprintf("(%d bytes)", len(b))
	case "full":
		return truncateBody(b)
	default: // masked
		return maskBody(b)
	}
}

// maskBody 对 JSON body 的敏感字段做脱敏（值替换为 ***），非 JSON 回退到截断展示
func maskBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		// 非 JSON，回退到截断展示
		return truncateBody(body)
	}
	maskMap(m)
	b, _ := json.Marshal(m)
	return truncateBody(b)
}

// maskMap 递归脱敏 map 中的敏感字段
func maskMap(m map[string]interface{}) {
	for k, v := range m {
		if isSensitiveKey(k) {
			m[k] = "***"
			continue
		}
		if nested, ok := v.(map[string]interface{}); ok {
			maskMap(nested)
		}
	}
}

// isSensitiveKey 判断字段名是否敏感（需脱敏）
func isSensitiveKey(key string) bool {
	lk := strings.ToLower(key)
	switch lk {
	case "password", "pwd", "secret", "token", "authorization", "oldpassword", "newpassword":
		return true
	}
	return false
}

// Logger 请求日志中间件：记录 request_id、method、path、query、status、duration、
// user_id、request body、response body 等，并按状态码分级（>=500 Error / >=400 Warn / 其余 Info）。
// bodyMode 控制 body 日志粒度：full（完整）/ masked（敏感字段脱敏）/ off（仅长度）。
// 同时通过 X-Request-ID 响应头回传 request_id，便于前端反馈定位。
func Logger(logger *zap.Logger, bodyMode string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 生成 request_id
			requestID := idgen.GenerateUuid()[:8]

			// 读取 request body（可重读，供后续 handler 正常解析）
			var reqBody []byte
			if r.Body != nil && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) {
				reqBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewReader(reqBody))
			}

			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			// 注入 reqInfo 到 context，供内层 Auth 写入 user_id
			ri := &reqInfo{requestID: requestID}
			ctx := context.WithValue(r.Context(), reqInfoKey, ri)
			ctx = context.WithValue(ctx, RequestIDKey, requestID)
			r = r.WithContext(ctx)

			// 回写 request_id 到响应头，前端可反馈此值用于定位
			w.Header().Set("X-Request-ID", requestID)

			next.ServeHTTP(rec, r)

			duration := time.Since(start)

			// 组装日志字段
			fields := []zap.Field{
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.Int("status", rec.status),
				zap.Duration("duration", duration),
				zap.String("ip", clientIP(r)),
				zap.String("ua", r.UserAgent()),
				zap.Int64("content_length", r.ContentLength),
				zap.String("user_id", ri.userID),
				zap.String("req_body", formatBody(reqBody, bodyMode)),
				zap.String("resp_body", formatBody(rec.response.Bytes(), bodyMode)),
			}

			level := "INFO"
			switch {
			case rec.status >= 500:
				level = "ERROR"
				logger.Error("request", fields...)
			case rec.status >= 400:
				level = "WARN"
				logger.Warn("request", fields...)
			default:
				logger.Info("request", fields...)
			}

			// 推送到日志查看器（SSE 网页）
			addLogEntry(LogEntry{
				Time:      start,
				Level:     level,
				RequestID: requestID,
				Method:    r.Method,
				Path:      r.URL.Path,
				Query:     r.URL.RawQuery,
				Status:    rec.status,
				Duration:  duration.String(),
				IP:        clientIP(r),
				UA:        r.UserAgent(),
				UserID:    ri.userID,
				ReqBody:   formatBody(reqBody, bodyMode),
				RespBody:  formatBody(rec.response.Bytes(), bodyMode),
			})
		})
	}
}

// clientIP 提取客户端 IP（优先 X-Forwarded-For / X-Real-IP）
func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// 取第一个（最原始的客户端）
		if idx := bytes.IndexByte([]byte(ip), ','); idx >= 0 {
			return ip[:idx]
		}
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// GetRequestID 从 context 中获取请求 ID
func GetRequestID(ctx context.Context) string {
	if rid, ok := ctx.Value(RequestIDKey).(string); ok {
		return rid
	}
	return ""
}

// setReqInfoUserID 供 Auth 中间件调用，把 user_id 写入 reqInfo（供外层 Logger 读取）
func setReqInfoUserID(ctx context.Context, userID string) {
	if ri, ok := ctx.Value(reqInfoKey).(*reqInfo); ok && ri != nil {
		ri.userID = userID
	}
}

// LogError 打一条带 request_id / method / path 的错误日志，供 handler 统一记录底层错误细节。
// 与 Logger 中间件的请求维度日志（记 resp_body）互补：本函数记原始 err，便于排查根因。
func LogError(r *http.Request, op string, err error) {
	zap.L().Error(op,
		zap.String("request_id", GetRequestID(r.Context())),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Error(err),
	)
}
