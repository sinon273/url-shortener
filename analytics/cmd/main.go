package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/segmentio/kafka-go"
	analytics_client "github.com/sinon273/url-shortener/analytics/internal/client"
	analytics_config "github.com/sinon273/url-shortener/analytics/internal/config"
	analytics_consumer "github.com/sinon273/url-shortener/analytics/internal/consumer"
	analytics_handler "github.com/sinon273/url-shortener/analytics/internal/handler"
	analytics_repository "github.com/sinon273/url-shortener/analytics/internal/repository"
	analytics_service "github.com/sinon273/url-shortener/analytics/internal/service"
	analyticsv1 "github.com/sinon273/url-shortener/shared/pkg/proto/analytics/v1"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
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
	chConn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.ClickHouseAddr},
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouseDB,
			Username: cfg.ClickHouseUser,
			Password: cfg.ClickHousePassword,
		},
	})
	if err != nil {
		fmt.Printf("failed to connect clickhouse: %v\n", err)
		os.Exit(1)
	}
	defer chConn.Close()

	if err = chConn.Ping(ctx); err != nil {
		fmt.Printf("failed to ping clickhouse: %v\n", err)
		os.Exit(1)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.KafkaBrokers,
		Topic:   "click_events",
		GroupID: "analytics-service",
	})
	defer reader.Close()

	urlConn, err := grpc.NewClient(
		cfg.URLServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Printf("failed to connect url service: %v\n", err)
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
		fmt.Printf("failed to ensure table: %v\n", err)
		os.Exit(1)
	}

	go func() {
		if err := consumer.Run(ctx); err != nil {
			fmt.Printf("failed to run consumer: %v\n", err)
		}
	}()

	grpcServer := grpc.NewServer()
	analyticsv1.RegisterAnalyticsServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+cfg.AnalyticsServicePort)
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		os.Exit(1)
	}
	defer lis.Close()

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	if err = grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
		os.Exit(1)
	}
}
