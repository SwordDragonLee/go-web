package service

import (
	"errors"
	"time"

	"github.com/SwordDragonLee/go-web/internal/domain/constant"
	"github.com/SwordDragonLee/go-web/internal/domain/dto"
	"github.com/SwordDragonLee/go-web/internal/domain/model"
	"github.com/SwordDragonLee/go-web/internal/repository"
	"github.com/SwordDragonLee/go-web/internal/utils"
)

// UserService 用户服务接口
type UserService interface {
	Register(req *dto.RegisterRequest) (*dto.UserInfo, error)
	Login(req *dto.LoginRequest, ip, userAgent string) (*dto.LoginResponse, error)
	GetUserInfo(userID uint) (*dto.UserInfo, error)
	UpdateProfile(userID uint, req *dto.UpdateProfileRequest) error
	ChangePassword(userID uint, req *dto.ChangePasswordRequest) error
	GetUserList(req *dto.UserListRequest) (*dto.UserListResponse, error)
	UpdateUserStatus(userID uint, status model.UserStatus) error
	DeleteUser(userID uint) error
	RefreshToken(userID uint, username, role string) (string, error)
	GetTokenExpireDuration() int64
	GetLoginLogs(userID uint, page, pageSize int) ([]*model.LoginLog, int64, error)
	SendVerificationCode(req *dto.SendVerificationCodeRequest) error
	VerifyEmail(req *dto.VerifyEmailRequest) error
	ResetPassword(req *dto.ResetPasswordRequest) error
	AssignRoles(userID uint, roleIDs []uint) error
	GetUserRoles(userID uint) ([]*dto.RoleInfo, error)
	GetUserPermissions(userID uint) (*dto.UserPermissionResponse, error)
}

type userService struct {
	userRepo             repository.UserRepository
	loginLogRepo         repository.LoginLogRepository
	verificationCodeRepo repository.VerificationCodeRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository, loginLogRepo repository.LoginLogRepository, verificationCodeRepo repository.VerificationCodeRepository) UserService {
	return &userService{
		userRepo:             userRepo,
		loginLogRepo:         loginLogRepo,
		verificationCodeRepo: verificationCodeRepo,
	}
}

// Register 用户注册
func (s *userService) Register(req *dto.RegisterRequest) (*dto.UserInfo, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New(constant.GetMessage(constant.CodeUserAlreadyExists))
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
			return nil, errors.New("邮箱已被注册")
		}
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		if _, err := s.userRepo.GetByPhone(req.Phone); err == nil {
			return nil, errors.New("手机号已被注册")
		}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Status:   model.UserStatusInactive, // 默认未激活，需要邮箱验证
		Role:     model.RoleUser,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	return s.toUserInfo(user), nil
}

// Login 用户登录
func (s *userService) Login(req *dto.LoginRequest, ip, userAgent string) (*dto.LoginResponse, error) {
	// 查找用户（支持用户名、邮箱、手机号登录）
	var user *model.User
	var err error
	// 尝试用户名
	if user, err = s.userRepo.GetByUsername(req.Username); err != nil {
		// 尝试邮箱
		if user, err = s.userRepo.GetByEmail(req.Username); err != nil {
			// 尝试手机号
			if user, err = s.userRepo.GetByPhone(req.Username); err != nil {
				// 记录登录失败日志
				s.recordLoginLog(0, req.Username, ip, userAgent, "failed", "用户不存在")
				return nil, errors.New(constant.GetMessage(constant.CodeUserNotFound))
			}
		}
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		// 记录登录失败日志
		s.recordLoginLog(user.ID, user.Username, ip, userAgent, "failed", "密码错误")
		return nil, errors.New(constant.GetMessage(constant.CodePasswordError))
	}

	// 检查用户状态
	if user.IsDisabled() {
		s.recordLoginLog(user.ID, user.Username, ip, userAgent, "failed", "用户已被禁用")
		return nil, errors.New(constant.GetMessage(constant.CodeUserDisabled))
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	// 更新最后登录信息
	s.userRepo.UpdateLastLogin(user.ID, ip)

	// 记录登录成功日志
	s.recordLoginLog(user.ID, user.Username, ip, userAgent, "success", "登录成功")

	// 返回登录响应
	return &dto.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: utils.GetTokenExpireDuration(),
		User:      s.toUserInfo(user),
	}, nil
}

// GetUserInfo 获取用户信息
func (s *userService) GetUserInfo(userID uint) (*dto.UserInfo, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New(constant.GetMessage(constant.CodeUserNotFound))
	}
	return s.toUserInfo(user), nil
}

