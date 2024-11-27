package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todolist/controller/code"
)

type ResponseData struct {
	Code code.ResCode `json:"code"`
	Msg  interface{}  `json:"msg"`
	Data interface{}  `json:"data"`
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code.CodeSuccess,
		Msg:  code.CodeSuccess.Msg(),
		Data: data,
	})
}

func ResponseError(c *gin.Context, code code.ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code code.ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
