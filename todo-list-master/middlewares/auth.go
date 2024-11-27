package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"todolist/controller"
	"todolist/controller/code"
	"todolist/utils"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 从请求体里面获取token
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			zap.L().Warn("Error Authorization")
			controller.ResponseError(c, code.CodeNeedLogin)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 && parts[0] != "Bearer" {
			zap.L().Warn("Error Authorization")
			controller.ResponseError(c, code.CodeInvalidToken)
			c.Abort()
			return
		}
		claims, err := utils.ParseAccessToken(parts[1])
		if err != nil {
			zap.L().Warn("Error Authorization")
			controller.ResponseError(c, code.CodeInvalidToken)
			c.Abort()
			return
		}
		userId, email := claims.UserId, claims.Email
		// 将当前请求的用户信息保存到请求的上下文c上
		c.Set(controller.CtxUserIDKey, userId)
		c.Set(controller.CtxUserEmail, email)
	}
}
