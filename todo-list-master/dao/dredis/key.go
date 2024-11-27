package dredis

import "todolist/utils"

// redis key注意使用命名空间的方式,方便查询和拆分
const (
	Prefix                   = "todoList:"        // 项目key前缀
	KeyRegisterCaptchaString = "Registercaptcha:" // 注册验证码
	KeyLoginCaptchaString    = "LoginCaptcha:"    // 登录验证码
	KeyCurrencyCaptchaString = "CurrencyCaptcha:" // 通用验证,预留
	KeyTaskAndDateMap        = "TaskAndDate:"     // task    zset 日期+taskid
	KeyTaskMap               = "TaskID:"          // task    uid：+ taskid
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}

// 转化验证码对应的redis key
func KeyRedis(captcha utils.CAPTCHA) string {
	switch captcha {
	case utils.RegisterCAPTCHA:
		return KeyRegisterCaptchaString
	case utils.LoginCAPTCHA:
		return KeyLoginCaptchaString
	default:
		return KeyCurrencyCaptchaString
	}
}
