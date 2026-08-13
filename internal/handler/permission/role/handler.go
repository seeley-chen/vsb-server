package role

import (
	"github.com/gorilla/mux"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
)

type Handler struct {
	roleSvc *svc.RoleService
	userSvc *svc.UserService
}

func NewHandler(roleSvc *svc.RoleService, userSvc *svc.UserService) *Handler {
	return &Handler{
		roleSvc: roleSvc,
		userSvc: userSvc,
	}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	sub := r.PathPrefix("/role").Subrouter()
	sub.HandleFunc("/create", h.Create).Methods("POST")
	sub.HandleFunc("/list", h.List).Methods("GET")
	sub.HandleFunc("/update/{id}", h.Update).Methods("PUT")
	sub.HandleFunc("/delete/{id}", h.Delete).Methods("DELETE")
}
