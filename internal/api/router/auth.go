package router

import (
	"github.com/SwordDragonLee/go-web/internal/api/handler"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 设置认证相关路由（无需认证）
func SetupAuthRoutes(v1 *gin.RouterGroup, userHandler *handler.UserHandler) {
	auth := v1.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)              // 注册
		auth.POST("/login", userHandler.Login)                    // 登录
		auth.POST("/send-code", userHandler.SendVerificationCode) // 发送验证码
		auth.POST("/verify-email", userHandler.VerifyEmail)       // 验证邮箱
		auth.POST("/reset-password", userHandler.ResetPassword)   // 重置密码
	}
}

