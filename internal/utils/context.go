package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// GetTokenFromContext 从 Context 获取 Token
func GetTokenFromContext(c *gin.Context) (string, error) {
	token, exists := c.Get("token")
	if !exists {
		return "", errors.New("Token不存在")
	}

	tokenString, ok := token.(string)
	if !ok {
		return "", errors.New("Token类型错误")
	}

	return tokenString, nil
}

// GetUserIDFromContext 从 Context 获取用户ID
func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("用户ID不存在")
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, errors.New("用户ID类型错误")
	}

	return id, nil
}

// GetUsernameFromContext 从 Context 获取用户名
func GetUsernameFromContext(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", errors.New("用户名不存在")
	}

	usernameString, ok := username.(string)
	if !ok {
		return "", errors.New("用户名类型错误")
	}

	return usernameString, nil
}

// GetRoleFromContext 从 Context 获取角色
func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")
	if !exists {
		return "", errors.New("角色不存在")
	}

	roleString, ok := role.(string)
	if !ok {
		return "", errors.New("角色类型错误")
	}

	return roleString, nil
}




