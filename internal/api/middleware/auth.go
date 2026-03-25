package middleware

import (
	"net/http"
	"strings"

	"github.com/SwordDragonLee/go-web/internal/domain/constant"
	"github.com/SwordDragonLee/go-web/internal/domain/dto"
	"github.com/SwordDragonLee/go-web/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, "缺少Authorization头"))
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, "Authorization格式错误"))
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, "Token为空"))
			c.Abort()
			return
		}

		// 解析Token
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, "Token无效或已过期"))
			c.Abort()
			return
		}

		// 将用户信息和Token存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("token", token) // 存储Token，方便后续使用

		c.Next()
	}
}

// RoleMiddleware 角色权限中间件，支持多个角色
// 参数 roles: 允许访问的角色列表，如果用户角色在列表中则允许访问
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, dto.Error(constant.CodeForbidden, "权限不足"))
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, dto.Error(constant.CodeForbidden, "角色信息格式错误"))
			c.Abort()
			return
		}

		// 检查用户角色是否在允许的角色列表中
		hasPermission := false
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, dto.Error(constant.CodeForbidden, "权限不足，需要以下角色之一: "+strings.Join(allowedRoles, ", ")))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件（向后兼容）
func AdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware("admin")
}
