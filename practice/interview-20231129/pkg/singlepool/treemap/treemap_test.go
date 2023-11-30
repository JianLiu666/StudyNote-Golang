package treemap

import (
	"interview20231129/model"
	"interview20231129/pkg/singlepool"
	"testing"

	"gotest.tools/assert"
)

func TestAddSinglePersonAndMatch_add_hit(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{Name: "boy1", Height: 180, Gender: 1, NumDates: 2},
	}

	groundtruth_boys := []*model.User{
		{UUID: "180-boy1", Name: "boy1", Height: 180, Gender: 1, NumDates: 2},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	for i, data := range sp.boys.Values() {
		boy := data.(*model.User)
		assert.DeepEqual(t, groundtruth_boys[i], boy)
	}
}

func TestAddSinglePersonAndMatch_add_missed(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{Name: "boy1", Height: 180, Gender: 1, NumDates: 2},
		{Name: "boy1", Height: 180, Gender: 1, NumDates: 20},
	}

	groundtruth_boys := []*model.User{
		{UUID: "180-boy1", Name: "boy1", Height: 180, Gender: 1, NumDates: 2},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	for i, data := range sp.boys.Values() {
		boy := data.(*model.User)
		assert.DeepEqual(t, groundtruth_boys[i], boy)
	}
}

func TestAddSinglePersonAndMatch_match_case1(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{Name: "boy1", Height: 180, Gender: 1, NumDates: 2},
		{Name: "boy2", Height: 180, Gender: 1, NumDates: 1},
		{Name: "boy3", Height: 180, Gender: 1, NumDates: 2},
		{Name: "girl1", Height: 155, Gender: 0, NumDates: 4},
	}

	groundtruth_boys := []*model.User{
		{UUID: "180-boy1", Name: "boy1", Height: 180, Gender: 1, NumDates: 1},
		{UUID: "180-boy3", Name: "boy3", Height: 180, Gender: 1, NumDates: 1},
	}
	groundtruth_girls := []*model.User{
		{UUID: "155-girl1", Name: "girl1", Height: 155, Gender: 0, NumDates: 1},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	for i, data := range sp.boys.Values() {
		boy := data.(*model.User)
		assert.DeepEqual(t, groundtruth_boys[i], boy)
	}

	for i, data := range sp.girls.Values() {
		girl := data.(*model.User)
		assert.DeepEqual(t, groundtruth_girls[i], girl)
	}
}

func TestAddSinglePersonAndMatch_match_case2(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{Name: "girl1", Height: 155, Gender: 0, NumDates: 2},
		{Name: "girl2", Height: 165, Gender: 0, NumDates: 2},
		{Name: "girl3", Height: 175, Gender: 0, NumDates: 2},
		{Name: "boy1", Height: 180, Gender: 1, NumDates: 4},
	}

	groundtruth_boys := []*model.User{
		{UUID: "180-boy1", Name: "boy1", Height: 180, Gender: 1, NumDates: 1},
	}
	groundtruth_girls := []*model.User{
		{UUID: "155-girl1", Name: "girl1", Height: 155, Gender: 0, NumDates: 1},
		{UUID: "165-girl2", Name: "girl2", Height: 165, Gender: 0, NumDates: 1},
		{UUID: "175-girl3", Name: "girl3", Height: 175, Gender: 0, NumDates: 1},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	for i, data := range sp.boys.Values() {
		boy := data.(*model.User)
		assert.DeepEqual(t, groundtruth_boys[i], boy)
	}

	for i, data := range sp.girls.Values() {
		girl := data.(*model.User)
		assert.DeepEqual(t, groundtruth_girls[i], girl)
	}
}

func TestAddSinglePersonAndMatch_match_case3(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{Name: "girl1", Height: 155, Gender: 0, NumDates: 2},
		{Name: "girl2", Height: 165, Gender: 0, NumDates: 2},
		{Name: "girl3", Height: 175, Gender: 0, NumDates: 2},
		{Name: "boy1", Height: 170, Gender: 1, NumDates: 4},
	}

	groundtruth_boys := []*model.User{
		{UUID: "170-boy1", Name: "boy1", Height: 170, Gender: 1, NumDates: 2},
	}
	groundtruth_girls := []*model.User{
		{UUID: "155-girl1", Name: "girl1", Height: 155, Gender: 0, NumDates: 1},
		{UUID: "165-girl2", Name: "girl2", Height: 165, Gender: 0, NumDates: 1},
		{UUID: "175-girl3", Name: "girl3", Height: 175, Gender: 0, NumDates: 2},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	for i, data := range sp.boys.Values() {
		boy := data.(*model.User)
		assert.DeepEqual(t, groundtruth_boys[i], boy)
	}

	for i, data := range sp.girls.Values() {
		girl := data.(*model.User)
		assert.DeepEqual(t, groundtruth_girls[i], girl)
	}
}

