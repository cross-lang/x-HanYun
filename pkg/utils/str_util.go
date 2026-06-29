package utils

import (
	"strings"
)

// IsEmpty 检查字符串是否为空
// str: 要检查的字符串
// 返回值: 如果字符串为空（长度为 0），则返回 true；否则返回 false
func IsEmpty(str string) bool {
	if str == "" || len(str) <= 0 {
		return true
	}
	return false
}

// IsNotEmpty 检查字符串是否不为空
// str: 要检查的字符串
// 返回值: 如果字符串不为空（长度大于 0），则返回 true；否则返回 false
func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

// IsBlank 检查字符串是否为空或仅包含空白字符
// str: 要检查的字符串
// 返回值: 如果字符串为空或仅包含空白字符，则返回 true；否则返回 false
func IsBlank(str string) bool {
	str = strings.Trim(str, " ")
	return IsEmpty(str)
}

// IsNotBlank 检查字符串是否不为空且不只是包含空白字符
// str: 要检查的字符串
// 返回值: 如果字符串不为空且不只是包含空白字符，则返回 true；否则返回 false
func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

// RemovePrefix 移除字符串的指定前缀
// url: 要处理的字符串
// prefix: 要移除的前缀
// 返回值: 如果字符串以指定前缀开头，则返回移除前缀后的字符串；否则返回原字符串
func RemovePrefix(url, prefix string) string {
	if strings.HasPrefix(url, prefix) {
		return url[len(prefix):]
	}
	return url
}

// IfString 根据给定的条件返回两个字符串中的一个
// condition: 用于判断的布尔条件
// trueVal: 当条件为 true 时返回的字符串
// falseVal: 当条件为 false 时返回的字符串
// 返回值: 根据条件返回 trueVal 或 falseVal
func IfString(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}
