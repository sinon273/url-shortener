package url_service

import (
	"context"
	"time"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
)

const linkCacheTTL = time.Hour

func (s *Service) Redirect(ctx context.Context, shortCode string) (string, error) {
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
	_ = s.cache.SetLink(ctx, shortCode, link.OriginalURL, linkCacheTTL)
	return link.OriginalURL, nil
}
