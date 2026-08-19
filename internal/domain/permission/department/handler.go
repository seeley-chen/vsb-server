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
	if response.FailIfValidation(w, err) {
		return
	}
	switch {
	case errors.Is(err, ErrDepartmentExists):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "department already exists")
	case errors.Is(err, ErrDepartmentNotFound):
		response.Fail(w, http.StatusNotFound, response.CodeNotFound, "department not found")
	default:
		middleware.LogError(r, op, err)
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
	}
}

// Create 创建部门
// @Summary 创建部门
// @Description 创建新部门，需要 Bearer Token
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param request body DepartmentRequest true "部门信息"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=DepartmentResponse} "创建成功"
// @Failure 400 {object} response.Response "参数错误或部门已存在"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/permission/department/create [post]
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

// List 获取部门列表
// @Summary 获取部门列表
// @Description 分页获取部门列表，支持按名称模糊搜索，需要 Bearer Token
// @Tags 部门管理
// @Produce json
// @Param pageIndex query int false "页码，默认 1"
// @Param pageSize query int false "每页数量，默认 20，上限 100"
// @Param name query string false "部门名称（模糊搜索）"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=response.PageData} "部门列表"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/permission/department/list [get]
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

// Update 更新部门
// @Summary 更新部门
// @Description 根据部门 ID 更新部门信息，空字段表示不修改，需要 Bearer Token
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id path string true "部门 ID"
// @Param request body DepartmentUpdateRequest true "更新信息"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=DepartmentResponse} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "部门不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/permission/department/update/{id} [put]
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

// Delete 删除部门
// @Summary 删除部门
// @Description 根据部门 ID 删除部门，需要 Bearer Token
// @Tags 部门管理
// @Produce json
// @Param id path string true "部门 ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "部门不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/permission/department/delete/{id} [delete]
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
