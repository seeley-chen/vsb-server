package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	_ "github.com/Vanselyn/vsb-server/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Vanselyn/vsb-server/config"
	"github.com/Vanselyn/vsb-server/internal/deps"
	"github.com/Vanselyn/vsb-server/internal/middleware"
	"github.com/Vanselyn/vsb-server/internal/module"
)

// New 创建并配置路由器
func New(cfg *config.Config, db *mongo.Database, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.Recover(logger))
	r.Use(middleware.MaxBodySize(1 << 20)) // 1MB
	r.Use(middleware.Logger(logger, cfg.LogBody))

	// 健康检查
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods("GET")

	d := &deps.Deps{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
		JWTExpire: cfg.GetJWTExpireDuration(),
		Logger:    logger,
	}

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.Auth(cfg.JWTSecret))

	for _, register := range module.All {
		register(d, r, protected)
	}

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
