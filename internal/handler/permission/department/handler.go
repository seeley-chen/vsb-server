package department

import (
	"github.com/gorilla/mux"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
)

type Handler struct {
	svc *svc.DepartmentService
}

func NewHandler(svc *svc.DepartmentService) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes 挂载 /department 子模块路由，完整路径示例：/api/permission/department/create
func (h *Handler) RegisterRoutes(r *mux.Router) {
	sub := r.PathPrefix("/department").Subrouter()
	sub.HandleFunc("/create", h.Create).Methods("POST")
	sub.HandleFunc("/list", h.List).Methods("GET")
	sub.HandleFunc("/update/{id}", h.Update).Methods("PUT")
	sub.HandleFunc("/delete/{id}", h.Delete).Methods("DELETE")
}
