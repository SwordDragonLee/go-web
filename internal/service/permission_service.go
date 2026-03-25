package service

import (
	"errors"

	"github.com/SwordDragonLee/go-web/internal/domain/dto"
	"github.com/SwordDragonLee/go-web/internal/domain/model"
	"github.com/SwordDragonLee/go-web/internal/repository"
)

type PermissionService interface {
	CreatePermission(req dto.PermissionRequest) (*dto.PermissionInfo, error)
	GetPermissionList(req *dto.PermissionListRequest) (*dto.PermissionListResponse, error)
	GetPermissionTree(status *int) (*dto.PermissionTreeResponse, error)
	UpdatePermission(id uint, req dto.PermissionRequest) (*dto.PermissionInfo, error)
	DeletePermission(id uint) error
	GetPermissionByID(id uint) (*dto.PermissionInfo, error)
}

type permissionService struct {
	permissionRepo repository.PermissionRepository
}

func NewPermissionService(permissionRepo repository.PermissionRepository) PermissionService {
	return &permissionService{permissionRepo: permissionRepo}
}

// CreatePermission 创建权限
func (s *permissionService) CreatePermission(req dto.PermissionRequest) (*dto.PermissionInfo, error) {
	// 检查权限编码是否已存在
	if _, err := s.permissionRepo.GetByCode(req.Code); err == nil {
		return nil, errors.New("权限编码已存在")
	}

	// 如果有父节点，检查父节点是否存在
	if req.ParentID > 0 {
		if _, err := s.permissionRepo.GetByID(req.ParentID); err != nil {
			return nil, errors.New("父权限不存在")
		}
	}

	permission := model.Permission{
		Name:      req.Name,
		Code:      req.Code,
		Type:      model.PermissionType(req.Type),
		Path:      req.Path,
		Method:    req.Method,
		ParentID:  req.ParentID,
		Sort:      req.Sort,
		Status:    model.PermissionStatus(req.Status),
		Remark:    req.Remark,
		Icon:      req.Icon,
		Component: req.Component,
		Redirect:  req.Redirect,
		Hidden:    req.Hidden,
		KeepAlive: req.KeepAlive,
	}

	createdPermission, err := s.permissionRepo.CreatePermission(permission)
	if err != nil {
		return nil, errors.New("创建权限失败")
	}

	return s.toPermissionInfo(&createdPermission, nil), nil
}

// GetPermissionList 获取权限列表
func (s *permissionService) GetPermissionList(req *dto.PermissionListRequest) (*dto.PermissionListResponse, error) {
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

	permissions, total, err := s.permissionRepo.List(page, pageSize, req.Keyword, req.Type, req.Status)
	if err != nil {
		return nil, errors.New("获取权限列表失败")
	}

	permissionInfos := make([]*dto.PermissionInfo, 0, len(permissions))
	for _, permission := range permissions {
		permissionInfos = append(permissionInfos, s.toPermissionInfo(permission, nil))
	}

	return &dto.PermissionListResponse{
		Total: total,
		List:  permissionInfos,
	}, nil
}

// GetPermissionTree 获取权限树
func (s *permissionService) GetPermissionTree(status *int) (*dto.PermissionTreeResponse, error) {
	permissions, err := s.permissionRepo.GetTree(status)
	if err != nil {
		return nil, errors.New("获取权限树失败")
	}

	// 构建树形结构
	tree := s.buildTree(permissions, status)

	return &dto.PermissionTreeResponse{
		List: tree,
	}, nil
}

// buildTree 递归构建权限树
func (s *permissionService) buildTree(permissions []*model.Permission, status *int) []*dto.PermissionInfo {
	tree := make([]*dto.PermissionInfo, 0)

	for _, permission := range permissions {
		info := s.toPermissionInfo(permission, nil)
		// 递归加载子节点
		children, _, _ := s.permissionRepo.List(1, 1000, "", nil, status)
		var childPermissions []*model.Permission
		for _, child := range children {
			if child.ParentID == permission.ID {
				childPermissions = append(childPermissions, child)
			}
		}
		if len(childPermissions) > 0 {
			info.Children = s.buildTree(childPermissions, status)
		}
		tree = append(tree, info)
	}

	return tree
}

// UpdatePermission 更新权限
func (s *permissionService) UpdatePermission(id uint, req dto.PermissionRequest) (*dto.PermissionInfo, error) {
	// 检查权限是否存在
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("权限不存在")
	}

	// 检查权限编码是否被其他权限占用
	if existingPermission, err := s.permissionRepo.GetByCode(req.Code); err == nil && existingPermission.ID != id {
		return nil, errors.New("权限编码已被其他权限使用")
	}

	// 如果有父节点，检查父节点是否存在且不是自己
	if req.ParentID > 0 {
		if req.ParentID == id {
			return nil, errors.New("不能将自己设置为父节点")
		}
		if _, err := s.permissionRepo.GetByID(req.ParentID); err != nil {
			return nil, errors.New("父权限不存在")
		}
	}

	// 更新权限信息
	permission.Name = req.Name
	permission.Code = req.Code
	permission.Type = model.PermissionType(req.Type)
	permission.Path = req.Path
	permission.Method = req.Method
	permission.ParentID = req.ParentID
	permission.Sort = req.Sort
	permission.Status = model.PermissionStatus(req.Status)
	permission.Remark = req.Remark
	permission.Icon = req.Icon
	permission.Component = req.Component
	permission.Redirect = req.Redirect
	permission.Hidden = req.Hidden
	permission.KeepAlive = req.KeepAlive

	if err := s.permissionRepo.Update(id, *permission); err != nil {
		return nil, errors.New("更新权限失败")
	}

	return s.toPermissionInfo(permission, nil), nil
}

// DeletePermission 删除权限
func (s *permissionService) DeletePermission(id uint) error {
	// 检查权限是否存在
	_, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return errors.New("权限不存在")
	}

	// TODO: 检查是否有子权限，如果有则不能删除
	// TODO: 检查权限是否被角色使用，如果被使用则不能删除

	if err := s.permissionRepo.Delete(id); err != nil {
		return errors.New("删除权限失败")
	}

	return nil
}

// GetPermissionByID 根据ID获取权限
func (s *permissionService) GetPermissionByID(id uint) (*dto.PermissionInfo, error) {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("权限不存在")
	}
	return s.toPermissionInfo(permission, nil), nil
}

// toPermissionInfo 转换为 PermissionInfo DTO
func (s *permissionService) toPermissionInfo(permission *model.Permission, children []*dto.PermissionInfo) *dto.PermissionInfo {
	info := &dto.PermissionInfo{
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
	}

	if children != nil {
		info.Children = children
	}

	return info
}
