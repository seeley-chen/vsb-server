package user

import (
	"github.com/gorilla/mux"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission/user"
)

type Handler struct {
	svc *svc.Service
}

func NewHandler(service *svc.Service) *Handler {
	return &Handler{svc: service}
}

// RegisterPublicRoutes 注册公开路由（无需鉴权）
func (h *Handler) RegisterPublicRoutes(r *mux.Router) {
	sub := r.PathPrefix("/api/permission/user").Subrouter()
	sub.HandleFunc("/login", h.Login).Methods("POST")
}

// RegisterProtectedRoutes 注册需要鉴权的路由（父路由已挂载 /api 前缀）
func (h *Handler) RegisterProtectedRoutes(r *mux.Router) {
	sub := r.PathPrefix("/permission/user").Subrouter()
	sub.HandleFunc("/create", h.Create).Methods("POST")
	sub.HandleFunc("/list", h.List).Methods("GET")
	sub.HandleFunc("/delete/{userId}", h.Delete).Methods("DELETE")
	sub.HandleFunc("/edit/{userId}", h.Edit).Methods("PUT")
}
