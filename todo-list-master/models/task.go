package models

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	Tid         int64          `json:"tid"`
	TaskContent string         `json:"task_content" gorm:"type:text"`
	Level       int            `json:"level"` // 一般（0） 重要（1）
	State       int            `json:"state"` //  未完成（0） 已完成（1）
	UserID      uint           `json:"userID,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
