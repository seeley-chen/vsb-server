package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string // 服务器端口
	MongoURI      string // MongoDB URI
	MongoDB       string // MongoDB 数据库名称
	JWTSecret     string // JWT 密钥
	JWTExpiration string // JWT 过期时间
}

func Load() *Config {
	// 加载 .env 文件（忽略错误，如果不存在则使用系统环境变量）
	_ = godotenv.Load()

	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		MongoURI:      getEnv("MONGODB_URI", ""),
		MongoDB:       getEnv("MONGODB_DB", "vsb"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
