package v1

import (
	"httpserver/models"
	"httpserver/pkg/e"
	"httpserver/pkg/setting"
	"httpserver/pkg/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]any)
	data := make(map[string]any)

	if name != "" {
		maps["name"] = name
	}

	if arg := c.Query("state"); arg != "" {
		state, _ := strconv.Atoi(arg)
		maps["state"] = state
	} else {
		maps["state"] = -1
	}

	code := e.SUCCESS

	data["lists"] = models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func AddTags(c *gin.Context) {

}

func EditTags(c *gin.Context) {

}

func DeleteTags(c *gin.Context) {

}
