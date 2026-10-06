package user_service

import (
	"context"

	"github.com/sinon273/url-shortener/user/internal/domain"
)

func (s *Service) GetLimit(ctx context.Context, userID string) (int, error) {
	limit, err := s.repo.GetSubscriptionLimit(ctx, userID)
	if err != nil {
		return 0, domain.ErrNotFound
	}
	return limit, nil
}
