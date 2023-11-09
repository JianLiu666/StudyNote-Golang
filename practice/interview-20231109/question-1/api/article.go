package api

import (
	"interview20231109/question-1/repository/article"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NOTE:
// You can use frameworks you are familiar with, such as gin, echo, iris and other.
//
// Design a function that can use paging query (ordered by the number of comments)
// default page = 1, size = 10
//
// example: GET /article?page=x&size=x
// result `{data: [{model.ARticle}]}`

func articleWithPagination(c *gin.Context) {
	page := 1
	if val := c.Query("page"); val != "" {
		page, _ = strconv.Atoi(val)
	}

	size := 10
	if val := c.Query("size"); val != "" {
		size, _ = strconv.Atoi(val)
	}

	result := article.ArticleWithPagination(page, size)

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
