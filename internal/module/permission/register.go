package permission

import (
	"github.com/gorilla/mux"

	"github.com/seeley-chen/vsb-server/internal/deps"
	permissionHandler "github.com/seeley-chen/vsb-server/internal/handler/permission"
	"github.com/seeley-chen/vsb-server/internal/handler/permission/department"
	"github.com/seeley-chen/vsb-server/internal/handler/permission/role"
	"github.com/seeley-chen/vsb-server/internal/handler/permission/user"
	permissionRepo "github.com/seeley-chen/vsb-server/internal/repository/permission"
	permissionSvc "github.com/seeley-chen/vsb-server/internal/service/permission"
)

// Register 组装 permission 大模块依赖并注册路由。
func Register(d *deps.Deps, _ *mux.Router, protected *mux.Router) {
	departmentRepo := permissionRepo.NewDepartmentRepo(d.DB)
	roleRepo := permissionRepo.NewRoleRepo(d.DB)
	userRepo := permissionRepo.NewUserRepo(d.DB)
	d.EnsureIndexes(departmentRepo.EnsureIndexes)
	d.EnsureIndexes(roleRepo.EnsureIndexes)

	departmentSvc := permissionSvc.NewDepartmentService(departmentRepo)
	roleSvc := permissionSvc.NewRoleService(roleRepo)
	userSvc := permissionSvc.NewUserService(userRepo, d.JWTSecret, d.JWTExpire)

	hdl := permissionHandler.NewHandler(
		department.NewHandler(departmentSvc),
		role.NewHandler(roleSvc, userSvc),
		user.NewHandler(userSvc),
	)
	hdl.RegisterRoutes(protected)
}
