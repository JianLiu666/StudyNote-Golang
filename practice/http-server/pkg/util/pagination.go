package util

import (
	"httpserver/pkg/setting"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) int {
	result := 0

	page := c.Query("page")
	if page == "" {
		return result
	}

	num, err := strconv.Atoi(page)
	if err != nil {
		return result
	}

	if num > 0 {
		result = (num - 1) * setting.AppSetting.PageSize
	}

	return result
}
