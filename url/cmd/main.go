package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
	url_cache "github.com/sinon273/url-shortener/url/internal/cache"
	url_client "github.com/sinon273/url-shortener/url/internal/client"
	url_config "github.com/sinon273/url-shortener/url/internal/config"
	url_handler "github.com/sinon273/url-shortener/url/internal/handler"
	url_logger "github.com/sinon273/url-shortener/url/internal/logger"
	url_metrics "github.com/sinon273/url-shortener/url/internal/metrics"
	url_producer "github.com/sinon273/url-shortener/url/internal/producer"
	url_repository "github.com/sinon273/url-shortener/url/internal/repository"
	url_service "github.com/sinon273/url-shortener/url/internal/service"
	url_tracing "github.com/sinon273/url-shortener/url/internal/tracing"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log/slog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := url_config.Load()
	if err != nil {
		fmt.Printf("loading config: %v\n", err)
		os.Exit(1)
	}

	logger, err := url_logger.Init(cfg.LogLevel)
	if err != nil {
		fmt.Println("failed to init logger:", err)
		os.Exit(1)
	}

	slog.SetDefault(logger)
	slog.Info("url service starting", "grpc_port", cfg.GRPCPort, "metrics_port", cfg.MetricsPort)

	shutdown, err := url_tracing.Init(ctx, cfg.OtelService, cfg.OtelExporter)
	if err != nil {
		slog.Error("failed to init tracing", "err", err)
		os.Exit(1)
	}
	defer shutdown(context.Background())

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		slog.Error("failed to connect mongo", "err", err)
		os.Exit(1)
	}
	defer mongoClient.Disconnect(ctx)

	if err = mongoClient.Ping(ctx, nil); err != nil {
		slog.Error("failed to ping mongo", "err", err)
		os.Exit(1)
	}
	coll := mongoClient.Database(cfg.MongoDB).Collection("links")

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer rdb.Close()

	if err = rdb.Ping(ctx).Err(); err != nil {
		slog.Error("failed to ping redis", "err", err)
		os.Exit(1)
	}

	userConn, err := grpc.NewClient(
		cfg.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		slog.Error("failed to connect user service", "err", err)
		os.Exit(1)
	}
	defer userConn.Close()

	kafkaWriter := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.KafkaBrokers...),
		Topic:                  "click_events",
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	producer := url_producer.NewProducer(kafkaWriter)

	userServiceClient := userv1.NewUserServiceClient(userConn)

	repo := url_repository.NewRepository(coll)
	cache := url_cache.NewCache(rdb)
	userClient := url_client.NewClient(userServiceClient)
	svc := url_service.NewService(repo, cache, userClient, cfg.ShortURLBase, producer)
	h := url_handler.NewHandler(svc)

	if err = repo.EnsureIndexes(ctx); err != nil {
		slog.Error("failed to ensure indexes", "err", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(url_metrics.UnaryServerInterceptor),
	)

	urlv1.RegisterURLServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		slog.Error("failed to listen", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	metricsServer := &http.Server{Addr: ":" + cfg.MetricsPort, Handler: mux}

	go func() {
		err := metricsServer.ListenAndServe()
		if err != nil {
			slog.Error("failed to start metrics server", "err", err)
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()
		slog.Info("shutting down url service")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricsServer.Shutdown(shutdownCtx)
		_ = kafkaWriter.Close()
		grpcServer.GracefulStop()
	}()

	if err = grpcServer.Serve(lis); err != nil {
		slog.Error("failed to serve", "err", err)
		os.Exit(1)
	}
}
