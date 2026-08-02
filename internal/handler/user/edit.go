package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/seeley-chen/vsb-server/internal/model"
	"github.com/seeley-chen/vsb-server/internal/service/user"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

type EditRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Account      string `json:"account"`
	Phone        string `json:"phone"`
	PasswordHash string `json:"passwordHash"`
}

// EditRequest 编辑用户请求
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
		Username:     req.Username,
		Email:        req.Email,
		Account:      req.Account,
		Phone:        req.Phone,
		PasswordHash: req.PasswordHash,
	})

	if err != nil {
		switch err {
		case user.ErrUserNotFound:
			response.Fail(w, http.StatusNotFound, response.CodeNotFound, err.Error())
		default:
			zap.L().Error("edit user failed", zap.Error(err))
			response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
		}
	}

	response.Success(w, nil)
}
