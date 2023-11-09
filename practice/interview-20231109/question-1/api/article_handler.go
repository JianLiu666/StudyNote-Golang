package api

import (
	"interview20231109/question-1/repository"
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

type articleHandler struct {
	repo repository.ArticleRepo
}

func NewArticleHandler(s *gin.Engine, repo repository.ArticleRepo) {
	handler := &articleHandler{
		repo: repo,
	}

	s.GET("/article", handler.getWithPagination)

	// e.g.
	// s.POST("/article", handler.create)
	// s.GET("/article/:id", handler.getByID)
	// s.DELETE("/article/:id", handler.delete)
}

func (a *articleHandler) getWithPagination(c *gin.Context) {
	page := 1
	if val := c.Query("page"); val != "" {
		page, _ = strconv.Atoi(val)
	}

	size := 10
	if val := c.Query("size"); val != "" {
		size, _ = strconv.Atoi(val)
	}

	result := a.repo.GetWithPagination(page, size)

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
