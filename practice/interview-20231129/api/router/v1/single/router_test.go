package single

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"interview20231129/model"
	"interview20231129/pkg/accessor"
	"interview20231129/pkg/singlepool/treemap"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/assert"
)

func setup() *gin.Engine {
	infra := accessor.Build()
	infra.InitSinglePool(context.Background(), treemap.NewTreemapSinglePool())

	router := gin.Default()
	api := router.Group("/api")
	NewSingleRouter(infra.SinglePool).Init(api)

	return router
}

func TestAddSinglePersonAndMatch_Success(t *testing.T) {
	router := setup()

	// prepare data
	user := &model.User{
		Name:     "boy",
		Height:   188,
		Gender:   1,
		NumDates: 10,
	}
	data, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
	}

	// testcase
	req := httptest.NewRequest(http.MethodPost, "/api/v1/singles", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestAddSinglePersonAndMatch_InvalidPaser(t *testing.T) {
	router := setup()

	// prepare data
	user := &model.User{
		Name:     "boy",
		Height:   188,
		Gender:   1,
		NumDates: 10,
	}
	data, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
	}
	data = data[:len(data)-1]

	// testcase
	req := httptest.NewRequest(http.MethodPost, "/api/v1/singles", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestAddSinglePersonAndMatch_DuplicatedUser(t *testing.T) {
	router := setup()

	// prepare data
	user := &model.User{
		Name:     "boy",
		Height:   188,
		Gender:   1,
		NumDates: 10,
	}
	data, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
	}

	// testcase
	req := httptest.NewRequest(http.MethodPost, "/api/v1/singles", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/singles", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRemoveSinglePerson_Success(t *testing.T) {
	router := setup()

	// prepare data
	user := &model.User{
		Name:     "boy",
		Height:   188,
		Gender:   1,
		NumDates: 10,
	}
	data, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
	}

	// testcase
	req := httptest.NewRequest(http.MethodPost, "/api/v1/singles", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/singles/%s", user.Name), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRemoveSinglePerson_UserNotFound(t *testing.T) {
	router := setup()

	// prepare data
	user := &model.User{
		Name:     "boy",
		Height:   188,
		Gender:   1,
		NumDates: 10,
	}

	// testcase
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/singles/%s", user.Name), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestQuerySinglePeople_Success(t *testing.T) {
	router := setup()

	// prepare data
	user := &model.User{
		UUID:     "188-boy",
		Name:     "boy",
		Height:   188,
		Gender:   1,
		NumDates: 10,
	}
	data, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
	}

	// testcase
	req := httptest.NewRequest(http.MethodPost, "/api/v1/singles", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/singles", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	result := []*model.User{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, len(result))
	assert.DeepEqual(t, user, result[0])
}
