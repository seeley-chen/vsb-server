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

func (h *Handler) RegisterPrivateRoutes(r *mux.Router) {
	sub := r.PathPrefix("/api/permission/department").Subrouter()
	sub.HandleFunc("/list", h.List).Methods("GET")
}
