package single

import (
	"interview20231129/pkg/e"
	"interview20231129/pkg/singlepool"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *singleRouter) QuerySinglePeople(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	minHeight, _ := strconv.Atoi(c.DefaultQuery("minHeight", "-1"))
	maxHeight, _ := strconv.Atoi(c.DefaultQuery("maxHeight", "-1"))
	gender, _ := strconv.Atoi(c.DefaultQuery("gender", "-1"))
	minNumDates, _ := strconv.Atoi(c.DefaultQuery("minNumDates", "-1"))
	maxNumDates, _ := strconv.Atoi(c.DefaultQuery("maxNumDates", "-1"))
	n, _ := strconv.Atoi(c.DefaultQuery("n", "10"))

	opts := &singlepool.QueryOpts{
		Name:        name,
		MinHeight:   minHeight,
		MaxHeight:   maxHeight,
		Gender:      gender,
		MinNumDates: minNumDates,
		MaxNumDates: maxNumDates,
	}

	res, code := s.singlePool.QuerySinglePeople(n, opts)
	if code != e.SUCCESS {
		c.String(http.StatusBadRequest, e.GetMsg(code))
	}

	c.JSON(http.StatusOK, res)
}
