package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gateway_client "github.com/sinon273/url-shortener/gateway/internal/client"
	gateway_config "github.com/sinon273/url-shortener/gateway/internal/config"
	gateway_handler "github.com/sinon273/url-shortener/gateway/internal/handler"
	gateway_middleware "github.com/sinon273/url-shortener/gateway/internal/middleware"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := gateway_config.Load()
	if err != nil {
		fmt.Println("failed to load config:", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	userConn, err := grpc.NewClient(cfg.UserServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("failed to connect user:", err)
		os.Exit(1)
	}
	urlConn, err := grpc.NewClient(cfg.URLServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("failed to connect url:", err)
		os.Exit(1)
	}
	defer userConn.Close()
	defer urlConn.Close()

	userClient := userv1.NewUserServiceClient(userConn)
	urlClient := urlv1.NewURLServiceClient(urlConn)

	client := gateway_client.NewClient(userClient, urlClient)
	authMW := gateway_middleware.NewAuth(userClient)
	h := gateway_handler.NewHandler(client)
	router := h.Router(authMW)

	srv := &http.Server{
		Addr:              ":" + cfg.GatewayPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		_ = srv.Shutdown(shutdownCtx)
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server error:", err)
		os.Exit(1)
	}
}
