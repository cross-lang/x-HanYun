package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"

	"x-HanYun/pkg/core/log"

	"go.uber.org/zap"
)

// RSA 填充模式常量
const (
	PKCS1v15 = "PKCS1v15"
	OAEP     = "OAEP"
)

var RSAPrivateKey *rsa.PrivateKey

func SetRSAPrivateKey(rsaPriKeyFile string) error {
	file, err := os.Open(rsaPriKeyFile)
	if err != nil {
		log.Error("<<<<<<<< Failed to open RSA private key file", zap.Error(err))
		return err
	}
	defer file.Close()

	rsaPriKeyData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error("<<<<<<<< Failed to read RSA private key file", zap.Error(err))
		return err
	}

	block, _ := pem.Decode(rsaPriKeyData)
	if block == nil {
		return errors.New("Failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	RSAPrivateKey = privateKey

	return nil
}

// RSAEncrypt 使用指定填充方式进行RSA加密
func RSAEncrypt(plaintext []byte, publicKey *rsa.PublicKey, padding string) ([]byte, error) {
	switch padding {
	case PKCS1v15:
		return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
	case OAEP:
		hash := crypto.SHA256
		return rsa.EncryptOAEP(hash.New(), rand.Reader, publicKey, plaintext, nil)
	default:
		return nil, errors.New("<<<<<<<< unsupported padding mode")
	}
}

// RSADecrypt 使用指定填充方式进行RSA解密
func RSADecrypt(ciphertext []byte, privateKey *rsa.PrivateKey, padding string) ([]byte, error) {
	switch padding {
	case PKCS1v15:
		return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	case OAEP:
		hash := crypto.SHA256
		return rsa.DecryptOAEP(hash.New(), rand.Reader, privateKey, ciphertext, nil)
	default:
		return nil, errors.New("<<<<<<<< unsupported padding mode")
	}
}
