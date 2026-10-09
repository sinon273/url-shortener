package url_service

import (
	"context"
	"testing"
	"time"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
	"github.com/sinon273/url-shortener/url/internal/service/mocks"
	"github.com/stretchr/testify/require"
)

func TestGetURL_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockUserClient := mocks.NewMockUserClient(t)
	svc := NewService(mockRepo, mockCache, mockUserClient, "https://short.url")

	link := url_domain.Link{
		ShortCode:   "abc1234",
		OriginalURL: "https://example.com",
		UserID:      "user-1",
		IsActive:    true,
	}

	mockRepo.EXPECT().
		GetByShortCode(ctx, "abc1234").
		Return(link, nil).
		Once()

	got, err := svc.GetURL(ctx, "abc1234")

	require.NoError(t, err)
	require.Equal(t, "https://example.com", got.OriginalURL)
}

func TestGetURL_Expired(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockUserClient := mocks.NewMockUserClient(t)
	svc := NewService(mockRepo, mockCache, mockUserClient, "https://short.url")

	past := time.Now().Add(-time.Hour)
	link := url_domain.Link{
		ShortCode: "abc1234",
		IsActive:  true,
		ExpiresAt: &past,
	}

	mockRepo.EXPECT().
		GetByShortCode(ctx, "abc1234").
		Return(link, nil).
		Once()

	_, err := svc.GetURL(ctx, "abc1234")

	require.ErrorIs(t, err, url_domain.ErrNotFound)
}
