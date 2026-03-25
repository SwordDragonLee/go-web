package router

import (
	"github.com/SwordDragonLee/go-web/internal/api/handler"
	"github.com/SwordDragonLee/go-web/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoleRoutes 设置角色相关路由（需要管理员权限）
func SetupRoleRoutes(v1 *gin.RouterGroup, roleHandler *handler.RoleHandler) {
	roles := v1.Group("/roles")
	roles.Use(middleware.AuthMiddleware())
	// roles.Use(middleware.RoleMiddleware("admin")) // 只允许管理员访问
	{
		roles.POST("", roleHandler.CreateRole)                 // 创建角色
		roles.GET("", roleHandler.GetRoleList)                 // 获取角色列表
		roles.POST("/:id/permissions", roleHandler.AssignPermissions) // 为角色分配权限
		roles.GET("/:id/permissions", roleHandler.GetRolePermissions)  // 获取角色的权限列表
	}
}
