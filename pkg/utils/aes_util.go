package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"x-HanYun/pkg/core/log"

	"go.uber.org/zap"
)

// 加密模式常量
const (
	ModeCBC = "CBC"
	ModeECB = "ECB"
	ModeGCM = "GCM"
)

// 填充模式常量
const (
	PKCS7Padding = "PKCS7"
	ZeroPadding  = "Zero"
)

// 编码模式常量
const (
	EncodingBase64 = "Base64"
	EncodingHex    = "Hex"
)

const (
	blockSize    = aes.BlockSize
	gcmNonceSize = 12 // GCM推荐使用12字节的nonce
)

// ConvertSKToKey 将 SK 转换为 key
func ConvertSKToKey(sk string) ([]byte, error) {
	// 验证密钥长度是否正确（32个字符，对应16字节）
	if len(sk) != 32 {
		return nil, fmt.Errorf("<<<<<<<< Invalid key string length: %d, expected 32", len(sk))
	}

	// 将十六进制字符串转换为字节数组
	keyBytes, err := hex.DecodeString(sk)
	if err != nil {
		return nil, fmt.Errorf("<<<<<<<< Failed to decode key: %v", err)
	}

	// AES-128需要16字节密钥
	if len(keyBytes) != 16 {
		return nil, fmt.Errorf("<<<<<<<< Invalid key bytes length: %d, expected 16", len(keyBytes))
	}

	return keyBytes, nil
}

func SimpleEncrypt(plaintext string, sk string) (string, error) {
	key, err7 := ConvertSKToKey(sk)
	if err7 != nil {
		return "", fmt.Errorf("<<<<<<<< ConvertSKToKey err: %v", err7)
	}
	keyStr := string(key)
	log.Info(">>>>>>>> key", zap.String("value", keyStr))

	// 生成指定长度的初始化向量（IV）
	iv := make([]byte, 16)
	if _, err3 := rand.Read(iv); err3 != nil {
		return "", fmt.Errorf("<<<<<<<< GenerateIV err: %v", err3)
	}
	ivStr := string(iv)
	log.Info(">>>>>>>> iv", zap.String("value", ivStr))

	// AES + CBC 加密, PKCS7方式填充，Base64编码
	ePlaintext, err2 := Encrypt(plaintext, keyStr, ivStr, ModeCBC, PKCS7Padding, EncodingBase64)
	if err2 != nil {
		log.Error("<<<<<<< To Encrypt err", zap.Error(err2))
		return "", err2
	}

	return ePlaintext, nil
}

func SimpleDecrypt(plaintext string, key string, iv string) (string, error) {
	log.Info(">>>>>>>> key", zap.String("value", key))
	log.Info(">>>>>>>> iv", zap.String("value", iv))

	// AES + CBC 加密, PKCS7方式填充，Base64编码
	dPlaintext, err2 := Decrypt(plaintext, key, iv, ModeCBC, PKCS7Padding, EncodingBase64)
	if err2 != nil {
		log.Error("<<<<<<< To Decrypt err", zap.Error(err2))
		return "", err2
	}

	return dPlaintext, nil
}

