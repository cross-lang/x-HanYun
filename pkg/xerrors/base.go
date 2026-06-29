package xerrors

import (
	"encoding/json"
)

// Error 定义错误接口
type Error interface {
	error
	Code() ErrorCode
	Message() string
	DebugInfo() string
	WithFields(fields ...Field) Error
}

// Field 定义错误详情字段
type Field struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// codeError 实现Error接口
type codeError struct {
	code  ErrorCode
	msg   string
	debug string
}

func (e *codeError) Error() string {
	return e.msg
}

func (e *codeError) Code() ErrorCode {
	return e.code
}

func (e *codeError) Message() string {
	return e.msg
}

func (e *codeError) DebugInfo() string {
	return e.debug
}

func (e *codeError) WithFields(fields ...Field) Error {
	if len(fields) == 0 {
		return e
	}

	debugFields := make(map[string]interface{})
	if e.debug != "" {
		_ = json.Unmarshal([]byte(e.debug), &debugFields)
	}

	for _, field := range fields {
		debugFields[field.Key] = field.Value
	}

	debugJSON, _ := json.Marshal(debugFields)
	return &codeError{
		code:  e.code,
		msg:   e.msg,
		debug: string(debugJSON),
	}
}

// Any 创建任意类型的错误字段
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// String 创建字符串类型的错误字段
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// NewError 创建带错误码的错误
func NewError(code ErrorCode, message string) Error {
	return &codeError{
		code: code,
		msg:  message,
	}
}
