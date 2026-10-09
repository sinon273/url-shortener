package url_cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
)

func (c *Cache) GetCounter(ctx context.Context, userID string) (int, error) {
	val, err := c.rdb.Get(ctx, counterKey(userID)).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, url_domain.ErrNotFound
		}
		return 0, fmt.Errorf("get counter from cache: %w", err)
	}
	return val, nil
}

func (c *Cache) SetCounter(ctx context.Context, userID string, counter int) error {
	err := c.rdb.Set(ctx, counterKey(userID), counter, 0).Err()
	if err != nil {
		return fmt.Errorf("set counter to cache: %w", err)
	}
	return nil
}

func (c *Cache) IncrCounter(ctx context.Context, userID string) (int, error) {
	val, err := c.rdb.Incr(ctx, counterKey(userID)).Result()
	if err != nil {
		return 0, fmt.Errorf("incr counter from cache: %w", err)
	}
	return int(val), nil
}

func (c *Cache) DecrCounter(ctx context.Context, userID string) error {
	err := c.rdb.Decr(ctx, counterKey(userID)).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return url_domain.ErrNotFound
		}
		return fmt.Errorf("decr counter from cache: %w", err)
	}
	return nil
}
