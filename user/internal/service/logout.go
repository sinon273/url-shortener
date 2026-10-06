package user_service

import (
	"context"
	"fmt"
)

func (s *Service) Logout(ctx context.Context, sessionID string) error {
	err := s.repo.RevokeSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	_ = s.cache.DelSession(ctx, sessionID)

	return nil
}
