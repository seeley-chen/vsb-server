package middleware

import "net/http"

// CORS 跨域中间件（允许所有来源）
// 必须包裹在 mux.Router 外层使用：gorilla/mux 的 r.Use 仅在路由匹配后执行，
// OPTIONS 预检请求无法匹配 POST-only 路由，会导致浏览器报跨域错误。
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