func TestRemoveSinglePerson_hit(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{Name: "boy1", Height: 170, Gender: 1, NumDates: 4},
		{Name: "girl1", Height: 175, Gender: 0, NumDates: 2},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	sp.RemoveSinglePerson(dataset[0].Name)

	assert.Equal(t, 0, len(sp.boys.Keys()))
	assert.Equal(t, 1, len(sp.girls.Keys()))
}

func TestRemoveSinglePerson_missed(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{Name: "boy1", Height: 170, Gender: 1, NumDates: 4},
		{Name: "girl1", Height: 175, Gender: 0, NumDates: 2},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	sp.RemoveSinglePerson("")

	assert.Equal(t, 1, len(sp.boys.Keys()))
	assert.Equal(t, 1, len(sp.girls.Keys()))
}

func TestQuerySinglePeople_case1(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{UUID: "170-boy1", Name: "boy1", Height: 170, Gender: 1, NumDates: 4},
		{UUID: "175-girl1", Name: "girl1", Height: 175, Gender: 0, NumDates: 2},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	result, _ := sp.QuerySinglePeople(4, &singlepool.QueryOpts{
		Name:        "",
		MinHeight:   -1,
		MaxHeight:   -1,
		Gender:      -1,
		MinNumDates: -1,
		MaxNumDates: -1,
	})

	for i, user := range result {
		assert.DeepEqual(t, dataset[i], user)
	}
}

func TestQuerySinglePeople_case2(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{UUID: "170-boy1", Name: "boy1", Height: 170, Gender: 1, NumDates: 4},
		{UUID: "175-girl1", Name: "girl1", Height: 175, Gender: 0, NumDates: 2},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	result, _ := sp.QuerySinglePeople(4, &singlepool.QueryOpts{
		Name:        "boy1",
		MinHeight:   -1,
		MaxHeight:   -1,
		Gender:      -1,
		MinNumDates: -1,
		MaxNumDates: -1,
	})

	for i, user := range result {
		assert.DeepEqual(t, dataset[i], user)
	}
}

func TestQuerySinglePeople_case3(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{UUID: "130-boy1", Name: "boy1", Height: 130, Gender: 1, NumDates: 1},
		{UUID: "140-boy2", Name: "boy2", Height: 140, Gender: 1, NumDates: 2},
		{UUID: "150-boy3", Name: "boy3", Height: 150, Gender: 1, NumDates: 3},
		{UUID: "160-boy4", Name: "boy4", Height: 160, Gender: 1, NumDates: 4},
		{UUID: "170-boy5", Name: "boy5", Height: 170, Gender: 1, NumDates: 5},
		{UUID: "180-girl1", Name: "girl1", Height: 180, Gender: 0, NumDates: 6},
		{UUID: "190-girl2", Name: "girl2", Height: 190, Gender: 0, NumDates: 7},
		{UUID: "200-girl3", Name: "girl3", Height: 200, Gender: 0, NumDates: 8},
		{UUID: "210-girl4", Name: "girl4", Height: 210, Gender: 0, NumDates: 9},
		{UUID: "220-girl5", Name: "girl5", Height: 220, Gender: 0, NumDates: 10},
	}
	groundtruth := []*model.User{
		{UUID: "160-boy4", Name: "boy4", Height: 160, Gender: 1, NumDates: 4},
		{UUID: "170-boy5", Name: "boy5", Height: 170, Gender: 1, NumDates: 5},
		{UUID: "180-girl1", Name: "girl1", Height: 180, Gender: 0, NumDates: 6},
		{UUID: "190-girl2", Name: "girl2", Height: 190, Gender: 0, NumDates: 7},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	result, _ := sp.QuerySinglePeople(5, &singlepool.QueryOpts{
		Name:        "",
		MinHeight:   160,
		MaxHeight:   190,
		Gender:      -1,
		MinNumDates: -1,
		MaxNumDates: -1,
	})

	for i, user := range result {
		assert.DeepEqual(t, groundtruth[i], user)
	}
}

