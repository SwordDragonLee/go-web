package model

import (
	"time"
)

// UserRoleLink 用户角色关联表（多对多关系，支持一个用户多个角色）
type UserRoleLink struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID uint `gorm:"type:int;not null;uniqueIndex:idx_user_role" json:"user_id"` // 用户ID
	RoleID uint `gorm:"type:int;not null;uniqueIndex:idx_user_role" json:"role_id"` // 角色ID
}

// TableName 指定表名
func (UserRoleLink) TableName() string {
	return "user_role"
}