// Encrypt 加密函数，支持CBC、ECB和GCM模式
func Encrypt(plaintext, key, iv string, mode, padding, encoding string) (string, error) {
	// 设置默认参数
	if mode == "" {
		mode = ModeCBC
	}
	if padding == "" && mode != ModeGCM {
		padding = PKCS7Padding
	}
	if encoding == "" {
		encoding = EncodingBase64
	}

	// 验证参数
	if mode != ModeCBC && mode != ModeECB && mode != ModeGCM {
		return "", fmt.Errorf("<<<<<<<< only CBC, ECB and GCM modes are supported, unsupported mode:%s", mode)
	}

	// GCM模式不需要填充
	if mode == ModeGCM && padding != "" {
		log.Error("<<<<<<<< GCM mode does not require padding, ignoring padding", zap.String("padding", padding))
		padding = ""
	} else if padding != PKCS7Padding && padding != ZeroPadding && padding != "" {
		return "", fmt.Errorf("<<<<<<<< only PKCS7 and ZeroPadding are supported, unsupported padding:%s", padding)
	}

	if encoding != EncodingBase64 && encoding != EncodingHex && encoding != "" {
		return "", fmt.Errorf("<<<<<<<< only base64 and hex encodings are supported, unsupported encoding:%s", encoding)
	}

	// 验证密钥长度
	keyBytes := []byte(key)
	keyLen := len(keyBytes)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return "", fmt.Errorf("<<<<<<<< key length must be 16, 24, or 32 bytes, current length:%d", keyLen)
	}

	// 验证IV长度
	var ivBytes []byte
	if mode == ModeGCM {
		// GCM模式需要12字节IV，如果未提供则生成随机IV
		if iv == "" {
			ivBytes = make([]byte, gcmNonceSize)
			if _, err := io.ReadFull(rand.Reader, ivBytes); err != nil {
				return "", err
			}
		} else {
			ivBytes = []byte(iv)
			if len(ivBytes) != gcmNonceSize {
				return "", fmt.Errorf("<<<<<<<< GCM mode requires IV length of %d bytes, current length:%d", gcmNonceSize, len(ivBytes))
			}
		}
	} else if mode == ModeCBC {
		// CBC模式需要16字节IV
		ivBytes = []byte(iv)
		if len(ivBytes) != blockSize {
			return "", fmt.Errorf("<<<<<<<< CBC mode requires IV length of %d bytes, current length:%d", blockSize, len(ivBytes))
		}
	}
	// ECB模式不需要IV

	// 创建加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)
	var ciphertext []byte

	// 根据模式进行加密
	switch mode {
	case ModeGCM:
		// GCM模式加密
		aesGCM, err5 := cipher.NewGCM(block)
		if err5 != nil {
			return "", err5
		}
		// GCM模式不需要填充，直接加密
		ciphertext = aesGCM.Seal(nil, ivBytes, plaintextBytes, nil)
		// 拼接IV和密文（包含认证标签）
		ciphertext = bytes.Join([][]byte{ivBytes, ciphertext}, nil)

	case ModeCBC:
		// CBC模式加密，需要填充
		switch padding {
		case PKCS7Padding:
			plaintextBytes = pkcs7Padding(plaintextBytes, blockSize)
		case ZeroPadding:
			plaintextBytes = zeroPadding(plaintextBytes, blockSize)
		}

		ciphertext = make([]byte, len(plaintextBytes))
		modeCBC := cipher.NewCBCEncrypter(block, ivBytes)
		modeCBC.CryptBlocks(ciphertext, plaintextBytes)

	case ModeECB:
		// ECB模式加密，需要填充
		switch padding {
		case PKCS7Padding:
			plaintextBytes = pkcs7Padding(plaintextBytes, blockSize)
		case ZeroPadding:
			plaintextBytes = zeroPadding(plaintextBytes, blockSize)
		}

		ciphertext = make([]byte, len(plaintextBytes))
		bs := block.BlockSize()
		// 分块加密
		for start := 0; start < len(plaintextBytes); start += bs {
			end := start + bs
			if end > len(plaintextBytes) {
				end = len(plaintextBytes)
			}
			block.Encrypt(ciphertext[start:end], plaintextBytes[start:end])
		}
	}

	// 编码
	var encoded string
	switch encoding {
	case EncodingBase64, "":
		encoded = base64.StdEncoding.EncodeToString(ciphertext)
	case EncodingHex:
		encoded = hex.EncodeToString(ciphertext)
	}

	return encoded, nil
}

