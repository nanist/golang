package dredis

import (
	"context"
	"go.uber.org/zap"
	"time"
	"todolist/controller/code"
	"todolist/utils"
)

func SetCaptcha(email string, code string, captch utils.CAPTCHA) (err error) {
	//	1.设置k-v
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := KeyRedis(captch) + email
	err = rdb.Set(ctx, getRedisKey(key), code, 5*60*time.Second).Err()
	if err != nil {
		zap.L().Error("Failed to set Key redis", zap.Error(err))
		return err
	}
	return
}

func GetCaptcha(email string, str string) (val string, err error) {
	//	1.获取k-v
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := getRedisKey(str) + email
	val, err = rdb.Get(ctx, key).Result()
	if err != nil {
		zap.L().Error("Failed to get Key redis", zap.Error(err))
		err = code.EerrorCaptchaNotExist
		return
	}
	return
}
