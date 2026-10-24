package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	analytics_client "github.com/sinon273/url-shortener/analytics/internal/client"
	analytics_config "github.com/sinon273/url-shortener/analytics/internal/config"
	analytics_consumer "github.com/sinon273/url-shortener/analytics/internal/consumer"
	analytics_handler "github.com/sinon273/url-shortener/analytics/internal/handler"
	analytics_logger "github.com/sinon273/url-shortener/analytics/internal/logger"
	analytics_metrics "github.com/sinon273/url-shortener/analytics/internal/metrics"
	analytics_repository "github.com/sinon273/url-shortener/analytics/internal/repository"
	analytics_service "github.com/sinon273/url-shortener/analytics/internal/service"
	analytics_tracing "github.com/sinon273/url-shortener/analytics/internal/tracing"
	analyticsv1 "github.com/sinon273/url-shortener/shared/pkg/proto/analytics/v1"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	batchSize     = 100
	flushInterval = 5 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := analytics_config.Load()
	if err != nil {
		fmt.Printf("loading config: %v\n", err)
		os.Exit(1)
	}

	logger, err := analytics_logger.Init(cfg.LogLevel)
	if err != nil {
		fmt.Printf("init logger: %v\n", err)
		os.Exit(1)
	}

	slog.SetDefault(logger)
	slog.Info("analytics starting", "port", cfg.AnalyticsServicePort, "metrics_port", cfg.MetricsPort)

	shutdown, err := analytics_tracing.Init(ctx, cfg.OtelService, cfg.OtelExporter)
	if err != nil {
		slog.Error("failed to init tracing", "err", err)
		os.Exit(1)
	}
	defer shutdown(context.Background())

	chConn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.ClickHouseAddr},
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouseDB,
			Username: cfg.ClickHouseUser,
			Password: cfg.ClickHousePassword,
		},
	})
	if err != nil {
		slog.Error("failed to connect clickhouse", "err", err)
		os.Exit(1)
	}
	defer chConn.Close()

	for i := 0; i < 10; i++ {
		if err = chConn.Ping(ctx); err == nil {
			break
		}
		slog.Warn("clickhouse not ready, retrying", "attempt", i, "err", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		slog.Error("failed to ping clickhouse after retries", "err", err)
		os.Exit(1)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.KafkaBrokers,
		Topic:       "click_events",
		GroupID:     "analytics-service",
		MaxAttempts: 10,
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		},
	})
	defer reader.Close()

	urlConn, err := grpc.NewClient(
		cfg.URLServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		slog.Error("failed to connect url service", "err", err)
		os.Exit(1)
	}
	defer urlConn.Close()

	urlServiceClient := urlv1.NewURLServiceClient(urlConn)

	repo := analytics_repository.NewRepository(chConn)
	urlClient := analytics_client.NewClient(urlServiceClient)
	svc := analytics_service.NewService(repo, urlClient)
	h := analytics_handler.NewHandler(svc)
	consumer := analytics_consumer.NewConsumer(reader, repo, batchSize, flushInterval)

	if err = repo.EnsureTable(ctx); err != nil {
		slog.Error("failed to ensure table", "err", err)
		os.Exit(1)
	}

	go func() {
		if err := consumer.Run(ctx); err != nil {
			slog.Error("failed to run consumer", "err", err)
		}
	}()

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // ← добавить
		grpc.UnaryInterceptor(analytics_metrics.UnaryServerInterceptor),
	)
	analyticsv1.RegisterAnalyticsServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+cfg.AnalyticsServicePort)
	if err != nil {
		slog.Error("failed to listen", "err", err)
		os.Exit(1)
	}
	defer lis.Close()

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	muxSrv := &http.Server{Addr: ":" + cfg.MetricsPort, Handler: mux}

	go func() {
		err := muxSrv.ListenAndServe()
		if err != nil {
			slog.Error("failed listen and server", "err", err)
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()
		slog.Info("shutting down analytics")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = muxSrv.Shutdown(shutdownCtx)
		grpcServer.GracefulStop()
	}()

	if err = grpcServer.Serve(lis); err != nil {
		slog.Error("failed to serve", "err", err)
		os.Exit(1)
	}
}
