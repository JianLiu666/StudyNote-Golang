package router

import "github.com/gin-gonic/gin"

type Router interface {
	Init(r *gin.RouterGroup)
}
