package url_cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
)

func (c *Cache) GetLimit(ctx context.Context, userID string) (int, error) {
	val, err := c.rdb.Get(ctx, limitKey(userID)).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, url_domain.ErrNotFound
		}
		return 0, fmt.Errorf("get limit from cache: %w", err)
	}
	return val, nil
}

func (c *Cache) SetLimit(ctx context.Context, userID string, limit int, ttl time.Duration) error {
	err := c.rdb.Set(ctx, limitKey(userID), limit, ttl).Err()
	if err != nil {
		return fmt.Errorf("set limit to cache: %w", err)
	}
	return nil
}
