package single

import (
	"interview20231129/model"
	"interview20231129/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (s *singleRouter) AddSinglePersonAndMatch(c *gin.Context) {
	user := &model.User{}

	if err := c.BindJSON(user); err != nil {
		logrus.Errorf("Failed to execute c.BindJSON: %v", err)
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}

	if code := s.singlePool.AddSinglePersonAndMatch(user); code != e.SUCCESS {
		c.String(http.StatusBadRequest, e.GetMsg(code))
		return
	}

	c.Status(http.StatusOK)
}
