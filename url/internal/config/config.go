package url_config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI        string
	MongoDB         string
	RedisAddr       string
	UserServiceAddr string
	GRPCPort        string
	ShortURLBase    string
	KafkaBrokers    []string
	MetricsPort     string
	OtelExporter    string
	OtelService     string
	LogLevel        string
}

func Load() (*Config, error) {
	if err := godotenv.Load("url/.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	grpcPort := os.Getenv("URL_SERVICE_PORT")
	shortURLBase := os.Getenv("SHORT_URL_BASE")
	kafkaBrokersRaw := os.Getenv("KAFKA_BROKERS")
	metricsPort := os.Getenv("METRICS_PORT")
	otelExporter := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	otelService := os.Getenv("OTEL_SERVICE_NAME")
	logLevel := os.Getenv("LOG_LEVEL")

	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI is required")
	}
	if mongoDB == "" {
		return nil, fmt.Errorf("MONGO_DB is required")
	}
	if redisHost == "" {
		return nil, fmt.Errorf("REDIS_HOST is required")
	}
	if redisPort == "" {
		return nil, fmt.Errorf("REDIS_PORT is required")
	}
	if userServiceAddr == "" {
		return nil, fmt.Errorf("USER_SERVICE_ADDR is required")
	}
	if grpcPort == "" {
		return nil, fmt.Errorf("URL_SERVICE_PORT is required")
	}
	if shortURLBase == "" {
		return nil, fmt.Errorf("SHORT_URL_BASE is required")
	}
	if kafkaBrokersRaw == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS is required")
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
		MongoURI:        mongoURI,
		MongoDB:         mongoDB,
		RedisAddr:       redisHost + ":" + redisPort,
		UserServiceAddr: userServiceAddr,
		GRPCPort:        grpcPort,
		ShortURLBase:    shortURLBase,
		KafkaBrokers:    strings.Split(kafkaBrokersRaw, ","),
		MetricsPort:     metricsPort,
		OtelExporter:    otelExporter,
		OtelService:     otelService,
		LogLevel:        logLevel,
	}, nil
}
