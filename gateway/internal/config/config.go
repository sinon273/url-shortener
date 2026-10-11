package gateway_config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GatewayPort     string
	UserServiceAddr string
	URLServiceAddr  string
}

func Load() (*Config, error) {
	if err := godotenv.Load("gateway/.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	gatewayPort := os.Getenv("GATEWAY_PORT")
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	urlServiceAddr := os.Getenv("URL_SERVICE_ADDR")

	if gatewayPort == "" {
		return nil, fmt.Errorf("GATEWAY_PORT is required")
	}
	if userServiceAddr == "" {
		return nil, fmt.Errorf("USER_SERVICE_ADDR is required")
	}
	if urlServiceAddr == "" {
		return nil, fmt.Errorf("URL_SERVICE_ADDR is required")
	}

	return &Config{
		GatewayPort:     gatewayPort,
		UserServiceAddr: userServiceAddr,
		URLServiceAddr:  urlServiceAddr,
	}, nil
}
