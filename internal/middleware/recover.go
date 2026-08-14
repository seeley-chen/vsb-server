package middleware

import (
	"net/http"

	"github.com/Vanselyn/vsb-server/pkg/response"
	"go.uber.org/zap"
)

// Recover 捕获 panic，避免单个请求导致进程崩溃
func Recover(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						zap.Any("panic", rec),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
					)
					response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
