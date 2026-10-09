package url_handler

import (
	"context"

	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
)

type Handler struct {
	urlv1.UnimplementedURLServiceServer
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type Service interface {
	CreateShortURL(ctx context.Context, originalURL, userID string, expiresIn *int64) (string, string, error)
	GetURL(ctx context.Context, shortCode string) (url_domain.Link, error)
	DeleteURL(ctx context.Context, shortCode, userID string) error
	ListUserURLs(ctx context.Context, userID string, limit int, cursor string) ([]url_domain.Link, string, error)
	Redirect(ctx context.Context, shortCode string) (string, error)
}
