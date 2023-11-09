package api

import (
	"interview20231109/question-1/pkg/orm"
	"interview20231109/question-1/repository/article"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	s := gin.New()

	s.GET("/ping", ping)

	NewArticleHandler(s, article.NewArticleService(orm.NewConnect()))

	return s
}

func ping(c *gin.Context) {
	c.String(200, "pong")
}
