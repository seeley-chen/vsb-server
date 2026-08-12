package user

import (
	"errors"
	"net/http"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req model.UserUpdateRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	userID, ok := response.PathVar(w, r, "id", "user id is required")
	if !ok {
		return
	}

	user, err := h.svc.UpdateUser(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, svc.ErrUserNotFound) {
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, "user not found")
			return
		}
		if errors.Is(err, svc.ErrUserNameEmpty) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user name is required")
			return
		}
		if errors.Is(err, svc.ErrUserAccountEmpty) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user account already exists")
			return
		}
		if errors.Is(err, svc.ErrInvalidPassword) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid password")
			return
		}
		if errors.Is(err, svc.ErrUserRoleEmpty) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user role is required")
			return
		}
		if errors.Is(err, svc.ErrUserDepartmentEmpty) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user department is required")
			return
		}
		zap.L().Error("update user failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(w, user)
}
