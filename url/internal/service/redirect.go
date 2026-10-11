package url_service

import (
	"context"
	"time"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
	url_producer "github.com/sinon273/url-shortener/url/internal/producer"
)

const linkCacheTTL = time.Hour

func (s *Service) Redirect(ctx context.Context, shortCode, ip, userAgent, referer string) (string, error) {
	originalURL, err := s.cache.GetLink(ctx, shortCode)
	if err == nil {
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
			//TODO: Log error(slog будет позже)
			_ = err
		}
	}()

	_ = s.cache.SetLink(ctx, shortCode, link.OriginalURL, linkCacheTTL)
	return link.OriginalURL, nil
}
