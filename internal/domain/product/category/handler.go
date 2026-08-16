package category

import (
	"errors"
	"net/http"

	"github.com/Vanselyn/vsb-server/internal/middleware"
	"github.com/Vanselyn/vsb-server/pkg/response"
	"github.com/gorilla/mux"
)

type Handler struct {
	svc *CategoryService
}

func NewCategoryHandler(svc *CategoryService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	sub := r.PathPrefix("/category").Subrouter()
	sub.HandleFunc("/create", h.Create).Methods("POST")
	sub.HandleFunc("/list", h.List).Methods("GET")
	sub.HandleFunc("/update/{id}", h.Update).Methods("PUT")
	sub.HandleFunc("/delete/{id}", h.Delete).Methods("DELETE")
}

func (h *Handler) handleErr(w http.ResponseWriter, r *http.Request, err error, op string) {
	switch {
	case errors.Is(err, ErrCategoryExists):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "category exists")

	case errors.Is(err, ErrCategoryNotFound):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "category not found")

	default:
		middleware.LogError(r, op, err)
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
	}

}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CategoryRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	category, err := h.svc.CreateCategory(r.Context(), &req)
	if err != nil {
		h.handleErr(w, r, err, "CreateCategory")
		return
	}
	response.Success(w, category)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	categories, total, err := h.svc.CategoryList(r.Context())
	if err != nil {
		h.handleErr(w, r, err, "CategoryList")
		return
	}
	response.Success(w, response.PageData{
		List:  categories,
		Total: total,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req CategoryUpdateRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	categoryID, ok := response.PathVar(w, r, "id", "category id is required")
	if !ok {
		return
	}

	category, err := h.svc.UpdateCategory(r.Context(), categoryID, &req)
	if err != nil {
		h.handleErr(w, r, err, "update category failed")
		return
	}

	response.Success(w, category)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	categoryID, ok := response.PathVar(w, r, "id", "category id is required")
	if !ok {
		return
	}

	err := h.svc.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		h.handleErr(w, r, err, "Delete category failed")
		return
	}

	response.Success(w, nil)
}
