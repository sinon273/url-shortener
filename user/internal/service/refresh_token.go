package user_service

import (
	"context"
	"fmt"

	"github.com/sinon273/url-shortener/user/internal/domain"
)

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	userID, sessionID, err := s.auth.Parse(refreshToken)
	if err != nil {
		return "", domain.ErrUnauthenticated
	}
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return "", domain.ErrUnauthenticated
	}
	if !session.IsAlive() {
		return "", domain.ErrUnauthenticated
	}
	accessToken, err := s.auth.SignAccess(userID, sessionID)
	if err != nil {
		return "", fmt.Errorf("sign access: %w", err)
	}
	_ = s.cache.SetSession(ctx, sessionID, userID, s.auth.AccessTTL())
	return accessToken, nil
}
