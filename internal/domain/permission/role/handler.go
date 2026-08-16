package role

import (
	"errors"
	"net/http"

	"github.com/Vanselyn/vsb-server/internal/domain/permission/user"
	"github.com/Vanselyn/vsb-server/internal/middleware"
	"github.com/Vanselyn/vsb-server/pkg/response"
	"github.com/gorilla/mux"
)

type Handler struct {
	roleSvc *RoleService
	userSvc *user.UserService
}

func NewHandler(roleSvc *RoleService, userSvc *user.UserService) *Handler {
	return &Handler{
		roleSvc: roleSvc,
		userSvc: userSvc,
	}
}

// RegisterRoutes 挂载 /role 子模块路由，完整路径示例：/api/permission/role/create
func (h *Handler) RegisterRoutes(r *mux.Router) {
	sub := r.PathPrefix("/role").Subrouter()
	sub.HandleFunc("/create", h.Create).Methods("POST")
	sub.HandleFunc("/list", h.List).Methods("GET")
	sub.HandleFunc("/update/{id}", h.Update).Methods("PUT")
	sub.HandleFunc("/delete/{id}", h.Delete).Methods("DELETE")
}

// handleErr 统一处理 service 层错误：已知业务错误映射为对应 HTTP 响应，
// 未知错误（default）通过 middleware.LogError 记录带 request_id 的底层错误日志后返回 500。
func (h *Handler) handleErr(w http.ResponseWriter, r *http.Request, err error, op string) {
	switch {
	case errors.Is(err, ErrRoleNameEmpty):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "role name is required")
	case errors.Is(err, ErrRoleExists):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "role already exists")
	case errors.Is(err, ErrRoleNotFound):
		response.Fail(w, http.StatusNotFound, response.CodeNotFound, "role not found")
	case errors.Is(err, ErrInvalidPermission):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid permission")
	case errors.Is(err, user.ErrUserNotFound):
		response.Fail(w, http.StatusNotFound, response.CodeNotFound, "user not found")
	default:
		middleware.LogError(r, op, err)
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req RoleRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	role, err := h.roleSvc.CreateRole(r.Context(), &req)
	if err != nil {
		h.handleErr(w, r, err, "create role failed")
		return
	}

	response.Success(w, role)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageIndex, pageSize := response.PageQuery(r)

	roles, total, err := h.roleSvc.GetRoleList(r.Context(), pageIndex, pageSize)
	if err != nil {
		h.handleErr(w, r, err, "list roles failed")
		return
	}

	response.Success(w, response.PageData{
		List:      roles,
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	roleId, ok := response.PathVar(w, r, "id", "role id is required")
	if !ok {
		return
	}

	var req RoleUpdateRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	role, err := h.roleSvc.UpdateRole(r.Context(), roleId, &req)
	if err != nil {
		h.handleErr(w, r, err, "update role failed")
		return
	}

	response.Success(w, role)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	roleId, ok := response.PathVar(w, r, "id", "role id is required")
	if !ok {
		return
	}

	if err := h.roleSvc.DeleteRole(r.Context(), roleId); err != nil {
		h.handleErr(w, r, err, "delete role failed")
		return
	}

	response.Success(w, nil)
}

// GetPrivileges 获取当前登录用户的合并权限
func (h *Handler) GetPrivileges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		response.Fail(w, http.StatusUnauthorized, response.CodeUnauthorized, "user not authenticated")
		return
	}

	u, err := h.userSvc.GetUserById(ctx, userID)
	if err != nil {
		h.handleErr(w, r, err, "get user by id failed")
		return
	}

	roleIds := make([]string, 0, 1)
	if u.RoleId != "" {
		roleIds = append(roleIds, u.RoleId)
	}

	roles, err := h.roleSvc.GetRolesByIds(ctx, roleIds)
	if err != nil {
		h.handleErr(w, r, err, "get roles by ids failed")
		return
	}

	permissions := h.roleSvc.MergePermissions(roles)

	response.Success(w, permissions)
}
