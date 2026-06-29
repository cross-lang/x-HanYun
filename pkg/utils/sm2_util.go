package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"x-HanYun/pkg/core/log"

	"github.com/tjfoc/gmsm/sm2"

	"go.uber.org/zap"
)

// GenerateKeyPair 生成随机秘钥对
func GenerateKeyPair() (privateKeyHex, publicKeyHex string, err error) {
	privateKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	// 私钥转换为16进制字符串
	privateKeyHex = hex.EncodeToString(privateKey.D.Bytes())

	// 公钥转换为04开头的未压缩格式
	publicKeyBytes := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	publicKeyHex = "04" + hex.EncodeToString(publicKeyBytes)

	log.Info(">>>>>>>> 公钥", zap.String("value", publicKeyHex))
	log.Info(">>>>>>>> 私钥", zap.String("value", privateKeyHex))

	return privateKeyHex, publicKeyHex, nil
}

// SM2Encrypt 数据加密,使用C1C2C3
func SM2Encrypt(publicKeyHex string, data []byte) (string, error) {
	if publicKeyHex == "" {
		return "", fmt.Errorf("<<<<<<<< public key is empty")
	}

	if len(data) == 0 {
		return "", fmt.Errorf("<<<<<<<< data is empty")
	}

	// 解析公钥
	publicKey, err := HexToPublicKey(publicKeyHex)
	if err != nil {
		return "", err
	}

	// 使用SM2加密
	cipherText, err := sm2.Encrypt(publicKey, data, rand.Reader, sm2.C1C2C3)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(cipherText), nil
}

// SM2Decrypt 数据解密，使用C1C2C3
func SM2Decrypt(privateKeyHex string, encryptedDataHex string) ([]byte, error) {
	if privateKeyHex == "" {
		return nil, fmt.Errorf("<<<<<<<< private key is empty")
	}

	if encryptedDataHex == "" {
		return nil, fmt.Errorf("<<<<<<<< encrypted data is empty")
	}

	// 解析私钥
	privateKey, err := HexToPrivateKey(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// 将加密数据从16进制转换为字节
	encryptedData, err := hex.DecodeString(encryptedDataHex)
	if err != nil {
		return nil, err
	}

	// 使用SM2解密
	plainText, err := sm2.Decrypt(privateKey, encryptedData, sm2.C1C2C3)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// HexToPublicKey 16进制字符串转换为公钥
func HexToPublicKey(publicKeyHex string) (*sm2.PublicKey, error) {
	if len(publicKeyHex) < 2 {
		return nil, fmt.Errorf("<<<<<<<< invalid public key length")
	}

	// 移除可能的04前缀（未压缩格式标识）
	if publicKeyHex[:2] == "04" {
		publicKeyHex = publicKeyHex[2:]
	}

	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, err
	}

	if len(publicKeyBytes) != 64 {
		return nil, fmt.Errorf("<<<<<<<< invalid public key length, expected 64 bytes, got %d", len(publicKeyBytes))
	}

	// 创建公钥
	x := new(big.Int).SetBytes(publicKeyBytes[:32])
	y := new(big.Int).SetBytes(publicKeyBytes[32:])

	publicKey := &sm2.PublicKey{
		Curve: sm2.P256Sm2(),
		X:     x,
		Y:     y,
	}

	return publicKey, nil
}

// HexToPrivateKey 16进制字符串转换为私钥
func HexToPrivateKey(privateKeyHex string) (*sm2.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	d := new(big.Int).SetBytes(privateKeyBytes)

	// 生成对应的公钥
	x, y := sm2.P256Sm2().ScalarBaseMult(d.Bytes())

	privateKey := &sm2.PrivateKey{
		PublicKey: sm2.PublicKey{
			Curve: sm2.P256Sm2(),
			X:     x,
			Y:     y,
		},
		D: d,
	}

	return privateKey, nil
}
