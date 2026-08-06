package department

import (
	"net/http"
	"strconv"

	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageIndex := 1
	pageSize := 20

	if p, err := strconv.Atoi(r.URL.Query().Get("pageIndex")); err == nil && p > 0 {
		pageIndex = p
	}
	if ps, err := strconv.Atoi(r.URL.Query().Get("pageSize")); err == nil && ps > 0 && ps <= 100 {
		pageSize = ps
	}

	departments, total, err := h.svc.GetDepartmentList(r.Context(), pageIndex, pageSize)
	if err != nil {
		zap.L().Error("list departments failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(w, response.PageData{
		List:      departments,
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
}
