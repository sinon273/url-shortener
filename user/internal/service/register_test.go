package user_service

import (
	"context"
	"testing"

	"github.com/sinon273/url-shortener/user/internal/domain"
	"github.com/sinon273/url-shortener/user/internal/service/mocks"
	"github.com/stretchr/testify/require"
)

func TestRegister_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockAuth := mocks.NewMockAuth(t)

	svc := NewService(mockRepo, mockCache, mockAuth)

	mockAuth.EXPECT().
		HashPassword("password123").
		Return("hashed", nil).
		Once()

	mockRepo.EXPECT().
		CreateUserWithSubscription(ctx, "test@example.com", "hashed").
		Return("user-id-123", nil).
		Once()

	userID, err := svc.Register(ctx, "test@example.com", "password123")

	require.NoError(t, err)
	require.Equal(t, "user-id-123", userID)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockAuth := mocks.NewMockAuth(t)

	svc := NewService(mockRepo, mockCache, mockAuth)

	mockAuth.EXPECT().
		HashPassword("password123").
		Return("hashed", nil).
		Once()

	mockRepo.EXPECT().
		CreateUserWithSubscription(ctx, "test@example.com", "hashed").
		Return("", domain.ErrEmailAlreadyExists).
		Once()

	userID, err := svc.Register(ctx, "test@example.com", "password123")

	require.ErrorIs(t, err, domain.ErrEmailAlreadyExists)
	require.Empty(t, userID)
}
