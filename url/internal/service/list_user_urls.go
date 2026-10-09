package url_service

import (
	"context"
	"fmt"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
)

func (s *Service) ListUserURLs(ctx context.Context, userID string, limit int, cursor string) ([]url_domain.Link, string, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	links, nextCursor, err := s.repo.ListByUser(ctx, userID, limit, cursor)
	if err != nil {
		return []url_domain.Link{}, "", fmt.Errorf("error list URLs user: %w", err)
	}
	return links, nextCursor, nil
}
