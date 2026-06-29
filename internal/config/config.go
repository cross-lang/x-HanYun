// Package config 提供应用程序配置管理功能
// 支持从配置文件加载配置，包括服务配置、数据库配置、日志配置等
package config

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"

	"x-HanYun/pkg/constants"
)

// Port 是服务的默认端口
const Port = 8888

// ESConf Elastic Search 配置结构体
// 包含连接 Elastic Search 所需的所有配置项
type ESConf struct {
	Addresses string `json:",default="`   // ES 服务地址，多个地址用逗号分隔
	Username  string `json:",default="`   // ES 认证用户名
	Password  string `json:",default="`   // ES 认证密码
}

// LoggerConfig 日志配置结构体
// 定义日志输出方式、级别和远程日志服务配置
type LoggerConfig struct {
	LogDir       string `json:",default=log"` // 日志文件存储目录
	Level        string `json:",default=info"` // 日志级别：debug, info, warn, error
	EnableRemote bool   `json:",default=false"` // 是否启用远程日志推送
	RemoteURL    string `json:",default="`    // 远程日志服务地址
}

// Config 结构体定义了应用的配置信息
type Config struct {
	Name     string        `json:",default=x-HanYun"` // 应用名称
	Host     string        `json:",default=0.0.0.0"`  // 应用主机
	Port     int           `json:",default=8888"`    // 应用端口
	Interval string        `json:",default=5s"`     // 调度任务间隔
	MySQLDSN string        `json:",default="`       // MySQL 数据库连接字符串
	ESConf   *ESConf       `json:",optional"`       // Elastic Search 配置
	Logger   *LoggerConfig `json:",optional"`       // 日志配置
}

// Init 初始化配置，从配置文件加载配置信息
func Init() Config {
	var configFile string
	flag.StringVar(&configFile, "f", "config.yaml", "the config file")
	flag.Parse()

	var cfg Config
	conf.MustLoad(configFile, &cfg)

	// 设置默认的应用名称
	if cfg.Name == "" {
		cfg.Name = constants.AppName
	}

	// MySQL DSN 需要添加时区参数
	if cfg.MySQLDSN != "" {
		cfg.MySQLDSN = cfg.MySQLDSN + "&parseTime=true&loc=Local"
	}

	// 如果没有配置日志，使用默认值
	if cfg.Logger == nil {
		cfg.Logger = &LoggerConfig{
			LogDir:       "log",
			Level:        "info",
			EnableRemote: false,
		}
	}

	fmt.Printf("配置文件: %s\n", configFile)
	fmt.Printf("应用名称: %s\n", cfg.Name)
	fmt.Printf("监听地址: %s:%d\n", cfg.Host, cfg.Port)

	return cfg
}
