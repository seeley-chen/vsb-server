package user

import (
	"errors"
	"net/http"

	"github.com/Vanselyn/vsb-server/pkg/response"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
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

// handleErr 统一处理 service 层错误
func (h *Handler) handleErr(w http.ResponseWriter, err error, op string) {
	switch {
	case errors.Is(err, ErrUserNameEmpty):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user name is required")
	case errors.Is(err, ErrUserAccountEmpty):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user account is required")
	case errors.Is(err, ErrInvalidPassword):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid password")
	case errors.Is(err, ErrUserRoleEmpty):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user role is required")
	case errors.Is(err, ErrUserDepartmentEmpty):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user department is required")
	case errors.Is(err, ErrUserAlreadyExists):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user account already exists")
	case errors.Is(err, ErrUserNotFound):
		response.Fail(w, http.StatusNotFound, response.CodeNotFound, "user not found")
	default:
		zap.L().Error(op, zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	user, err := h.svc.CreateUser(r.Context(), &req)
	if err != nil {
		h.handleErr(w, err, "create user failed")
		return
	}

	response.Success(w, user)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageIndex, pageSize := response.PageQuery(r)
	q := r.URL.Query()
	username := q.Get("username")
	account := q.Get("account")
	email := q.Get("email")

	users, total, err := h.svc.GetUserList(r.Context(), pageIndex, pageSize, username, account, email)
	if err != nil {
		zap.L().Error("list users failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(w, response.PageData{
		List:      users,
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
}

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
		h.handleErr(w, err, "update user failed")
		return
	}

	response.Success(w, user)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := response.PathVar(w, r, "id", "user id is required")
	if !ok {
		return
	}

	if err := h.svc.DeleteUser(r.Context(), userID); err != nil {
		h.handleErr(w, err, "delete user failed")
		return
	}

	response.Success(w, nil)
}
