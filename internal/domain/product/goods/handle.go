package goods

import (
	"errors"
	"net/http"

	"github.com/Vanselyn/vsb-server/internal/middleware"
	"github.com/Vanselyn/vsb-server/internal/models"
	"github.com/Vanselyn/vsb-server/pkg/i18n"
	"github.com/Vanselyn/vsb-server/pkg/response"
	"github.com/gorilla/mux"
)

type Handler struct {
	svc *GoodsService
}

func NewGoodsHandler(svc *GoodsService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	sub := r.PathPrefix("/goods").Subrouter()
	sub.HandleFunc("/create", h.Create).Methods("POST")
	sub.HandleFunc("/list", h.List).Methods("GET")
	sub.HandleFunc("/update/{id}", h.Update).Methods("PUT")
	sub.HandleFunc("/delete/{id}", h.Delete).Methods("DELETE")
}

func (h *Handler) handleErr(w http.ResponseWriter, r *http.Request, err error, op string) {
	if response.FailIfValidation(w, err) {
		return
	}

	switch {
	case errors.Is(err, ErrGoodsNotFound):
		response.Fail(w, http.StatusNotFound, response.CodeNotFound, "goods not found")

	case errors.Is(err, ErrGoodsSkuAlreadyExists):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "goods sku already exists")

	default:
		middleware.LogError(r, op, err)
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
	}

}

// Create 创建商品
// @Summary 创建商品
// @Description 创建新商品，需要 Bearer Token
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param request body GoodsRequest true "商品信息"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=GoodsResponse} "创建成功"
// @Failure 400 {object} response.Response "参数错误或 SKU 已存在"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/product/goods/create [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req GoodsRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	goods, err := h.svc.CreateGoods(r.Context(), &req)
	if err != nil {
		h.handleErr(w, r, err, "create goods failed")
		return
	}

	response.Success(w, goods)

}

// List 获取商品列表
// @Summary 获取商品列表
// @Description 分页获取商品列表，支持 SKU 模糊搜索、商品名称模糊搜索、分类 ID 精确匹配
// @Tags 商品管理
// @Produce json
// @Param pageIndex query int false "页码，默认 1"
// @Param pageSize query int false "每页数量，默认 20，上限 100"
// @Param sku query string false "SKU（模糊搜索）"
// @Param name query string false "商品名称（模糊搜索）"
// @Param categoryId query string false "分类 ID（精确匹配）"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=response.PageData} "商品列表"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/product/goods/list [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageIndex, pageSize := response.PageQuery(r)
	q := r.URL.Query()
	params := GoodsQueryParams{
		Sku:        q.Get("sku"),
		Name:       i18n.Locale{i18n.ZhCN: q.Get("name")},
		CategoryID: q.Get("categoryId"),
		Status:     models.StatusEnum(q.Get("status")),
	}

	goodsList, total, err := h.svc.GetGoodsList(r.Context(), pageIndex, pageSize, params)
	if err != nil {
		h.handleErr(w, r, err, "get goods list failed")
		return
	}

	response.Success(w, response.PageData{
		List:      goodsList,
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
}

// Update 更新商品
// @Summary 更新商品
// @Description 根据商品 ID 更新商品信息，空字段表示不修改，需要 Bearer Token
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param id path string true "商品 ID"
// @Param request body GoodsUpdateRequest true "商品信息"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=GoodsResponse} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/product/goods/update/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req GoodsUpdateRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	goodsID, ok := response.PathVar(w, r, "id", "goods id is required")
	if !ok {
		return
	}

	goods, err := h.svc.UpdateGoods(r.Context(), goodsID, &req)
	if err != nil {
		h.handleErr(w, r, err, "update goods failed")
		return
	}

	response.Success(w, goods)

}

// Delete 删除商品
// @Summary 删除商品
// @Description 根据商品 ID 删除商品，需要 Bearer Token
// @Tags 商品管理
// @Produce json
// @Param id path string true "商品 ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	goodsID, ok := response.PathVar(w, r, "id", "goods id is required")
	if !ok {
		return
	}

	err := h.svc.DeleteGoods(r.Context(), goodsID)
	if err != nil {
		h.handleErr(w, r, err, "delete goods failed")
		return
	}

	response.Success(w, nil)
}
