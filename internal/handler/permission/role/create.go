package role

import (
	"encoding/json"
	"errors"
	"net/http"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.RoleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "请求体解析失败")
		return
	}

	role, err := h.svc.CreateRole(r.Context(), &req)
	if err != nil {
		if errors.Is(err, svc.ErrRoleNameEmpty) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "角色名不能为空")
			return
		}

		if errors.Is(err, svc.ErrRoleExist) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "角色名已存在")
			return
		}

		zap.L().Error("create role failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "创建角色失败")
		return
	}

	response.Success(w, role)
}
