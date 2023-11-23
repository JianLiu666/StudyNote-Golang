package model

type List struct {
	Key         string `json:"key"`
	NextPageKey string `json:"nextPageKey"`
}

type Page struct {
	Key         string     `json:"key"`
	NextPageKey string     `json:"nextPageKey"`
	Articles    *[]Article `json:"articles"`
}

type Article struct {
	ID int `json:"id"`
}
