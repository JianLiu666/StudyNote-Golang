package redis

import (
	"context"
	"encoding/json"
	"interview20231116/model"
	"interview20231116/pkg/config"
	"interview20231116/pkg/e"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"gotest.tools/assert"
)

func setup() (*redisClient, *miniredis.Miniredis) {
	mockServer, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	conf := config.NewFromDefault()
	conf.Redis.Address = mockServer.Addr()

	client := newRedisClient(context.Background(), &conf.Redis)

	return client, mockServer
}

func TestSetPageToListHead_Success(t *testing.T) {
	client, _ := setup()

	// prepare data
	listKey := "list-uuid"
	page := &model.Page{
		Key:         "page-uuid",
		NextPageKey: "",
		Articles:    []*model.Article{{ID: 1}, {ID: 2}},
	}

	// testcase
	code := client.SetPageToListHead(context.Background(), listKey, page)
	if code != e.SUCCESS {
		t.Error(e.GetMsg(code))
	}

	res, code := client.GetPage(context.Background(), page.Key)
	if code != e.SUCCESS {
		t.Error(e.GetMsg(code))
	}

	assert.DeepEqual(t, page, res)
}

func TestGetListHead_Success(t *testing.T) {
	client, server := setup()

	// prepare data
	listKey := "list-uuid"
	pageKey := "page-uuid"

	// init redis
	server.HSet(client.conf.ListCollectionName, listKey, pageKey)

	// testcase
	res, code := client.GetListHead(context.Background(), listKey)
	if code != e.SUCCESS {
		t.Error(e.GetMsg(code))
	}

	assert.Equal(t, pageKey, res)
}

func TestGetListHead_Fail_DataNotFound(t *testing.T) {
	client, _ := setup()

	// prepare data
	listKey := "list-uuid"

	// testcase
	_, code := client.GetListHead(context.Background(), listKey)
	if code != e.ERROR_DATA_NOT_FOUND {
		t.Error(e.GetMsg(code))
	}
}

func TestGetPage_Success(t *testing.T) {
	client, server := setup()

	// prepare data
	page := &model.Page{
		Key:         "page-uuid",
		NextPageKey: "",
		Articles:    []*model.Article{{ID: 1}, {ID: 2}},
	}
	pageJSON, err := json.Marshal(page)
	if err != nil {
		t.Error(err)
	}

	// init redis
	server.Set(client.genPageKey(page.Key), string(pageJSON))

	// testcase
	res, code := client.GetPage(context.Background(), page.Key)
	if code != e.SUCCESS {
		t.Error(e.GetMsg(code))
	}

	assert.DeepEqual(t, page, res)
}

func TestGetPage_Fail_DataNotFund(t *testing.T) {
	client, _ := setup()

	// prepare data
	page := &model.Page{
		Key:         "page1-uuid",
		NextPageKey: "",
		Articles:    []*model.Article{{ID: 1}, {ID: 2}},
	}

	// testcase
	_, code := client.GetPage(context.Background(), page.Key)
	if code != e.ERROR_DATA_NOT_FOUND {
		t.Error(e.GetMsg(code))
	}
}
