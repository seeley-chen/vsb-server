package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Vanselyn/vsb-server/pkg/jwt"
	"github.com/Vanselyn/vsb-server/pkg/response"
)

// 注意：contextKey / UserIDKey / RequestIDKey 已在 logger.go 统一定义，避免重复。

// Auth JWT 鉴权中间件
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Fail(w, http.StatusUnauthorized, response.CodeUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Fail(w, http.StatusUnauthorized, response.CodeUnauthorized, "invalid authorization header format")
				return
			}

			tokenString := parts[1]
			userId, err := jwt.GetUserIdFromToken(tokenString, jwtSecret)
			if err != nil {
				response.Fail(w, http.StatusUnauthorized, response.CodeUnauthorized, "invalid or expired token")
				return
			}

			// 将 userId 放入 context（供 handler 取用）
			ctx := context.WithValue(r.Context(), UserIDKey, userId)
			// 同步回写 reqInfo，让外层 Logger 能在请求结束时记录 user_id
			setReqInfoUserID(ctx, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID 从 context 中获取用户 ID
func GetUserID(ctx context.Context) string {
	if uid, ok := ctx.Value(UserIDKey).(string); ok {
		return uid
	}
	return ""
}
