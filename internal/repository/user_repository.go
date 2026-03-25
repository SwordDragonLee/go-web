package repository

import (
	"time"

	"github.com/SwordDragonLee/go-web/internal/domain/model"
	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByPhone(phone string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	List(page, pageSize int, keyword string, status *int, role string) ([]*model.User, int64, error)
	UpdateLastLogin(userID uint, ip string) error
	AssignRoles(userID uint, roleIDs []uint) error
	GetRoles(userID uint) ([]*model.Role, error)
	RemoveRoles(userID uint) error
	GetUserPermissions(userID uint) ([]*model.Permission, error)
}

// db 是数据库连接实例
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号获取用户
func (r *userRepository) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

/**
 * 查询用户列表
 * 查询条件：
 * keyword: '%admin%'
 * status: 1
 * role: 'admin'
 * page: 1
 * pageSize: 10
 * 查询语句：
 * SELECT * FROM `user`
 * WHERE (username LIKE '%admin%'
 * OR email LIKE '%admin%'
 * OR phone LIKE '%admin%'
 * OR nickname LIKE '%admin%')
 * AND status = 1
 * AND role = 'admin'
 * ORDER BY created_at DESC
 * LIMIT 10
 * OFFSET 10
 */
func (r *userRepository) List(page, pageSize int, keyword string, status *int, role string) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.Model(&model.User{})
	// 关键词搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ? OR nickname LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 角色筛选
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateLastLogin 更新最后登录信息
func (r *userRepository) UpdateLastLogin(userID uint, ip string) error {
	now := time.Now()
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"last_login_at": &now,
		"last_login_ip": ip,
	}).Error
}

// AssignRoles 为用户分配角色
func (r *userRepository) AssignRoles(userID uint, roleIDs []uint) error {
	// 先删除该用户的所有角色
	if err := r.RemoveRoles(userID); err != nil {
		return err
	}

	// 批量插入新的角色关联
	if len(roleIDs) == 0 {
		return nil
	}

	var userRoles []model.UserRoleLink
	for _, roleID := range roleIDs {
		userRoles = append(userRoles, model.UserRoleLink{
			UserID: userID,
			RoleID: roleID,
		})
	}

	return r.db.Create(&userRoles).Error
}

// GetRoles 获取用户的角色列表
func (r *userRepository) GetRoles(userID uint) ([]*model.Role, error) {
	var roles []*model.Role

	err := r.db.
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}

// RemoveRoles 移除用户的所有角色
func (r *userRepository) RemoveRoles(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.UserRoleLink{}).Error
}

// GetUserPermissions 获取用户的所有权限（通过角色）
func (r *userRepository) GetUserPermissions(userID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission

	err := r.db.
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Distinct("permissions.*").
		Order("permissions.sort ASC, permissions.id ASC").
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}
