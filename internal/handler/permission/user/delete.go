package user

import (
	"errors"
	"net/http"

	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := response.PathVar(w, r, "id", "user id is required")
	if !ok {
		return
	}

	if err := h.svc.DeleteUser(r.Context(), userID); err != nil {
		if errors.Is(err, svc.ErrUserNotFound) {
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, "user not found")
			return
		}
		zap.L().Error("delete user failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(w, nil)
}
