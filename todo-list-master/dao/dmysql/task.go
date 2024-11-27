package dmysql

import (
	"fmt"
	"todolist/models"
	"todolist/utils"
)

// AddTask 增加task
func AddTask(task *models.Task) (*models.Task, error) {
	// 1. 生成id
	tid := utils.GenID()
	// 2. 创建
	task = &models.Task{Tid: tid, Level: task.Level, State: 0, UserID: task.UserID, TaskContent: task.TaskContent}
	err := db.Create(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

// QueryTasksByUidAndDate 通过的uid和日期查询
func QueryTasksByUidAndDate(uid uint, date *models.ParamDate) (tasks []models.Task, err error) {
	// 1. 解析日期
	startDate := date.StartDate
	endDate := date.EndDate
	// 2. 查询数据库
	err = db.Select("tid", "level", "state", "user_id", "task_content", "created_at").Where("user_id = ? and created_at >= ?  and  created_at < ?", uid, startDate, endDate).Find(&tasks).Error
	return
}

func DelTasks(taskIDs []int64) (err error) {
	fmt.Println("-------删除数据库-----")
	fmt.Println(taskIDs)
	err = db.Where("tid in ?", taskIDs).Delete(&models.Task{}).Error
	fmt.Println(err)
	fmt.Println("-------删除数据库结束-----")
	return
}

// QueryTasksIdByUid 根据用户id 查询tids
func QueryTasksIdByUid(uid uint) (tids []models.Task, err error) {
	err = db.Model(&models.Task{}).Debug().Select("tid", "created_at").Where("user_id = ?", uid).Find(&tids).Error
	return
}

// QueryTasksIdByTid 根据tid 返回详细task
func QueryTasksIdByTid(tids []string) (task []models.Task, err error) {
	err = db.Debug().Where("tid in ?", tids).Find(&task).Error
	return
}

func UpdateTask(tid int64, task map[string]interface{}) (err error) {
	err = db.Debug().Model(&models.Task{}).Where("tid = ?", tid).Updates(task).Error
	return
}
