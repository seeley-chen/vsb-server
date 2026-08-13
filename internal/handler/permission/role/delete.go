package role

import (
	"errors"
	"net/http"

	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	roleId, ok := response.PathVar(w, r, "id", "role id is required")
	if !ok {
		return
	}

	if err := h.roleSvc.DeleteRole(r.Context(), roleId); err != nil {
		if errors.Is(err, svc.ErrRoleNotFound) {
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, "role not found")
			return
		}

		zap.L().Error("delete role failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(w, nil)
}
