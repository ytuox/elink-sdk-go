package serial

import (
	"errors"
	"sync"
)

// SerialManager 管理串口的打开和关闭
type SerialManager struct {
	sync.Mutex  // 互斥锁，用于保证唯一性
	openedPorts map[string]interface{}
}

// 创建一个新的SerialManager实例
func NewSerialManager() *SerialManager {
	return &SerialManager{
		openedPorts: make(map[string]interface{}),
	}
}

// AddPort 添加一个打开的串口
func (sm *SerialManager) AddPort(address string, port interface{}) error {
	sm.Lock()
	defer sm.Unlock()

	if _, exists := sm.openedPorts[address]; exists {
		return errors.New("serial port already exists")
	}

	sm.openedPorts[address] = port
	return nil
}

// GetPort 获取已经打开的串口
func (sm *SerialManager) GetPort(address string) (interface{}, bool) {
	sm.Lock()
	defer sm.Unlock()

	port, ok := sm.openedPorts[address]
	return port, ok
}

// ClosePort 关闭一个已经打开的串口
func (sm *SerialManager) DelPort(address string) error {
	sm.Lock()
	defer sm.Unlock()

	delete(sm.openedPorts, address)
	return nil
}
