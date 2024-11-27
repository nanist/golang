package controller

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"todolist/controller/code"
	"todolist/dao/dmysql"
	"todolist/logic"
	"todolist/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// GetTestByUseridHandler 查询
func GetTestByUseridHandler(c *gin.Context) {
	ids := c.Query("id")

	fmt.Println("=======" + ids)

	intNum, _ := strconv.Atoi(ids) //字符串转int
	var uid uint = uint(intNum)    // int转uint显式转换

	tasks, err := dmysql.QueryTasksIdByUid(uid)
	if err != nil {
		zap.L().Error("Error GetTasks about logic", zap.Error(err))
		ResponseError(c, code.CodeServerBusy)
		return
	}

	tasks[0].TaskContent = "测试内容" //修改返回值

	b, _ := json.Marshal(tasks) //结构体转json
	s3 := string(b)
	fmt.Println("tasks=======" + s3) //json转字符串

	data := map[string]interface{}{
		"user_id": uid,
		"tasks":   tasks,
	}
	jsonBytes, _ := json.Marshal(data)
	m := string(jsonBytes)
	fmt.Println("data=======" + m) //map转字符串
	ResponseSuccess(c, data)
	return
}

// GetTestListHandler 查询列表
func GetTestListHandler(c *gin.Context) {
	ids := c.Query("id")
	fmt.Println("=======" + ids)
	arr := strings.Split(ids, ",") //字符串转数组

	tasks, err := dmysql.QueryTasksIdByTid(arr)
	if err != nil {
		zap.L().Error("Error GetTasks about logic", zap.Error(err))
		ResponseError(c, code.CodeServerBusy)
		return
	}

	ResponseSuccess(c, tasks)
	return
}

// AddTestHanler 新增
func AddTestHanler(c *gin.Context) {

	var uid uint = 1                            //创建人
	p := new(models.ParamTask)                  //new一个空的结构体
	if err := c.ShouldBindJSON(p); err != nil { //使用 ShouldBindJSON() 或 BindJSON() 解析 JSON 数据并绑定到结构体
		zap.L().Error("GetTask with invalid param:ParamDate", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeServerBusy)
		return
	}
	logic.AddTask(uid, p)

	ResponseSuccess(c, nil)
	return
}

// UptateTestHandler 更新
func UptateTestHandler(c *gin.Context) {

	var uid uint = 1 //修改人
	p := new(models.ParamTask)

	if err := c.ShouldBindJSON(p); err != nil { //使用 ShouldBindJSON() 或 BindJSON() 解析 JSON 数据并绑定到结构体
		zap.L().Error("GetTask with invalid param:ParamTask", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeServerBusy)
		return
	}

	logic.UpdateTask(uid, p)

	ResponseSuccess(c, nil)
	return
}

// DeleteTestHandler 批量删除test
func DeleteTestHandler(c *gin.Context) {
	var uid uint = 1 //修改人

	// 2. 校验参数
	p := new(models.ParamTaskIDs)

	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("DeleteTasks with invalid param:ParamTaskIDs", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeServerBusy)
		return
	}
	fmt.Println(p.TaskIDs) //获取id
	logic.DelTasks(uid, p)
	fmt.Println("-----------")
	ResponseSuccess(c, nil)
	return
}
