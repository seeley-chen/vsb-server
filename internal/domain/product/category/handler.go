package category

import (
	"errors"
	"net/http"

	"github.com/Vanselyn/vsb-server/internal/middleware"
	"github.com/Vanselyn/vsb-server/pkg/response"
	"github.com/gorilla/mux"
)

type Handler struct {
	svc *CategoryService
}

func NewCategoryHandler(svc *CategoryService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	sub := r.PathPrefix("/category").Subrouter()
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
	case errors.Is(err, ErrCategoryExists):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "category exists")

	case errors.Is(err, ErrCategoryNotFound):
		response.Fail(w, http.StatusBadRequest, response.CodeBadRequest, "category not found")

	default:
		middleware.LogError(r, op, err)
		response.Fail(w, http.StatusInternalServerError, response.CodeInternalError)
	}

}

// Create 创建商品分类
// @Summary 创建商品分类
// @Description 创建新商品分类，需要 Bearer Token
// @Tags 商品分类
// @Accept json
// @Produce json
// @Param request body CategoryRequest true "分类信息"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=CategoryResponse} "创建成功"
// @Failure 400 {object} response.Response "参数错误或分类已存在"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/product/category/create [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CategoryRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	category, err := h.svc.CreateCategory(r.Context(), &req)
	if err != nil {
		h.handleErr(w, r, err, "CreateCategory")
		return
	}
	response.Success(w, category)
}

// List 获取商品分类列表
// @Summary 获取商品分类列表
// @Description 获取全部商品分类，需要 Bearer Token
// @Tags 商品分类
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=response.PageData{list=[]CategoryResponse}} "分类列表"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/product/category/list [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	categories, total, err := h.svc.CategoryList(r.Context())
	if err != nil {
		h.handleErr(w, r, err, "CategoryList")
		return
	}
	response.Success(w, response.PageData{
		List:  categories,
		Total: total,
	})
}

// Update 更新商品分类
// @Summary 更新商品分类
// @Description 根据分类 ID 更新商品分类，空字段表示不修改，需要 Bearer Token
// @Tags 商品分类
// @Accept json
// @Produce json
// @Param id path string true "分类 ID"
// @Param request body CategoryUpdateRequest true "更新信息"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=CategoryResponse} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/product/category/update/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req CategoryUpdateRequest
	if !response.BindJSON(w, r, &req) {
		return
	}

	categoryID, ok := response.PathVar(w, r, "id", "category id is required")
	if !ok {
		return
	}

	category, err := h.svc.UpdateCategory(r.Context(), categoryID, &req)
	if err != nil {
		h.handleErr(w, r, err, "update category failed")
		return
	}

	response.Success(w, category)
}

// Delete 删除商品分类
// @Summary 删除商品分类
// @Description 根据分类 ID 删除商品分类，需要 Bearer Token
// @Tags 商品分类
// @Produce json
// @Param id path string true "分类 ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/product/category/delete/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	categoryID, ok := response.PathVar(w, r, "id", "category id is required")
	if !ok {
		return
	}

	err := h.svc.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		h.handleErr(w, r, err, "Delete category failed")
		return
	}

	response.Success(w, nil)
}
