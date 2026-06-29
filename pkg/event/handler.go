package event

import (
	"context"
	"errors"
	"fmt"
	"x-HanYun/pkg/constants"
	"x-HanYun/pkg/core/log"
	"x-HanYun/pkg/utils"

	"go.uber.org/zap"
)

// 处理方法
type handlerFunc func(ctx context.Context, data *Entity, msg string) error

// Process 事件处理程序
type Process struct {
	handlerMap map[string]handlerFunc
	ak         string
	sk         string
}

func NewProcess(ak, sk string) *Process {
	return &Process{
		handlerMap: make(map[string]handlerFunc),
		ak:         ak,
		sk:         sk,
	}
}

// RegisterHandler 注册处理器
func (p *Process) RegisterHandler(name string, handler handlerFunc) {
	p.handlerMap[name] = handler
}

// Handle 处理事件
func (p *Process) Handle(ctx context.Context, data *Entity) error {
	// 验签
	if !p.VerifySignature(data) {
		return fmt.Errorf("<<<<<<<< Verification Signature Failed")
	}

	// 解密
	cipher := utils.CalcMd5(p.sk)
	EncryptedDataStr, err := utils.Decrypt(data.Data, cipher, data.Nonce, constants.ModeCBC, constants.PKCS7Padding, constants.EncodingBase64)
	if err != nil {
		return fmt.Errorf("decrypt:%v", err)
	}

	key := data.Topic + "." + data.Operation
	log.WithContext(ctx).Debug("Event.Topic", zap.String("topic", data.Topic))
	log.WithContext(ctx).Debug("Event.Operation", zap.String("operation", data.Operation))
	log.WithContext(ctx).Debug("Event.Time", zap.Int64("time", data.Time))
	log.WithContext(ctx).Debug("Event.Nonce", zap.String("nonce", data.Nonce))
	log.WithContext(ctx).Debug("Event.Signature", zap.String("signature", data.Signature))
	log.WithContext(ctx).Debug("Event.EncryptedData", zap.String("encryptedData", EncryptedDataStr))

	// 获取处理方法
	fn, ok := p.handlerMap[key]
	if !ok {
		return errors.New("<<<<<<<< No Handler Exists for " + key)
	}

	// 调用处理方法
	return fn(ctx, data, EncryptedDataStr)
}

func (p *Process) VerifySignature(data *Entity) bool {
	message := fmt.Sprintf("%s:%s:%s:%d:%s", p.ak, data.Topic, data.Nonce, data.Time, data.Data)
	hmacSha256 := utils.CalcHMACSha256(message, p.sk)
	return hmacSha256 == data.Signature
}
