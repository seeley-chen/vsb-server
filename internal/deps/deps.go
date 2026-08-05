package deps

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Deps 各模块共享的依赖
type Deps struct {
	DB        *mongo.Database
	JWTSecret string
	JWTExpire time.Duration
	Logger    *zap.Logger
}

// EnsureIndexes 在超时上下文中创建索引，失败时仅记录警告
func (d *Deps) EnsureIndexes(fn func(context.Context) error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := fn(ctx); err != nil {
		d.Logger.Warn("failed to ensure indexes", zap.Error(err))
	}
}
