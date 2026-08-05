package department

import (
	"github.com/gorilla/mux"

	"github.com/seeley-chen/vsb-server/internal/deps"
	departmentHandler "github.com/seeley-chen/vsb-server/internal/handler/permission/department"
	departmentRepo "github.com/seeley-chen/vsb-server/internal/repository/permission/department"
	departmentSvc "github.com/seeley-chen/vsb-server/internal/service/permission/department"
)

// Register 初始化部门模块并注册路由
func Register(d *deps.Deps, public, protected *mux.Router) {
	repo := departmentRepo.NewRepository(d.DB)
	d.EnsureIndexes(repo.EnsureIndexes)

	svc := departmentSvc.NewService(repo, d.JWTSecret, d.JWTExpire)
	hdl := departmentHandler.NewHandler(svc)

	hdl.RegisterPublicRoutes(public)
	hdl.RegisterProtectedRoutes(protected)
}
