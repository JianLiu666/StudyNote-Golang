package v1

import (
	"httpserver/models"
	"httpserver/pkg/e"
	"httpserver/pkg/setting"
	"httpserver/pkg/util"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/validation"
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
	name := c.Query("name")
	state, _ := strconv.Atoi(c.DefaultQuery("state", "0"))
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名稱不能為空")
	valid.MaxSize(name, 100, "name").Message("名稱最長為100字符")
	valid.Required(createdBy, "created_by").Message("創建人不能為空")
	valid.MaxSize(createdBy, 100, "created_by").Message("創建人最長為100字符")
	valid.Range(state, 0, 1, "state").Message("狀態只允許0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func EditTags(c *gin.Context) {

}

func DeleteTags(c *gin.Context) {

}
