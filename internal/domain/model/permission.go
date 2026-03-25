package model

import (
	"time"
)

// PermissionType 权限类型
type PermissionType int

const (
	PermissionTypeDir  PermissionType = 1 // 目录（文件夹，用于分组，无实际页面）
	PermissionTypeMenu PermissionType = 2 // 菜单（有实际页面的菜单项）
	PermissionTypeAPI  PermissionType = 3 // API权限
	PermissionTypeBtn  PermissionType = 4 // 按钮权限
)

// PermissionStatus 权限状态
type PermissionStatus int

const (
	PermissionStatusDisabled PermissionStatus = 0 // 禁用
	PermissionStatusActive   PermissionStatus = 1 // 激活
)

// Permission 权限模型
type Permission struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 基本信息
	Name     string           `gorm:"type:varchar(50);not null" json:"name"`              // 权限名称
	Code     string           `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"` // 权限编码（唯一）
	Type     PermissionType   `gorm:"type:tinyint;not null" json:"type"`                  // 权限类型：1-目录，2-菜单，3-API，4-按钮
	Path     string           `gorm:"type:varchar(255)" json:"path"`                      // 权限路径（API路径或菜单路径）
	Method   string           `gorm:"type:varchar(10)" json:"method"`                     // HTTP方法（GET/POST/PUT/DELETE等）
	ParentID uint             `gorm:"default:0;index" json:"parent_id"`                   // 父权限ID（用于权限树）
	Sort     int              `gorm:"default:0" json:"sort"`                              // 排序
	Status   PermissionStatus `gorm:"type:tinyint;default:1;not null" json:"status"`      // 状态
	Remark   string           `gorm:"type:varchar(255)" json:"remark"`                    // 备注

	// 菜单/目录相关字段（Type=1目录或Type=2菜单时使用）
	Icon      string `gorm:"type:varchar(100)" json:"icon"`      // 图标（如：el-icon-user, ant-design:user-outlined）
	Component string `gorm:"type:varchar(255)" json:"component"` // 前端组件路径（目录类型为空，菜单类型为实际组件路径）
	Redirect  string `gorm:"type:varchar(255)" json:"redirect"`  // 重定向路径（目录类型通常指向第一个子菜单）
	Hidden    bool   `gorm:"default:0" json:"hidden"`            // 是否隐藏（true-隐藏，false-显示）
	KeepAlive bool   `gorm:"default:1" json:"keep_alive"`        // 是否缓存页面（目录类型不使用，菜单类型使用）
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permission"
}

// IsActive 检查权限是否激活
func (p *Permission) IsActive() bool {
	return p.Status == PermissionStatusActive
}
