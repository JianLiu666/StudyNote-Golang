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

	maps["state"] = -1
	if arg := c.Query("state"); arg != "" {
		state, _ := strconv.Atoi(arg)
		maps["state"] = state
	}

	data["lists"] = models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": data,
	})
}

type AddTagsForm struct {
	Name      string `json:"name" binding:"required,max=100"`
	State     int    `json:"state" binding:"gte=0,lte=1"`
	CreatedBy string `json:"created_by" binding:"required,max=100"`
}

func AddTags(c *gin.Context) {
	var form AddTagsForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  e.GetMsg(e.INVALID_PARAMS),
			"date": make(map[string]string),
		})
		return
	}

	code := e.ERROR_EXIST_TAG
	if !models.ExistTagByName(form.Name) {
		code = e.SUCCESS
		models.AddTag(form.Name, form.State, form.CreatedBy)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func EditTags(c *gin.Context) {
	// id, _ := strconv.Atoi(c.Param("id"))
	// name := c.Query("name")
	// modifiedBy := c.Query("modified_by")

	// valid := validation.Validation{}

}

func DeleteTags(c *gin.Context) {

}
