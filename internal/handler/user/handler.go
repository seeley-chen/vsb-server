package user

import (
	"github.com/gorilla/mux"
	svc "github.com/seeley-chen/vsb-server/internal/service/user"
)

type Handler struct {
	svc *svc.Service
}

func NewHandler(service *svc.Service) *Handler {
	return &Handler{svc: service}
}

// RegisterPublicRoutes 注册公开路由（无需鉴权）
func (h *Handler) RegisterPublicRoutes(r *mux.Router) {
	sub := r.PathPrefix("/api/user").Subrouter()
	sub.HandleFunc("/register", h.Register).Methods("POST")
	sub.HandleFunc("/login", h.Login).Methods("POST")
}

// RegisterProtectedRoutes 注册需要鉴权的路由
func (h *Handler) RegisterProtectedRoutes(r *mux.Router) {
	sub := r.PathPrefix("/api/user").Subrouter()
	sub.HandleFunc("/list", h.List).Methods("GET")
}
