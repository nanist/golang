package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"todolist/controller/code"
	"todolist/logic"
	"todolist/models"
)

// SignUpHadnler 用户注册
func SignUpHadnler(c *gin.Context) {
	//	 1. 校验参数
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 校验参数失败，直接返回
		zap.L().Error("SignUp with invalid param:ParamSignUp", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeInvalidParam)
		return
	}
	// 2. 调用注册逻辑
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("SignUp logic error", zap.Error(err))
		switch err {
		case code.EerrorCaptchaNotExist:
			ResponseError(c, code.CodeCaptchaNotExist)
		case code.EerrorCaptchaAtypism:
			ResponseError(c, code.CodeCaptchaAtypism)
		case code.ErrorUserExist:
			ResponseError(c, code.CodeEmailExist)
		default:
			ResponseError(c, code.CodeServerBusy)
		}
		return
	}
	// 3. 返回响应
	// 这里可以增加其他逻辑，如查询用户id并返回
	ResponseSuccess(c, nil)
}

// AccessTokenHandler 获取AccessToken
func AccessTokenHandler(c *gin.Context) {
	// 1.解析参数
	p := new(models.ParamTokens)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("AccessToken with invalid param:ParamTokens")
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeInvalidParam)
		return
	}
	// 2. 调用生成token的逻辑
	tokenMap, err := logic.AccessToken(p)
	if err != nil {
		if errors.Is(err, code.EerrorInvalidRefreshToken) {
			ResponseError(c, code.CodeInvalidRefreshToken)
			return
		}
		ResponseError(c, code.CodeServerBusy)
		return
	}
	ResponseSuccess(c, tokenMap)
}

// LoginHandler 用户登录
func LoginHandler(c *gin.Context) {
	// 1.校验参数
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("AccessToken with invalid param:ParamLogin")
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeInvalidParam)
		return
	}
	// 2. 调用生成登录逻辑
	data, err := logic.Login(p)
	if err != nil {
		zap.L().Error(" Error from check user", zap.String("email", p.Email), zap.Error(err))
		switch err {
		case code.EerrorCaptchaAtypism:
			ResponseError(c, code.CodeCaptchaAtypism)
		case code.ErrorUserExist:
			ResponseError(c, code.CodeUserNotExist)
		case code.ErrorInvalidPassword:
			ResponseError(c, code.CodeInvalidPassword)
		}
		return
	}
	ResponseSuccess(c, data)
}
