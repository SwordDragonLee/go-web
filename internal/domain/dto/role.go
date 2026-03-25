package dto

import "time"

type RoleRequest struct {
	Name   string `json:"name" validate:"required"`
	Code   string `json:"code" validate:"required"`
	Status int    `json:"status" validate:"required"`
}

type RoleInfo struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`              // 页码
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"` // 每页数量
	Keyword  string `form:"keyword"`                                     // 搜索关键词
}

// RoleListResponse 角色列表响应
type RoleListResponse struct {
	Total int64       `json:"total"` // 总数
	List  []*RoleInfo `json:"list"`  // 列表
}
