package role

import (
	"github.com/gorilla/mux"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
)

type Handler struct {
	svc *svc.RoleService
}

func NewHandler(svc *svc.RoleService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	sub := r.PathPrefix("/role").Subrouter()
	sub.HandleFunc("/create", h.Create).Methods("POST")
	sub.HandleFunc("/list", h.List).Methods("GET")
	sub.HandleFunc("/update/{id}", h.Update).Methods("PUT")
	sub.HandleFunc("/delete/{id}", h.Delete).Methods("DELETE")
}
