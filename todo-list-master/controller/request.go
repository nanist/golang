package controller

import (
	"github.com/gin-gonic/gin"
	"todolist/controller/code"
)

const CtxUserIDKey = "userID"
const CtxUserEmail = "userEmail"

func getCurrentUserID(c *gin.Context) (userID uint, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = code.ErrorUserNotLogin
		return
	}
	userID, ok = uid.(uint)
	if !ok {
		err = code.ErrorUserNotLogin
		return
	}
	return

}
