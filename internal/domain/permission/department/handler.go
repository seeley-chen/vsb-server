package department

import (
	"errors"
	"net/http"

	"github.com/Vanselyn/vsb-server/internal/middleware"
	"github.com/Vanselyn/vsb-server/pkg/response"
	"github.com/gorilla/mux"
)

type Handler struct {
	svc *DepartmentService
}

func NewHandler(svc *DepartmentService) *Handler {
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

// handleErr 统一处理 service 层错误：已知业务错误映射为对应 HTTP 响应，
// 未知错误（default）通过 middleware.LogError 记录带 request_id 的底层错误日志后返回 500。
func (h *Handler) handleErr(w http.ResponseWriter, r *http.Request, err error, op string) {
	switch {
	case errors.Is(err, ErrDepartmentNameEmpty):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "department name is required")
	case errors.Is(err, ErrDepartmentExists):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "department already exists")
	case errors.Is(err, ErrDepartmentNotFound):
		response.Fail(w, http.StatusNotFound, response.CodeNotFound, "department not found")
	default:
		middleware.LogError(r, op, err)
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req DepartmentRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	department, err := h.svc.CreateDepartment(r.Context(), &req)
	if err != nil {
		h.handleErr(w, r, err, "create department failed")
		return
	}

	response.Success(w, department)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageIndex, pageSize := response.PageQuery(r)
	q := r.URL.Query()
	name := q.Get("name")

	departments, total, err := h.svc.GetDepartmentList(r.Context(), pageIndex, pageSize, name)
	if err != nil {
		h.handleErr(w, r, err, "list departments failed")
		return
	}

	response.Success(w, response.PageData{
		List:      departments,
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req DepartmentUpdateRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	departmentID, ok := response.PathVar(w, r, "id", "department id is required")
	if !ok {
		return
	}

	department, err := h.svc.UpdateDepartment(r.Context(), departmentID, &req)
	if err != nil {
		h.handleErr(w, r, err, "update department failed")
		return
	}

	response.Success(w, department)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	departmentID, ok := response.PathVar(w, r, "id", "department id is required")
	if !ok {
		return
	}

	if err := h.svc.DeleteDepartment(r.Context(), departmentID); err != nil {
		h.handleErr(w, r, err, "delete department failed")
		return
	}

	response.Success(w, nil)
}
