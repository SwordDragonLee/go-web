package model

import (
	"time"

	"gorm.io/gorm"
)

// UserStatus 用户状态
type UserStatus int

const (
	UserStatusActive   UserStatus = 1 // 激活
	UserStatusInactive UserStatus = 0 // 未激活
	UserStatusDisabled UserStatus = 2 // 禁用
)

// UserRole 用户角色
type UserRole string

const (
	RoleAdmin UserRole = "admin" // 管理员
	RoleUser  UserRole = "user"  // 普通用户
	RoleGuest UserRole = "guest" // 访客
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 基本信息
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"` // 用户名
	Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`            // 邮箱
	Phone    string `gorm:"type:varchar(20);uniqueIndex" json:"phone"`             // 手机号
	Password string `gorm:"type:varchar(255);not null" json:"-"`                   // 密码（加密后）

	// 用户信息
	Nickname string     `gorm:"type:varchar(50)" json:"nickname"`              // 昵称
	Avatar   string     `gorm:"type:varchar(255)" json:"avatar"`               // 头像URL
	Status   UserStatus `gorm:"type:tinyint;default:0;not null" json:"status"` // 状态
	Role     UserRole   `gorm:"type:varchar(20);default:'user'" json:"role"`   // 角色

	// 扩展信息
	LastLoginAt *time.Time `gorm:"type:datetime" json:"last_login_at"`    // 最后登录时间
	LastLoginIP string     `gorm:"type:varchar(50)" json:"last_login_ip"` // 最后登录IP
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsDisabled 检查用户是否被禁用
func (u *User) IsDisabled() bool {
	return u.Status == UserStatusDisabled
}
