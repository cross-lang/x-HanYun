package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// CalcMd5 ①对字符串进行 MD5 哈希计算；②将结果转换为十六进制字符串
// s: 待计算哈希的字符串
// 返回值: MD5 哈希后的十六进制字符串
func CalcMd5(s string) string {
	// ①对字符串进行 MD5 哈希计算
	h := md5.New()
	h.Write([]byte(s))

	// ②将结果转换为十六进制字符串
	return hex.EncodeToString(h.Sum(nil))
}

// CalcHMACSha256 ①对字符串进行 HMAC-SHA256 哈希计算；②将结果进行 Base64 编码
// message: 待签名的消息
// secret: 签名使用的密钥
// 返回值: 签名并编码后的字符串
func CalcHMACSha256(message string, secret string) string {
	// ①对字符串进行  HMAC-SHA256 哈希计算
	if len(secret) == 0 {
		return ""
	}
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))

	// ②将结果进行 Base64 编码
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
