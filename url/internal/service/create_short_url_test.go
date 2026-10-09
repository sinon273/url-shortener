package url_service

import (
	"context"
	"errors"
	"testing"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
	"github.com/sinon273/url-shortener/url/internal/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURL_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockUserClient := mocks.NewMockUserClient(t)

	svc := NewService(mockRepo, mockCache, mockUserClient, "https://short.url")

	userID := "user-1"

	mockCache.EXPECT().
		GetLimit(ctx, userID).
		Return(1000, nil).
		Once()

	mockCache.EXPECT().
		GetCounter(ctx, userID).
		Return(5, nil).
		Once()

	mockCache.EXPECT().
		IncrCounter(ctx, userID).
		Return(6, nil).
		Once()

	mockRepo.EXPECT().
		Create(ctx, mock.Anything).
		Return(nil).
		Once()

	shortCode, shortURL, err := svc.CreateShortURL(ctx, "https://example.com", userID, nil)

	require.NoError(t, err)
	require.Len(t, shortCode, 7)
	require.Equal(t, "https://short.url/"+shortCode, shortURL)
}

func TestCreateShortURL_LimitExceeded(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockUserClient := mocks.NewMockUserClient(t)

	svc := NewService(mockRepo, mockCache, mockUserClient, "https://short.url")

	userID := "user-1"

	mockCache.EXPECT().
		GetLimit(ctx, userID).
		Return(3, nil).
		Once()

	mockCache.EXPECT().
		GetCounter(ctx, userID).
		Return(3, nil).
		Once()

	mockCache.EXPECT().
		IncrCounter(ctx, userID).
		Return(4, nil).
		Once()

	mockCache.EXPECT().
		DecrCounter(ctx, userID).
		Return(nil).
		Once()

	_, _, err := svc.CreateShortURL(ctx, "https://example.com", userID, nil)

	require.ErrorIs(t, err, url_domain.ErrResourceExhausted)
}

func TestCreateShortURL_MongoInsertFails(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockUserClient := mocks.NewMockUserClient(t)

	svc := NewService(mockRepo, mockCache, mockUserClient, "https://short.url")

	userID := "user-1"

	mockCache.EXPECT().
		GetLimit(ctx, userID).
		Return(1000, nil).
		Once()

	mockCache.EXPECT().
		GetCounter(ctx, userID).
		Return(5, nil).
		Once()

	mockCache.EXPECT().
		IncrCounter(ctx, userID).
		Return(6, nil).
		Once()

	mockRepo.EXPECT().
		Create(ctx, mock.Anything).
		Return(errors.New("mongo down")).
		Once()

	mockCache.EXPECT().
		DecrCounter(ctx, userID).
		Return(nil).
		Once()

	_, _, err := svc.CreateShortURL(ctx, "https://example.com", userID, nil)

	require.Error(t, err)
	require.NotErrorIs(t, err, url_domain.ErrResourceExhausted)
}
