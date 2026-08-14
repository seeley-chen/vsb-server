package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// LogEntry 日志查看器展示的一条请求日志
type LogEntry struct {
	Time      time.Time `json:"time"`
	Level     string    `json:"level"` // INFO / WARN / ERROR
	RequestID string    `json:"request_id"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Query     string    `json:"query"`
	Status    int       `json:"status"`
	Duration  string    `json:"duration"`
	IP        string    `json:"ip"`
	UA        string    `json:"ua"`
	UserID    string    `json:"user_id"`
	ReqBody   string    `json:"req_body"`
	RespBody  string    `json:"resp_body"`
}

// logBuffer 环形缓冲 + SSE 广播
type logBuffer struct {
	mu      sync.Mutex
	entries []LogEntry
	max     int
	next    int // 环形写入位置
	full    bool
	subs    map[chan LogEntry]struct{}
}

var defaultBuffer *logBuffer

// InitLogViewer 初始化日志查看器缓冲区
func InitLogViewer(maxSize int) {
	defaultBuffer = &logBuffer{
		entries: make([]LogEntry, maxSize),
		max:     maxSize,
		subs:    make(map[chan LogEntry]struct{}),
	}
}

// addLogEntry 写入缓冲区并广播给所有 SSE 订阅者
func addLogEntry(e LogEntry) {
	if defaultBuffer == nil {
		return
	}
	defaultBuffer.add(e)
}

func (b *logBuffer) add(e LogEntry) {
	b.mu.Lock()
	b.entries[b.next] = e
	b.next = (b.next + 1) % b.max
	if b.next == 0 {
		b.full = true
	}
	subs := make([]chan LogEntry, 0, len(b.subs))
	for ch := range b.subs {
		subs = append(subs, ch)
	}
	b.mu.Unlock()

	// 非阻塞广播，慢客户端丢弃消息
	for _, ch := range subs {
		select {
		case ch <- e:
		default:
		}
	}
}

// snapshot 返回当前缓冲区的历史日志（按时间顺序）
func (b *logBuffer) snapshot() []LogEntry {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.full {
		return append([]LogEntry(nil), b.entries[:b.next]...)
	}
	result := make([]LogEntry, 0, b.max)
	result = append(result, b.entries[b.next:]...)
	result = append(result, b.entries[:b.next]...)
	return result
}

func (b *logBuffer) subscribe() chan LogEntry {
	ch := make(chan LogEntry, 64)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *logBuffer) unsubscribe(ch chan LogEntry) {
	b.mu.Lock()
	delete(b.subs, ch)
	b.mu.Unlock()
	close(ch)
}

// LogViewerPage 返回日志查看器 HTML 页面
func LogViewerPage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(logViewerHTML))
}

// LogViewerStream SSE 端点：先推送历史日志，再实时推送新日志
func LogViewerStream(w http.ResponseWriter, r *http.Request) {
	if defaultBuffer == nil {
		http.Error(w, "log viewer not initialized", http.StatusServiceUnavailable)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // nginx 不缓冲

	// 推送历史日志
	for _, e := range defaultBuffer.snapshot() {
		writeSSE(w, e)
	}
	flusher.Flush()

	// 订阅实时日志
	ch := defaultBuffer.subscribe()
	defer defaultBuffer.unsubscribe(ch)

	ticker := time.NewTicker(15 * time.Second) // 心跳保活
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case e := <-ch:
			writeSSE(w, e)
			flusher.Flush()
		case <-ticker.C:
			// 心跳注释，防止连接超时
			_, _ = fmt.Fprintf(w, ": ping\n\n")
			flusher.Flush()
		}
	}
}

// writeSSE 写一条 SSE 事件
func writeSSE(w http.ResponseWriter, e LogEntry) {
	data, _ := json.Marshal(e)
	_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
}
