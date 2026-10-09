package url_service

import (
	"context"
	"testing"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
	"github.com/sinon273/url-shortener/url/internal/service/mocks"
	"github.com/stretchr/testify/require"
)

func TestDeleteURL_PermissionDenied(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockUserClient := mocks.NewMockUserClient(t)
	svc := NewService(mockRepo, mockCache, mockUserClient, "https://short.url")

	link := url_domain.Link{
		ShortCode: "abc1234",
		UserID:    "owner-1",
	}

	mockRepo.EXPECT().
		GetByShortCode(ctx, "abc1234").
		Return(link, nil).
		Once()

	err := svc.DeleteURL(ctx, "abc1234", "attacker-2")

	require.ErrorIs(t, err, url_domain.ErrPermissionDenied)
}

func TestDeleteURL_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockCache := mocks.NewMockCache(t)
	mockUserClient := mocks.NewMockUserClient(t)
	svc := NewService(mockRepo, mockCache, mockUserClient, "https://short.url")

	link := url_domain.Link{
		ShortCode: "abc1234",
		UserID:    "user-1",
	}

	mockRepo.EXPECT().
		GetByShortCode(ctx, "abc1234").
		Return(link, nil).
		Once()

	mockRepo.EXPECT().
		DeleteByShortCode(ctx, "abc1234").
		Return(nil).
		Once()

	mockCache.EXPECT().
		DelLink(ctx, "abc1234").
		Return(nil).
		Once()

	mockCache.EXPECT().
		DecrCounter(ctx, "user-1").
		Return(nil).
		Once()

	err := svc.DeleteURL(ctx, "abc1234", "user-1")

	require.NoError(t, err)
}
