package router

import (
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/seeley-chen/vsb-server/config"
	userHandler "github.com/seeley-chen/vsb-server/internal/handler/user"
	"github.com/seeley-chen/vsb-server/internal/middleware"
	repo "github.com/seeley-chen/vsb-server/internal/repository/user"
	svc "github.com/seeley-chen/vsb-server/internal/service/user"
)

// New 创建并配置路由器
func New(cfg *config.Config, db *mongo.Database, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	// 全局中间件
	r.Use(middleware.CORS)
	r.Use(middleware.Logger(logger))

	// 初始化用户模块
	userRepo := repo.NewRepository(db)
	userSvc := svc.NewService(userRepo, cfg.JWTSecret, cfg.GetJWTExpireDuration())
	userHdl := userHandler.NewHandler(userSvc)

	// 注册公开路由（无需鉴权）
	userHdl.RegisterPublicRoutes(r)

	// 注册受保护路由（需要 JWT 鉴权）
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.Auth(cfg.JWTSecret))
	userHdl.RegisterProtectedRoutes(protected)

	return r
}
