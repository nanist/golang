package code

import "errors"

var (
	ErrorUserExist            = errors.New("用户已存在")
	ErrorUserNotExist         = errors.New("用户不存在")
	ErrorUserNotLogin         = errors.New("用户未登录")
	ErrorInvalidPassword      = errors.New("账号或密码错误")
	EerrorCaptchaAtypism      = errors.New("验证码不一致")
	EerrorCaptchaNotExist     = errors.New("验证码已过期")
	EerrorEmailTexNeedNone    = errors.New("邮箱的Text必须为空,请检查传入的参数")
	EerrorInvalidAccessToken  = errors.New("无效access token")
	EerrorInvalidRefreshToken = errors.New("无效refresh token")
	ErrorRedisKeyNotExist     = errors.New("redis key 不存在")
)
