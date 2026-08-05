package user

import (
	"github.com/gorilla/mux"

	"github.com/seeley-chen/vsb-server/internal/deps"
	userHandler "github.com/seeley-chen/vsb-server/internal/handler/permission/user"
	userRepo "github.com/seeley-chen/vsb-server/internal/repository/permission/user"
	userSvc "github.com/seeley-chen/vsb-server/internal/service/permission/user"
)

// Register 初始化用户模块并注册路由
func Register(d *deps.Deps, public, protected *mux.Router) {
	repo := userRepo.NewRepository(d.DB)
	d.EnsureIndexes(repo.EnsureIndexes)

	svc := userSvc.NewService(repo, d.JWTSecret, d.JWTExpire)
	hdl := userHandler.NewHandler(svc)

	hdl.RegisterPublicRoutes(public)
	hdl.RegisterProtectedRoutes(protected)
}
