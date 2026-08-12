package user

import (
	"errors"
	"net/http"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	svc "github.com/seeley-chen/vsb-server/internal/service/permission"
	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.UserLoginRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	if req.Account == "" || req.Password == "" {
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "account and password are required")
		return
	}

	login, err := h.svc.Login(r.Context(), &req)
	if err != nil {
		if errors.Is(err, svc.ErrInvalidAccount) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid account")
			return
		}

		if errors.Is(err, svc.ErrInvalidPassword) {
			response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "invalid password")
			return
		}

		zap.L().Error("login failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(w, login)
}
