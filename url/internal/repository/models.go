package url_repository

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type LinkModel struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	ShortCode   string        `bson:"short_code"`
	OriginalURL string        `bson:"original_url"`
	UserID      string        `bson:"user_id"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
	ExpiresAt   *time.Time    `bson:"expires_at,omitempty"`
	IsActive    bool          `bson:"is_active"`
}
