package login

import (
	"errors"
	"net/http"

	"github.com/Vanselyn/vsb-server/internal/middleware"
	"github.com/Vanselyn/vsb-server/pkg/response"
	"github.com/gorilla/mux"
)

type Handler struct {
	svc *LoginService
}

func NewHandler(svc *LoginService) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes 挂载 /login 路由，完整路径示例：/api/login
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/login", h.Login).Methods("POST")
}

// handleErr 统一处理 service 层错误：已知业务错误映射为对应 HTTP 响应，
// 未知错误（default）通过 middleware.LogError 记录带 request_id 的底层错误日志后返回 500。
func (h *Handler) handleErr(w http.ResponseWriter, r *http.Request, err error, op string) {
	switch {
	case errors.Is(err, ErrInvalidAccount), errors.Is(err, ErrInvalidPassword):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid account or password")
	default:
		middleware.LogError(r, op, err)
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	login, err := h.svc.Login(r.Context(), &req)
	if err != nil {
		h.handleErr(w, r, err, "login failed")
		return
	}

	response.Success(w, login)
}
