package router

import (
	"time"

	"github.com/SwordDragonLee/go-web/internal/api/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由（主入口，组装所有模块路由）
func SetupRouter(userHandler *handler.UserHandler, roleHandler *handler.RoleHandler, permissionHandler *handler.PermissionHandler) *gin.Engine {
	r := gin.Default()

	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9000"}, // 前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 允许携带凭证（Cookie、Authorization 等）
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1路由组
	v1 := r.Group("/api/v1")
	{
		// 注册各个模块的路由
		SetupAuthRoutes(v1, userHandler)       // 认证相关路由（无需认证）
		SetupUserRoutes(v1, userHandler)       // 用户相关路由（需要认证）
		SetupRoleRoutes(v1, roleHandler)       // 角色相关路由（需要管理员权限）
		SetupAdminRoutes(v1, userHandler)      // 管理员相关路由（需要管理员权限）
		SetupPermissionRoutes(v1, permissionHandler) // 权限相关路由（需要管理员权限）
	}

	return r
}
