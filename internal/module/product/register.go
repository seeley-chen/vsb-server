package product

import (
	"github.com/Vanselyn/vsb-server/internal/deps"
	"github.com/Vanselyn/vsb-server/internal/domain/product/category"
	"github.com/Vanselyn/vsb-server/internal/domain/product/goods"
	"github.com/gorilla/mux"
)

func Register(d *deps.Deps, _ *mux.Router, protected *mux.Router) {
	categoryRepo := category.NewCategoryRepo(d.DB)
	goodsRepo := goods.NewGoodsRepo(d.DB)

	d.EnsureIndexes(categoryRepo.EnsureIndexes)
	d.EnsureIndexes(goodsRepo.EnsureIndexes)

	categorySvc := category.NewCategoryService(categoryRepo)
	goodSvc := goods.NewGoodsService(goodsRepo)

	sub := protected.PathPrefix("/product").Subrouter()
	category.NewCategoryHandler(categorySvc).RegisterRoutes(sub)
	goods.NewGoodsHandler(goodSvc).RegisterRoutes(sub)
}
