package code

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeEmailExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
	CodeCaptchaAtypism
	CodeCaptchaNotExist
	CodeEmailTexNeedNone
	CodeInvalidAccessToken
	CodeInvalidRefreshToken
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:             "success",
	CodeInvalidParam:        "请求参数错误",
	CodeEmailExist:          "邮箱已存在",
	CodeUserNotExist:        "用户名不存在",
	CodeInvalidPassword:     "邮箱或密码错误",
	CodeServerBusy:          "服务繁忙",
	CodeNeedLogin:           "需要登录",
	CodeInvalidToken:        "无效的token",
	CodeCaptchaAtypism:      "验证码不一致",
	CodeCaptchaNotExist:     "验证码已过期",
	CodeEmailTexNeedNone:    "邮箱内容需要为空",
	CodeInvalidAccessToken:  "无效access token",
	CodeInvalidRefreshToken: "无效refresh token",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if ok {
		return msg
	}
	return codeMsgMap[CodeServerBusy]
}
