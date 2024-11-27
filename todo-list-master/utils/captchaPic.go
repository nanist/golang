package utils

import (
	"github.com/mojocn/base64Captcha"
)

type Captchabase64 struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

// 这里可以设置，过期时间和存储数量的等
var store = base64Captcha.DefaultMemStore
var captchaIint Captchabase64

func init() {
	var driverMath = base64Captcha.DriverMath{
		Height: 80,
		Width:  240,
		//NoiseCount: 6,
		//ShowLineOptions: 4,
	}
	captchaIint.CaptchaType = "math"
	captchaIint.DriverDigit = base64Captcha.DefaultDriverDigit
	captchaIint.DriverMath = &driverMath
}

func GetcaptchaPic() (id string, b64s string, err error) {
	var driver base64Captcha.Driver
	switch captchaIint.CaptchaType {
	case "math":
		driver = captchaIint.DriverMath.ConvertFonts()
	case "audio":
		driver = captchaIint.DriverAudio // 未自定义,可自定义
	case "string":
		driver = captchaIint.DriverString.ConvertFonts() // 未自定义,可自定义
	case "chinese":
		driver = captchaIint.DriverChinese.ConvertFonts()
	default:
		driver = captchaIint.DriverDigit
	}
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err = captcha.Generate()
	if err != nil {
		return
	}
	return
}

func VerfiyCaptcha(id, value string) bool {
	if store.Verify(id, value, true) {
		return true
	}
	return false
}
