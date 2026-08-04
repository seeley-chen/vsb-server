package department

import (
	"encoding/json"
	"errors"
	"net/http"

	model "github.com/seeley-chen/vsb-server/internal/model/permission"
	"github.com/seeley-chen/vsb-server/internal/service/permission/department"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

type DepartmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req DepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid request request")
		return
	}

	// 验证请求参数
	if req.Name == "" {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "account is required")
		return
	}

	newDepartment, err := h.svc.Create(r.Context(), model.DepartmentRequest{
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		if errors.Is(err, department.ErrDepartmentExists) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "department already exists")
		}

		zap.L().Error("Create failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "internal server error")

		return
	}

	response.Success(w, newDepartment)
}
