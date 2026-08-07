package permission

import (
	"github.com/gorilla/mux"

	"github.com/seeley-chen/vsb-server/internal/handler/permission/department"
	"github.com/seeley-chen/vsb-server/internal/handler/permission/role"
)

// Handler 权限大模块入口，统一挂载 /permission 前缀下的子模块路由。
type Handler struct {
	Department *department.Handler
	Role       *role.Handler
}

func NewHandler(department *department.Handler, role *role.Handler) *Handler {
	return &Handler{
		Department: department,
		Role:       role,
	}
}

// RegisterRoutes 在 protected 路由（/api）下注册权限模块全部接口。
func (h *Handler) RegisterRoutes(r *mux.Router) {
	sub := r.PathPrefix("/permission").Subrouter()
	h.Department.RegisterRoutes(sub)
	h.Role.RegisterRoutes(sub)
}
