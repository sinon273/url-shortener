package user_service

import (
	"context"
	"fmt"
	"log/slog"
)

func (s *Service) Logout(ctx context.Context, sessionID string) error {
	err := s.repo.RevokeSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	err = s.cache.DelSession(ctx, sessionID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete session cache", "err", err, "session_id", sessionID)
	}
	return nil
}
