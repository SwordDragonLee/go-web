package handler

import (
	"net/http"
	"strconv"

	"github.com/SwordDragonLee/go-web/internal/domain/constant"
	"github.com/SwordDragonLee/go-web/internal/domain/dto"
	"github.com/SwordDragonLee/go-web/internal/domain/model"
	"github.com/SwordDragonLee/go-web/internal/service"
	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 200 {object} dto.Response
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	//ShouldBindJSON自动将 HTTP 请求的参数（JSON / 表单 / URL 查询等）绑定到 Go 结构体 / 变量，简化手动解析请求参数的繁琐操作。
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	userInfo, err := h.userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(userInfo))
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} dto.Response
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	// 获取客户端IP和User-Agent
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	loginResp, err := h.userService.Login(&req, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(loginResp))
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Router /api/v1/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, constant.GetMessage(constant.CodeUnauthorized)))
		return
	}

	userInfo, err := h.userService.GetUserInfo(userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(userInfo))
}

// UpdateProfile 更新个人信息
// @Summary 更新个人信息
// @Description 更新当前登录用户的个人信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "更新信息"
// @Success 200 {object} dto.Response
// @Router /api/v1/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, constant.GetMessage(constant.CodeUnauthorized)))
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	if err := h.userService.UpdateProfile(userID.(uint), &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMessage("更新成功", nil))
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "密码信息"
// @Success 200 {object} dto.Response
// @Router /api/v1/user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, constant.GetMessage(constant.CodeUnauthorized)))
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	if err := h.userService.ChangePassword(userID.(uint), &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMessage("密码修改成功", nil))
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表（管理员功能）
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态筛选"
// @Param role query string false "角色筛选"
// @Success 200 {object} dto.Response
// @Router /api/v1/users [get]
func (h *UserHandler) GetUserList(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	resp, err := h.userService.GetUserList(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(resp))
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新用户状态（管理员功能）
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param status body int true "状态：0-未激活，1-激活，2-禁用"
// @Success 200 {object} dto.Response
// @Router /api/v1/users/{id}/status [put]
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, "无效的用户ID"))
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	if err := h.userService.UpdateUserStatus(uint(id), model.UserStatus(req.Status)); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMessage("状态更新成功", nil))
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除用户（管理员功能）
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} dto.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, "无效的用户ID"))
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMessage("删除成功", nil))
}

// RefreshToken 刷新Token
// @Summary 刷新Token
// @Description 刷新当前用户的Token
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Router /api/v1/user/refresh-token [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, constant.GetMessage(constant.CodeUnauthorized)))
		return
	}

	username, _ := c.Get("username")
	role, _ := c.Get("role")

	// 重新生成Token
	token, err := h.userService.RefreshToken(userID.(uint), username.(string), role.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"token":      token,
		"token_type": "Bearer",
		"expires_in": h.userService.GetTokenExpireDuration(),
	}))
}

// GetLoginLogs 获取登录日志
// @Summary 获取登录日志
// @Description 获取当前用户的登录日志
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} dto.Response
// @Router /api/v1/user/login-logs [get]
func (h *UserHandler) GetLoginLogs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, constant.GetMessage(constant.CodeUnauthorized)))
		return
	}

	var req struct {
		Page     int `form:"page" binding:"omitempty,min=1"`
		PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	// 设置默认值
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	logs, total, err := h.userService.GetLoginLogs(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"total": total,
		"list":  logs,
	}))
}

// SendVerificationCode 发送验证码
// @Summary 发送验证码
// @Description 发送邮箱验证码（用于邮箱验证或密码重置）
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.SendVerificationCodeRequest true "验证码请求"
// @Success 200 {object} dto.Response
// @Router /api/v1/auth/send-code [post]
func (h *UserHandler) SendVerificationCode(c *gin.Context) {
	var req dto.SendVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	if err := h.userService.SendVerificationCode(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMessage("验证码已发送", nil))
}

// VerifyEmail 验证邮箱
// @Summary 验证邮箱
// @Description 验证用户邮箱
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.VerifyEmailRequest true "验证请求"
// @Success 200 {object} dto.Response
// @Router /api/v1/auth/verify-email [post]
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	if err := h.userService.VerifyEmail(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMessage("邮箱验证成功", nil))
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Description 通过邮箱验证码重置密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "重置密码请求"
// @Success 200 {object} dto.Response
// @Router /api/v1/auth/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	if err := h.userService.ResetPassword(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMessage("密码重置成功", nil))
}

// AssignRoles 为用户分配角色
// @Summary 为用户分配角色
// @Description 为用户分配一个或多个角色
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AssignRoleRequest true "角色分配请求"
// @Success 200 {object} dto.Response
// @Router /api/v1/users/roles [post]
func (h *UserHandler) AssignRoles(c *gin.Context) {
	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	if err := h.userService.AssignRoles(req.UserID, req.RoleIDs); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMessage("角色分配成功", nil))
}

// GetUserRoles 获取用户的角色列表
// @Summary 获取用户的角色列表
// @Description 获取指定用户的角色列表
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} dto.Response
// @Router /api/v1/users/{id}/roles [get]
func (h *UserHandler) GetUserRoles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, "无效的用户ID"))
		return
	}

	roles, err := h.userService.GetUserRoles(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(roles))
}

// GetUserPermissions 获取用户的所有权限
// @Summary 获取用户的所有权限
// @Description 获取指定用户的所有权限（通过角色合并）
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} dto.Response
// @Router /api/v1/users/{id}/permissions [get]
func (h *UserHandler) GetUserPermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, "无效的用户ID"))
		return
	}

	permissions, err := h.userService.GetUserPermissions(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(constant.CodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(permissions))
}
