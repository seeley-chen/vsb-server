package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client // 客户端
var DB *mongo.Database   // 数据库

func Connect(uri, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 设置超时时间
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri)) // 连接数据库
	if err != nil {
		return err // 如果连接失败，则返回错误
	}

	// Ping 测试连接
	if err = client.Ping(ctx, nil); err != nil {
		return err // 如果ping失败，则返回错误
	}

	Client = client              // 设置客户端
	DB = client.Database(dbName) // 设置数据库
	return nil                   // 返回nil
}

func Disconnect() {
	// 断开连接
	if Client != nil {
		_ = Client.Disconnect(context.Background()) // 断开连接
	}
}

// 获取数据库
func GetDB() *mongo.Database {
	return DB // 获取数据库
}
