package dmysql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
	"todolist/controller/code"
	"todolist/models"
)

// MD5加密的密钥
const secret = "todolist.com"

// CheckUserExist 邮箱是否已经被注册
func CheckUserExist(email string) (err error) {
	user := &models.User{}
	var count int64
	db.Model(user).Where("email = ?", email).Count(&count)
	if count > 0 {
		err = code.ErrorUserExist
	}
	return
}

// InsertUser 新增用户
func InsertUser(user *models.User) (err error) {
	// 1. 对密码进行加密
	password := encryptPassword(user.Password)
	// 2. 创建用户
	user = &models.User{Email: user.Email, Name: user.Name, Phone: user.Phone, Password: password}
	err = db.Create(user).Error
	if err != nil {
		return err
	}
	return
}

// encryptPassword 密码加密
func encryptPassword(originalPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(originalPassword)))
}

func Login(user *models.User) (userMap map[string]interface{}, err error) {
	userObject := new(models.User)
	err = db.Select("id", "email", "password", "name").Where("email = ?", user.Email).First(userObject).Error
	// 判断用户是否存在
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = code.ErrorUserNotExist
			return nil, err
		}
	}
	// 检查密码是否正确

	if encryptPassword(user.Password) != userObject.Password {
		err = code.ErrorInvalidPassword
		return nil, err
	}
	userMap = map[string]interface{}{
		"id":    userObject.ID,
		"email": userObject.Email,
		"name":  userObject.Name,
	}
	return
}
