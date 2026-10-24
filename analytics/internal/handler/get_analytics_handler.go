package analytics_handler

import (
	"context"

	analyticsv1 "github.com/sinon273/url-shortener/shared/pkg/proto/analytics/v1"
)

func (h *Handler) GetAnalytics(ctx context.Context, req *analyticsv1.GetAnalyticsRequest) (*analyticsv1.GetAnalyticsResponse, error) {
	shortCode := req.ShortCode
	userID := req.UserId
	days := int(req.Days)

	a, err := h.service.GetAnalytics(ctx, shortCode, userID, days)
	if err != nil {
		return nil, mapError(ctx, err)
	}
	clicksByDay := make([]*analyticsv1.DayClicks, len(a.ClicksByDay))
	for i, d := range a.ClicksByDay {
		clicksByDay[i] = &analyticsv1.DayClicks{
			Date:   d.Date,
			Clicks: d.Clicks,
		}
	}

	referer := make([]*analyticsv1.RefererClicks, len(a.TopReferers))
	for i, r := range a.TopReferers {
		referer[i] = &analyticsv1.RefererClicks{
			Referer: r.Referer,
			Clicks:  r.Clicks,
		}
	}

	response := analyticsv1.GetAnalyticsResponse{
		ShortCode:   shortCode,
		TotalClicks: a.TotalClicks,
		ClicksByDay: clicksByDay,
		TopReferers: referer,
	}

	return &response, nil
}
