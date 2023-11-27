package page

import (
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
	NewPageRouter(infra.KvStore).Init(api)

	return app, infra, redisServer
}

func TestGetPage_Success(t *testing.T) {
	app, infra, redisServer := setup()

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
	redisServer.Set(infra.Config.Redis.PageKeyPrefix+page.Key, string(pageJSON))

	// testcase
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/page?pageKey=%v", page.Key), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	res := &getPageResponse{}
	err = json.Unmarshal(respBody, res)
	if err != nil {
		t.Error(err)
	}

	assert.DeepEqual(t, page.NextPageKey, res.NextPageKey)
	assert.DeepEqual(t, page.Articles, res.Articles)
}

func TestGetPage_BadRequest_Query(t *testing.T) {
	app, _, _ := setup()

	// testcase
	req := httptest.NewRequest(fiber.MethodGet, "/api/v1/page", nil)
	resp, err := app.Test(req)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetPage_BadRequest_DataNotFound(t *testing.T) {
	app, _, _ := setup()

	// testcase
	req := httptest.NewRequest(fiber.MethodGet, "/api/v1/page?pageKey=nil", nil)
	resp, err := app.Test(req)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
