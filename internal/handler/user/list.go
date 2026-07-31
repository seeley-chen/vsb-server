package user

import (
	"net/http"
	"strconv"

	"github.com/seeley-chen/vsb-server/pkg/response"
	"go.uber.org/zap"
)

// List 获取用户列表（需鉴权）
// @Summary 获取用户列表
// @Description 分页获取用户列表，需要 Bearer Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认20"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=response.PageData} "用户列表"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/user/list [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	pageSize := 20

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
		pageSize = ps
	}

	users, total, err := h.svc.ListUsers(r.Context(), page, pageSize)
	if err != nil {
		zap.L().Error("list users failed", zap.Error(err))
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
		return
	}

	data := response.PageData{
		List:      users,
		Total:     total,
		PageIndex: page,
		PageSize:  pageSize,
	}
	response.Success(w, data)
}
