package serialport

import (
	"errors"
	"time"
)

func calculateTimeout(baudRate int, dataBits int, stopBits int, parity string, packetSize int) (time.Duration, error) {
	// 参数有效性检查
	if baudRate <= 0 || dataBits < 5 || dataBits > 8 || stopBits < 1 || stopBits > 2 || (parity != "N" && parity != "E" && parity != "O") || packetSize <= 0 {
		return 0, errors.New("invalid input parameters")
	}

	// 计算每个字符的传输时间
	byteTime := 10 * time.Second / time.Duration(baudRate)
	charTime := (time.Duration(dataBits) + 1) * byteTime

	// 根据校验位类型增加传输时间
	if parity != "N" {
		charTime += byteTime
	}

	// 加上停止位的传输时间
	stopBitTime := time.Duration(stopBits) * byteTime
	charTime += stopBitTime

	// 计算总的传输时间
	totalTime := charTime * time.Duration(packetSize)

	// 处理可能的溢出情况
	if totalTime < 0 {
		return 0, errors.New("timeout value exceeds maximum")
	}

	return totalTime, nil
}
