package product

import (
	"github.com/Vanselyn/vsb-server/internal/deps"
	"github.com/Vanselyn/vsb-server/internal/domain/product/category"
	"github.com/gorilla/mux"
)

func Register(d *deps.Deps, _ *mux.Router, protected *mux.Router) {
	categoryRepo := category.NewCategoryRepo(d.DB)

	d.EnsureIndexes(categoryRepo.EnsureIndexes)

	categorySvc := category.NewCategoryService(categoryRepo)

	sub := protected.PathPrefix("/product").Subrouter()
	category.NewCategoryHandler(categorySvc).RegisterRoutes(sub)
}
