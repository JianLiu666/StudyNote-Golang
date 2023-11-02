package v1

import (
	"httpserver/models"
	"httpserver/pkg/app"
	"httpserver/pkg/e"
	"httpserver/pkg/setting"
	"httpserver/pkg/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	maps := make(map[string]any)
	data := make(map[string]any)

	name := c.Query("name")
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

	app.Response(c, http.StatusOK, e.SUCCESS, data)
}

type AddTagsForm struct {
	Name      string `json:"name"       valid:"Required;MaxSize(100)"`
	State     int    `json:"state"      valid:"Range(0,1)"`
	CreatedBy string `json:"created_by" valid:"Required;MaxSize(100)"`
}

func AddTags(c *gin.Context) {
	var form AddTagsForm

	if httpCode, errCode := app.BindAndValid(c, &form); errCode != e.SUCCESS {
		app.Response(c, httpCode, errCode, nil)
		return
	}

	exists, err := models.ExistTagByName(form.Name)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		app.Response(c, http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	if err := models.AddTag(form.Name, form.State, form.CreatedBy); err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, nil)
}

type EditTagsForm struct {
	ID         int `json:"" valid:`
	Name       string
	ModifiedBy string
	State      int
}

func EditTags(c *gin.Context) {
	// id, _ := strconv.Atoi(c.Param("id"))
	// name := c.Query("name")
	// modifiedBy := c.Query("modified_by")

	// valid := validation.Validation{}

}

func DeleteTags(c *gin.Context) {

}
