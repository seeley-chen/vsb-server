package module

import (
	"github.com/gorilla/mux"

	"github.com/Vanselyn/vsb-server/internal/deps"
	"github.com/Vanselyn/vsb-server/internal/module/login"
	"github.com/Vanselyn/vsb-server/internal/module/permission"
)

// Registrar 大模块注册函数：完成依赖注入并挂载路由。
// 每个大模块（permission、login、announcement 等）在 internal/module/<name>/register.go 中实现。
type Registrar func(d *deps.Deps, public, protected *mux.Router)

// All 所有大模块注册表，新增模块只需追加一行。
var All = []Registrar{
	permission.Register,
	login.Register,
}
