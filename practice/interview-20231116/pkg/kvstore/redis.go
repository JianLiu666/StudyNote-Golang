package kvstore

import (
	"context"
	"encoding/json"
	"interview20231116/model"
	"interview20231116/pkg/e"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
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

func (c *redisClient) SetPageToListHead(ctx context.Context, listKey string, page *model.Page) e.CODE {
	script := `
		local listKey = KEYS[1]
		local oldPageKey = redis.call('HGET', 'list', listKey)
		
		if not oldPageKey then
			return -1
		end

		local page = cjson.decode(ARGV[1])
		page.nextPageKey = oldPageKey

		redis.call('HSET', 'list', listKey, page.key)

		local pageJSON = cjson.encode(page)
		redis.call('SET', 'page/' .. page.key, pageJSON)

		return 1
	`

	page.Key = xid.New().String()
	pageJSON, err := json.Marshal(page)
	if err != nil {
		logrus.Errorf("failed to execute json.Marshal: %v", err)
		return e.ERROR_MARSHAL
	}

	res, err := c.conn.Eval(ctx, script, []string{listKey}, pageJSON).Result()
	if err != nil {
		logrus.Errorf("failed to execute redis command Eval: %v", err)
		return e.ERROR_REDIS_COMMAND
	}

	if res.(int64) == -1 {
		logrus.Errorf("list key not found: %v", listKey)
		return e.ERROR_DATA_NOT_FOUND
	}

	return e.SUCCESS
}

func (c *redisClient) GetListHead(ctx context.Context, listKey string) (string, e.CODE) {
	// TODO: remove hardcore string
	res, err := c.conn.HGet(ctx, "list", listKey).Result()
	if err != nil {
		logrus.Errorf("failed to execute redis command HGet: %v", err)
		return "", e.ERROR_REDIS_COMMAND
	}

	return res, e.SUCCESS
}

func (c *redisClient) GetPage(ctx context.Context, pageKey string) (*model.Page, e.CODE) {
	res, err := c.conn.Get(ctx, "page/"+pageKey).Result()
	if err != nil {
		logrus.Errorf("failed to execute redis command Get: %v", err)
		return nil, e.ERROR_REDIS_COMMAND
	}

	page := &model.Page{}
	err = json.Unmarshal([]byte(res), page)
	if err != nil {
		return nil, e.ERROR_UNMARSHAL
	}

	return page, e.SUCCESS
}
