package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/golang/snappy"
)

func HmacSha256(secret []byte, data string) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha1(secret []byte, data string) string {
	h := hmac.New(sha1.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padLen := blockSize - len(src)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padding...)
}

func AesCbcBase64(src, productSecret string) (string, error) {
	if src == "" || productSecret == "" {
		return "", fmt.Errorf("加密参数错误")
	}
	
	// 截取 productSecret 前 16 位作为密钥
	key := []byte(productSecret)[:16]
	// 以长度 16 的字符 "0" 作为偏移量
	iv := bytes.Repeat([]byte("0"), 16)

	data := []byte(src)

	// 使用 AES-CBC 加密
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 对补全后的数据进行加密
	blockSize := block.BlockSize()
	data = PKCS5Padding(data, blockSize)
	cryptData := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cryptData, data)

	// 进行 base64 编码
	return base64.StdEncoding.EncodeToString(cryptData), nil
}

// 使用 AES-GCM 进行加密
func EncryptAESGCM(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 生成随机的 nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

// 使用 AES-GCM 进行解密
func DecryptAESGCM(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 提取 nonce
	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("invalid ciphertext")
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// 使用 Snappy 进行压缩
func CompressSnappy(data []byte) []byte {
	return snappy.Encode(nil, data)
}

// 使用 Snappy 进行解压缩
func DecompressSnappy(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