// Decrypt 解密函数，支持CBC、ECB和GCM模式
func Decrypt(ciphertext, key, iv string, mode, padding, encoding string) (string, error) {
	// 设置默认参数
	if mode == "" {
		mode = ModeCBC
	}
	if padding == "" && mode != ModeGCM {
		padding = PKCS7Padding
	}
	if encoding == "" {
		encoding = EncodingBase64
	}

	// 验证参数
	if mode != ModeCBC && mode != ModeECB && mode != ModeGCM {
		return "", fmt.Errorf("<<<<<<<< only CBC, ECB and GCM modes are supported, unsupported mode:%s", mode)
	}

	// GCM模式不需要填充
	if mode == ModeGCM && padding != "" {
		padding = ""
	} else if padding != PKCS7Padding && padding != ZeroPadding && padding != "" {
		return "", fmt.Errorf("<<<<<<<< only PKCS7 and ZeroPadding are supported, unsupported padding:%s", padding)
	}

	if encoding != EncodingBase64 && encoding != EncodingHex && encoding != "" {
		return "", fmt.Errorf("<<<<<<<< only base64 and hex encodings are supported, unsupported encoding:%s", encoding)
	}

	// 验证密钥长度
	keyBytes := []byte(key)
	keyLen := len(keyBytes)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return "", fmt.Errorf("<<<<<<<< key length must be 16, 24, or 32 bytes, current length:%d", keyLen)
	}

	// 解码
	var ciphertextBytes []byte
	var err error
	switch encoding {
	case EncodingBase64, "":
		ciphertextBytes, err = base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			return "", err
		}
	case EncodingHex:
		ciphertextBytes, err = hex.DecodeString(ciphertext)
		if err != nil {
			return "", err
		}
	}

	// 创建解密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	var plaintextBytes []byte
	var ivBytes []byte

	// 根据不同模式进行解密
	switch mode {
	case ModeGCM:
		// GCM模式解密
		// 分离IV和密文（包含认证标签）
		if len(ciphertextBytes) < gcmNonceSize {
			return "", fmt.Errorf("<<<<<<<< ciphertext too short for GCM mode")
		}

		ivBytes = ciphertextBytes[:gcmNonceSize]
		ciphertextBytes = ciphertextBytes[gcmNonceSize:]

		aesGCM, err := cipher.NewGCM(block)
		if err != nil {
			return "", err
		}

		// GCM解密同时验证认证标签
		plaintextBytes, err = aesGCM.Open(nil, ivBytes, ciphertextBytes, nil)
		if err != nil {
			return "", fmt.Errorf("<<<<<<<< GCM decryption failed:%v", err)
		}

	case ModeCBC:
		// 验证密文长度
		if len(ciphertextBytes)%blockSize != 0 {
			return "", fmt.Errorf("<<<<<<<< ciphertext length must be a multiple of block size (%d bytes) for CBC mode, current length:%d", blockSize, len(ciphertextBytes))
		}

		ivBytes = []byte(iv)
		if len(ivBytes) != blockSize {
			return "", fmt.Errorf("<<<<<<<< CBC mode requires IV length of %d bytes, current length:%d", blockSize, len(ivBytes))
		}

		plaintextBytes = make([]byte, len(ciphertextBytes))
		modeCBC := cipher.NewCBCDecrypter(block, ivBytes)
		modeCBC.CryptBlocks(plaintextBytes, ciphertextBytes)

		// 去除填充
		switch padding {
		case PKCS7Padding:
			plaintextBytes = pkcs7UnPadding(plaintextBytes)
		case ZeroPadding:
			plaintextBytes = zeroUnPadding(plaintextBytes)
		}

	case ModeECB:
		// 验证密文长度
		if len(ciphertextBytes)%blockSize != 0 {
			return "", fmt.Errorf("<<<<<<<< ciphertext length must be a multiple of block size (%d bytes) for ECB mode, current length:%d", blockSize, len(ciphertextBytes))
		}

		plaintextBytes = make([]byte, len(ciphertextBytes))
		bs := block.BlockSize()
		for start := 0; start < len(ciphertextBytes); start += bs {
			block.Decrypt(plaintextBytes[start:start+bs], ciphertextBytes[start:start+bs])
		}

		// 去除填充
		switch padding {
		case PKCS7Padding:
			plaintextBytes = pkcs7UnPadding(plaintextBytes)
		case ZeroPadding:
			plaintextBytes = zeroUnPadding(plaintextBytes)
		}
	}

	return string(plaintextBytes), nil
}

// pkcs7Padding PKCS#7填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7UnPadding PKCS#7去填充
func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return nil
	}
	unpadding := int(data[length-1])

	// 验证填充字节是否符合PKCS#7规范
	if unpadding < 1 || unpadding > blockSize {
		return nil
	}

	for i := length - unpadding; i < length; i++ {
		if data[i] != byte(unpadding) {
			return nil
		}
	}

	return data[:(length - unpadding)]
}

// zeroPadding 零填充
func zeroPadding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	if padding == 0 {
		return data
	}
	return append(data, bytes.Repeat([]byte{0}, padding)...)
}

// zeroUnPadding 零去填充
func zeroUnPadding(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}
	for len(data) > 0 && data[len(data)-1] == 0 {
		data = data[:len(data)-1]
	}
	return data
}
