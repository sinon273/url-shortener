package analytics_repository

import (
	"context"
	"fmt"

	analytics_domain "github.com/sinon273/url-shortener/analytics/internal/domain"
)

func (r *Repository) TotalClicks(ctx context.Context, shortCode string, days int) (int64, error) {
	query := `
	SELECT toInt64(count()) FROM click_events
	WHERE short_code = ? AND clicked_at >= now() - INTERVAL ? DAY
	`

	var total int64
	err := r.conn.QueryRow(ctx, query, shortCode, days).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("get total clicks: %w", err)
	}
	return total, nil
}

func (r *Repository) ClicksByDay(ctx context.Context, shortCode string, days int) ([]analytics_domain.DayClicks, error) {
	query := `
	SELECT formatDateTime(toDate(clicked_at), '%Y-%m-%d') AS d, toInt64(count())
	FROM click_events
	WHERE short_code = ? AND clicked_at >= now() - INTERVAL ? DAY
	GROUP BY d
	ORDER BY d DESC
	`

	rows, err := r.conn.Query(ctx, query, shortCode, days)
	if err != nil {
		return nil, fmt.Errorf("get clicks by day: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var result []analytics_domain.DayClicks
	for rows.Next() {
		var item analytics_domain.DayClicks
		if err := rows.Scan(&item.Date, &item.Clicks); err != nil {
			return nil, fmt.Errorf("scan day clicks: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return result, nil
}

func (r *Repository) TopReferers(ctx context.Context, shortCode string, days int) ([]analytics_domain.RefererClicks, error) {

	query := `
	SELECT referer, toInt64(count()) AS c
	FROM click_events
	WHERE short_code = ? AND clicked_at >= now() - INTERVAL ? DAY
	GROUP BY referer
	ORDER BY c DESC
	LIMIT 10
	`

	rows, err := r.conn.Query(ctx, query, shortCode, days)
	if err != nil {
		return nil, fmt.Errorf("get top referers: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var result []analytics_domain.RefererClicks
	for rows.Next() {
		var item analytics_domain.RefererClicks
		if err := rows.Scan(&item.Referer, &item.Clicks); err != nil {
			return nil, fmt.Errorf("scan referer clicks: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return result, nil
}
