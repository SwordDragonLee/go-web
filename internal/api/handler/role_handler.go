package handler

import (
	"fmt"
	"net/http"

	"github.com/SwordDragonLee/go-web/internal/domain/constant"
	"github.com/SwordDragonLee/go-web/internal/domain/dto"
	"github.com/gin-gonic/gin"

	"github.com/SwordDragonLee/go-web/internal/service"
)

type RoleHandler struct {
	roleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}
	roleInfo, err := h.roleService.CreateRole(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, roleInfo)
}

func (h *RoleHandler) GetRoleList(c *gin.Context) {
	var req dto.RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	resp, err := h.roleService.GetRoleList(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(resp))
}

// AssignPermissions 为角色分配权限
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	var req dto.AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	if err := h.roleService.AssignPermissions(req.RoleID, req.PermissionIDs); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}

// GetRolePermissions 获取角色的权限列表
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	id := c.Param("id")
	var roleID uint
	if _, err := fmt.Sscanf(id, "%d", &roleID); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, "无效的角色ID"))
		return
	}

	permissions, err := h.roleService.GetRolePermissions(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(permissions))
}
