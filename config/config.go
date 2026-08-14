package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	MongoURI           string
	MongoDB            string
	JWTSecret          string
	JWTExpiration      string // 原始字符串，如 "24h"
	LogLevel           string
	CORSAllowedOrigins []string
}

func Load() *Config {
	_ = godotenv.Load() // 忽略错误，允许从系统环境变量读取

	return &Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		MongoURI:           getEnv("MONGODB_URI", ""),
		MongoDB:            getEnv("MONGODB_DB", "vsb"),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpiration:      getEnv("JWT_EXPIRATION", "24h"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		CORSAllowedOrigins: parseCSV(getEnv("CORS_ALLOWED_ORIGINS", "")),
	}
}

// Validate 校验必填配置项
func (c *Config) Validate() error {
	if c.MongoURI == "" {
		return errors.New("MONGODB_URI is required")
	}
	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}
	if _, err := time.ParseDuration(c.JWTExpiration); err != nil {
		return fmt.Errorf("invalid JWT_EXPIRATION %q: %w", c.JWTExpiration, err)
	}
	return nil
}

// GetJWTExpireDuration 将字符串解析为 time.Duration（Validate 已保证合法性）
func (c *Config) GetJWTExpireDuration() time.Duration {
	d, _ := time.ParseDuration(c.JWTExpiration)
	return d
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	for _, part := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
