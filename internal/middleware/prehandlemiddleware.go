// Package middleware 提供 HTTP 中间件功能
// 包含请求预处理、认证、日志记录等中间件实现
package middleware

import (
	"net/http"

	"x-HanYun/pkg/utils"
)

// PreHandleMiddleware 预处理中间件
// 用于在请求处理前提取和设置请求上下文信息
// 包括请求来源、请求ID、真实IP等信息，便于后续处理和日志追踪
type PreHandleMiddleware struct{}

// NewPreHandleMiddleware 创建一个新的预处理中间件实例
// 返回值是预处理中间件指针
func NewPreHandleMiddleware() *PreHandleMiddleware {
	return &PreHandleMiddleware{}
}

// Handle 实现中间件的处理逻辑
// 从请求头中提取上下文信息并设置到 context.Context 中
// 参数 next 是下一个处理器
// 返回值是包装后的处理器
func (m *PreHandleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取并设置请求来源信息
		// X-Origin 头通常由网关或代理设置，用于标识请求的原始来源
		c2 := utils.SetXOrigin(r.Context(), r.Header.Get("X-Origin"))

		// 从请求头中获取并设置请求ID
		// X-Request-Id 用于请求追踪，便于在日志中关联同一请求的所有记录
		c2 = utils.SetReqId(c2, r.Header.Get("X-Request-Id"))

		// 从请求头中获取并设置请求真实IP
		// X-Real-Ip 通常由反向代理设置，包含客户端的真实IP地址
		c2 = utils.SetReqRealIP(c2, r.Header.Get("X-Real-Ip"))

		// 继续执行下一个处理器，并传递更新后的上下文
		next(w, r.WithContext(c2))
	}
}
