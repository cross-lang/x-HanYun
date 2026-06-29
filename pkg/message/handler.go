package message

import (
	"context"
	"errors"
	"x-HanYun/pkg/core/log"

	"go.uber.org/zap"
)

// 处理方法
type handlerFunc func(ctx context.Context, data *Entity, msg *Message) error

// Process 消息处理程序
type Process struct {
	handlerMap map[string]handlerFunc
}

func NewProcess() *Process {
	return &Process{
		handlerMap: make(map[string]handlerFunc),
	}
}

// RegisterHandler 注册处理器
func (p *Process) RegisterHandler(name string, handler handlerFunc) {
	p.handlerMap[name] = handler
}

// Handle 处理事件
func (p *Process) Handle(ctx context.Context, data *Entity) error {
	if len(data.Msgs) < 1 {
		return errors.New("<<<<<<<< Msgs are empty")
	}
	key := data.Source + "." + data.Key

	log.WithContext(ctx).Debug(">>>>>>>> NotificationData.Source", zap.String("source", data.Source))
	log.WithContext(ctx).Debug(">>>>>>>> NotificationData.Key", zap.String("key", data.Key))
	log.WithContext(ctx).Debug(">>>>>>>> NotificationData.ToUsers", zap.Any("toUsers", data.ToUsers))
	log.WithContext(ctx).Debug(">>>>>>>> NotificationData.MsgTime", zap.Int64("msgTime", data.MsgTime))
	log.WithContext(ctx).Debug(">>>>>>>> NotificationData.Extra", zap.String("extra", data.Extra))
	log.WithContext(ctx).Debug(">>>>>>>> NotificationData.MsgId", zap.Int64("msgId", data.MsgId))
	log.WithContext(ctx).Debug(">>>>>>>> Message.Terminals", zap.Any("terminals", data.Msgs[0].Terminals))
	log.WithContext(ctx).Debug(">>>>>>>> Message.Body", zap.Any("body", data.Msgs[0].Body))
	log.WithContext(ctx).Debug(">>>>>>>> Message.Category", zap.String("category", data.Msgs[0].Category))
	log.WithContext(ctx).Debug(">>>>>>>> Message.Template", zap.String("template", data.Msgs[0].Template))
	log.WithContext(ctx).Debug(">>>>>>>> Message.Version", zap.Int64("version", data.Msgs[0].Version))
	log.WithContext(ctx).Debug(">>>>>>>> Message.Ext", zap.String("ext", data.Msgs[0].Ext))
	log.WithContext(ctx).Debug(">>>>>>>> Message.Nopop", zap.Bool("nopop", data.Msgs[0].Nopop))

	// 获取处理方法
	fn, ok := p.handlerMap[key]
	if !ok {
		return errors.New("<<<<<<<< No Handler Exists for " + key)
	}

	// 调用处理方法
	return fn(ctx, data, data.Msgs[0])
}
