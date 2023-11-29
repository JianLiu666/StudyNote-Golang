package singlepool

import (
	"interview20231129/model"
	"testing"

	"gotest.tools/assert"
)

func TestAddSinglePersonAndMatch_add_hit(t *testing.T) {
	sp := NewSinglePool()

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
	sp := NewSinglePool()

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
	sp := NewSinglePool()

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
	sp := NewSinglePool()

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
	sp := NewSinglePool()

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
	sp := NewSinglePool()

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
	sp := NewSinglePool()

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
