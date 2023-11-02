package routers

import (
	v1 "httpserver/routers/api/v1"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", ping)

	apiv1 := r.Group("/api/v1")
	apiv1.GET("/tags", v1.GetTags)
	apiv1.POST("/tags", v1.AddTag)
	apiv1.PUT("/tags/:id", v1.EditTag)
	apiv1.DELETE("/tags/:id", v1.DeleteTags)

	return r
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
