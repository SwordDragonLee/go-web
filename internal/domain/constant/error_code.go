package constant

// 错误码定义
const (
	// 通用错误码
	CodeSuccess        = 200  // 成功
	CodeBadRequest     = 400  // 请求参数错误
	CodeUnauthorized   = 401  // 未授权
	CodeForbidden      = 403  // 禁止访问
	CodeNotFound       = 404  // 资源不存在
	CodeInternalError  = 500  // 服务器内部错误

	// 业务错误码
	CodeUserNotFound      = 1001 // 用户不存在
	CodeUserAlreadyExists = 1002 // 用户已存在
	CodePasswordError     = 1003 // 密码错误
	CodeUserDisabled      = 1004 // 用户已被禁用
	CodeUserInactive      = 1005 // 用户未激活
	CodeTokenInvalid      = 1006 // Token无效
	CodeTokenExpired      = 1007 // Token已过期
	CodeOldPasswordError  = 1008 // 旧密码错误
	CodeEmailNotVerified  = 1009 // 邮箱未验证
)

// ErrorMessages 错误消息映射
var ErrorMessages = map[int]string{
	CodeSuccess:        "操作成功",
	CodeBadRequest:     "请求参数错误",
	CodeUnauthorized:   "未授权，请先登录",
	CodeForbidden:      "禁止访问",
	CodeNotFound:       "资源不存在",
	CodeInternalError:  "服务器内部错误",
	CodeUserNotFound:   "用户不存在",
	CodeUserAlreadyExists: "用户已存在",
	CodePasswordError:     "密码错误",
	CodeUserDisabled:      "用户已被禁用",
	CodeUserInactive:      "用户未激活",
	CodeTokenInvalid:      "Token无效",
	CodeTokenExpired:      "Token已过期",
	CodeOldPasswordError:  "旧密码错误",
	CodeEmailNotVerified:  "邮箱未验证",
}

// GetMessage 获取错误消息
func GetMessage(code int) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

