package single

import (
	"interview20231129/api/router"
	"interview20231129/pkg/singlepool"

	"github.com/gin-gonic/gin"
)

type singleRouter struct {
	singlePool singlepool.SinglePool
}

func NewSingleRouter(singlePool singlepool.SinglePool) router.Router {
	return &singleRouter{
		singlePool: singlePool,
	}
}

func (s *singleRouter) Init(r *gin.RouterGroup) {
	v1 := r.Group("/v1")

	v1.GET("/singles", s.QuerySinglePeople)
	v1.POST("/singles", s.AddSinglePersonAndMatch)
	v1.DELETE("/singles/:name", s.RemoveSinglePerson)
}
