package main

import (
	"context"
	"flag"
	"redisstream"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

const POOL_SIZE int = 100
const NUM_CONSUMERS int = 1
const NUM_MESSAGES int = 10000
const STREAM_COUNT int64 = 1
const STREAM_TOPIC string = "guchat:mq"
const CONSUMER_GROUP string = "guchat-group"

var addr = flag.String("addr", ":6379", "redis address")

func init() {
	// enable logger modules
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05-07:00",
	})
}

func main() {
	flag.Parse()

	ctx := context.Background()

	client := redisstream.InitRedisClient(ctx, *addr, POOL_SIZE)
	defer client.Close()

	createGroup(ctx, client)

	var wg sync.WaitGroup

	elasped := time.Now()
	wg.Add(NUM_CONSUMERS)
	for i := 0; i < NUM_CONSUMERS; i++ {
		go readMessages(ctx, &wg, client)
		// go readMessagesByGroup(ctx, &wg, client)
	}
	wg.Wait()
	logrus.Printf("Time elasped: %v ms", time.Now().Sub(elasped).Milliseconds())
}

func createGroup(ctx context.Context, client *redis.Client) {
	err := client.XGroupCreate(ctx, STREAM_TOPIC, CONSUMER_GROUP, "0").Err()
	if err != nil {
		logrus.Errorln(err)
	}
}

func readMessages(ctx context.Context, wg *sync.WaitGroup, client *redis.Client) {
	defer wg.Done()

	args := &redis.XReadArgs{
		Streams: []string{STREAM_TOPIC, "0"},
		Count:   STREAM_COUNT,
		Block:   time.Duration(100) * time.Millisecond,
	}
	for i := 0; i < NUM_MESSAGES; i++ {
		_, err := client.XRead(ctx, args).Result()
		if err != nil {
			logrus.Errorln(err)
		}
	}
}

func readMessagesByGroup(ctx context.Context, wg *sync.WaitGroup, client *redis.Client) {
	defer wg.Done()

	uuid := xid.New().String()
	args := &redis.XReadGroupArgs{
		Group:    CONSUMER_GROUP,
		Consumer: uuid,
		Streams:  []string{STREAM_TOPIC, ">"},
		Count:    STREAM_COUNT,
		Block:    0,
		NoAck:    true,
	}

	for i := 0; i < NUM_MESSAGES; i++ {
		_, err := client.XReadGroup(ctx, args).Result()
		if err != nil {
			logrus.Fatalln(err)
		}
	}
}
