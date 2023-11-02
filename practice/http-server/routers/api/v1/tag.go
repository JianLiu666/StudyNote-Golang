package v1

import (
	"httpserver/models"
	"httpserver/pkg/app"
	"httpserver/pkg/e"
	"httpserver/pkg/setting"
	"httpserver/pkg/util"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	maps := make(map[string]any)

	name := c.Query("name")
	if name != "" {
		maps["name"] = name
	}

	maps["state"] = -1
	if arg := c.Query("state"); arg != "" {
		state, _ := strconv.Atoi(arg)
		maps["state"] = state
	}

	lists, err := models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := models.GetTagTotal(maps)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, map[string]any{
		"lists": lists,
		"total": count,
	})
}

type AddTagForm struct {
	Name      string `json:"name"       valid:"Required;MaxSize(100)"`
	State     int    `json:"state"      valid:"Range(0,1)"`
	CreatedBy string `json:"created_by" valid:"Required;MaxSize(100)"`
}

func AddTag(c *gin.Context) {
	var form AddTagForm

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

type EditTagForm struct {
	ID         int    `json:"id"          valid:"Required;Min(1)"`
	Name       string `json:"name"        valid:"Required;MaxSize(100)"`
	State      int    `json:"state"       valid:"Range(0,1)"`
	ModifiedBy string `json:"modified_by" valid:"Required;MaxSize(100)"`
}

func EditTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	form := EditTagForm{
		ID: id,
	}

	if httpCode, errCode := app.BindAndValid(c, &form); errCode != e.SUCCESS {
		app.Response(c, httpCode, errCode, nil)
		return
	}

	exists, err := models.ExistTagByID(form.ID)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := models.EditTag(form.ID, form); err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}
	app.Response(c, http.StatusOK, e.SUCCESS, nil)
}

func DeleteTags(c *gin.Context) {
	valid := validation.Validation{}
	id, _ := strconv.Atoi(c.Param("id"))
	valid.Min(id, 1, "id").Message("ID必須大於0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	exists, err := models.ExistTagByID(id)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := models.DeleteTag(id); err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, nil)
}
