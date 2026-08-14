package login

import (
	"github.com/gorilla/mux"

	"github.com/Vanselyn/vsb-server/internal/deps"
	loginhttp "github.com/Vanselyn/vsb-server/internal/domain/login"
	"github.com/Vanselyn/vsb-server/internal/domain/permission/user"
)

// Register 组装 login 模块依赖并注册路由。
// 登录为公开路由，无需鉴权，路径为 POST /api/login。
func Register(d *deps.Deps, public, _ *mux.Router) {
	userRepo := user.NewUserRepo(d.DB)
	svc := loginhttp.NewLoginService(userRepo, d.JWTSecret, d.JWTExpire)
	loginhttp.NewHandler(svc).RegisterRoutes(public.PathPrefix("/api").Subrouter())
}
