// Package config 提供应用配置管理
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置结构体
type Config struct {
	AppPort        string
	AppMode        string
	JWTSecret      string
	JWTExpireHours int
	DBPath         string
	AdminUsername  string
	AdminPassword  string
}

// Global 全局配置实例
var Global *Config

// Load 加载配置 (优先从 .env 文件, 其次从环境变量, 最后使用默认值)
func Load() *Config {
	// 尝试加载 .env 文件, 不存在也不报错
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		AppMode:        getEnv("APP_MODE", "debug"),
		JWTSecret:      getEnv("JWT_SECRET", "go-cargo-default-secret-key-2024"),
		JWTExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 72),
		DBPath:         getEnv("DB_PATH", "./data/cargo.db"),
		AdminUsername:  getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", "admin123"),
	}

	Global = cfg
	return cfg
}

// getEnv 获取环境变量，支持默认值
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvInt 获取整数类型的环境变量
func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
