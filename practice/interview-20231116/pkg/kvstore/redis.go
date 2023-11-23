package kvstore

import (
	"context"
	"encoding/json"
	"interview20231116/model"
	"interview20231116/pkg/config"
	"interview20231116/pkg/e"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var _ KvStore = (*redisClient)(nil)

type redisClient struct {
	conn *redis.Client
	conf *config.RedisOpts
}

func NewRedisClient(ctx context.Context, conf *config.RedisOpts) KvStore {
	conn := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.DB,
		PoolSize: conf.PoolSize,
	})

	ct, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if _, err := conn.Ping(ct).Result(); err != nil {
		logrus.Panicf("failed to ping redis server: %v", err)
	}

	return &redisClient{
		conn: conn,
		conf: conf,
	}
}

func (c *redisClient) Shutdown(ctx context.Context) {
	if err := c.conn.Close(); err != nil {
		logrus.Errorf("failed to close redis client: %v", err)
	}
}

func (c *redisClient) SetPageToListHead(ctx context.Context, listKey string, page *model.Page) e.CODE {
	script := `
		local hashKey = KEYS[1]
		local listKey = KEYS[2]
		local pageKey = KEYS[3]
		
		-- step.1 如果 hash key 跟 list key 不存在的話, 給定初始值
		redis.call('HSETNX', hashKey, listKey, '')
		
		-- step.2 取出 list 的 head page
		local nextPageKey = redis.call('HGET', hashKey, listKey)

		-- step.3 更新 linked list
		local page = cjson.decode(ARGV[1])
		page.nextPageKey = nextPageKey

		redis.call('HSET', hashKey, listKey, page.key)

		-- step.4 寫入 page
		local pageJSON = cjson.encode(page)
		redis.call('SET', pageKey, pageJSON)

		return 1
	`

	pageJSON, err := json.Marshal(page)
	if err != nil {
		logrus.Errorf("failed to execute json.Marshal: %v", err)
		return e.ERROR_MARSHAL
	}

	keys := []string{c.conf.ListCollectionName, listKey, c.genPageKey(page.Key)}
	res, err := c.conn.Eval(ctx, script, keys, pageJSON).Result()
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
	res, err := c.conn.HGet(ctx, c.conf.ListCollectionName, listKey).Result()
	if err != nil {
		logrus.Errorf("failed to execute redis command HGet: %v", err)
		return "", e.ERROR_REDIS_COMMAND
	}

	return res, e.SUCCESS
}

func (c *redisClient) GetPage(ctx context.Context, pageKey string) (*model.Page, e.CODE) {
	res, err := c.conn.Get(ctx, c.genPageKey(pageKey)).Result()
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

func (c *redisClient) genPageKey(pageKey string) string {
	return c.conf.PageKeyPrefix + pageKey
}
