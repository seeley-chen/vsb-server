package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	MongoURI      string
	MongoDB       string
	JWTSecret     string
	JWTExpiration string // 原始字符串，如 "24h"
}

func Load() *Config {
	_ = godotenv.Load() // 忽略错误，允许从系统环境变量读取

	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		MongoURI:      getEnv("MONGODB_URI", ""),
		MongoDB:       getEnv("MONGODB_DB", "vsb"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),
	}
}

// GetJWTExpireDuration 将字符串解析为 time.Duration
func (c *Config) GetJWTExpireDuration() time.Duration {
	d, err := time.ParseDuration(c.JWTExpiration)
	if err != nil {
		return 24 * time.Hour // 解析失败默认 24h
	}
	return d
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
