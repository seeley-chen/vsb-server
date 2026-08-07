package role

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
	roleId := mux.Vars(r)["id"]
	if roleId == "" {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "角色ID必传")
		return
	}

	var req model.RoleUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "请求体解析失败")
		return
	}

	role, err := h.svc.UpdateRole(r.Context(), roleId, &req)
	if err != nil {
		if errors.Is(err, svc.ErrRoleNotFound) {
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, "角色不存在")
			return
		}
		zap.L().Error("update role failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "更新角色失败")
		return
	}

	response.Success(w, role)
}
