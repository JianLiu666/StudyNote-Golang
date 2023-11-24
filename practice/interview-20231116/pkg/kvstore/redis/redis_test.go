package redis

import (
	"context"
	"encoding/json"
	"interview20231116/model"
	"interview20231116/pkg/config"
	"interview20231116/pkg/e"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

var _redisClient redisClient
var _redisMock redismock.ClientMock

func TestMain(m *testing.M) {
	setup()

	m.Run()
}

func setup() {
	conn, mock := redismock.NewClientMock()

	_redisClient = redisClient{
		conn: conn,
		conf: &config.NewFromDefault().Redis,
	}

	_redisMock = mock
}

func TestGetListHead_Success(t *testing.T) {
	listKey := "list-uuid"
	pageKey := "page-uuid"

	_redisMock.ExpectHGet(_redisClient.conf.ListCollectionName, listKey).SetVal(pageKey)

	res, code := _redisClient.GetListHead(context.Background(), listKey)
	if code != e.SUCCESS {
		t.Error(e.GetMsg(code))
	}

	assert.Equal(t, pageKey, res)
}

func TestGetListHead_Fail_DataNotFound(t *testing.T) {
	listKey := "list-uuid"

	_redisMock.ExpectHGet(_redisClient.conf.ListCollectionName, listKey).RedisNil()

	_, code := _redisClient.GetListHead(context.Background(), listKey)
	if code != e.ERROR_DATA_NOT_FOUND {
		t.Error(e.GetMsg(code))
	}
}

func TestGetPage_Success(t *testing.T) {
	page := &model.Page{
		Key:         "page-uuid",
		NextPageKey: "",
		Articles:    []*model.Article{{ID: 1}, {ID: 2}},
	}
	pageJSON, err := json.Marshal(page)
	if err != nil {
		t.Error(err)
	}

	_redisMock.ExpectGet(_redisClient.genPageKey(page.Key)).SetVal(string(pageJSON))

	res, code := _redisClient.GetPage(context.Background(), page.Key)
	if code != e.SUCCESS {
		t.Error(e.GetMsg(code))
	}

	assert.Equal(t, page, res)
}

func TestGetPage_Fail_DataNotFund(t *testing.T) {
	page := &model.Page{
		Key:         "page1-uuid",
		NextPageKey: "",
		Articles:    []*model.Article{{ID: 1}, {ID: 2}},
	}

	_redisMock.ExpectGet(_redisClient.genPageKey(page.Key)).RedisNil()

	_, code := _redisClient.GetPage(context.Background(), page.Key)
	if code != e.ERROR_DATA_NOT_FOUND {
		t.Error(e.GetMsg(code))
	}
}
