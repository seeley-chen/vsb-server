package department

import (
	"net/http"
	"strconv"

	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageIndexStr := r.URL.Query().Get("pageIndex")
	pageSizeStr := r.URL.Query().Get("pageSize")

	pageIndex := 1
	pageSize := 20

	if p, err := strconv.Atoi(pageIndexStr); err == nil && p > 0 {
		pageIndex = p
	}
	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
		pageSize = ps
	}

	departments, total, err := h.svc.ListDepartments(r.Context(), pageIndex, pageSize)
	if err != nil {
		zap.L().Error("list departments failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
		return
	}

	data := response.PageData{
		List:      departments,
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}

	response.Success(w, data)
}
