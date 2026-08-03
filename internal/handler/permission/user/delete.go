package user

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/seeley-chen/vsb-server/internal/service/permission/user"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

// Delete 删除用户
// @Summary 删除用户
// @Description 根据 userId 删除用户，需要 Bearer Token
// @Tags 用户管理
// @Produce json
// @Param userId path string true "用户 ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/permission/user/delete/{userId} [delete]
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
		if errors.Is(err, user.ErrUserNotFound) {
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, err.Error())
			return
		}
		zap.L().Error("delete user failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
		return
	}

	response.Success(w, nil)
}
