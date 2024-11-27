package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	//Uid      string `json:"uid" gorm:"primary_key"`  // 这里可以根据需要增加雪花id
	Name     string `json:"name"`
	Phone    string `json:"phone" gorm:"index"`
	Email    string `json:"email" gorm:"index,unique"`
	Password string `json:"password"`
	Task     []Task `json:"task"`
}
