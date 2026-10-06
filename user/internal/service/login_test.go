package user_service

import (
	"context"
	"testing"

	"github.com/sinon273/url-shortener/user/internal/domain"
	user_domain "github.com/sinon273/url-shortener/user/internal/domain/user"
	"github.com/sinon273/url-shortener/user/internal/service/mocks"
	"github.com/stretchr/testify/require"
)

func TestLogin_InvalidPassword(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockAuth := mocks.NewMockAuth(t)

	svc := NewService(mockRepo, mockCache, mockAuth)

	mockRepo.EXPECT().
		GetUserByEmail(ctx, "test@example.com").
		Return(user_domain.User{
			ID:           "user-id-123",
			Email:        "test@example.com",
			PasswordHash: "hashed",
		}, nil).
		Once()

	mockAuth.EXPECT().
		ComparePassword("hashed", "wrong-password").
		Return(domain.ErrUnauthenticated).
		Once()

	access, refresh, err := svc.Login(ctx, "test@example.com", "wrong-password")

	require.ErrorIs(t, err, domain.ErrUnauthenticated)
	require.Empty(t, access)
	require.Empty(t, refresh)
}
