package handler

import (
	"fmt"
	"net/http"

	"github.com/SwordDragonLee/go-web/internal/domain/constant"
	"github.com/SwordDragonLee/go-web/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"github.com/SwordDragonLee/go-web/internal/service"
)

type PermissionHandler struct {
	permissionService service.PermissionService
}

func NewPermissionHandler(permissionService service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService}
}

// CreatePermission 创建权限
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req dto.PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	permissionInfo, err := h.permissionService.CreatePermission(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(permissionInfo))
}

// GetPermissionList 获取权限列表
func (h *PermissionHandler) GetPermissionList(c *gin.Context) {
	var req dto.PermissionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	resp, err := h.permissionService.GetPermissionList(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(resp))
}

// GetPermissionTree 获取权限树
func (h *PermissionHandler) GetPermissionTree(c *gin.Context) {
	// 从查询参数获取状态筛选
	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		s := 0
		if statusStr == "1" {
			s = 1
		}
		status = &s
	}

	resp, err := h.permissionService.GetPermissionTree(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(resp))
}

// GetPermissionDetail 获取权限详情
func (h *PermissionHandler) GetPermissionDetail(c *gin.Context) {
	id := c.Param("id")
	var permissionID uint
	if _, err := fmt.Sscanf(id, "%d", &permissionID); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, "无效的权限ID"))
		return
	}

	permissionInfo, err := h.permissionService.GetPermissionByID(permissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(constant.CodeNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(permissionInfo))
}

// UpdatePermission 更新权限
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	var permissionID uint
	if _, err := fmt.Sscanf(id, "%d", &permissionID); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, "无效的权限ID"))
		return
	}

	var req dto.PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	permissionInfo, err := h.permissionService.UpdatePermission(permissionID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(permissionInfo))
}

// DeletePermission 删除权限
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id := c.Param("id")
	var permissionID uint
	if _, err := fmt.Sscanf(id, "%d", &permissionID); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, "无效的权限ID"))
		return
	}

	if err := h.permissionService.DeletePermission(permissionID); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
