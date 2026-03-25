package router

import (
	"github.com/SwordDragonLee/go-web/internal/api/handler"
	"github.com/SwordDragonLee/go-web/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// SetupAdminRoutes 设置管理员相关路由（需要管理员权限）
func SetupAdminRoutes(v1 *gin.RouterGroup, userHandler *handler.UserHandler) {
	admin := v1.Group("/users")
	admin.Use(middleware.AuthMiddleware())
	// admin.Use(middleware.RoleMiddleware("admin")) // 只允许管理员访问
	{
		admin.GET("", userHandler.GetUserList)                 // 获取用户列表
		admin.PUT("/:id/status", userHandler.UpdateUserStatus) // 更新用户状态
		admin.DELETE("/:id", userHandler.DeleteUser)           // 删除用户
		admin.POST("/roles", userHandler.AssignRoles)          // 为用户分配角色
		admin.GET("/:id/roles", userHandler.GetUserRoles)      // 获取用户的角色列表
		admin.GET("/:id/permissions", userHandler.GetUserPermissions) // 获取用户的权限列表
	}
}
