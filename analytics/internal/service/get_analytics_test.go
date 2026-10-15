package analytics_service

import (
	"context"
	"errors"
	"testing"

	analytics_domain "github.com/sinon273/url-shortener/analytics/internal/domain"
	"github.com/sinon273/url-shortener/analytics/internal/service/mocks"
	"github.com/stretchr/testify/require"
)

func TestGetAnalytics_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockUserClient := mocks.NewMockUserClient(t)

	svc := NewService(mockRepo, mockUserClient)

	shortCode := "abc1234"
	userID := "user-1"

	mockUserClient.EXPECT().
		GetLinkOwner(ctx, shortCode).
		Return(userID, nil).
		Once()

	mockRepo.EXPECT().
		TotalClicks(ctx, shortCode, 30).
		Return(int64(100), nil).
		Once()

	mockRepo.EXPECT().
		ClicksByDay(ctx, shortCode, 30).
		Return([]analytics_domain.DayClicks{
			{Date: "2026-10-15", Clicks: 60},
			{Date: "2026-10-14", Clicks: 40},
		}, nil).
		Once()

	mockRepo.EXPECT().
		TopReferers(ctx, shortCode, 30).
		Return([]analytics_domain.RefererClicks{
			{Referer: "https://t.me", Clicks: 50},
		}, nil).
		Once()

	result, err := svc.GetAnalytics(ctx, shortCode, userID, 0)

	require.NoError(t, err)
	require.Equal(t, shortCode, result.ShortCode)
	require.Equal(t, int64(100), result.TotalClicks)
	require.Len(t, result.ClicksByDay, 2)
	require.Len(t, result.TopReferers, 1)
}

func TestGetAnalytics_PermissionDenied(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockUserClient := mocks.NewMockUserClient(t)

	svc := NewService(mockRepo, mockUserClient)

	shortCode := "abc1234"

	mockUserClient.EXPECT().
		GetLinkOwner(ctx, shortCode).
		Return("owner-1", nil).
		Once()

	_, err := svc.GetAnalytics(ctx, shortCode, "attacker-2", 30)

	require.ErrorIs(t, err, analytics_domain.ErrPermissionDenied)
}

func TestGetAnalytics_LinkNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(t)
	mockUserClient := mocks.NewMockUserClient(t)

	svc := NewService(mockRepo, mockUserClient)

	shortCode := "abc1234"

	mockUserClient.EXPECT().
		GetLinkOwner(ctx, shortCode).
		Return("", analytics_domain.ErrNotFound).
		Once()

	_, err := svc.GetAnalytics(ctx, shortCode, "user-1", 30)

	require.Error(t, err)
	require.True(t, errors.Is(err, analytics_domain.ErrNotFound))
}