// UpdateProfile 更新个人信息
func (s *userService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New(constant.GetMessage(constant.CodeUserNotFound))
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	return s.userRepo.Update(user)
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New(constant.GetMessage(constant.CodeUserNotFound))
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errors.New(constant.GetMessage(constant.CodeOldPasswordError))
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

// GetUserList 获取用户列表
func (s *userService) GetUserList(req *dto.UserListRequest) (*dto.UserListResponse, error) {
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

	users, total, err := s.userRepo.List(page, pageSize, req.Keyword, req.Status, req.Role)
	if err != nil {
		return nil, errors.New("获取用户列表失败")
	}

	userInfos := make([]*dto.UserInfo, 0, len(users))
	for _, user := range users {
		userInfos = append(userInfos, s.toUserInfo(user))
	}

	return &dto.UserListResponse{
		Total: total,
		List:  userInfos,
	}, nil
}

// UpdateUserStatus 更新用户状态
func (s *userService) UpdateUserStatus(userID uint, status model.UserStatus) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New(constant.GetMessage(constant.CodeUserNotFound))
	}

	user.Status = status
	return s.userRepo.Update(user)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(userID uint) error {
	return s.userRepo.Delete(userID)
}

// toUserInfo 转换为UserInfo DTO
func (s *userService) toUserInfo(user *model.User) *dto.UserInfo {
	userInfo := &dto.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Status:   int(user.Status),
		Role:     string(user.Role),
	}

	if user.LastLoginAt != nil {
		lastLoginAt := user.LastLoginAt.Format(time.RFC3339)
		userInfo.LastLoginAt = &lastLoginAt
	}

	return userInfo
}

// RefreshToken 刷新Token
func (s *userService) RefreshToken(userID uint, username, role string) (string, error) {
	// 验证用户是否存在
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", errors.New(constant.GetMessage(constant.CodeUserNotFound))
	}

	// 生成新Token
	token, err := utils.GenerateToken(userID, username, role)
	if err != nil {
		return "", errors.New("生成Token失败")
	}

	return token, nil
}

// GetTokenExpireDuration 获取Token过期时间
func (s *userService) GetTokenExpireDuration() int64 {
	return utils.GetTokenExpireDuration()
}

// GetLoginLogs 获取登录日志
func (s *userService) GetLoginLogs(userID uint, page, pageSize int) ([]*model.LoginLog, int64, error) {
	return s.loginLogRepo.List(userID, page, pageSize)
}

// SendVerificationCode 发送验证码
func (s *userService) SendVerificationCode(req *dto.SendVerificationCodeRequest) error {
	codeType := model.VerificationCodeType(req.Type)

	// 检查邮箱是否存在（根据类型）
	if codeType == model.CodeTypePasswordReset {
		_, err := s.userRepo.GetByEmail(req.Email)
		if err != nil {
			return errors.New("邮箱未注册")
		}
	} else if codeType == model.CodeTypeEmailVerification {
		_, err := s.userRepo.GetByEmail(req.Email)
		if err == nil {
			return errors.New("邮箱已被注册")
		}
	}

	// 检查是否在1分钟内发送过验证码
	latestCode, err := s.verificationCodeRepo.GetLatestByEmailAndType(req.Email, codeType)
	if err == nil && time.Since(latestCode.CreatedAt) < time.Minute {
		return errors.New("验证码发送过于频繁，请稍后再试")
	}

	// 生成验证码
	code := utils.GenerateVerificationCode()

	// 创建验证码记录
	verificationCode := &model.VerificationCode{
		Email:     req.Email,
		Code:      code,
		Type:      codeType,
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10分钟过期
		Used:      false,
	}

	if err := s.verificationCodeRepo.Create(verificationCode); err != nil {
		return errors.New("创建验证码失败")
	}

	// TODO: 这里应该发送邮件，目前只是记录日志
	// 实际项目中应该集成邮件服务（如SendGrid、阿里云邮件等）
	// sendEmail(req.Email, code, codeType)

	return nil
}

