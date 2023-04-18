package main

import (
	"context"
	"flag"
	"fmt"
	"redisstream"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const POOL_SIZE int = 100
const NUM_CONSUMERS int = 100
const NUM_MESSAGES int = 1000
const STREAM_COUNT int64 = 1

var addr = flag.String("addr", ":6379", "redis address")

func main() {
	flag.Parse()

	ctx := context.Background()

	client := redisstream.InitRedisClient(ctx, *addr, POOL_SIZE)
	defer client.Close()

	var wg sync.WaitGroup

	elasped := time.Now()
	wg.Add(NUM_CONSUMERS)
	for i := 0; i < NUM_CONSUMERS; i++ {
		go readMessages(ctx, &wg, client)
	}
	wg.Wait()
	fmt.Printf("Time elasped: %vms\n", time.Now().Sub(elasped).Milliseconds())
}

func readMessages(ctx context.Context, wg *sync.WaitGroup, client *redis.Client) {
	defer wg.Done()

	args := &redis.XReadArgs{
		Streams: []string{"guchat:mq", "0"},
		Count:   STREAM_COUNT,
		Block:   time.Duration(2) * time.Second,
	}
	for i := 0; i < NUM_MESSAGES; i++ {
		_, err := client.XRead(ctx, args).Result()
		if err != nil {
			panic(err)
		}
		// fmt.Println(res)
	}
}
