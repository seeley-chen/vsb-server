package department

import (
	"errors"
	"net/http"

	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	departmentID, ok := response.PathVar(w, r, "id", "department id is required")
	if !ok {
		return
	}

	if err := h.svc.DeleteDepartment(r.Context(), departmentID); err != nil {
		if errors.Is(err, svc.ErrDepartmentNotFound) {
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, "department not found")
			return
		}
		zap.L().Error("delete department failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(w, nil)
}
