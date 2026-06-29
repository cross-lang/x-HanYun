package utils

import (
	"encoding/json"

	"x-HanYun/pkg/core/log"

	"go.uber.org/zap"
)

// ToString 将一个对象序列化为 JSON 字符串
// obj: 要进行序列化的对象，类型为 interface{}，表示可以传入任意类型的对象
// 返回值: 序列化后的 JSON 字符串和可能出现的错误
// 若序列化失败，会使用 log 记录错误信息，包括错误详情和传入的对象
func ToString(obj interface{}) (string, error) {
	str, err := json.Marshal(obj)
	if err != nil {
		log.Error("<<<<<<<< Failed to Marshal json obj", zap.Error(err), zap.Any("obj", obj))
		return "", err
	}
	return string(str), nil
}

// ToByteSlice 将一个对象序列化为 JSON 字节切片
// obj: 要进行序列化的对象，类型为 interface{}，可以传入任意类型的对象
// 返回值: 序列化后的 JSON 字节切片和可能出现的错误
// 若序列化失败，会使用 logx 记录错误信息，包括错误详情和传入的对象
func ToByteSlice(obj interface{}) ([]byte, error) {
	byteSlice, err := json.Marshal(obj)
	if err != nil {
		log.Error("<<<<<<<< Failed to Marshal json obj", zap.Error(err), zap.Any("obj", obj))
		return []byte(nil), err
	}
	return byteSlice, nil
}

// ToObject 将 JSON 字符串反序列化为指定对象
// str: 要进行反序列化的 JSON 字符串
// obj: 用于存储反序列化结果的对象指针，类型为 interface{}，调用时需传入对象的指针
// 返回值: 可能出现的错误
// 若反序列化失败，会使用 logx 记录错误信息，包括错误详情和传入的 JSON 字符串
func ToObject(str string, obj interface{}) error {
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		log.Error("<<<<<<<< Failed to Unmarshal json str", zap.Error(err), zap.String("str", str))
		return err
	}
	return nil
}
