package url_service

import (
	"context"
	"fmt"
	"log/slog"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
)

func (s *Service) DeleteURL(ctx context.Context, shortCode, userID string) error {
	link, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return fmt.Errorf("get link by short code: %w", err)
	}

	if link.UserID != userID {
		return url_domain.ErrPermissionDenied
	}
	err = s.repo.DeleteByShortCode(ctx, shortCode)
	if err != nil {
		return fmt.Errorf("delete link by short code: %w", err)
	}
	err = s.cache.DelLink(ctx, shortCode)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete link", "err", err, "short_code", shortCode)
	}
	err = s.cache.DecrCounter(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decrCounter", "err", err, "user_id", userID)
	}
	return nil
}
