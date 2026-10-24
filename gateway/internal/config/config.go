package gateway_config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GatewayPort          string
	UserServiceAddr      string
	URLServiceAddr       string
	AnalyticsServiceAddr string
	MetricsPort          string
	OtelExporter         string
	OtelService          string
	LogLevel             string
}

func Load() (*Config, error) {
	if err := godotenv.Load("gateway/.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	gatewayPort := os.Getenv("GATEWAY_PORT")
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	urlServiceAddr := os.Getenv("URL_SERVICE_ADDR")
	analyticsServiceAddr := os.Getenv("ANALYTICS_SERVICE_ADDR")
	metricsPort := os.Getenv("METRICS_PORT")
	otelExporter := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	otelService := os.Getenv("OTEL_SERVICE_NAME")
	logLevel := os.Getenv("LOG_LEVEL")

	if gatewayPort == "" {
		return nil, fmt.Errorf("GATEWAY_PORT is required")
	}
	if userServiceAddr == "" {
		return nil, fmt.Errorf("USER_SERVICE_ADDR is required")
	}
	if urlServiceAddr == "" {
		return nil, fmt.Errorf("URL_SERVICE_ADDR is required")
	}
	if analyticsServiceAddr == "" {
		return nil, fmt.Errorf("ANALYTICS_SERVICE_ADDR is required")
	}
	if metricsPort == "" {
		return nil, fmt.Errorf("METRICS_PORT is required")
	}
	if otelService == "" {
		return nil, fmt.Errorf("OTEL_SERVICE_NAME is required")
	}
	if otelExporter == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT is required")
	}
	if logLevel == "" {
		return nil, fmt.Errorf("LOG_LEVEL is required")
	}

	return &Config{
		GatewayPort:          gatewayPort,
		UserServiceAddr:      userServiceAddr,
		URLServiceAddr:       urlServiceAddr,
		AnalyticsServiceAddr: analyticsServiceAddr,
		MetricsPort:          metricsPort,
		OtelExporter:         otelExporter,
		OtelService:          otelService,
		LogLevel:             logLevel,
	}, nil
}
