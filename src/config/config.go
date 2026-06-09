package config

import (
	"os"
	"time"
)

type Config struct {
	Port        string
	JWTSecret   string
	TokenExpiry time.Duration
	AppName     string
	AppVersion  string
	DatabaseURL string
	RedisAddr   string
	RedisPass   string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		TokenExpiry: 24 * time.Hour,
		AppName:     "todo-app",
		AppVersion:  "1.0.0",
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@127.0.0.1:5432/todoapp?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:   getEnv("REDIS_PASS", ""),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
