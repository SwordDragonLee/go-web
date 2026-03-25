package repository

import (
	"github.com/SwordDragonLee/go-web/internal/domain/model"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	CreatePermission(permission model.Permission) (model.Permission, error)
	GetByID(id uint) (*model.Permission, error)
	GetByCode(code string) (*model.Permission, error)
	List(page, pageSize int, keyword string, permissionType *int, status *int) ([]*model.Permission, int64, error)
	GetTree(status *int) ([]*model.Permission, error)
	Update(id uint, permission model.Permission) error
	Delete(id uint) error
	GetByRoleID(roleID uint) ([]*model.Permission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

// CreatePermission 创建权限
func (r *permissionRepository) CreatePermission(permission model.Permission) (model.Permission, error) {
	err := r.db.Create(&permission).Error
	return permission, err
}

// GetByID 根据ID获取权限
func (r *permissionRepository) GetByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetByCode 根据权限编码获取权限
func (r *permissionRepository) GetByCode(code string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// List 获取权限列表（分页）
func (r *permissionRepository) List(page, pageSize int, keyword string, permissionType *int, status *int) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	query := r.db.Model(&model.Permission{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 类型筛选
	if permissionType != nil {
		query = query.Where("type = ?", *permissionType)
	}

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("sort ASC, id ASC").Offset(offset).Limit(pageSize).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// GetTree 获取权限树（按照 sort 排序）
func (r *permissionRepository) GetTree(status *int) ([]*model.Permission, error) {
	var permissions []*model.Permission

	query := r.db.Where("parent_id = ?", 0) // 只获取根节点
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Order("sort ASC, id ASC").Find(&permissions).Error; err != nil {
		return nil, err
	}

	// 递归加载子节点
	for _, permission := range permissions {
		r.loadChildren(permission, status)
	}

	return permissions, nil
}

// loadChildren 递归加载子权限
func (r *permissionRepository) loadChildren(permission *model.Permission, status *int) {
	var children []*model.Permission
	query := r.db.Where("parent_id = ?", permission.ID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Order("sort ASC, id ASC").Find(&children).Error; err != nil {
		return
	}

	// 递归加载子节点的子节点
	for _, child := range children {
		r.loadChildren(child, status)
	}

	// 这里不能直接设置，因为 model.Permission 没有 Children 字段
	// 需要在 DTO 层面处理树形结构
}

// GetByRoleID 根据角色ID获取权限列表
func (r *permissionRepository) GetByRoleID(roleID uint) ([]*model.Permission, error) {
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

// Update 更新权限
func (r *permissionRepository) Update(id uint, permission model.Permission) error {
	return r.db.Model(&model.Permission{}).Where("id = ?", id).Updates(&permission).Error
}

// Delete 删除权限（软删除）
func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Permission{}, id).Error
}
