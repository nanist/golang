package router

import (
	"time"
	"todolist/controller"
	"todolist/logger"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/", func(c *gin.Context) {
		controller.ResponseSuccess(c, gin.H{
			"time": time.Now(),
		})
	})
	v1 := r.Group("/api/v1")
	v1.GET("/captcha", controller.CaptchaHandler)              // 获取邮箱验证码   用于注册
	v1.GET("/captcha/pic", controller.CaptchaPicHandler)       // 获取图片验证码   用于登录
	v1.POST("/captcha/pic", controller.CheckCaptchaPicHandler) // 获取图片验证码校验   用于登录
	v1.POST("/signup", controller.SignUpHadnler)               // 注册
	v1.POST("/login", controller.LoginHandler)                 // 登录
	v1.GET("/accesstoken", controller.AccessTokenHandler)      // 刷新accesstoken
	//v1.Use(middlewares.JWTAuthMiddleware())//中间件

	//RESTful风格的优势在于其简洁性和可扩展性。通过使用标准的HTTP方法，RESTful风格使得接口的设计更加清晰，易于理解和维护。
	{
		v1.POST("/task", controller.AddTaskHanler)          // 增加task
		v1.GET("/tasks", controller.GetTestByUseridHandler) // 获取tasks
		v1.DELETE("/tasks", controller.DeleteTaskHandler)   // 删除tasks
		v1.PUT("/tasks", controller.UptateTaskHandler)      // 更新tasks
	}
	{
		v1.POST("/test", controller.AddTestHanler)         // 增加test
		v1.GET("/test", controller.GetTestByUseridHandler) // 根据userid获取test对象
		v1.GET("/testList", controller.GetTestListHandler) // 获取test列表
		v1.PUT("/test", controller.UptateTestHandler)      // 更新test
		v1.DELETE("/test", controller.DeleteTestHandler)   // 批量删除test
	}
	return r
}
