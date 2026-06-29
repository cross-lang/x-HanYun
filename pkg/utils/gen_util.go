package utils

import (
	cryptoRand "crypto/rand"
	"errors"
	"math/big"
	mathRand "math/rand"
	"time"
)

// GenerateSalt 生成指定长度的随机盐值
func GenerateSalt(length int) string {
	if length <= 0 {
		return ""
	}

	// 初始化随机数生成器
	mathRand.Seed(time.Now().UnixNano())

	// 定义盐值字符集
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"

	// 生成随机盐值
	salt := make([]byte, length)
	for i := range salt {
		salt[i] = charset[mathRand.Intn(len(charset))]
	}

	return string(salt)
}

// GeneratePassword 生成指定长度的安全随机密码
func GeneratePassword(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("密码长度必须大于0")
	}

	// 定义密码字符集（包含大小写字母、数字和特殊字符）
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	charsetLen := big.NewInt(int64(len(charset)))

	// 生成随机密码
	password := make([]byte, length)
	for i := range password {
		// 使用crypto/rand生成安全随机数
		idx, err := cryptoRand.Int(cryptoRand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		password[i] = charset[idx.Int64()]
	}

	return string(password), nil
}

// GenerateIV 生成指定长度的初始化向量（IV）
// 用于对称加密算法（如AES-CBC、AES-GCM）
func GenerateIV(length int) ([]byte, error) {
	if length <= 0 {
		return nil, errors.New("IV长度必须大于0")
	}

	iv := make([]byte, length)
	_, err := cryptoRand.Read(iv)
	if err != nil {
		return nil, err
	}

	return iv, nil
}

// GenerateNonce 生成指定长度的临时随机数（nonce）
// 用于一次性加密、消息认证或防止重放攻击
func GenerateNonce(length int) ([]byte, error) {
	if length <= 0 {
		return nil, errors.New("nonce长度必须大于0")
	}

	nonce := make([]byte, length)
	_, err := cryptoRand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}
