package dmysql

import (
	"fmt"
	"time"
	"todolist/models"
	"todolist/setting"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitClient(mcg *setting.MysqlConfig) (err error) {
	//var dsn = "root:root@tcp(127.0.0.1:3306)/todolist?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mcg.User, mcg.Password, mcg.Host, mcg.Port, mcg.Db)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 设置日志模式为Info，将会打印出所有SQL语句
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	sqlDB, _ := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	//automigrate是gorm库中一个非常重要的功能,它可以自动创建数据库表和对应的字段,无需手动编写SQL语句。
	err = db.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		return err
	}
	return
}
