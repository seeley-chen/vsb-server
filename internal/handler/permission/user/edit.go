package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	model "github.com/seeley-chen/vsb-server/internal/model/permission"
	"github.com/seeley-chen/vsb-server/internal/service/permission/user"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

type EditRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// Edit 编辑用户
// @Summary 编辑用户
// @Description 根据 userId 更新用户信息，需要 Bearer Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param userId path string true "用户 ID"
// @Param request body EditRequest true "更新信息"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/permission/user/edit/{userId} [put]
func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	if userId == "" {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "user_id is required")
		return
	}

	var req EditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	err := h.svc.UpdateUserById(r.Context(), userId, model.User{
		Username: req.Username,
		Email:    req.Email,
		Account:  req.Account,
		Phone:    req.Phone,
	}, req.Password)

	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, err.Error())
			return
		}
		zap.L().Error("edit user failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
		return
	}

	response.Success(w, nil)
}
