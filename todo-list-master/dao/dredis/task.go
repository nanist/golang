package dredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"todolist/models"
)

const (
	KeyTaskAndDateMapExpireDuration = time.Second * 60 * 60
	KeyTaskMapExpireDuration        = time.Second * 60 * 60
)

// ZAddTasks 将task写入redis uid: zset(时间戳)   task_id:{task具体的内容}
func ZAddTask(task *models.Task) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// 创建一个pipline
	pipe := rdb.Pipeline()
	suid := strconv.FormatUint(uint64(task.UserID), 10)
	key := getRedisKey(KeyTaskAndDateMap + "uid_" + suid)
	stid := strconv.FormatInt(task.Tid, 10)
	timestamp := time.Now().Unix()
	// 写入uid_日期:{ task_id:task_id}
	dataSet := &redis.Z{
		float64(timestamp),
		stid,
	}
	// 判断键存不存在，不存在的话，则不写入redis中，等查询的时候自己写入
	flag, err := CheckKeyTotalExist(task.UserID)
	if err == nil && flag == true {
		pipe.ZAdd(ctx, key, dataSet)
		pipe.Expire(ctx, key, KeyTaskAndDateMapExpireDuration)
	}

	// 写入task_id:{task具体的内容}
	key = getRedisKey(KeyTaskMap + "uid_" + suid + ":" + stid)
	dataMap := map[string]string{
		"uid":     suid,
		"tid":     stid,
		"level":   strconv.Itoa(task.Level),
		"state":   strconv.Itoa(task.State),
		"content": task.TaskContent,
	}
	pipe.HSet(ctx, key, dataMap)
	pipe.Expire(ctx, key, KeyTaskMapExpireDuration)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return
	}
	return
}

// ZaddTotalTask 将用户所有的task写入到redis中
func ZaddTotalTask(uid uint, tids []models.Task) (err error) {
	suid := strconv.FormatUint(uint64(uid), 10)
	key := getRedisKey(KeyTaskAndDateMap + "uid_" + suid)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// 创建一个pipline
	pipe := rdb.Pipeline()
	for _, task := range tids {
		stid := strconv.FormatInt(task.Tid, 10)
		dataSet := &redis.Z{
			float64(task.CreatedAt.Unix()),
			stid,
		}
		pipe.ZAdd(ctx, key, dataSet)
	}
	pipe.Expire(ctx, key, KeyTaskAndDateMapExpireDuration)
	_, err = pipe.Exec(ctx)
	return
}

// ZAddTasks 添加task到redis
func ZAddTasks(tasks []models.Task) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// 创建一个pipline
	pipe := rdb.Pipeline()
	for _, task := range tasks {
		suid := strconv.FormatUint(uint64(task.UserID), 10)
		key := getRedisKey(KeyTaskAndDateMap + "uid_" + suid)
		stid := strconv.FormatInt(task.Tid, 10)
		timestamp := task.CreatedAt.Unix()
		// 写入uid_日期:{ task_id:task_id}
		dataSet := &redis.Z{
			float64(timestamp),
			stid,
		}
		pipe.ZAdd(ctx, key, dataSet)
		key = getRedisKey(KeyTaskMap + "uid_" + suid + ":" + stid)
		dataMap := map[string]string{
			"uid":     suid,
			"tid":     stid,
			"level":   strconv.Itoa(task.Level),
			"state":   strconv.Itoa(task.State),
			"content": task.TaskContent,
		}
		pipe.HSet(ctx, key, dataMap)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return
	}
	return
}

// ZrangTasks 根据score（时间戳）返回数据  查询到需要返回那些todo task
func ZrangTasks(uid uint, date *models.ParamDate) (vals []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	suid := strconv.FormatUint(uint64(uid), 10)
	key := getRedisKey(KeyTaskAndDateMap + "uid_" + suid)
	// 1. 解析日期数据
	loc, _ := time.LoadLocation("Local")
	startDate, _ := time.ParseInLocation("2006-1-02", date.StartDate, loc)
	endDate, _ := time.ParseInLocation("2006-1-02", date.EndDate, loc)
	// 查询对应的改日期范围内的taskId
	opt := &redis.ZRangeBy{
		Min: strconv.FormatUint(uint64(startDate.Unix()), 10),
		Max: strconv.FormatUint(uint64(endDate.Unix()), 10),
	}
	vals = rdb.ZRangeByScore(ctx, key, opt).Val()
	return vals

}

