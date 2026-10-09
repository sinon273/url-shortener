package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
	url_cache "github.com/sinon273/url-shortener/url/internal/cache"
	url_client "github.com/sinon273/url-shortener/url/internal/client"
	url_config "github.com/sinon273/url-shortener/url/internal/config"
	url_handler "github.com/sinon273/url-shortener/url/internal/handler"
	url_repository "github.com/sinon273/url-shortener/url/internal/repository"
	url_service "github.com/sinon273/url-shortener/url/internal/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := url_config.Load()
	if err != nil {
		fmt.Println("loading config: %w", err)
		os.Exit(1)
	}

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		fmt.Println("failed to connect mongo:", err)
		os.Exit(1)
	}
	defer mongoClient.Disconnect(ctx)

	if err = mongoClient.Ping(ctx, nil); err != nil {
		fmt.Println("failed to ping mongo:", err)
		os.Exit(1)
	}
	coll := mongoClient.Database(cfg.MongoDB).Collection("links")

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer rdb.Close()

	if err = rdb.Ping(ctx).Err(); err != nil {
		fmt.Println("failed to ping redis:", err)
		os.Exit(1)
	}

	userConn, err := grpc.NewClient(
		cfg.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Println("failed to connect user service:", err)
		os.Exit(1)
	}
	defer userConn.Close()

	userServiceClient := userv1.NewUserServiceClient(userConn)

	repo := url_repository.NewRepository(coll)
	cache := url_cache.NewCache(rdb)
	userClient := url_client.NewClient(userServiceClient)
	svc := url_service.NewService(repo, cache, userClient, cfg.ShortURLBase)
	h := url_handler.NewHandler(svc)

	if err = repo.EnsureIndexes(ctx); err != nil {
		fmt.Println("failed to ensure indexes:", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	urlv1.RegisterURLServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		fmt.Println("failed to listen:", err)
		os.Exit(1)
	}

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	if err = grpcServer.Serve(lis); err != nil {
		fmt.Println("failed to serve:", err)
		os.Exit(1)
	}
}
