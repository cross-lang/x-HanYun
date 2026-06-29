package utils

import (
	"context"
)

// 定义上下文键常量
const (
	XOriginKey   = "x-origin"   // 请求来源信息键
	ReqIdKey     = "req-id"     // 请求ID键
	ReqRealIPKey = "req-real-ip" // 请求真实IP键
)

// SetXOrigin 设置请求来源信息到上下文
// 参数 ctx 是上下文
// 参数 origin 是请求来源信息
// 返回值是更新后的上下文
func SetXOrigin(ctx context.Context, origin string) context.Context {
	return context.WithValue(ctx, XOriginKey, origin)
}

// GetXOrigin 从上下文获取请求来源信息
// 参数 ctx 是上下文
// 返回值是请求来源信息，如果不存在则返回空字符串
func GetXOrigin(ctx context.Context) string {
	if origin, ok := ctx.Value(XOriginKey).(string); ok {
		return origin
	}
	return ""
}

// SetReqId 设置请求ID到上下文
// 参数 ctx 是上下文
// 参数 reqId 是请求ID
// 返回值是更新后的上下文
func SetReqId(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, ReqIdKey, reqId)
}

// GetReqId 从上下文获取请求ID
// 参数 ctx 是上下文
// 返回值是请求ID，如果不存在则返回空字符串
func GetReqId(ctx context.Context) string {
	if reqId, ok := ctx.Value(ReqIdKey).(string); ok {
		return reqId
	}
	return ""
}

// SetReqRealIP 设置请求真实IP到上下文
// 参数 ctx 是上下文
// 参数 realIP 是请求真实IP
// 返回值是更新后的上下文
func SetReqRealIP(ctx context.Context, realIP string) context.Context {
	return context.WithValue(ctx, ReqRealIPKey, realIP)
}

// GetReqRealIP 从上下文获取请求真实IP
// 参数 ctx 是上下文
// 返回值是请求真实IP，如果不存在则返回空字符串
func GetReqRealIP(ctx context.Context) string {
	if realIP, ok := ctx.Value(ReqRealIPKey).(string); ok {
		return realIP
	}
	return ""
}
