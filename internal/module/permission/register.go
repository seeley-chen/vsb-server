package permission

import (
	"github.com/gorilla/mux"

	"github.com/Vanselyn/vsb-server/internal/deps"
	"github.com/Vanselyn/vsb-server/internal/domain/permission/department"
	"github.com/Vanselyn/vsb-server/internal/domain/permission/role"
	"github.com/Vanselyn/vsb-server/internal/domain/permission/user"
)

// Register 组装 permission 大模块依赖并注册路由（均需鉴权）。
func Register(d *deps.Deps, _ *mux.Router, protected *mux.Router) {
	departmentRepo := department.NewDepartmentRepo(d.DB)
	roleRepo := role.NewRoleRepo(d.DB)
	userRepo := user.NewUserRepo(d.DB)
	d.EnsureIndexes(departmentRepo.EnsureIndexes)
	d.EnsureIndexes(roleRepo.EnsureIndexes)

	departmentSvc := department.NewDepartmentService(departmentRepo)
	roleSvc := role.NewRoleService(roleRepo)
	userSvc := user.NewUserService(userRepo)

	sub := protected.PathPrefix("/permission").Subrouter()
	department.NewHandler(departmentSvc).RegisterRoutes(sub)
	roleHdl := role.NewHandler(roleSvc, userSvc)
	roleHdl.RegisterRoutes(sub)
	user.NewHandler(userSvc).RegisterRoutes(sub)
	sub.HandleFunc("/privileges", roleHdl.GetPrivileges).Methods("GET")
}
