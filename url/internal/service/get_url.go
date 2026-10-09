package url_service

import (
	"context"
	"fmt"
	"time"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
)

func (s *Service) GetURL(ctx context.Context, shortCode string) (url_domain.Link, error) {
	link, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return url_domain.Link{}, fmt.Errorf("get by short code: %w", err)
	}

	if !link.IsActive {
		return url_domain.Link{}, url_domain.ErrNotFound
	}

	if link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now()) {
		return url_domain.Link{}, url_domain.ErrNotFound
	}

	return link, nil
}
