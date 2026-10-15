package analytics_domain

import "time"

type ClickEvent struct {
	ShortCode string    `json:"short_code"`
	UserID    string    `json:"user_id"`
	ClickedAt time.Time `json:"clicked_at"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
}

type DayClicks struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

type RefererClicks struct {
	Referer string `json:"referer"`
	Clicks  int64  `json:"clicks"`
}

type Analytics struct {
	ShortCode   string          `json:"short_code"`
	TotalClicks int64           `json:"total_clicks"`
	ClicksByDay []DayClicks     `json:"clicks_by_day"`
	TopReferers []RefererClicks `json:"top_referers"`
}
