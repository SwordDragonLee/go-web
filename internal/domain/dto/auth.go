package dto

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名
	Email    string `json:"email" binding:"required,email"`           // 邮箱
	Phone    string `json:"phone" binding:"required,len=11"`          // 手机号
	Password string `json:"password" binding:"required,min=6,max=20"` // 密码
	Nickname string `json:"nickname"`                                 // 昵称（可选）
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名/邮箱/手机号
	Password string `json:"password" binding:"required"` // 密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string    `json:"token"`      // JWT Token
	TokenType string    `json:"token_type"` // Token类型
	ExpiresIn int64     `json:"expires_in"` // 过期时间（秒）
	User      *UserInfo `json:"user"`       // 用户信息
}

// UserInfo 用户信息
type UserInfo struct {
	ID          uint    `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Nickname    string  `json:"nickname"`
	Avatar      string  `json:"avatar"`
	Status      int     `json:"status"`
	Role        string  `json:"role"`
	LastLoginAt *string `json:"last_login_at"`
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像URL
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`              // 旧密码
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"` // 新密码
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`               // 邮箱
	Code        string `json:"code" binding:"required"`                      // 验证码
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"` // 新密码
}

// SendVerificationCodeRequest 发送验证码请求
type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required,email"`                                  // 邮箱
	Type  string `json:"type" binding:"required,oneof=email_verification password_reset"` // 类型
}

// VerifyEmailRequest 验证邮箱请求
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"` // 邮箱
	Code  string `json:"code" binding:"required"`        // 验证码
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`              // 页码
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"` // 每页数量
	Keyword  string `form:"keyword"`                                     // 搜索关键词
	Status   *int   `form:"status"`                                      // 状态筛选
	Role     string `form:"role"`                                        // 角色筛选
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64       `json:"total"` // 总数
	List  []*UserInfo `json:"list"`  // 列表
}
