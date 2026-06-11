package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv           string
	AppPort          string
	PlatformDomain   string
	DatabaseURL      string
	RedisURL         string
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	StoragePath      string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env")

	accessExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_EXPIRY: %w", err)
	}

	refreshExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRY: %w", err)
	}

	cfg := &Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		AppPort:          getEnv("APP_PORT", "8080"),
		PlatformDomain:   getEnv("PLATFORM_DOMAIN", "localhost"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://owncommerce:owncommerce@localhost:5432/owncommerce?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:        getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTAccessExpiry:  accessExpiry,
		JWTRefreshExpiry: refreshExpiry,
		StoragePath:      getEnv("STORAGE_PATH", "../../storage"),
	}

	if cfg.JWTSecret == "dev-secret-change-me" && cfg.AppEnv == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
