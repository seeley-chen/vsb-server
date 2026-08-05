package module

import (
	"github.com/gorilla/mux"

	"github.com/seeley-chen/vsb-server/internal/deps"
	"github.com/seeley-chen/vsb-server/internal/module/department"
	"github.com/seeley-chen/vsb-server/internal/module/user"
)

// Registrar 模块注册函数：完成依赖注入并挂载路由
type Registrar func(d *deps.Deps, public, protected *mux.Router)

// All 所有模块注册表，新增模块只需追加一行
var All = []Registrar{
	user.Register,
	department.Register,
}
