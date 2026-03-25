package model

import (
	"time"
)

// RolePermission 角色权限关联表（多对多关系）
type RolePermission struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	RoleID       uint `gorm:"type:int;not null;uniqueIndex:idx_role_permission" json:"role_id"`       // 角色ID
	PermissionID uint `gorm:"type:int;not null;uniqueIndex:idx_role_permission" json:"permission_id"` // 权限ID
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permission"
}
