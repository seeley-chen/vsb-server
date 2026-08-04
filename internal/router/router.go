package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	_ "github.com/seeley-chen/vsb-server/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/seeley-chen/vsb-server/config"
	departmentHandler "github.com/seeley-chen/vsb-server/internal/handler/permission/department"
	userHandler "github.com/seeley-chen/vsb-server/internal/handler/permission/user"
	"github.com/seeley-chen/vsb-server/internal/middleware"
	departmentRepo "github.com/seeley-chen/vsb-server/internal/repository/permission/department"
	userRepo "github.com/seeley-chen/vsb-server/internal/repository/permission/user"
	userSvc "github.com/seeley-chen/vsb-server/internal/service/permission/user"

	deptSvc "github.com/seeley-chen/vsb-server/internal/service/permission/department"
)

// New 创建并配置路由器
func New(cfg *config.Config, db *mongo.Database, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.Recover(logger))
	r.Use(middleware.MaxBodySize(1 << 20)) // 1MB
	r.Use(middleware.Logger(logger))

	// 健康检查
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods("GET")

	// 初始化用户模块
	userRepo := userRepo.NewRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := userRepo.EnsureIndexes(ctx); err != nil {
		logger.Warn("failed to ensure indexes", zap.Error(err))
	}

	userSvc := userSvc.NewService(userRepo, cfg.JWTSecret, cfg.GetJWTExpireDuration())
	userHdl := userHandler.NewHandler(userSvc)

	departmentRepo := departmentRepo.NewRepository(db)

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := departmentRepo.EnsureIndexes(ctx); err != nil {
		logger.Warn("failed to ensure indexes", zap.Error(err))
	}

	departmentSvc := deptSvc.NewService(departmentRepo, cfg.JWTSecret, cfg.GetJWTExpireDuration())
	departmentHdl := departmentHandler.NewHandler(departmentSvc)
	departmentHdl.RegisterPublicRoutes(r)

	// 注册公开路由（无需鉴权）
	userHdl.RegisterPublicRoutes(r)

	// 注册受保护路由（需要 JWT 鉴权）
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.Auth(cfg.JWTSecret))
	userHdl.RegisterProtectedRoutes(protected)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
