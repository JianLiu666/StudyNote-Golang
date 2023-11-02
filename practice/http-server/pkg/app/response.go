package app

import (
	"httpserver/pkg/e"

	"github.com/gin-gonic/gin"
)

type RespBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Response(c *gin.Context, httpCode, errCode int, data any) {
	c.JSON(httpCode, RespBody{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
}
