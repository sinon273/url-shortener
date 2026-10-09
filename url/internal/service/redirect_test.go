package url_service

import (
	"context"
	"testing"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
	"github.com/sinon273/url-shortener/url/internal/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRedirect_CacheHit(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockUserClient := mocks.NewMockUserClient(t)
	svc := NewService(mockRepo, mockCache, mockUserClient, "https://short.url")

	mockCache.EXPECT().
		GetLink(ctx, "abc1234").
		Return("https://example.com", nil).
		Once()

	url, err := svc.Redirect(ctx, "abc1234")

	require.NoError(t, err)
	require.Equal(t, "https://example.com", url)
}

func TestRedirect_CacheMiss_MongoHit(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockUserClient := mocks.NewMockUserClient(t)
	svc := NewService(mockRepo, mockCache, mockUserClient, "https://short.url")

	link := url_domain.Link{
		ShortCode:   "abc1234",
		OriginalURL: "https://example.com",
		IsActive:    true,
	}

	mockCache.EXPECT().
		GetLink(ctx, "abc1234").
		Return("", url_domain.ErrNotFound).
		Once()

	mockRepo.EXPECT().
		GetByShortCode(ctx, "abc1234").
		Return(link, nil).
		Once()

	mockCache.EXPECT().
		SetLink(ctx, "abc1234", "https://example.com", mock.Anything).
		Return(nil).
		Once()

	url, err := svc.Redirect(ctx, "abc1234")

	require.NoError(t, err)
	require.Equal(t, "https://example.com", url)
}
