package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/seeley-chen/vsb-server/internal/service/user"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

// DeleteRequest 删除用户请求
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// 从url当中获取userId,例如/api/user/delete/1234567890
	vars := mux.Vars(r)

	userId := vars["userId"]

	if userId == "" {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user_id is required")
		return
	}

	err := h.svc.DeleteUserById(r.Context(), userId)
	if err != nil {
		switch err {
		case user.ErrUserNotFound:
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, err.Error())
		default:
			zap.L().Error("delete user failed", zap.Error(err))
			response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
		}
		return
	}

	response.Success(w, nil)
}
