package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/seeley-chen/vsb-server/internal/service/permission/user"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

// LoginRequest 登录请求体
type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// LoginResponse 登录响应体
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码登录，返回 JWT token 和用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=LoginResponse} "登录成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/permission/user/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	if req.Account == "" || req.Password == "" {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "account and password are required")
		return
	}

	token, userObj, err := h.svc.Login(r.Context(), req.Account, req.Password)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrInvalidPassword) {
			response.Fail(w, http.StatusUnauthorized, response.CodeUnauthorized, "invalid account or password")
			return
		}
		zap.L().Error("login failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
		return
	}

	resp := LoginResponse{
		Token: token,
		User:  userObj,
	}
	response.Success(w, resp)
}
