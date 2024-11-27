package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"time"
	"todolist/controller/code"
)

var (
	// 用于签名的字符串
	CustomSecret = []byte("todolist")
)

const (
	AccessTokenExpireDuration  = time.Second * 60 * 60 * 24 * 7
	RefreshTokenExpireDuration = time.Second * 60 * 60 * 24 * 7
	AccessTokenTime            = 60 * 60 * 24 * 7
	RefreshTokenTime           = 60 * 60 * 24 * 7
)

type UserClaims struct {
	// 可根据需要自行添加字段
	UserId               uint   `json:"user_id"`
	Email                string `json:"email"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

// GensToken 生成refresh token access token
func GensTokens(userId uint, email string) (tokens map[string]map[string]interface{}, err error) {
	// 创建一个声明
	accessClaims := UserClaims{
		userId,
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpireDuration)),
			Issuer:    "to-do-list", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	accessToken, err := token.SignedString(CustomSecret)
	if err != nil {
		zap.L().Error("Failed to get access token", zap.Error(err))
		return
	}
	// refreshToken 不需要存储任何自定义数据
	refreshClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpireDuration)),
		Issuer:    "to-do-list", // 签发人
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := token.SignedString(CustomSecret)
	if err != nil {
		zap.L().Error("Failed to get refresh token ", zap.Error(err))
		return
	}
	tokens = map[string]map[string]interface{}{
		"access_token": {
			"access_token": accessToken,
			"expires_in":   AccessTokenTime,
		},
		"refresh_token": {
			"refresh_token": refreshToken,
			"expires_in":    RefreshTokenTime,
		},
	}
	return
}

// GenAccessToken 生成  access token
func GenAccessToken(userId uint, email string) (tokenMap map[string]interface{}, err error) {
	// 创建一个声明
	accessClaims := UserClaims{
		userId,
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpireDuration)),
			Issuer:    "to-do-list", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	accessToken, err := token.SignedString(CustomSecret)
	if err != nil {
		zap.L().Error("Failed to get accessToken", zap.Error(err))
	}
	tokenMap = map[string]interface{}{
		"access_token": accessToken,
		"expires_in":   AccessTokenTime,
	}
	return
}

// ParseToken 解析JWT
func ParseAccessToken(tokenStr string) (claims *UserClaims, err error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		zap.L().Error("Invalid AccessToken", zap.Error(err))
		return
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	err = code.EerrorInvalidAccessToken
	return
}

// refreshToken 刷新access token
func RefreshToken(accessToken, refreshToken string) (tokenMap map[string]interface{}, err error) {
	//  1. 判断 refreshToken 是正确的且未过期
	_, err = jwt.Parse(refreshToken, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		err = code.EerrorInvalidRefreshToken
		return
	}
	userClaims := &UserClaims{}
	// 2. accessToken中解析出 cliams 数据
	_, err = jwt.ParseWithClaims(accessToken, userClaims, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	// 3. 判断错误类型
	validationError := err.(*jwt.ValidationError)
	if validationError.Errors == jwt.ValidationErrorExpired {
		return GenAccessToken(userClaims.UserId, userClaims.Email)
	}
	return
}
