package department

import (
	"net/http"

	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageIndex, pageSize := response.PageQuery(r)

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
