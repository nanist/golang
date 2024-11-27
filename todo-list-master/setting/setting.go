package setting

import (
	"fmt"
	"github.com/spf13/viper"
	"path"
	"runtime"
)

// 自定义logerr，其中zapcore.Core需要三个配置——Encoder，WriteSyncer，LogLevel
// 使用Lumberjack进行日志分割

type AppConfig struct {
	LogConfig   *LogConfig   `json:"logConfig" mapstructure:"log"`
	MysqlConfig *MysqlConfig `json:"mysqlConfig" mapstructure:"mysql"`
	EmailConfig *EmailConfig `json:"emailConfig" mapstructure:"email"`
	RedisConfig *RedisConfig `json:"redisConfig" mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `json:"level" mapstructure:"level"`
	Filename   string `json:"filename" mapstructure:"filename"`
	MaxSize    int    `json:"maxsize" mapstructure:"maxsize"`
	MaxAge     int    `json:"max_age" mapstructure:"max_age"`
	MaxBackups int    `json:"max_backups" mapstructure:"max_backups"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Db           string `mapstructure:"db"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	Port     int    `mapstructure:"port"`
	PoolSize int    `mapstructure:"pool_size"`
}
type EmailConfig struct {
	Addr     string `mapstructure:"addr"`
	Identity string `mapstructure:"identity"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
}

var Appconfig = new(AppConfig)

func Init() (err error) {
	// 获取当前文件的路径
	_, filename, _, _ := runtime.Caller(0)
	rootPath := path.Dir(path.Dir(filename))
	viper.SetConfigFile(rootPath + "/conf/Appconfig.yaml") // 指定配置文件路径

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
			fmt.Printf("未找到配置文件,err:%v", err)
		} else {
			// 配置文件被找到，但产生了另外的错误
			fmt.Printf("未知错误,err:%v", err)
		}
	}
	if err = viper.Unmarshal(Appconfig); err != nil {
		fmt.Printf("viper.Unmarshal failed,err:%v", err)
	}
	return err
}
