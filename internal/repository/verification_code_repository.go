package repository

import (
	"time"

	"github.com/SwordDragonLee/go-web/internal/domain/model"
	"gorm.io/gorm"
)

// VerificationCodeRepository 验证码仓储接口
type VerificationCodeRepository interface {
	Create(code *model.VerificationCode) error
	GetLatestByEmailAndType(email string, codeType model.VerificationCodeType) (*model.VerificationCode, error)
	MarkAsUsed(id uint) error
	DeleteExpired() error
}

type verificationCodeRepository struct {
	db *gorm.DB
}

// NewVerificationCodeRepository 创建验证码仓储
func NewVerificationCodeRepository(db *gorm.DB) VerificationCodeRepository {
	return &verificationCodeRepository{db: db}
}

// Create 创建验证码
func (r *verificationCodeRepository) Create(code *model.VerificationCode) error {
	return r.db.Create(code).Error
}

// GetLatestByEmailAndType 获取最新的验证码
func (r *verificationCodeRepository) GetLatestByEmailAndType(email string, codeType model.VerificationCodeType) (*model.VerificationCode, error) {
	var code model.VerificationCode
	err := r.db.Where("email = ? AND type = ?", email, codeType).
		Order("created_at DESC").
		First(&code).Error
	if err != nil {
		return nil, err
	}
	return &code, nil
}

// MarkAsUsed 标记验证码为已使用
func (r *verificationCodeRepository) MarkAsUsed(id uint) error {
	now := time.Now()
	return r.db.Model(&model.VerificationCode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"used":    true,
			"used_at": &now,
		}).Error
}

// DeleteExpired 删除过期的验证码
func (r *verificationCodeRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&model.VerificationCode{}).Error
}
