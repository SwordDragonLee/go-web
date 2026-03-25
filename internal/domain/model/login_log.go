package model

import (
	"time"

	"gorm.io/gorm"
)

// LoginLog 登录日志
type LoginLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID    uint   `gorm:"type:int;not null;index" json:"user_id"`    // 用户ID
	Username  string `gorm:"type:varchar(50);not null" json:"username"`  // 用户名
	IP        string `gorm:"type:varchar(50)" json:"ip"`                 // IP地址
	UserAgent string `gorm:"type:varchar(255)" json:"user_agent"`       // User Agent
	Status    string `gorm:"type:varchar(20);not null" json:"status"`    // 登录状态：success/failed
	Message   string `gorm:"type:varchar(255)" json:"message"`           // 消息
}

// TableName 指定表名
func (LoginLog) TableName() string {
	return "login_logs"
}

