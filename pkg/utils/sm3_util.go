package utils

import (
	"crypto/rand"
	"encoding/hex"

	"x-HanYun/pkg/core/log"

	"github.com/tjfoc/gmsm/sm3"

	"go.uber.org/zap"
)

// SM3Hash 使用SM3算法计算数据的哈希值
func SM3Hash(data []byte) string {
	hash := sm3.Sm3Sum(data)
	return hex.EncodeToString(hash)
}

// SM3GenerateSalt 生成指定长度的随机盐值
func SM3GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		log.Error("<<<<<<<< Failed to generate salt", zap.Error(err))
		return "", err
	}
	log.Info(">>>>>>>> Generated salt", zap.String("value", hex.EncodeToString(salt)))
	return hex.EncodeToString(salt), nil
}

// SM3HashWithSalt 使用SM3算法计算数据和盐值的哈希值
func SM3HashWithSalt(data []byte, salt string) string {
	saltBytes, _ := hex.DecodeString(salt)
	combined := append(data, saltBytes...)
	hash := sm3.Sm3Sum(combined)
	return hex.EncodeToString(hash)
}
