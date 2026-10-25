package analytics_config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBrokers         []string
	ClickHouseAddr       string
	ClickHouseDB         string
	ClickHouseUser       string
	ClickHousePassword   string
	URLServiceAddr       string
	AnalyticsServicePort string
	MetricsPort          string
	OtelExporter         string
	OtelService          string
	LogLevel             string
}

func Load() (*Config, error) {
	if err := godotenv.Load("analytics/.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	kafkaBrokersRaw := os.Getenv("KAFKA_BROKERS")
	clickHouseAddr := os.Getenv("CLICKHOUSE_ADDR")
	clickHouseDb := os.Getenv("CLICKHOUSE_DB")
	clickHouseUser := os.Getenv("CLICKHOUSE_USER")
	clickHousePassword := os.Getenv("CLICKHOUSE_PASSWORD")
	urlServiceAddr := os.Getenv("URL_SERVICE_ADDR")
	analyticsServicePort := os.Getenv("ANALYTICS_SERVICE_PORT")
	metricsPort := os.Getenv("METRICS_PORT")
	otelExporter := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	otelService := os.Getenv("OTEL_SERVICE_NAME")
	logLevel := os.Getenv("LOG_LEVEL")

	if kafkaBrokersRaw == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS is required")
	}
	if clickHouseAddr == "" {
		return nil, fmt.Errorf("CLICKHOUSE_ADDR is required")
	}
	if clickHouseDb == "" {
		return nil, fmt.Errorf("CLICKHOUSE_DB is required")
	}
	if clickHouseUser == "" {
		return nil, fmt.Errorf("CLICKHOUSE_USER is required")
	}
	if urlServiceAddr == "" {
		return nil, fmt.Errorf("URL_SERVICE_ADDR is required")
	}
	if analyticsServicePort == "" {
		return nil, fmt.Errorf("ANALYTICS_SERVICE_PORT is required")
	}
	if metricsPort == "" {
		return nil, fmt.Errorf("METRICS_PORT is required")
	}
	if otelExporter == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT is required")
	}
	if otelService == "" {
		return nil, fmt.Errorf("OTEL_SERVICE_NAME is required")
	}
	if logLevel == "" {
		return nil, fmt.Errorf("LOG_LEVEL is required")
	}

	return &Config{
		KafkaBrokers:         strings.Split(kafkaBrokersRaw, ","),
		ClickHouseAddr:       clickHouseAddr,
		ClickHouseDB:         clickHouseDb,
		ClickHouseUser:       clickHouseUser,
		ClickHousePassword:   clickHousePassword,
		URLServiceAddr:       urlServiceAddr,
		AnalyticsServicePort: analyticsServicePort,
		MetricsPort:          metricsPort,
		OtelService:          otelService,
		OtelExporter:         otelExporter,
		LogLevel:             logLevel,
	}, nil
}
