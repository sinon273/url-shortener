package url_domain

import (
	"time"
)

type Link struct {
	ShortCode   string
	OriginalURL string
	UserID      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   *time.Time
	IsActive    bool
}
