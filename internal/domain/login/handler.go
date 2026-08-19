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
	if response.FailIfValidation(w, err) {
		return
	}
	switch {
	case errors.Is(err, ErrInvalidAccount), errors.Is(err, ErrInvalidPassword):
		response.Fail(w, http.StatusBadRequest, response.CodeNotFound, "invalid account or password")
	default:
		middleware.LogError(r, op, err)
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用账号和密码登录，返回 JWT token 和用户信息
// @Tags 登录
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=LoginResponse} "登录成功"
// @Failure 400 {object} response.Response "参数错误或账号密码错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/login [post]
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
