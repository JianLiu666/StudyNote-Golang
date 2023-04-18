package main

import (
	"context"
	"redisstream"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

const POOL_SIZE int = 100
const NUM_MESSAGES int = 10000

func init() {
	// enable logger modules
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05-07:00",
	})
}

func main() {
	ctx := context.Background()

	client := redisstream.InitRedisClient(ctx, ":6379", POOL_SIZE)
	defer client.Close()

	for i := 0; i < NUM_MESSAGES; i++ {
		args := redis.XAddArgs{
			Stream: "guchat:mq",
			Values: map[string]interface{}{
				"timestamp": time.Now().UnixMicro(),
				"field1":    true,
				"field2":    100,
				"field3":    "string",
				"field4":    "stringstringstringstringstringstring",
				"field5":    "stringstringstringstringstringstringstringstringstring",
			},
		}
		client.XAdd(ctx, &args)
	}
}
