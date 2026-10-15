package analytics_service

import (
	"context"

	analytics_domain "github.com/sinon273/url-shortener/analytics/internal/domain"
)

type Service struct {
	repo      Repository
	urlClient UserClient
}

func NewService(repo Repository, urlClient UserClient) *Service {
	return &Service{
		repo:      repo,
		urlClient: urlClient,
	}
}

type Repository interface {
	TotalClicks(ctx context.Context, shortCode string, days int) (int64, error)
	//сколько кликов было по дням(день+клики)
	ClicksByDay(ctx context.Context, shortCode string, days int) ([]analytics_domain.DayClicks, error)
	//колько кликов из каждого источника (источник+клики)
	TopReferers(ctx context.Context, shortCode string, days int) ([]analytics_domain.RefererClicks, error)
}

type UserClient interface {
	GetLinkOwner(ctx context.Context, shortCode string) (string, error)
}
