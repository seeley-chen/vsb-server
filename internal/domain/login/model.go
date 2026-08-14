package login

import "github.com/Vanselyn/vsb-server/internal/domain/permission/user"

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string             `json:"token"`
	User  *user.UserResponse `json:"user"`
}
