package login

import (
	"errors"
	"net/http"

	"github.com/Vanselyn/vsb-server/pkg/response"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
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

// handleErr 统一处理 service 层错误
func (h *Handler) handleErr(w http.ResponseWriter, err error, op string) {
	switch {
	case errors.Is(err, ErrInvalidAccount), errors.Is(err, ErrInvalidPassword):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid account or password")
	default:
		zap.L().Error(op, zap.Error(err))
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
		h.handleErr(w, err, "login failed")
		return
	}

	response.Success(w, login)
}
