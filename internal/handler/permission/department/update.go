package department

import (
	"errors"
	"net/http"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req model.DepartmentUpdateRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	departmentID, ok := response.PathVar(w, r, "id", "department id is required")
	if !ok {
		return
	}

	department, err := h.svc.UpdateDepartment(r.Context(), departmentID, &req)
	if err != nil {
		if errors.Is(err, svc.ErrDepartmentNotFound) {
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, "department not found")
			return
		}
		zap.L().Error("update department failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(w, department)
}
