package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"todolist/controller/code"
	"todolist/logic"
	"todolist/models"
)

// AddTaskHanler 添加单条task
func AddTaskHanler(c *gin.Context) {
	// 1.  判断是否已登录
	uid, err := checkLogin(c)
	if err != nil {
		zap.L().Error("Error form checkLogin", zap.Error(err))
		return
	}
	// 2. 校验参数
	p := new(models.ParamTask)
	if err = c.ShouldBindJSON(p); err != nil {
		zap.L().Warn("AddTask with invalid param:ParamTask", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeServerBusy)
		return
	}
	// 3. 调用添加task逻辑
	if err = logic.AddTask(uid, p); err != nil {
		zap.L().Error("Error AddTask about logic", zap.Error(err))
		ResponseError(c, code.CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
	return
}

// DeleteTaskHandler 删除tasks
func DeleteTaskHandler(c *gin.Context) {
	// 1.  判断是否已登录
	uid, err := checkLogin(c)
	if err != nil {
		zap.L().Error("Error form checkLogin", zap.Error(err))
		return
	}
	// 2. 校验参数
	p := new(models.ParamTaskIDs)
	if err = c.ShouldBindJSON(p); err != nil {
		zap.L().Error("DeleteTasks with invalid param:ParamTaskIDs", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeServerBusy)
		return
	}
	// 3. 调用删除tasks逻辑
	err = logic.DelTasks(uid, p)
	if err != nil {
		zap.L().Error("Error from DelTasks", zap.Error(err))
		ResponseError(c, code.CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
	return
}

// GetTasksHandler 查询tasks
func GetTasksHandler(c *gin.Context) {
	// 1.  判断是否已登录
	uid, err := checkLogin(c)
	if err != nil {
		zap.L().Error("Error form checkLogin", zap.Error(err))
		return
	}
	// 2. 校验参数
	p := new(models.ParamDate)
	if err = c.ShouldBindJSON(p); err != nil {
		zap.L().Error("GetTask with invalid param:ParamDate", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeServerBusy)
		return
	}
	// 3. 调用查询task逻辑
	tasks, err := logic.GetTasks(uid, p)
	if err != nil {
		zap.L().Error("Error GetTasks about logic", zap.Error(err))
		ResponseError(c, code.CodeServerBusy)
		return
	}
	data := map[string]interface{}{
		"user_id": uid,
		"date":    p,
		"tasks":   tasks,
	}
	ResponseSuccess(c, data)
	return
}

// PutTaskHanler 修改tasks
func UptateTaskHandler(c *gin.Context) {
	// 1.  判断是否已登录
	uid, err := checkLogin(c)
	if err != nil {
		zap.L().Error("Error form checkLogin", zap.Error(err))
		return
	}
	// 2. 校验参数
	p := new(models.ParamTask)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("GetTask with invalid param:ParamTask", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, code.CodeServerBusy)
		return
	}
	// 3. 调用更新task逻辑
	if err = logic.UpdateTask(uid, p); err != nil {
		zap.L().Error("Error AddTask about logic", zap.Error(err))
		ResponseError(c, code.CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
	return
}

// 检查是否登录
func checkLogin(c *gin.Context) (uid uint, err error) {
	uid, err = getCurrentUserID(c)
	if err != nil {
		zap.L().Warn("Failed  to get user id", zap.Error(err))
		if errors.Is(err, code.ErrorUserNotLogin) {
			ResponseError(c, code.CodeNeedLogin)
			return
		}
		ResponseError(c, code.CodeServerBusy)
		return
	}
	return
}
