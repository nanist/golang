package utils

import (
	"fmt"
	"github.com/jordan-wright/email"
	"go.uber.org/zap"
	"net/smtp"
	"net/textproto"
	"todolist/controller/code"
	"todolist/setting"
)

type EmailParameter struct {
	To      []string
	From    string
	Subject string
	Text    string
	HTML    string
	Code    string
}

type CAPTCHA int64
type CaptchaMap map[string]string

const (
	RegisterCAPTCHA CAPTCHA = 2000 + iota
	LoginCAPTCHA
	CurrencyCAPTCHA
)

var CaptchaMsgMap = map[CAPTCHA]map[string]string{
	RegisterCAPTCHA: CaptchaMap{"Text": "to-do-list的注册验证码：%s,\n验证码五分钟内有效,请抓紧注册。", "Subject": "to-do-list注册验证码"},
	LoginCAPTCHA:    CaptchaMap{"Text": "to-do-list的登录验证码：%s,\n验证码五分钟内有效,请抓紧登录。", "Subject": "to-do-list登录验证码"},
	CurrencyCAPTCHA: CaptchaMap{"Text": "to-do-list的登录验证码:%s,\n验证码五分钟内有效。", "Subject": "to-do-list验证码"},
}

func (captcha CAPTCHA) Msg() CaptchaMap {
	msg, ok := CaptchaMsgMap[captcha]
	if !ok {
		return CaptchaMsgMap[CurrencyCAPTCHA]
	}
	return msg
}

func SendCAPTCHA(ep *EmailParameter, codeType CAPTCHA) (err error) {
	var e = &email.Email{}
	if ep.Text == "" {
		text := fmt.Sprintf(codeType.Msg()["Text"], ep.Code)
		e = &email.Email{
			To:      ep.To,
			From:    "to-do-list <xxxxxx@qq.com>",
			Subject: codeType.Msg()["Subject"],
			Text:    []byte(text),
			//HTML:    []byte(ep.HTML),
			Headers: textproto.MIMEHeader{},
		}
	} else {
		err = code.EerrorEmailTexNeedNone
		zap.L().Error("The text of the mailbox must be empty", zap.Error(err))
		return
	}
	//设置服务器相关的配置
	err = e.Send(setting.Appconfig.EmailConfig.Addr, smtp.PlainAuth(setting.Appconfig.EmailConfig.Identity, setting.Appconfig.EmailConfig.Username, setting.Appconfig.EmailConfig.Password, setting.Appconfig.EmailConfig.Host))
	if err != nil {
		zap.L().Error("Failed to send Email", zap.String("msg", "邮箱发送失败"), zap.Error(err))
		return
	}
	return
}
