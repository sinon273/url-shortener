package url_service

import (
	"context"
	"log/slog"
	"time"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
	url_metrics "github.com/sinon273/url-shortener/url/internal/metrics"
	url_producer "github.com/sinon273/url-shortener/url/internal/producer"
)

const linkCacheTTL = time.Hour

func (s *Service) Redirect(ctx context.Context, shortCode, ip, userAgent, referer string) (string, error) {
	originalURL, err := s.cache.GetLink(ctx, shortCode)
	if err == nil {
		url_metrics.RedirectsTotal.Inc()
		return originalURL, nil
	}

	link, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	if !link.IsActive {
		return "", url_domain.ErrNotFound
	}

	if link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now()) {
		return "", url_domain.ErrNotFound
	}

	if s.producer != nil {
		go func() {
			event := url_producer.ClickEvent{
				ShortCode: shortCode,
				UserID:    link.UserID,
				ClickedAt: time.Now().UTC(),
				IP:        ip,
				UserAgent: userAgent,
				Referer:   referer,
			}

			if err := s.producer.PublishClick(context.Background(), event); err != nil {
				slog.ErrorContext(context.Background(), "failed to publish click event", "err", err)
			}
		}()
	}

	err = s.cache.SetLink(ctx, shortCode, link.OriginalURL, linkCacheTTL)
	if err != nil {
		slog.ErrorContext(ctx, "failed to cache link", "err", err)
	}
	url_metrics.RedirectsTotal.Inc()

	return link.OriginalURL, nil
}
