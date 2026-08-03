package user

import (
	"encoding/json"
	"net/http"

	"github.com/seeley-chen/vsb-server/internal/model"
	"github.com/seeley-chen/vsb-server/internal/service/user"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

// RegisterRequest 注册请求体
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Phone    string `json:"phone"`
}

// Create 创建用户
// @Description 使用用户名、密码、登录类型注册新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=model.User} "注册成功"
// @Failure 400 {object} response.Response "参数错误或登录类型不支持"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/user/register [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	// 简单参数校验
	if req.Account == "" || req.Password == "" {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "account and password are required")
		return
	}

	username := req.Username
	if username == "" {
		username = req.Account
	}

	newUser, err := h.svc.Create(r.Context(), model.RegisterRequest(req))

	if err != nil {
		switch err {
		case user.ErrAccountExists:
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		default:
			zap.L().Error("register failed", zap.Error(err))
			response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
		}
		return
	}

	response.Success(w, newUser)
}
