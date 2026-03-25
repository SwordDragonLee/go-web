package repository

import (
	"github.com/SwordDragonLee/go-web/internal/domain/model"
	"gorm.io/gorm"
)

// LoginLogRepository 登录日志仓储接口
type LoginLogRepository interface {
	Create(log *model.LoginLog) error
	List(userID uint, page, pageSize int) ([]*model.LoginLog, int64, error)
}

type loginLogRepository struct {
	db *gorm.DB
}

// NewLoginLogRepository 创建登录日志仓储
func NewLoginLogRepository(db *gorm.DB) LoginLogRepository {
	return &loginLogRepository{db: db}
}

// Create 创建登录日志
func (r *loginLogRepository) Create(log *model.LoginLog) error {
	return r.db.Create(log).Error
}

// List 获取登录日志列表
func (r *loginLogRepository) List(userID uint, page, pageSize int) ([]*model.LoginLog, int64, error) {
	var logs []*model.LoginLog
	var total int64

	query := r.db.Model(&model.LoginLog{}).Where("user_id = ?", userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

