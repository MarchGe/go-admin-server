package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// EncryptString 加密。
// key的长度16、24、32字节，分别对应AES-128、AES-192、AES-256加密算法，
// data是需要加密的数据，
// additionalData是用于进行身份验证的数据，按需使用，可以不用传。
func EncryptString(key, data, additionalData string) (string, error) {
	encryptBytes, err := Encrypt([]byte(key), []byte(data), []byte(additionalData))
	if err != nil {
		return "", fmt.Errorf("encrypt data error, %w", err)
	}
	encodeStr := base64.RawURLEncoding.EncodeToString(encryptBytes)
	return encodeStr, nil
}

// Encrypt 加密。
// key的长度16、24、32字节，分别对应AES-128、AES-192、AES-256加密算法，
// data是需要加密的数据，
// additionalData是用于进行身份验证的数据，按需使用，可以不用传。
func Encrypt(key, data, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher error, %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new GCM error, %w", err)
	}
	nonce := RandomBytes(gcm.NonceSize())
	bytes := gcm.Seal(nil, nonce, data, additionalData)
	return append(nonce, bytes...), nil
}

// DecryptString 解密。
// key的长度16、24、32字节，分别对应AES-128、AES-192、AES-256加密算法，
// data是需要解密的数据，
// additionalData是用于进行身份验证的数据，按需使用，可以不用传。
func DecryptString(key, data, additionalData string) (string, error) {
	decodeBytes, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("cipher text base64 decode error, %w", err)
	}
	decryptBytes, err := Decrypt([]byte(key), decodeBytes, []byte(additionalData))
	if err != nil {
		return "", err
	}
	return string(decryptBytes), nil
}

// Decrypt 解密。
// key的长度16、24、32字节，分别对应AES-128、AES-192、AES-256加密算法，
// data是需要解密的数据，
// additionalData是用于进行身份验证的数据，按需使用，可以不用传。
func Decrypt(key, data, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher error, %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new GCM error, %w", err)
	}
	nonceSize := gcm.NonceSize()
	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	bytes, err := gcm.Open(nil, nonce, cipherData, additionalData)
	if err != nil {
		return nil, fmt.Errorf("decrypt data error, %w", err)
	}
	return bytes, nil
}
