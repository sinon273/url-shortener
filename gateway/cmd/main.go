package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	gateway_client "github.com/sinon273/url-shortener/gateway/internal/client"
	gateway_config "github.com/sinon273/url-shortener/gateway/internal/config"
	gateway_handler "github.com/sinon273/url-shortener/gateway/internal/handler"
	gateway_logger "github.com/sinon273/url-shortener/gateway/internal/logger"
	gateway_middleware "github.com/sinon273/url-shortener/gateway/internal/middleware"
	gateway_tracing "github.com/sinon273/url-shortener/gateway/internal/tracing"
	analyticsv1 "github.com/sinon273/url-shortener/shared/pkg/proto/analytics/v1"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := gateway_config.Load()
	if err != nil {
		fmt.Println("failed to load config:", err)
		os.Exit(1)
	}

	logger, err := gateway_logger.Init(cfg.LogLevel)
	if err != nil {
		fmt.Println("failed to init logger:", err)
		os.Exit(1)
	}
	slog.SetDefault(logger)
	slog.Info("gateway starting", "port", cfg.GatewayPort, "metrics_port", cfg.MetricsPort)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	shutdown, err := gateway_tracing.Init(ctx, cfg.OtelService, cfg.OtelExporter)
	if err != nil {
		slog.Error("failed to init tracing", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown tracing", "err", err)
		}
	}()

	userConn, err := grpc.NewClient(
		cfg.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		slog.Error("failed to connect user", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := userConn.Close(); err != nil {
			slog.Error("failed to close user conn", "err", err)
		}
	}()

	urlConn, err := grpc.NewClient(
		cfg.URLServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		slog.Error("failed to connect url", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := urlConn.Close(); err != nil {
			slog.Error("failed to close url conn", "err", err)
		}
	}()

	analyticsConn, err := grpc.NewClient(
		cfg.AnalyticsServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		slog.Error("failed to connect analytics", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := analyticsConn.Close(); err != nil {
			slog.Error("failed to close analytics conn", "err", err)
		}
	}()

	userClient := userv1.NewUserServiceClient(userConn)
	urlClient := urlv1.NewURLServiceClient(urlConn)
	analyticsClient := analyticsv1.NewAnalyticsServiceClient(analyticsConn)

	client := gateway_client.NewClient(userClient, urlClient, analyticsClient)
	authMW := gateway_middleware.NewAuth(userClient)
	h := gateway_handler.NewHandler(client)
	router := h.Router(authMW)

	handler := otelhttp.NewHandler(router, "gateway")

	srv := &http.Server{
		Addr:              ":" + cfg.GatewayPort,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
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
			slog.Error("metrics server error", "err", err)
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()
		slog.Info("shutting down gateway")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		_ = metricsSrv.Shutdown(shutdownCtx)
		_ = srv.Shutdown(shutdownCtx)
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
