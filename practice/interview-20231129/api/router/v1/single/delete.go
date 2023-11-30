package single

import (
	"interview20231129/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *singleRouter) RemoveSinglePerson(c *gin.Context) {
	name := c.Param("name")

	if code := s.singlePool.RemoveSinglePerson(name); code != e.SUCCESS {
		c.String(http.StatusBadRequest, e.GetMsg(code))
		return
	}

	c.Status(http.StatusOK)
}
