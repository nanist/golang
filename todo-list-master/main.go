package main

import (
	"fmt"
	"go.uber.org/zap"
	"todolist/controller"
	"todolist/dao/dmysql"
	"todolist/dao/dredis"
	"todolist/logger"
	"todolist/router"
	"todolist/setting"
	"todolist/utils"
)

func main() {
	err := setting.Init()
	if err != nil {
		fmt.Printf("err:%v", err)
	}
	// 初始化日志
	logger.InitLogger(setting.Appconfig.LogConfig, "dev")
	defer zap.L().Sync()

	// 初始化Snowflake
	err = utils.InitSnowflake("2023-01-10")
	if err != nil {
		zap.L().Error("Error from Snowflake", zap.Error(err))
		return
	}
	// 初始化mysql
	err = dmysql.InitClient(setting.Appconfig.MysqlConfig)
	if err != nil {
		zap.L().Error("Error from init mysql", zap.Error(err))
		return
	}
	// 初始化redis
	err = dredis.InitClient(setting.Appconfig.RedisConfig)
	if err != nil {
		zap.L().Error("Error from init redis", zap.Error(err))
		return
	}
	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v\n", err)
		return
	}
	r := router.SetupRouter()
	r.Run("127.0.0.1:8080")
}
