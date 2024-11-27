package logic

import (
	"go.uber.org/zap"
	"todolist/controller/code"
	"todolist/dao/dmysql"
	"todolist/dao/dredis"
	"todolist/models"
	"todolist/utils"
)

// SignUp 注册逻辑
func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 从redis中拿验证码并进行校验
	val, err := dredis.GetCaptcha(p.Email, dredis.KeyRegisterCaptchaString)
	if err != nil {
		return err
	} else if p.Code != val {
		err = code.EerrorCaptchaAtypism
		return
	}
	// 2. 校验用户是否存在
	err = dmysql.CheckUserExist(p.Email)
	if err != nil {
		return
	}
	// 3. 创建用户
	user := &models.User{Email: p.Email, Name: p.Name, Phone: p.Phone, Password: p.Password}
	err = dmysql.InsertUser(user)
	if err != nil {
		return
	}
	return
}

// AccessToken 刷新accessToken逻辑
func AccessToken(p *models.ParamTokens) (tokenMap map[string]interface{}, err error) {
	tokenMap, err = utils.RefreshToken(p.AccessToken, p.RefreshToken)
	if err != nil {
		return
	}
	return
}

// Login 登录逻辑
func Login(p *models.ParamLogin) (data map[string]interface{}, err error) {

	// 1. 校验图形验证码
	flag := VerfiyCaptcha(p.CaptchaPic.Pid, p.CaptchaPic.Value)
	if !flag {
		zap.L().Error("Failed to VerfiyCaptcha", zap.String("msg", "验证校验失败"))
		err = code.EerrorCaptchaAtypism
		return
	}
	// 2. 校验用户
	user := &models.User{
		Email:    p.Email,
		Password: p.Password,
	}
	userData, err := dmysql.Login(user)
	if err != nil {
		zap.L().Error("Failed to find  user by emial", zap.Error(err))
		return
	}

	// 3. 生成 refresh token
	tokens, err := utils.GensTokens(userData["id"].(uint), userData["email"].(string))
	data = map[string]interface{}{
		"user_info": userData,
		"tokens":    tokens,
	}
	return
}
func VerfiyCaptcha(id, value string) bool {
	if utils.VerfiyCaptcha(id, value) {
		return true
	}
	return false
}
