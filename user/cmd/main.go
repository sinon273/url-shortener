package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
	user_auth "github.com/sinon273/url-shortener/user/internal/auth"
	user_redis_cache "github.com/sinon273/url-shortener/user/internal/cache"
	user_config "github.com/sinon273/url-shortener/user/internal/config"
	user_transport_grpc "github.com/sinon273/url-shortener/user/internal/handler"
	user_postgres_repository "github.com/sinon273/url-shortener/user/internal/repository"
	user_service "github.com/sinon273/url-shortener/user/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := user_config.Load()
	if err != nil {
		fmt.Println("failed to load config:", err)
		os.Exit(1)
	}

	migration, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		fmt.Println("failed to init migrate:", err)
		os.Exit(1)
	}
	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("failed to apply migrations:", err)
		os.Exit(1)
	}
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Println("failed to init postgres pool:", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		fmt.Println("failed to ping postgres:", err)
		os.Exit(1)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer rdb.Close()

	if err = rdb.Ping(ctx).Err(); err != nil {
		fmt.Println("failed to ping redis:", err)
		os.Exit(1)
	}
	repo := user_postgres_repository.New(pool)
	cache := user_redis_cache.New(rdb)
	authn := user_auth.New(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	svc := user_service.NewService(repo, cache, authn)
	h := user_transport_grpc.New(svc)

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		fmt.Println("failed to listen:", err)
		os.Exit(1)
	}

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Println("failed to serve grpc:", err)
	}
}
