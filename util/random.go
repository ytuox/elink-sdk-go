package util

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

// 获取随机数 纯数字
func Random(n int) string {
	charset := []byte("0123456789")
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	result := make([]byte, n)
	for i := range result {
		idx := r.Intn(len(charset))
		result[i] = charset[idx]
	}
	return string(result)
}

// 获取随机数 纯文字
func GetRandomString(n int) string {
	charset := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	for i := range result {
		idx := r.Intn(len(charset))
		result[i] = charset[idx]
	}
	return string(result)
}

// 获取随机数 数字和文字
func GetRandomBoth(n int) string {
	charset := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	result := make([]byte, n)
	for i := range result {
		idx := r.Intn(len(charset))
		result[i] = charset[idx]
	}
	return string(result)
}

// 获取随机数 纯数字
func RandomNum() string {
	return fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
}

func GetRandomBase64(n int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	buf := make([]byte, n)
	r.Read(buf)
	return base64.StdEncoding.EncodeToString(buf)
}
