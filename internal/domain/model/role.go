package model

import (
	"time"
)

// RoleStatus 角色状态
type RoleStatus int

const (
	RoleStatusDisabled RoleStatus = 0 // 禁用
	RoleStatusActive   RoleStatus = 1 // 激活
)

// Role 角色模型
type Role struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 基本信息
	Name   string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"` // 角色名
	Code   string     `gorm:"type:varchar(100);uniqueIndex" json:"code"`         // 角色编码
	Status RoleStatus `gorm:"type:tinyint;default:0;not null" json:"status"`     // 状态
}

// TableName 指定表名
func (Role) TableName() string {
	return "role"
}
