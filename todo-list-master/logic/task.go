package logic

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"todolist/dao/dmysql"
	"todolist/dao/dredis"
	"todolist/models"
)

func AddTask(uid uint, taskParam *models.ParamTask) (err error) {
	// 1. 写入mysql
	task := &models.Task{
		TaskContent: taskParam.TaskContent,
		State:       taskParam.State,
		Level:       taskParam.Level,
		UserID:      uid,
	}
	task, err = dmysql.AddTask(task)
	if err != nil {
		zap.L().Error("Faild to insert task in mysql", zap.Error(err))
		return
	}
	// 写入redis
	err = dredis.ZAddTask(task)
	if err != nil {
		zap.L().Error("Faild to insert task in redis", zap.Error(err))
		return
	}
	return
}

func GetTasks(uid uint, date *models.ParamDate) (data []map[string]string, err error) {
	// 1. 判断键存不存在
	flag, _ := dredis.CheckKeyTotalExist(uid)
	if err == nil && flag == false {
		// 不存在从mysql中读取用户的的所有todo task
		tids, _ := dmysql.QueryTasksIdByUid(uid)
		// 写入redis
		err = dredis.ZaddTotalTask(uid, tids)
		if err != nil {
			return
		}
	}
	// 2. 先从reids中获取需要返回那些task 数据
	tasks := dredis.ZrangTasks(uid, date)
	redisNotExist := make([]string, 0, 30)
	for _, task := range tasks {
		tid, _ := strconv.ParseInt(task, 10, 64)
		flag, _ = dredis.CheckTaskKeyExist(uid, tid)
		if err == nil && flag == false {
			redisNotExist = append(redisNotExist, task)
		}
	}
	fmt.Printf("redisNotExist:%v", redisNotExist)
	// 长度大于则需要redis中key不全
	if len(redisNotExist) > 0 {
		// mysql中查询到task
		taskList, err := dmysql.QueryTasksIdByTid(redisNotExist)
		if err != nil {
			zap.L().Warn("Error from mysql", zap.Error(err))
			return nil, err
		}
		// 写入到redis中
		err = dredis.HSetTasks(uid, taskList)
		if err != nil {
			zap.L().Warn("Error from redis", zap.Error(err))
			return nil, err
		}
	}
	// 从redis中获取要返回的值
	data, err = dredis.HGetAllTasks(uid, tasks)
	if err != nil {
		zap.L().Error("Faild to get task in redis", zap.Error(err))
		return
	}
	return
}

func DelTasks(uid uint, tids *models.ParamTaskIDs) (err error) {
	// 先删除mysql的数据,然后删除redis的数据
	// 1. 删除mysql的数据
	err = dmysql.DelTasks(tids.TaskIDs)
	if err != nil {
		zap.L().Error("Faild to delete tasks in mysql", zap.Error(err))
		return
	}
	// 2. 删除redis的数据
	// 删除KeyTaskAndDateMap中的数据
	err = dredis.DelKeyTaskAndDateMapTasks(uid, tids.TaskIDs)
	if err != nil {
		zap.L().Error("Faild to delete tasks in redis", zap.Error(err))
		return
	}
	//  删除KeyTaskMap中的数据
	err = dredis.DelTasks(uid, tids.TaskIDs)
	if err != nil {
		zap.L().Error("Faild to delete tasks in redis", zap.Error(err))
		return
	}
	return
}

func UpdateTask(uid uint, taskParam *models.ParamTask) (err error) {
	// 修改mysql,删除redis中的task,查询逻辑已经做了没有的时候添加所以不需要修改redis
	// 1. 修改task
	task := map[string]interface{}{
		"Level":       taskParam.Level,
		"State":       taskParam.State,
		"TaskContent": taskParam.TaskContent,
	}
	err = dmysql.UpdateTask(taskParam.Tid, task)
	if err != nil {
		zap.L().Error("Faild to update task in mysql", zap.Error(err))
		return
	}
	// 2. 删除redis中的key
	tids := []int64{taskParam.Tid}
	err = dredis.DelTasks(uid, tids)
	if err != nil {
		zap.L().Error("Faild to del task in redis", zap.Error(err))
		return
	}
	return
}
