package department

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	departmentID := mux.Vars(r)["id"]
	if departmentID == "" {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "department id is required")
		return
	}

	var req model.DepartmentUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid request")
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
