package log

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"x-HanYun/internal/config"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig 日志配置结构体
type LoggerConfig struct {
	LogDir       string
	Level        string
	EnableRemote bool
	RemoteURL    string
}

// Logger 全局日志实例
var Logger *zap.Logger

// RemoteCore 自定义日志核心，支持本地输出和远程推送
type RemoteCore struct {
	core         zapcore.Core // 底层zap核心
	enableRemote bool         // 是否启用远程推送
	remoteURL    string       // 远程日志服务地址
}

// Enabled 检查日志级别是否启用
func (r *RemoteCore) Enabled(level zapcore.Level) bool {
	return r.core.Enabled(level)
}

// With 添加结构化字段到日志核心
func (r *RemoteCore) With(fields []zapcore.Field) zapcore.Core {
	return &RemoteCore{
		core:         r.core.With(fields),
		enableRemote: r.enableRemote,
		remoteURL:    r.remoteURL,
	}
}

// Check 检查日志条目是否应该被记录
func (r *RemoteCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if r.Enabled(ent.Level) {
		return ce.AddCore(ent, r)
	}
	return ce
}

// Write 写入日志条目，同时处理远程推送
func (r *RemoteCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// 先写入本地
	err := r.core.Write(entry, fields)
	if err != nil {
		return err
	}

	// 异步推送到远程日志服务
	// 注释掉远程推送代码，因为当前还没有配置远程日志主机
	// if r.enableRemote && r.remoteURL != "" {
	// 	go r.pushRemote(entry, fields)
	// }
	return nil
}

// Sync 同步刷新日志缓冲区
func (r *RemoteCore) Sync() error {
	return r.core.Sync()
}

// pushRemote 推送日志到远程服务
func (r *RemoteCore) pushRemote(entry zapcore.Entry, fields []zapcore.Field) {
	defer func() {
		if err := recover(); err != nil {
			Error("远程日志推送发生panic", zap.Any("error", err))
		}
	}()

	// 构建日志消息
	logMsg := map[string]interface{}{
		"level":     entry.Level.String(),
		"message":   entry.Message,
		"timestamp": entry.Time.UnixNano() / int64(time.Millisecond),
		"logger":    entry.LoggerName,
		"caller":    entry.Caller.String(),
		"stack":     entry.Stack,
		"fields":    fields,
	}

	// 转换为JSON
	jsonData, err := json.Marshal(logMsg)
	if err != nil {
		Error("日志消息序列化失败", zap.Error(err))
		return
	}

	// 发送HTTP请求
	req, err := http.NewRequest("POST", r.remoteURL, strings.NewReader(string(jsonData)))
	if err != nil {
		Error("创建远程日志请求失败", zap.Error(err))
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// 设置超时时间
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		Error("远程日志推送失败", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		Error("远程日志服务返回错误状态码", zap.Int("statusCode", resp.StatusCode))
	}
}

// levelSet 日志级别映射表
var levelSet = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}

// InitLogger 初始化日志模块
func InitLogger(conf config.LoggerConfig) error {
	// 参数验证
	if conf.LogDir == "" {
		conf.LogDir = "./logs"
	}

	// 确保日志目录存在
	if err := os.MkdirAll(conf.LogDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 配置文件轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   conf.LogDir + "/" + time.Now().Format("2006-01-02") + ".log",
		MaxSize:    100,  // 单个文件最大100MB
		MaxBackups: 30,   // 保留30个备份文件
		MaxAge:     7,    // 保留7天
		Compress:   true, // 启用压缩
	}

	// 配置输出目标：同时输出到文件和标准输出
	writeSyncer := zapcore.AddSync(io.MultiWriter(lumberjackLogger, os.Stdout))

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000")) // 精确到毫秒
	}
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 解析日志级别
	level, ok := levelSet[strings.ToLower(conf.Level)]
	if !ok {
		level = zapcore.InfoLevel
		Error("未知的日志级别, 使用默认级别: info", zap.String("level", conf.Level))
	}

	// 创建基础日志核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		level,
	)

	// 创建全局日志实例
	Logger = zap.New(&RemoteCore{
		core:         core,
		enableRemote: conf.EnableRemote,
		remoteURL:    conf.RemoteURL,
	}, zap.AddCaller(), zap.AddCallerSkip(1))

	Logger.Info("日志模块初始化完成",
		zap.String("logDir", conf.LogDir),
		zap.String("level", level.String()),
		zap.Bool("enableRemote", conf.EnableRemote),
		zap.String("remoteURL", conf.RemoteURL),
	)

	return nil
}

// TraceIDKey 上下文追踪ID的键名
const TraceIDKey = "trace_id"

// TraceIDFromContext 从 Context 获取 trace_id
func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if v := ctx.Value(TraceIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// ContextWithTraceID 在 Context 中设置 trace_id
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithContext 获取带有上下文追踪信息的日志实例
func WithContext(ctx context.Context) *zap.Logger {
	if Logger == nil {
		// 如果Logger未初始化，返回默认logger避免nil panic
		return zap.NewNop()
	}

	traceID := TraceIDFromContext(ctx)
	if traceID != "" {
		return Logger.With(zap.String("trace_id", traceID))
	}
	return Logger
}

// Debug 记录调试级别日志
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

// Info 记录信息级别日志
func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

// Warn 记录警告级别日志
func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

// Error 记录错误级别日志
func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

// Panic 记录恐慌级别日志并触发panic
func Panic(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Panic(msg, fields...)
	}
}

// Fatal 记录致命级别日志并退出程序
func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}

// Sync 同步刷新所有日志缓冲区
func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}
