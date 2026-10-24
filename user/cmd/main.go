package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
	user_auth "github.com/sinon273/url-shortener/user/internal/auth"
	user_redis_cache "github.com/sinon273/url-shortener/user/internal/cache"
	user_config "github.com/sinon273/url-shortener/user/internal/config"
	user_transport_grpc "github.com/sinon273/url-shortener/user/internal/handler"
	user_logger "github.com/sinon273/url-shortener/user/internal/logger"
	user_metrics "github.com/sinon273/url-shortener/user/internal/metrics"
	user_postgres_repository "github.com/sinon273/url-shortener/user/internal/repository"
	user_service "github.com/sinon273/url-shortener/user/internal/service"
	user_tracing "github.com/sinon273/url-shortener/user/internal/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := user_config.Load()
	if err != nil {
		fmt.Println("failed to load config:", err)
		os.Exit(1)
	}

	logger, err := user_logger.Init(cfg.LogLevel)
	if err != nil {
		fmt.Println("failed to init logger:", err)
		os.Exit(1)
	}
	slog.SetDefault(logger)
	slog.Info("user service starting", "grpc_port", cfg.GRPCPort, "metrics_port", cfg.MetricsPort)
	shutdown, err := user_tracing.Init(ctx, cfg.OtelService, cfg.OtelExporter)
	if err != nil {
		slog.Error("failed to init tracing: ", "err", err)
		os.Exit(1)
	}
	defer shutdown(context.Background())

	migration, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to init migrate: ", "err", err)
		os.Exit(1)
	}
	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("failed to apply migrations: ", "err", err)
		os.Exit(1)
	}
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to init postgres pool: ", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		slog.Error("failed to ping postgres: ", "err", err)
		os.Exit(1)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer rdb.Close()

	if err = rdb.Ping(ctx).Err(); err != nil {
		slog.Error("failed to ping redis: ", "err", err)
		os.Exit(1)
	}
	repo := user_postgres_repository.New(pool)
	cache := user_redis_cache.New(rdb)
	authn := user_auth.New(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	svc := user_service.NewService(repo, cache, authn)
	h := user_transport_grpc.New(svc)

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(user_metrics.UnaryServerInterceptor),
	)
	userv1.RegisterUserServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		slog.Error("failed to listen: ", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	metricsSrv := &http.Server{Addr: ":" + cfg.MetricsPort, Handler: mux}

	go func() {
		err := metricsSrv.ListenAndServe()
		if err != nil {
			slog.Error("failed to serve metrics: ", "err", err)
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()
		slog.Info("shutting down user service")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricsSrv.Shutdown(shutdownCtx)
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("failed to serve grpc: ", "err", err)
	}
}
