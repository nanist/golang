package logic

import (
	"go.uber.org/zap"
	"todolist/dao/dredis"
	"todolist/utils"
)

// SetCaptcha 验证码逻辑
func SetCaptcha(email string, captch utils.CAPTCHA) (err error) {
	// 1. 获取一个随机验证码
	code := utils.Getcode(6)
	ep := &utils.EmailParameter{
		To:   []string{email},
		Code: code,
	}
	// 2. 发送邮件
	if err = utils.SendCAPTCHA(ep, captch); err != nil {
		zap.L().Error("Failed to send mail", zap.Error(err))
		return
	}
	// 写入redis
	if err = dredis.SetCaptcha(email, code, captch); err != nil {
		zap.L().Error("Failed to set key to  redis", zap.Error(err))
		return
	}
	return
}

func SetCaptchaPic() (captchId string, b64s string, err error) {
	// 1. 生成图形验证码,返回id和图片
	captchId, b64s, err = utils.GetcaptchaPic()
	if err != nil {
		zap.L().Error("Failed to get a GetcaptchaPic", zap.Error(err))
	}
	return
}
