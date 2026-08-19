package goods

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	svc *GoodsService
}

func NewGoodsHandler(svc *GoodsService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	sub := r.PathPrefix("/goods").Subrouter()
	sub.HandleFunc("/create", h.Create).Methods("POST")
	sub.HandleFunc("/list", h.List).Methods("GET")
	sub.HandleFunc("/update/{id}", h.Update).Methods("PUT")
	sub.HandleFunc("/delete/{id}", h.Delete).Methods("DELETE")
}

func (*Handler) Create(w http.ResponseWriter, r *http.Request) {}
func (*Handler) List(w http.ResponseWriter, r *http.Request)   {}
func (*Handler) Update(w http.ResponseWriter, r *http.Request) {}
func (*Handler) Delete(w http.ResponseWriter, r *http.Request) {}
