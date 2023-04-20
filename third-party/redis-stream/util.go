package redisstream

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedisClient(ctx context.Context, addr string, poolSize int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		PoolSize:     poolSize,
	})

	return client
}
