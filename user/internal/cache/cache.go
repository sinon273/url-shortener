package user_redis_cache

import "github.com/redis/go-redis/v9"

type Cache struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Cache {
	return &Cache{rdb: rdb}
}
