package main

import (
	"interview20231109/question-1/api"
	"interview20231109/question-1/pkg/orm"
	"interview20231109/question-1/repository/article"
	"net/http"
)

// Top Articles
//
// In this challenge, you have to design a RESTful API,
// through the data in the database to design a query with pagination function, ordered by the number of comments.
//
// Database and repository
// Table name is `article`
// Each article record has the following schema:
//  - id: integer: paimary key, not null
//  - title: varchar(50): the title of the article, not null
//  - url: varchar(50): the URL of the article, not null
//  - author: varchar(50): the username of the author of the article, not null
//  - num_comments: the number of comments the article has, not null
//  - created_at: the date and time when record was created, not null
//
// Use the sqlite database, example `db := orm.NewConnect()`
// and you have to implement the database interface and complete the api
//
// Testing(optional) if there is enough time
// Add test cases to ensure correctness of `ArticleWithPagination` function

func init() {
	article.Init(orm.NewConnect())
}

func main() {
	s := &http.Server{
		Addr:    ":8800",
		Handler: api.InitRouter(),
	}

	s.ListenAndServe()
}
