package role

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	roleId := mux.Vars(r)["id"]

	if roleId == "" {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "角色ID必传")
		return
	}

	if err := h.svc.DeleteRole(r.Context(), roleId); err != nil {
		if errors.Is(err, svc.ErrRoleNotFound) {
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, "未找到角色")
			return
		}

		zap.L().Error("delete role failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(w, nil)
}