func TestQuerySinglePeople_case4(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{UUID: "130-boy1", Name: "boy1", Height: 130, Gender: 1, NumDates: 1},
		{UUID: "140-boy2", Name: "boy2", Height: 140, Gender: 1, NumDates: 2},
		{UUID: "150-boy3", Name: "boy3", Height: 150, Gender: 1, NumDates: 3},
		{UUID: "160-boy4", Name: "boy4", Height: 160, Gender: 1, NumDates: 4},
		{UUID: "170-boy5", Name: "boy5", Height: 170, Gender: 1, NumDates: 5},
		{UUID: "180-girl1", Name: "girl1", Height: 180, Gender: 0, NumDates: 6},
		{UUID: "190-girl2", Name: "girl2", Height: 190, Gender: 0, NumDates: 7},
		{UUID: "200-girl3", Name: "girl3", Height: 200, Gender: 0, NumDates: 8},
		{UUID: "210-girl4", Name: "girl4", Height: 210, Gender: 0, NumDates: 9},
		{UUID: "220-girl5", Name: "girl5", Height: 220, Gender: 0, NumDates: 10},
	}
	groundtruth := []*model.User{
		{UUID: "160-boy4", Name: "boy4", Height: 160, Gender: 1, NumDates: 4},
		{UUID: "170-boy5", Name: "boy5", Height: 170, Gender: 1, NumDates: 5},
		{UUID: "180-girl1", Name: "girl1", Height: 180, Gender: 0, NumDates: 6},
		{UUID: "190-girl2", Name: "girl2", Height: 190, Gender: 0, NumDates: 7},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	result, _ := sp.QuerySinglePeople(5, &singlepool.QueryOpts{
		Name:        "",
		MinHeight:   -1,
		MaxHeight:   -1,
		Gender:      -1,
		MinNumDates: 4,
		MaxNumDates: 7,
	})

	for i, user := range result {
		assert.DeepEqual(t, groundtruth[i], user)
	}
}

func TestQuerySinglePeople_case5(t *testing.T) {
	sp := newTreemapSinglePool()

	dataset := []*model.User{
		{UUID: "130-boy1", Name: "boy1", Height: 130, Gender: 1, NumDates: 1},
		{UUID: "140-boy2", Name: "boy2", Height: 140, Gender: 1, NumDates: 2},
		{UUID: "150-boy3", Name: "boy3", Height: 150, Gender: 1, NumDates: 3},
		{UUID: "160-boy4", Name: "boy4", Height: 160, Gender: 1, NumDates: 4},
		{UUID: "170-boy5", Name: "boy5", Height: 170, Gender: 1, NumDates: 5},
		{UUID: "180-girl1", Name: "girl1", Height: 180, Gender: 0, NumDates: 6},
		{UUID: "190-girl2", Name: "girl2", Height: 190, Gender: 0, NumDates: 7},
		{UUID: "200-girl3", Name: "girl3", Height: 200, Gender: 0, NumDates: 8},
		{UUID: "210-girl4", Name: "girl4", Height: 210, Gender: 0, NumDates: 9},
		{UUID: "220-girl5", Name: "girl5", Height: 220, Gender: 0, NumDates: 10},
	}
	groundtruth := []*model.User{
		{UUID: "130-boy1", Name: "boy1", Height: 130, Gender: 1, NumDates: 1},
		{UUID: "140-boy2", Name: "boy2", Height: 140, Gender: 1, NumDates: 2},
		{UUID: "150-boy3", Name: "boy3", Height: 150, Gender: 1, NumDates: 3},
		{UUID: "160-boy4", Name: "boy4", Height: 160, Gender: 1, NumDates: 4},
		{UUID: "170-boy5", Name: "boy5", Height: 170, Gender: 1, NumDates: 5},
	}

	for _, data := range dataset {
		sp.AddSinglePersonAndMatch(data)
	}

	result, _ := sp.QuerySinglePeople(6, &singlepool.QueryOpts{
		Name:        "",
		MinHeight:   -1,
		MaxHeight:   -1,
		Gender:      1,
		MinNumDates: -1,
		MaxNumDates: -1,
	})

	for i, user := range result {
		assert.DeepEqual(t, groundtruth[i], user)
	}
}
