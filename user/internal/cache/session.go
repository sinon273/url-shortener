package user_redis_cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sinon273/url-shortener/user/internal/domain"
)

func (c *Cache) SetSession(ctx context.Context, sessionID string, userID string, ttl time.Duration) error {
	key := "session:" + sessionID
	err := c.rdb.Set(ctx, key, userID, ttl).Err()
	if err != nil {
		return fmt.Errorf("set session: %w", err)
	}
	return nil
}

func (c *Cache) GetSession(ctx context.Context, sessionID string) (string, error) {
	key := "session:" + sessionID
	userID, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrNotFound
		}
		return "", fmt.Errorf("get session: %w", err)
	}
	return userID, nil
}

func (c *Cache) DelSession(ctx context.Context, sessionID string) error {
	key := "session:" + sessionID
	err := c.rdb.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("del session: %w", err)
	}
	return nil
}
