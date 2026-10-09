package url_service

import (
	"context"
	"time"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
)

type Service struct {
	repo         Repository
	cache        Cache
	userClient   UserClient
	shortURLBase string
}

func NewService(repo Repository, cache Cache, client UserClient, shortURLBase string) *Service {
	return &Service{
		repo:         repo,
		cache:        cache,
		userClient:   client,
		shortURLBase: shortURLBase,
	}
}

type Repository interface {
	Create(ctx context.Context, link url_domain.Link) error
	GetByShortCode(ctx context.Context, shortCode string) (url_domain.Link, error)
	DeleteByShortCode(ctx context.Context, shortCode string) error
	ListByUser(ctx context.Context, userID string, limit int, cursor string) ([]url_domain.Link, string, error)
	CountActiveByUser(ctx context.Context, userID string) (int, error)
}

type Cache interface {
	GetLink(ctx context.Context, shortCode string) (string, error)
	SetLink(ctx context.Context, shortCode, originalURL string, ttl time.Duration) error
	DelLink(ctx context.Context, shortCode string) error

	GetCounter(ctx context.Context, userID string) (int, error)
	SetCounter(ctx context.Context, userID string, counter int) error
	IncrCounter(ctx context.Context, userID string) (int, error)
	DecrCounter(ctx context.Context, userID string) error

	GetLimit(ctx context.Context, userID string) (int, error)
	SetLimit(ctx context.Context, userID string, limit int, ttl time.Duration) error
}

type UserClient interface {
	GetLimit(ctx context.Context, userID string) (int, error)
}
