package dto

import "time"

// PermissionRequest 权限请求
type PermissionRequest struct {
	Name      string `json:"name" binding:"required"` // 权限名称
	Code      string `json:"code" binding:"required"` // 权限编码
	Type      int    `json:"type" binding:"required"` // 权限类型：1-目录，2-菜单，3-API，4-按钮
	Path      string `json:"path"`                    // 权限路径
	Method    string `json:"method"`                  // HTTP方法
	ParentID  uint   `json:"parent_id"`               // 父权限ID
	Sort      int    `json:"sort"`                    // 排序
	Status    int    `json:"status"`                  // 状态
	Remark    string `json:"remark"`                  // 备注
	Icon      string `json:"icon"`                    // 图标（目录和菜单类型使用）
	Component string `json:"component"`               // 前端组件路径（目录类型为空，菜单类型使用）
	Redirect  string `json:"redirect"`                // 重定向路径（目录和菜单类型使用）
	Hidden    bool   `json:"hidden"`                  // 是否隐藏菜单
	KeepAlive bool   `json:"keep_alive"`              // 是否缓存页面
}

// PermissionInfo 权限信息
type PermissionInfo struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Code      string            `json:"code"`
	Type      int               `json:"type"`
	Path      string            `json:"path"`
	Method    string            `json:"method"`
	ParentID  uint              `json:"parent_id"`
	Sort      int               `json:"sort"`
	Status    int               `json:"status"`
	Remark    string            `json:"remark"`
	Icon      string            `json:"icon"`       // 图标
	Component string            `json:"component"`  // 前端组件路径
	Redirect  string            `json:"redirect"`   // 重定向路径
	Hidden    bool              `json:"hidden"`     // 是否隐藏
	KeepAlive bool              `json:"keep_alive"` // 是否缓存
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Children  []*PermissionInfo `json:"children,omitempty"` // 子权限列表（用于树形结构）
}

// PermissionListRequest 权限列表请求
type PermissionListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`              // 页码
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"` // 每页数量
	Keyword  string `form:"keyword"`                                     // 搜索关键词
	Type     *int   `form:"type"`                                        // 权限类型筛选
	Status   *int   `form:"status"`                                      // 状态筛选
}

// PermissionListResponse 权限列表响应
type PermissionListResponse struct {
	Total int64             `json:"total"` // 总数
	List  []*PermissionInfo `json:"list"`  // 列表
}

// PermissionTreeResponse 权限树响应
type PermissionTreeResponse struct {
	List []*PermissionInfo `json:"list"` // 树形列表
}

// AssignPermissionRequest 分配权限请求
type AssignPermissionRequest struct {
	RoleID        uint   `json:"role_id" binding:"required"`        // 角色ID
	PermissionIDs []uint `json:"permission_ids" binding:"required"` // 权限ID列表
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`  // 用户ID
	RoleIDs []uint `json:"role_ids" binding:"required"` // 角色ID列表
}

// UserPermissionResponse 用户权限响应
type UserPermissionResponse struct {
	UserID      uint              `json:"user_id"`
	Username    string            `json:"username"`
	Roles       []*RoleInfo       `json:"roles"`       // 用户角色列表
	Permissions []*PermissionInfo `json:"permissions"` // 用户权限列表（合并所有角色的权限）
}