// VerifyEmail 验证邮箱
func (s *userService) VerifyEmail(req *dto.VerifyEmailRequest) error {
	// 获取最新的验证码
	code, err := s.verificationCodeRepo.GetLatestByEmailAndType(req.Email, model.CodeTypeEmailVerification)
	if err != nil {
		return errors.New("验证码不存在或已过期")
	}

	// 检查验证码是否有效
	if !code.IsValid() {
		return errors.New("验证码已使用或已过期")
	}

	// 验证验证码
	if code.Code != req.Code {
		return errors.New("验证码错误")
	}

	// 标记验证码为已使用
	if err := s.verificationCodeRepo.MarkAsUsed(code.ID); err != nil {
		return errors.New("验证失败")
	}

	// 激活用户
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.Status = model.UserStatusActive
	return s.userRepo.Update(user)
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(req *dto.ResetPasswordRequest) error {
	// 获取最新的验证码
	code, err := s.verificationCodeRepo.GetLatestByEmailAndType(req.Email, model.CodeTypePasswordReset)
	if err != nil {
		return errors.New("验证码不存在或已过期")
	}

	// 检查验证码是否有效
	if !code.IsValid() {
		return errors.New("验证码已使用或已过期")
	}

	// 验证验证码
	if code.Code != req.Code {
		return errors.New("验证码错误")
	}

	// 标记验证码为已使用
	if err := s.verificationCodeRepo.MarkAsUsed(code.ID); err != nil {
		return errors.New("验证失败")
	}

	// 获取用户
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

// recordLoginLog 记录登录日志
func (s *userService) recordLoginLog(userID uint, username, ip, userAgent, status, message string) {
	log := &model.LoginLog{
		UserID:    userID,
		Username:  username,
		IP:        ip,
		UserAgent: userAgent,
		Status:    status,
		Message:   message,
	}
	_ = s.loginLogRepo.Create(log) // 忽略错误，不影响主流程
}

// AssignRoles 为用户分配角色
func (s *userService) AssignRoles(userID uint, roleIDs []uint) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 分配角色
	if err := s.userRepo.AssignRoles(userID, roleIDs); err != nil {
		return errors.New("分配角色失败")
	}

	// 注意：这里更新了用户的角色，但不修改 User 表中的 Role 字段
	// User.Role 字段保留作为主要角色或向后兼容
	_ = user

	return nil
}

// GetUserRoles 获取用户的角色列表
func (s *userService) GetUserRoles(userID uint) ([]*dto.RoleInfo, error) {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	roles, err := s.userRepo.GetRoles(userID)
	if err != nil {
		return nil, errors.New("获取角色列表失败")
	}

	roleInfos := make([]*dto.RoleInfo, 0, len(roles))
	for _, role := range roles {
		roleInfos = append(roleInfos, &dto.RoleInfo{
			ID:        role.ID,
			Name:      role.Name,
			Code:      role.Code,
			Status:    int(role.Status),
			CreatedAt: role.CreatedAt,
			UpdatedAt: role.UpdatedAt,
		})
	}

	return roleInfos, nil
}

// GetUserPermissions 获取用户的所有权限
func (s *userService) GetUserPermissions(userID uint) (*dto.UserPermissionResponse, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 获取用户的角色
	roles, err := s.userRepo.GetRoles(userID)
	if err != nil {
		return nil, errors.New("获取角色列表失败")
	}

	// 转换为 RoleInfo
	roleInfos := make([]*dto.RoleInfo, 0, len(roles))
	for _, role := range roles {
		roleInfos = append(roleInfos, &dto.RoleInfo{
			ID:        role.ID,
			Name:      role.Name,
			Code:      role.Code,
			Status:    int(role.Status),
			CreatedAt: role.CreatedAt,
			UpdatedAt: role.UpdatedAt,
		})
	}

	// 获取用户的权限
	permissions, err := s.userRepo.GetUserPermissions(userID)
	if err != nil {
		return nil, errors.New("获取权限列表失败")
	}

	// 转换为 PermissionInfo
	permissionInfos := make([]*dto.PermissionInfo, 0, len(permissions))
	for _, permission := range permissions {
		permissionInfos = append(permissionInfos, &dto.PermissionInfo{
			ID:        permission.ID,
			Name:      permission.Name,
			Code:      permission.Code,
			Type:      int(permission.Type),
			Path:      permission.Path,
			Method:    permission.Method,
			ParentID:  permission.ParentID,
			Sort:      permission.Sort,
			Status:    int(permission.Status),
			Remark:    permission.Remark,
			Icon:      permission.Icon,
			Component: permission.Component,
			Redirect:  permission.Redirect,
			Hidden:    permission.Hidden,
			KeepAlive: permission.KeepAlive,
			CreatedAt: permission.CreatedAt,
			UpdatedAt: permission.UpdatedAt,
		})
	}

	return &dto.UserPermissionResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Roles:       roleInfos,
		Permissions: permissionInfos,
	}, nil
}