func DelTasks(uid uint, taskIDs []int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	pipe := rdb.Pipeline()
	suid := strconv.FormatUint(uint64(uid), 10)
	stids := make([]string, 0, 30)
	for _, tid := range taskIDs {
		stid := strconv.FormatUint(uint64(tid), 10)
		key := getRedisKey(KeyTaskMap + "uid_" + suid + ":" + stid)
		stids = append(stids, key)
	}
	// 删除KeyTaskMap中的数据
	pipe.Del(ctx, stids...)
	_, err = pipe.Exec(ctx)
	return
}

// 删除KeyTaskAndDateMap中的部分key
func DelKeyTaskAndDateMapTasks(uid uint, taskIDs []int64) (err error) {
	// 1. 删除KeyTaskAndDateMap中的数据
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	suid := strconv.FormatUint(uint64(uid), 10)
	key := getRedisKey(KeyTaskAndDateMap + "uid_" + suid)
	itaskIDs := make([]interface{}, 0, 30)
	for _, tid := range taskIDs {
		itaskIDs = append(itaskIDs, tid)
	}
	pipe := rdb.Pipeline()
	pipe.ZRem(ctx, key, itaskIDs...)
	_, err = pipe.Exec(ctx)
	return
}

func HSetTasks(uid uint, tasks []models.Task) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// 创建一个pipline
	pipe := rdb.Pipeline()
	suid := strconv.FormatUint(uint64(uid), 10)
	for _, task := range tasks {
		stid := strconv.FormatInt(task.Tid, 10)
		key := getRedisKey(KeyTaskMap + "uid_" + suid + ":" + stid)
		dataMap := map[string]string{
			"uid":     suid,
			"tid":     stid,
			"level":   strconv.Itoa(task.Level),
			"state":   strconv.Itoa(task.State),
			"content": task.TaskContent,
		}
		pipe.HSet(ctx, key, dataMap)
		pipe.Expire(ctx, key, KeyTaskMapExpireDuration)
	}
	_, err = pipe.Exec(ctx)
	return
}

func HGetAllTasks(uid uint, tids []string) (tasks []map[string]string, err error) {
	suid := strconv.FormatUint(uint64(uid), 10)
	// 查询到所有的数据
	taskIds := make([]string, 0, 50)
	for _, item := range tids {
		key := getRedisKey(KeyTaskMap + "uid_" + suid + ":" + item)
		taskIds = append(taskIds, key)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// 创建一个pipline
	pipe := rdb.Pipeline()
	for _, item := range taskIds {
		pipe.HGetAll(ctx, item)
	}
	cmdr, err := pipe.Exec(ctx)
	if err != nil {
		return
	}
	tasks = make([]map[string]string, 0, 30)
	// 获取执行结果
	for _, item := range cmdr {
		tasks = append(tasks, item.(*redis.StringStringMapCmd).Val())
	}
	return
}

// CheckKeyTotalExist 查询用户下的total key存在不
func CheckKeyTotalExist(uid uint) (flag bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	suid := strconv.FormatUint(uint64(uid), 10)
	key := getRedisKey(KeyTaskAndDateMap + "uid_" + suid)
	res, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if res == 1 {
		return true, err
	}
	return false, err
}

func CheckTaskKeyExist(uid uint, tid int64) (flag bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	suid := strconv.FormatUint(uint64(uid), 10)
	stid := strconv.FormatInt(tid, 10)
	key := getRedisKey(KeyTaskMap + "uid_" + suid + ":" + stid)
	res, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if res == 1 {
		return true, err
	}
	return false, err
}
