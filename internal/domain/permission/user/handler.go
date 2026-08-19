package user

import (
	"errors"
	"net/http"

	"github.com/Vanselyn/vsb-server/internal/middleware"
	"github.com/Vanselyn/vsb-server/internal/models"
	"github.com/Vanselyn/vsb-server/pkg/response"
	"github.com/gorilla/mux"
)

type Handler struct {
	svc *UserService
}

func NewHandler(svc *UserService) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes 挂载 /user 子模块路由，完整路径示例：/api/permission/user/create
func (h *Handler) RegisterRoutes(r *mux.Router) {
	sub := r.PathPrefix("/user").Subrouter()
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
	case errors.Is(err, ErrUserAlreadyExists):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user account already exists")
	case errors.Is(err, ErrUserNotFound):
		response.Fail(w, http.StatusNotFound, response.CodeNotFound, "user not found")
	default:
		middleware.LogError(r, op, err)
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
	}
}

// Create 创建用户
// @Summary 创建用户
// @Description 创建新用户，需要 Bearer Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body UserRequest true "用户信息"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=UserResponse} "创建成功"
// @Failure 400 {object} response.Response "参数错误或账号已存在"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/permission/user/create [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	user, err := h.svc.CreateUser(r.Context(), &req)
	if err != nil {
		h.handleErr(w, r, err, "create user failed")
		return
	}

	response.Success(w, user)
}

// List 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表，支持 username/account 模糊搜索、email 精确匹配，需要 Bearer Token
// @Tags 用户管理
// @Produce json
// @Param pageIndex query int false "页码，默认 1"
// @Param pageSize query int false "每页数量，默认 20，上限 100"
// @Param username query string false "用户名（模糊搜索）"
// @Param account query string false "账号（模糊搜索）"
// @Param email query string false "邮箱（精确匹配）"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=response.PageData} "用户列表"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/permission/user/list [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageIndex, pageSize := response.PageQuery(r)
	q := r.URL.Query()
	params := UserQueryParams{
		Username:     q.Get("username"),
		Account:      q.Get("account"),
		Email:        q.Get("email"),
		RoleID:       q.Get("roleId"),
		DepartmentID: q.Get("departmentId"),
		Status:       models.StatusEnum(q.Get("status")),
	}

	users, total, err := h.svc.GetUserList(r.Context(), pageIndex, pageSize, params)
	if err != nil {
		h.handleErr(w, r, err, "list users failed")
		return
	}

	response.Success(w, response.PageData{
		List:      users,
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
}

// Update 更新用户
// @Summary 更新用户
// @Description 根据用户 ID 更新用户信息，空字段表示不修改，需要 Bearer Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Param request body UserUpdateRequest true "更新信息"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=UserResponse} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/permission/user/update/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req UserUpdateRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	userID, ok := response.PathVar(w, r, "id", "user id is required")
	if !ok {
		return
	}

	user, err := h.svc.UpdateUser(r.Context(), userID, &req)
	if err != nil {
		h.handleErr(w, r, err, "update user failed")
		return
	}

	response.Success(w, user)
}

// Delete 删除用户
// @Summary 删除用户
// @Description 根据用户 ID 删除用户，需要 Bearer Token
// @Tags 用户管理
// @Produce json
// @Param id path string true "用户 ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/permission/user/delete/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := response.PathVar(w, r, "id", "user id is required")
	if !ok {
		return
	}

	if err := h.svc.DeleteUser(r.Context(), userID); err != nil {
		h.handleErr(w, r, err, "delete user failed")
		return
	}

	response.Success(w, nil)
}
