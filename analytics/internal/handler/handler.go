package analytics_handler

import (
	"context"

	analytics_domain "github.com/sinon273/url-shortener/analytics/internal/domain"
	analytics_v1 "github.com/sinon273/url-shortener/shared/pkg/proto/analytics/v1"
)

type Handler struct {
	analytics_v1.UnimplementedAnalyticsServiceServer
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type Service interface {
	GetAnalytics(ctx context.Context, shortCode, userID string, days int) (analytics_domain.Analytics, error)
}
