package kvstore

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var _ KvStore = (*redisClient)(nil)

type redisClient struct {
	conn *redis.Client
}

func NewRedisClient(ctx context.Context, addr, password string, db int) KvStore {
	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 10, // TODO: remove magic number
	})

	ct, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if _, err := conn.Ping(ct).Result(); err != nil {
		logrus.Panicf("failed to ping redis server: %v", err)
	}

	return &redisClient{
		conn: conn,
	}
}

func (c *redisClient) Shutdown(ctx context.Context) {
	if err := c.conn.Close(); err != nil {
		logrus.Errorf("failed to close redis client: %v", err)
	}
}
