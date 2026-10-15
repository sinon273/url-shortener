//go:build integration

package analytics_test

import (
	"context"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	analytics_domain "github.com/sinon273/url-shortener/analytics/internal/domain"
	analytics_repository "github.com/sinon273/url-shortener/analytics/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupClickHouse(t *testing.T) (clickhouse.Conn, func()) {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "clickhouse/clickhouse-server:24.8",
		ExposedPorts: []string{"9000/tcp"},
		Env: map[string]string{
			"CLICKHOUSE_USER":                      "default",
			"CLICKHOUSE_PASSWORD":                  "clickhouse",
			"CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT": "1",
		},
		WaitingFor: wait.ForListeningPort("9000/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "9000")
	addr := host + ":" + port.Port()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "clickhouse",
		},
	})
	require.NoError(t, err)
	require.NoError(t, conn.Ping(ctx))

	cleanup := func() {
		_ = conn.Close()
		_ = container.Terminate(ctx)
	}
	return conn, cleanup
}

func TestRepository_InsertAndGetAnalytics(t *testing.T) {
	ctx := context.Background()
	conn, cleanup := setupClickHouse(t)
	defer cleanup()

	repo := analytics_repository.NewRepository(conn)
	require.NoError(t, repo.EnsureTable(ctx))

	shortCode := "abc1234"
	userID := "user-1"
	now := time.Now()

	events := []analytics_domain.ClickEvent{
		{ShortCode: shortCode, UserID: userID, ClickedAt: now, IP: "1.1.1.1", UserAgent: "curl", Referer: "https://t.me"},
		{ShortCode: shortCode, UserID: userID, ClickedAt: now, IP: "1.1.1.2", UserAgent: "curl", Referer: "https://t.me"},
		{ShortCode: shortCode, UserID: userID, ClickedAt: now, IP: "1.1.1.3", UserAgent: "curl", Referer: "https://google.com"},
	}
	require.NoError(t, repo.InsertBatch(ctx, events))

	total, err := repo.TotalClicks(ctx, shortCode, 30)
	require.NoError(t, err)
	require.Equal(t, int64(3), total)

	byDay, err := repo.ClicksByDay(ctx, shortCode, 30)
	require.NoError(t, err)
	require.Len(t, byDay, 1)
	require.Equal(t, int64(3), byDay[0].Clicks)

	referers, err := repo.TopReferers(ctx, shortCode, 30)
	require.NoError(t, err)
	require.Len(t, referers, 2)
	require.Equal(t, "https://t.me", referers[0].Referer)
	require.Equal(t, int64(2), referers[0].Clicks)
}
