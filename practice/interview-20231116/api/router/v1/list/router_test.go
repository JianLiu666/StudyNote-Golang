package list

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"interview20231116/model"
	"interview20231116/pkg/accessor"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gofiber/fiber/v2"
	"gotest.tools/assert"
)

func setup() (*fiber.App, *accessor.Accessor, *miniredis.Miniredis) {
	redisServer, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	infra := accessor.Build()
	infra.Config.Redis.Address = redisServer.Addr()

	infra.InitKvStore(context.Background())

	app := fiber.New()
	api := app.Group("/api")
	NewListRouter(infra.KvStore).Init(api)

	return app, infra, redisServer
}

func TestGetList_Success(t *testing.T) {
	app, infra, redisServer := setup()

	// prepare data
	listKey := "list-uuid"
	pageKey := "page-uuid"

	// init redis
	redisServer.HSet(infra.Config.Redis.ListCollectionName, listKey, pageKey)

	// testcase
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/list?listKey=%v", listKey), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	res := &getListResponse{}
	err = json.Unmarshal(respBody, res)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, pageKey, res.NextPageKey)
}

func TestGetList_BadRequest_Query(t *testing.T) {
	app, _, _ := setup()

	// testcase
	req := httptest.NewRequest(fiber.MethodGet, "/api/v1/list", nil)
	resp, err := app.Test(req)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetList_BadRequest_DataNotFound(t *testing.T) {
	app, _, _ := setup()

	// testcase
	req := httptest.NewRequest(fiber.MethodGet, "/api/v1/list?listKey=nil", nil)
	resp, err := app.Test(req)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSetList_Success(t *testing.T) {
	app, _, _ := setup()

	// prepare data
	payload := setListPayload{
		ListKey:  "list-uuid",
		Articles: []*model.Article{{ID: 1}, {ID: 2}},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	// testcase
	req := httptest.NewRequest(fiber.MethodPost, "/api/v1/list", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestSetList_InternalServerError(t *testing.T) {
	app, _, _ := setup()

	// prepare data
	payload := setListPayload{
		ListKey:  "list-uuid",
		Articles: []*model.Article{{ID: 1}, {ID: 2}},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	// testcase
	req := httptest.NewRequest(fiber.MethodPost, "/api/v1/list", bytes.NewReader(data))
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
