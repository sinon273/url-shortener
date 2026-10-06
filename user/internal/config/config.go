package user_config

import (
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
}

func Load() (*Config, error) {
	_ = godotenv.Load("user/.env")

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisAddr:   redisHost + ":" + redisPort,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		GRPCPort:    os.Getenv("USER_SERVICE_PORT"),
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
