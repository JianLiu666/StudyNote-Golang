package single

import (
	"interview20231129/api/router"

	"github.com/gin-gonic/gin"
)

type singleRouter struct {
}

func NewSingleRouter() router.Router {
	return &singleRouter{}
}

func (s *singleRouter) Init(r *gin.RouterGroup) {
	v1 := r.Group("/v1")

	v1.POST("/singles", s.AddSinglePersonAndMatch)
	v1.DELETE("/singles/:name", s.RemoveSinglePerson)
}
