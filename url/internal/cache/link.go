package url_cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
)

func (c *Cache) GetLink(ctx context.Context, shortCode string) (string, error) {
	val, err := c.rdb.Get(ctx, linkKey(shortCode)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", url_domain.ErrNotFound
		}
		return "", fmt.Errorf("get link from cache: %w", err)
	}
	return val, nil
}

func (c *Cache) SetLink(ctx context.Context, shortCode, originalURL string, ttl time.Duration) error {
	err := c.rdb.Set(ctx, linkKey(shortCode), originalURL, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) DelLink(ctx context.Context, shortCode string) error {
	err := c.rdb.Del(ctx, linkKey(shortCode)).Err()
	if err != nil {
		return err
	}
	return nil
}
