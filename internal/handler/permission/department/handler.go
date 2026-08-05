package department

import (
	"github.com/gorilla/mux"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission/department"
)

type Handler struct {
	svc *svc.Service
}

func NewHandler(service *svc.Service) *Handler {
	return &Handler{svc: service}
}

func (h *Handler) RegisterPublicRoutes(r *mux.Router) {
	sub := r.PathPrefix("/api/permission/department").Subrouter()
	sub.HandleFunc("/create", h.Create).Methods("POST")
}

// RegisterProtectedRoutes 注册需要鉴权的路由（父路由已挂载 /api 前缀）
func (h *Handler) RegisterProtectedRoutes(r *mux.Router) {
	sub := r.PathPrefix("/permission/department").Subrouter()
	sub.HandleFunc("/list", h.List).Methods("GET")
}
