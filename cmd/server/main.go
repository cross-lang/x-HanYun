// Package main 是 x-HanYun 项目的入口包
// 提供基于 GoZero 框架的 RESTful API 服务
package main

import (
	"flag"

	"x-HanYun/pkg/core/log"
	"x-HanYun/pkg/swagger"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"

	"x-HanYun/internal/config"
	"x-HanYun/internal/handler"
	"x-HanYun/internal/svc"
	"x-HanYun/pkg/tasks"
	"x-HanYun/pkg/xerrors"

	"go.uber.org/zap"
)

// main 函数是项目的入口点
// 负责初始化配置、数据库连接、路由注册和服务启动
func main() {
	// 解析命令行参数
	flag.Parse()

	// 初始化配置
	// 从配置文件中读取服务配置、数据库配置等
	cfg := config.Init()

	// 初始化 Elastic Search（可选）
	// 取消注释以下代码以启用 ES 功能
	// es.Init(cfg.ESConf.Addresses, cfg.ESConf.Username, cfg.ESConf.Password)

	// 初始化 MySQL（可选）
	// 取消注释以下代码以启用 MySQL 功能
	// model.Init(cfg.MySQLDSN)

	// 创建 REST 服务器实例
	server := rest.MustNewServer(rest.RestConf{
		Host: cfg.Host,
		Port: cfg.Port,
	})
	// 确保服务器在程序退出时停止
	defer server.Stop()

	// 创建服务上下文
	// 服务上下文包含配置、数据库连接、缓存等共享资源
	sc := svc.NewServiceContext(cfg)

	// 注册路由
	// 将 API 路由与对应的 Handler 绑定
	handler.RegisterHandlers(server, sc)

	// 注册 Swagger UI 路由
	// 提供在线 API 文档和测试功能
	swagger.RegisterRoutes(server)

	// 设置全局错误处理器
	// 统一处理业务错误和系统错误，返回标准化的错误响应
	httpx.SetErrorHandlerCtx(xerrors.HTTPErrorHandler)

	// 初始化日志模块
	// 日志输出到文件和标准输出，支持日志轮转
	if err := log.InitLogger(*cfg.Logger); err != nil {
		panic("日志初始化失败: " + err.Error())
	}
	defer log.Sync()

	// 启动定时任务（每天凌晨 00:00:00 执行）
	// 使用 cron 表达式定义执行时间（支持秒级：秒 分 时 日 月 星期）
	go tasks.StartScheduledTask("0 0 0 * * *", func() {
		log.Info(">>>>>>>> 定时任务开始执行")
		// 这里可以添加具体的定时任务逻辑
		log.Info(">>>>>>>> 定时任务执行成功")
	})

	// 启动周期任务（每 24 小时执行一次）
	// 支持指定首次执行时间
	go tasks.StartPeriodicTask("24h", "", func() {
		log.Info(">>>>>>>> 周期任务开始执行")
		// 这里可以添加具体的周期任务逻辑
		log.Info(">>>>>>>> 周期任务执行成功")
	})

	// 启动服务
	// 服务将在配置的地址和端口上监听请求
	log.Info(">>>>>>>> 服务器启动中", zap.String("host", cfg.Host), zap.Int("port", cfg.Port))
	server.Start()
}
