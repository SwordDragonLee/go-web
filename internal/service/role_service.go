package service

import (
	"errors"

	"github.com/SwordDragonLee/go-web/internal/domain/dto"
	"github.com/SwordDragonLee/go-web/internal/domain/model"
	"github.com/SwordDragonLee/go-web/internal/repository"
)

type RoleService interface {
	CreateRole(roleRequest dto.RoleRequest) (*dto.RoleInfo, error)
	GetRoleList(req *dto.RoleListRequest) (*dto.RoleListResponse, error)
	AssignPermissions(roleID uint, permissionIDs []uint) error
	GetRolePermissions(roleID uint) ([]*dto.PermissionInfo, error)
}

type roleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}
func (s *roleService) CreateRole(req dto.RoleRequest) (*dto.RoleInfo, error) {
	// 检查角色名是否已存在
	if _, err := s.roleRepo.GetByName(req.Name); err == nil {
		return nil, errors.New("角色名已存在")
	}

	// 检查角色编码是否已存在
	if _, err := s.roleRepo.GetByCode(req.Code); err == nil {
		return nil, errors.New("角色编码已存在")
	}

	role := model.Role{
		Name:   req.Name,
		Code:   req.Code,
		Status: model.RoleStatus(req.Status),
	}
	createdRole, err := s.roleRepo.CreateRole(role)
	if err != nil {
		return nil, errors.New("创建角色失败")
	}
	return s.toRoleInfo(&createdRole), nil
}

func (s *roleService) GetRoleList(req *dto.RoleListRequest) (*dto.RoleListResponse, error) {
	// 设置默认值
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	roles, total, err := s.roleRepo.List(page, pageSize, req.Keyword)
	if err != nil {
		return nil, errors.New("获取角色列表失败")
	}
	roleInfos := make([]*dto.RoleInfo, 0, len(roles))
	for _, role := range roles {
		roleInfos = append(roleInfos, s.toRoleInfo(role))
	}
	return &dto.RoleListResponse{
		Total: total,
		List:  roleInfos,
	}, nil
}

// toRoleInfo 转换为RoleInfo DTO
func (s *roleService) toRoleInfo(role *model.Role) *dto.RoleInfo {
	return &dto.RoleInfo{
		ID:        role.ID,
		Name:      role.Name,
		Code:      role.Code,
		Status:    int(role.Status),
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

// AssignPermissions 为角色分配权限
func (s *roleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// 检查角色是否存在
	if _, err := s.roleRepo.GetByID(roleID); err != nil {
		return errors.New("角色不存在")
	}

	return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}

// GetRolePermissions 获取角色的权限列表
func (s *roleService) GetRolePermissions(roleID uint) ([]*dto.PermissionInfo, error) {
	// 检查角色是否存在
	if _, err := s.roleRepo.GetByID(roleID); err != nil {
		return nil, errors.New("角色不存在")
	}

	permissions, err := s.roleRepo.GetPermissions(roleID)
	if err != nil {
		return nil, errors.New("获取权限列表失败")
	}

	permissionInfos := make([]*dto.PermissionInfo, 0, len(permissions))
	for _, permission := range permissions {
		permissionInfos = append(permissionInfos, &dto.PermissionInfo{
			ID:        permission.ID,
			Name:      permission.Name,
			Code:      permission.Code,
			Type:      int(permission.Type),
			Path:      permission.Path,
			Method:    permission.Method,
			ParentID:  permission.ParentID,
			Sort:      permission.Sort,
			Status:    int(permission.Status),
			Remark:    permission.Remark,
			Icon:      permission.Icon,
			Component: permission.Component,
			Redirect:  permission.Redirect,
			Hidden:    permission.Hidden,
			KeepAlive: permission.KeepAlive,
			CreatedAt: permission.CreatedAt,
			UpdatedAt: permission.UpdatedAt,
		})
	}

	return permissionInfos, nil
}
