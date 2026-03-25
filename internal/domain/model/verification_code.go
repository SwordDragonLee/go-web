package model

import (
	"time"

	"gorm.io/gorm"
)

// VerificationCodeType 验证码类型
type VerificationCodeType string

const (
	CodeTypeEmailVerification VerificationCodeType = "email_verification" // 邮箱验证
	CodeTypePasswordReset     VerificationCodeType = "password_reset"     // 密码重置
)

// VerificationCode 验证码模型
type VerificationCode struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email     string               `gorm:"type:varchar(100);not null;index" json:"email"`  // 邮箱
	Code      string               `gorm:"type:varchar(10);not null" json:"code"`          // 验证码
	Type      VerificationCodeType `gorm:"type:varchar(50);not null;index" json:"type"`    // 类型
	ExpiresAt time.Time            `gorm:"type:datetime;not null;index" json:"expires_at"` // 过期时间
	Used      bool                 `gorm:"type:tinyint;default:0;not null" json:"used"`    // 是否已使用
	UsedAt    *time.Time           `gorm:"type:datetime" json:"used_at"`                   // 使用时间
}

// TableName 指定表名
func (VerificationCode) TableName() string {
	return "verification_codes"
}

// IsExpired 检查验证码是否过期
func (v *VerificationCode) IsExpired() bool {
	return time.Now().After(v.ExpiresAt)
}

// IsValid 检查验证码是否有效（未使用且未过期）
func (v *VerificationCode) IsValid() bool {
	return !v.Used && !v.IsExpired()
}
