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

	return &Config{
		KafkaBrokers:         strings.Split(kafkaBrokersRaw, ","),
		ClickHouseAddr:       clickHouseAddr,
		ClickHouseDB:         clickHouseDb,
		ClickHouseUser:       clickHouseUser,
		ClickHousePassword:   clickHousePassword,
		URLServiceAddr:       urlServiceAddr,
		AnalyticsServicePort: analyticsServicePort,
	}, nil
}
