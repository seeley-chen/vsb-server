// @title VSB Server API
// @version 1.0
// @description VSB 服务端 API 文档
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer JWT token，格式：Bearer {token}
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vanselyn/vsb-server/config"
	"github.com/Vanselyn/vsb-server/internal/database"
	"github.com/Vanselyn/vsb-server/internal/middleware"
	"github.com/Vanselyn/vsb-server/internal/router"
	"github.com/Vanselyn/vsb-server/pkg/logger"
)

func main() {
	// 1. 加载配置
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	// 2. 初始化日志
	logger.Init(cfg.LogLevel)
	defer logger.Sync()

	log.Printf("config loaded: port=%s, db=%s", cfg.ServerPort, cfg.MongoDB)

	// 3. 连接 MongoDB
	client, db, err := database.Connect(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("failed to connect mongodb: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = client.Disconnect(ctx)
	}()
	log.Println("mongodb connected")

	// 4. 创建路由器
	r := router.New(cfg, db, logger.Log)

	// 5. 启动 HTTP 服务（优雅关停）
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	server := &http.Server{
		Addr:              addr,
		Handler:           middleware.CORS(cfg.CORSAllowedOrigins)(r),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 vsb-server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-done
	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	log.Println("bye 👋")
}
