package permission

import (
	"github.com/gorilla/mux"

	"github.com/seeley-chen/vsb-server/internal/deps"
	permissionHandler "github.com/seeley-chen/vsb-server/internal/handler/permission"
	"github.com/seeley-chen/vsb-server/internal/handler/permission/department"
	permissionRepo "github.com/seeley-chen/vsb-server/internal/repository/permission"
	permissionSvc "github.com/seeley-chen/vsb-server/internal/service/permission"
)

// Register 组装 permission 大模块依赖并注册路由。
func Register(d *deps.Deps, _ *mux.Router, protected *mux.Router) {
	departmentRepo := permissionRepo.NewDepartmentRepo(d.DB)
	departmentSvc := permissionSvc.NewDepartmentService(departmentRepo)

	hdl := permissionHandler.NewHandler(department.NewHandler(departmentSvc))
	hdl.RegisterRoutes(protected)
}
