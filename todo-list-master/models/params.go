package models

import (
	"todolist/utils"
)

// ParamEmail 发送邮箱参数
type ParamEmail struct {
	Email string `json:"email" binding:"required,email"` // 邮箱
	Type  int    `json:"type" binding:"oneof=0 1 2"`     // 通用(0) 注册（1） 登录（2）
}

// ParamSignUp 注册参数
type ParamSignUp struct {
	Name       string `json:"name" binding:"required"`                         // 用户名称
	Phone      string `json:"phone"`                                           // 用户手机号,预留
	Email      string `json:"email" binding:"required,email"`                  // 邮箱
	Password   string `json:"password" binding:"required"`                     // 密码
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` // 确认密码
	Code       string `json:"code" binding:"required"`                         // 验证码
}

// ParamTokens 刷新token参数
type ParamTokens struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ParamLogin 登录参数
type ParamLogin struct {
	Email      string      `json:"email" binding:"required,email"`
	Password   string      `json:"password" binding:"required"`
	CaptchaPic *CaptchaPic `json:"captcha_pic" binding:"required"` // 不使用验证码的话，可以去掉
}

// CaptchaPic 图形验证码校验参数
type CaptchaPic struct {
	Pid   string `json:"pid" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// ParamTask  task参数
type ParamTask struct {
	Tid         int64  `json:"tid"`
	Level       int    `json:"level"` // 一般（0） 重要（1）
	State       int    `json:"state"` //  未完成（0） 已完成（1）
	TaskContent string `json:"task_content"`
	UserID      uint   `json:"uid"`
}

type ParamDate struct {
	StartDate string `json:"start_date"  binding:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date"  binding:"required,datetime=2006-01-02"`
}

type ParamTaskIDs struct {
	TaskIDs []int64 `json:"task_ids" binding:"required"`
}

// GetCaptchType 根据type返回验证码类型
func (e *ParamEmail) GetCaptchType() utils.CAPTCHA {
	switch e.Type {
	case 1:
		return utils.RegisterCAPTCHA
	case 2:
		return utils.LoginCAPTCHA
	default:
		return utils.CurrencyCAPTCHA
	}
}
