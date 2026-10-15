package analytics_service

import (
	"context"
	"fmt"

	analytics_domain "github.com/sinon273/url-shortener/analytics/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetAnalytics(ctx context.Context, shortCode, userID string, days int) (analytics_domain.Analytics, error) {
	if days <= 0 {
		days = 30
	}
	if days > 365 {
		days = 365
	}
	owner, err := s.urlClient.GetLinkOwner(ctx, shortCode)
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			return analytics_domain.Analytics{}, analytics_domain.ErrUnavailable
		}
		return analytics_domain.Analytics{}, fmt.Errorf("get link owner: %w", err)
	}
	if owner != userID {
		return analytics_domain.Analytics{}, analytics_domain.ErrPermissionDenied
	}
	totalClicks, err := s.repo.TotalClicks(ctx, shortCode, days)
	if err != nil {
		return analytics_domain.Analytics{}, fmt.Errorf("get total clicks: %w", err)
	}
	clicksByDay, err := s.repo.ClicksByDay(ctx, shortCode, days)
	if err != nil {
		return analytics_domain.Analytics{}, fmt.Errorf("get clicks by day: %w", err)
	}
	topReferers, err := s.repo.TopReferers(ctx, shortCode, days)
	if err != nil {
		return analytics_domain.Analytics{}, fmt.Errorf("get top referers: %w", err)
	}

	return analytics_domain.Analytics{
		ShortCode:   shortCode,
		TotalClicks: totalClicks,
		ClicksByDay: clicksByDay,
		TopReferers: topReferers,
	}, nil
}
