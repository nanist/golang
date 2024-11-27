package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"todolist/controller/code"
	"todolist/logic"
	"todolist/models"
	"todolist/utils"
)

// CaptchaHandler 邮箱验证码
func CaptchaHandler(c *gin.Context) {
	// 1. 校验参数,接受需要发送的邮箱
	email := new(models.ParamEmail)
	if err := c.ShouldBindJSON(email); err != nil {
		// 校验参数失败，直接返回
		zap.L().Error("Get Captcha with invalid param:email", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeInvalidParam)
		return
	}
	//	2. 调用生成验证码逻辑
	if err := logic.SetCaptcha(email.Email, email.GetCaptchType()); err != nil {
		zap.L().Error("Failed to set  Captcha", zap.Error(err))
		switch err {
		case code.EerrorEmailTexNeedNone:
			ResponseError(c, code.CodeEmailTexNeedNone)
		case code.EerrorCaptchaNotExist:
			ResponseError(c, code.CodeCaptchaNotExist)
		}
		ResponseError(c, code.CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// 图形验证码
func CaptchaPicHandler(c *gin.Context) {
	captchId, b64s, err := logic.SetCaptchaPic()
	if err != nil {
		zap.L().Error("Failed to get a SetCaptchaPic", zap.Error(err))
		ResponseError(c, code.CodeServerBusy)
		return
	}
	data := map[string]string{
		"captchId": captchId,
		"b64s":     b64s,
	}
	ResponseSuccess(c, data)
}

func CheckCaptchaPicHandler(c *gin.Context) {
	// 1. 校验参数
	p := new(models.CaptchaPic)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Get Captcha with invalid param:CaptchaPic", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeServerBusy)
		return

	}
	// 2. 调用验证逻辑
	var data string
	if utils.VerfiyCaptcha(p.Pid, p.Value) {
		data = "验证成功"
	} else {
		data = "验证失败"
	}
	ResponseSuccess(c, data)
}
