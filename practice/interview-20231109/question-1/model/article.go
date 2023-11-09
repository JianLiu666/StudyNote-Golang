package model

import "time"

// Article
type Article struct {
	Author      string    `json:"author"`
	CreatedTime time.Time `json:"created_time"`
	ID          int       `json:"id"`
	NumComments int       `json:"num_comments"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
}
