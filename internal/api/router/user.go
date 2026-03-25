package router

import (
	"github.com/SwordDragonLee/go-web/internal/api/handler"
	"github.com/SwordDragonLee/go-web/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户相关路由（需要认证）
func SetupUserRoutes(v1 *gin.RouterGroup, userHandler *handler.UserHandler) {
	user := v1.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/profile", userHandler.GetProfile)          // 获取个人信息
		user.PUT("/profile", userHandler.UpdateProfile)       // 更新个人信息
		user.PUT("/password", userHandler.ChangePassword)     // 修改密码
		user.GET("/login-logs", userHandler.GetLoginLogs)     // 获取登录日志
		user.POST("/refresh-token", userHandler.RefreshToken) // 刷新Token
	}
}
