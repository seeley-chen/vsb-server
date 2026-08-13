package role

import (
	"errors"
	"net/http"

	"github.com/seeley-chen/vsb-server/internal/middleware"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) GetPrivileges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		response.Fail(w, http.StatusUnauthorized, response.CodeUnauthorized, "user not authenticated")
		return
	}

	user, err := h.userSvc.GetUserById(ctx, userID)
	if err != nil {
		if errors.Is(err, svc.ErrUserNotFound) {
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, "user not found")
			return
		}
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		zap.L().Error("get user by id failed", zap.Error(err))
		return
	}

	roleIds := make([]string, 0, 1)
	if user.RoleId != "" {
		roleIds = append(roleIds, user.RoleId)
	}

	roles, err := h.roleSvc.GetRolesByIds(ctx, roleIds)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		zap.L().Error("get roles by ids failed", zap.Error(err))
		return
	}

	permissions := h.roleSvc.MergePermissions(roles)

	response.Success(w, permissions)
}
