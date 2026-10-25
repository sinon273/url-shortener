package user_config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	RedisAddr       string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	GRPCPort        string
	MetricsPort     string
	OtelExporter    string
	OtelService     string
	LogLevel        string
}

func Load() (*Config, error) {
	if err := godotenv.Load("user/.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	cfg := &Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		RedisAddr:    redisHost + ":" + redisPort,
		JWTSecret:    os.Getenv("JWT_SECRET"),
		GRPCPort:     os.Getenv("USER_SERVICE_PORT"),
		MetricsPort:  os.Getenv("METRICS_PORT"),
		OtelExporter: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OtelService:  os.Getenv("OTEL_SERVICE_NAME"),
		LogLevel:     os.Getenv("LOG_LEVEL"),
	}

	if cfg.MetricsPort == "" {
		return nil, fmt.Errorf("METRICS_PORT is required")
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.GRPCPort == "" {
		return nil, fmt.Errorf("USER_SERVICE_PORT is required")
	}
	if cfg.OtelExporter == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT is required")
	}
	if cfg.OtelService == "" {
		return nil, fmt.Errorf("OTEL_SERVICE_NAME is required")
	}
	if cfg.LogLevel == "" {
		return nil, fmt.Errorf("LOG_LEVEL is required")
	}
	if redisHost == "" {
		return nil, fmt.Errorf("REDIS_HOST is required")
	}
	if redisPort == "" {
		return nil, fmt.Errorf("REDIS_PORT is required")
	}

	var err error
	cfg.AccessTokenTTL, err = time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		return nil, fmt.Errorf("parse ACCESS_TOKEN_TTL: %w", err)
	}

	cfg.RefreshTokenTTL, err = time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, fmt.Errorf("parse REFRESH_TOKEN_TTL: %w", err)
	}

	return cfg, nil
}
