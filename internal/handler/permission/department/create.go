package department

import (
	"errors"
	"net/http"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.DepartmentRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	department, err := h.svc.CreateDepartment(r.Context(), &req)
	if err != nil {
		if errors.Is(err, svc.ErrDepartmentNameEmpty) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "department name is required")
			return
		}
		if errors.Is(err, svc.ErrDepartmentExists) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "department already exists")
			return
		}
		zap.L().Error("create department failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(w, department)
}
