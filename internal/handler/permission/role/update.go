package role

import (
	"errors"
	"net/http"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	roleId, ok := response.PathVar(w, r, "id", "role id is required")
	if !ok {
		return
	}

	var req model.RoleUpdateRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	role, err := h.svc.UpdateRole(r.Context(), roleId, &req)
	if err != nil {
		if errors.Is(err, svc.ErrRoleNotFound) {
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, "role not found")
			return
		}
		zap.L().Error("update role failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "update role failed")
		return
	}

	response.Success(w, role)
}
