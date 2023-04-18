package redisstream

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedisClient(ctx context.Context, addr string, poolSize int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     poolSize,
	})

	return client
}
