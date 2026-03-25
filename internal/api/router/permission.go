package router

import (
	"github.com/SwordDragonLee/go-web/internal/api/handler"
	"github.com/SwordDragonLee/go-web/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// SetupPermissionRoutes 设置权限相关路由（需要管理员权限）
func SetupPermissionRoutes(v1 *gin.RouterGroup, permissionHandler *handler.PermissionHandler) {
	permissions := v1.Group("/permissions")
	permissions.Use(middleware.AuthMiddleware())
	// permissions.Use(middleware.RoleMiddleware("admin")) // 只允许管理员访问
	{
		permissions.POST("", permissionHandler.CreatePermission)       // 创建权限
		permissions.GET("", permissionHandler.GetPermissionList)       // 获取权限列表
		permissions.GET("/tree", permissionHandler.GetPermissionTree)  // 获取权限树
		permissions.GET("/:id", permissionHandler.GetPermissionDetail) // 获取权限详情
		permissions.PUT("/:id", permissionHandler.UpdatePermission)    // 更新权限
		permissions.DELETE("/:id", permissionHandler.DeletePermission) // 删除权限
	}
}
