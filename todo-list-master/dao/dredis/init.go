package dredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"todolist/setting"
)

var rdb *redis.Client

func InitClient(rcg *setting.RedisConfig) (err error) {

	rdb = redis.NewClient(&redis.Options{
		Addr:     rcg.Host + ":" + strconv.Itoa(rcg.Port),
		Password: rcg.Password, // no password set
		DB:       rcg.Db,       // use default DB
		PoolSize: 100,          // 连接池大小
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = rdb.Ping(ctx).Result()
	return
}
