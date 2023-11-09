package api

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	s := gin.New()

	s.GET("/ping", ping)
	s.GET("/article", articleWithPagination)

	return s
}
