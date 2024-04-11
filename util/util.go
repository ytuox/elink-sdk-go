package util

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"
)

// StringDecoder 解码 JSON 字符串到目标结构体
func StringDecoder(data string, v interface{}) error {
	d := json.NewDecoder(strings.NewReader(data))
	d.UseNumber()
	return d.Decode(v)
}

// ByteDecoder 解码 JSON 字节切片到目标结构体
func ByteDecoder(data []byte, v interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	return d.Decode(v)
}

// ByteEncoder 编码目标结构体为 JSON 字节切片
func ByteEncoder(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// InterfaceDecoder 解码任意类型的数据到目标结构体
func InterfaceDecoder(data, v interface{}) error {
	byteData, err := ByteEncoder(data)
	if err != nil {
		return err
	}
	return ByteDecoder(byteData, v)
}

// StringEncoder 编码目标结构体为 JSON 字符串
func StringEncoder(data interface{}) string {
	byteData, err := ByteEncoder(data)
	if err != nil {
		return ""
	}
	return string(byteData)
}

// 读文件
func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

// 获取当前时间戳
func GetTimeStamp() int64 {
	now := time.Now().Local()
	return now.UnixNano() / 1e6
}
