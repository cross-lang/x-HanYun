// Package svc 提供服务上下文管理功能
// 服务上下文包含应用程序运行时所需的所有共享资源和配置
package svc

import (
	"github.com/zeromicro/go-zero/rest"

	"x-HanYun/internal/config"
	"x-HanYun/internal/middleware"
)

// ServiceContext 服务上下文结构体
// 存储应用程序运行时所需的所有共享资源，包括配置、数据库连接、中间件等
// 所有 Handler 和 Logic 层都可以通过 ServiceContext 访问这些共享资源
type ServiceContext struct {
	Config              config.Config   // 应用配置，包含服务、数据库、日志等配置
	PreHandleMiddleware rest.Middleware // 预处理中间件，用于请求前处理（如设置上下文信息）
}

// NewServiceContext 创建一个新的服务上下文实例
// 参数 c 是应用配置，从配置文件或环境变量加载
// 返回值是服务上下文指针，包含所有初始化后的共享资源
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		PreHandleMiddleware: middleware.NewPreHandleMiddleware().Handle,
	}
}
