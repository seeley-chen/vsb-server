package user

import (
	"net/http"

	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageIndex, pageSize := response.PageQuery(r)
	q := r.URL.Query()
	username := q.Get("username")
	account := q.Get("account")
	email := q.Get("email")

	users, total, err := h.svc.GetUserList(r.Context(), pageIndex, pageSize, username, account, email)
	if err != nil {
		zap.L().Error("list users failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(w, response.PageData{
		List:      users,
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
}
