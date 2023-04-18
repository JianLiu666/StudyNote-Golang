package main

import (
	"context"
	"redisstream"
	"time"

	"github.com/go-redis/redis/v8"
)

const POOL_SIZE int = 100
const NUM_MESSAGES int = 1000

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
