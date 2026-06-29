package xerrors

// ErrorCode 定义错误码结构
type ErrorCode int

// 错误码格式：前三位(HTTP响应码) + 中间两位(模块编号) + 后面三位(模块内部业务编号)
const (
	// 认证相关错误 (401xxxxx)
	_ ErrorCode = 40101000 + iota
	ErrNotLogin
	ErrAccessNotPermission
	ErrDetailCurrentSession
	ErrCaptcha
	ErrInvalidUserPassword
	ErrInvalidUser

	// InternalError 内部错误 (500xxxxx)
	InternalError = 50001001 + iota
)

// ErrorMap 定义错误码到错误信息的映射
var ErrorMap = map[ErrorCode]string{
	ErrNotLogin:             "用户未登录",
	ErrAccessNotPermission:  "用户无权访问",
	ErrDetailCurrentSession: "获取当前会话信息错误",
	InternalError:           "内部错误",
	ErrCaptcha:              "验证码错误",
	ErrInvalidUser:          "用户名错误",
	ErrInvalidUserPassword:  "用户名或者密码错误",
}

// GetError 根据错误码获取错误
func GetError(code ErrorCode) Error {
	msg, exists := ErrorMap[code]
	if !exists {
		msg = "未知错误"
	}
	return NewError(code, msg)
}
