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

	"github.com/seeley-chen/vsb-server/config"
	"github.com/seeley-chen/vsb-server/internal/database"
	"github.com/seeley-chen/vsb-server/internal/router"
	"github.com/seeley-chen/vsb-server/pkg/logger"
)

func main() {
	// 1. 加载配置
	cfg := config.Load()

	// 2. 初始化日志
	logger.Init("debug") // 开发阶段用 debug，生产可改为 info
	defer logger.Sync()

	log.Printf("config loaded: port=%s, db=%s", cfg.ServerPort, cfg.MongoDB)

	// 3. 连接 MongoDB
	if err := database.Connect(cfg.MongoURI, cfg.MongoDB); err != nil {
		log.Fatalf("failed to connect mongodb: %v", err)
	}
	defer database.Disconnect()
	log.Println("mongodb connected")

	// 4. 创建路由器
	r := router.New(cfg, database.GetDB(), logger.Log)

	// 5. 启动 HTTP 服务（优雅关停）
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
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
