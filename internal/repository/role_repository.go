package repository

import (
	"github.com/SwordDragonLee/go-web/internal/domain/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	CreateRole(role model.Role) (model.Role, error)
	GetByName(name string) (*model.Role, error)
	GetByCode(code string) (*model.Role, error)
	GetByID(id uint) (*model.Role, error)
	List(page, pageSize int, keyword string) ([]*model.Role, int64, error)
	AssignPermissions(roleID uint, permissionIDs []uint) error
	GetPermissions(roleID uint) ([]*model.Permission, error)
	RemovePermissions(roleID uint) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// CreateRole 创建角色
func (r *roleRepository) CreateRole(role model.Role) (model.Role, error) {
	err := r.db.Create(&role).Error
	return role, err
}

// GetByName 根据角色名获取角色
func (r *roleRepository) GetByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("name = ?", name).First(&role).Error // First(&role) 会执行查询并将结果填充到 role 中
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByCode 根据角色编码获取角色
func (r *roleRepository) GetByCode(code string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) List(page, pageSize int, keyword string) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	query := r.db.Model(&model.Role{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// GetByID 根据ID获取角色
func (r *roleRepository) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// AssignPermissions 为角色分配权限
func (r *roleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// 先删除该角色的所有权限
	if err := r.RemovePermissions(roleID); err != nil {
		return err
	}

	// 批量插入新的权限关联
	if len(permissionIDs) == 0 {
		return nil
	}

	var rolePermissions []model.RolePermission
	for _, permissionID := range permissionIDs {
		rolePermissions = append(rolePermissions, model.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		})
	}

	return r.db.Create(&rolePermissions).Error
}

// GetPermissions 获取角色的权限列表
func (r *roleRepository) GetPermissions(roleID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission

	err := r.db.
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.sort ASC, permissions.id ASC").
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// RemovePermissions 移除角色的所有权限
func (r *roleRepository) RemovePermissions(roleID uint) error {
	return r.db.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error
}
