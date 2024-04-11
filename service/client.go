/*******************************************************************************
 * Copyright 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package service

import (
	"github.com/ytuox/elink-sdk-go/common"
	"github.com/ytuox/elink-sdk-go/interfaces"
	"github.com/ytuox/elink-sdk-go/internal/logger"
	"github.com/ytuox/elink-sdk-go/model"
)

// Start 启动驱动
func (d *PluginService) Start(plugin interfaces.Plugin) error {
	return d.start(plugin)
}

// Stop 停止驱动
func (d *PluginService) Stop() error {
	return d.stop()
}

// GetLogger 获取日志接口
func (d *PluginService) GetLogger() logger.Logger {
	return d.logger
}

// GetCustomParam 获取自定义参数
func (d *PluginService) GetPluginParam() string {
	return d.cfg.PluginParam
}

// Online 设备与平台建立连接
func (d *PluginService) Online(deviceId string) error {
	return d.connectIotPlatform(deviceId)
}

// Offline 设备与平台断开连接
func (d *PluginService) Offline(deviceId string) error {
	return d.disconnectIotPlatform(deviceId)
}

// GetConnectStatus 获取设备连接状态
func (d *PluginService) GetConnectStatus(deviceId string) (common.DeviceConnectStatus, error) {
	return d.getConnectStatus(deviceId)
}

// CreateDevice 创建设备
func (d *PluginService) CreateDevice(device model.AddDevice) (model.Device, error) {
	return d.createDevice(device)
}

// GetDeviceList 获取所有的设备
func (d *PluginService) GetDeviceList() []model.Device {
	return d.getDeviceList()
}

// GetDeviceList 获取所有的设备
func (d *PluginService) GetDeviceListByAppId(id string) []model.Device {
	return d.getDeviceListByAppId(id)
}

// GetDeviceList 获取所有的设备
func (d *PluginService) GetDeviceListByUserId(id string) ([]model.Device, error) {
	return d.getDeviceListByUserId(id)
}

// GetDeviceById 通过设备id获取设备详情
func (d *PluginService) GetDeviceById(deviceId string) (model.Device, bool) {
	return d.getDeviceById(deviceId)
}

// GetDeviceById 通过设备id获取设备属性影子
func (d *PluginService) GetDevicePropertyShadow(deviceId, identifier string) ([]model.PropertyShadowData, error) {
	return d.getDevicePropertyShadow(deviceId, identifier)
}

// GetDeviceById 通过设备id获取设备服务影子
func (d *PluginService) GetDeviceServiceShadow(deviceId, identifier string) ([]model.ServiceShadowData, error) {
	return d.getDeviceServiceShadow(deviceId, identifier)
}

// ProductList 获取当前实例下的所有产品
func (d *PluginService) ProductList() map[string]model.Product {
	return d.productCache.All()
}

// GetProductById 根据产品id获取产品信息
func (d *PluginService) GetProductById(productId string) (model.Product, bool) {
	return d.productCache.SearchById(productId)
}

// GetProductProperties 根据产品id获取产品所有属性信息
func (d *PluginService) GetProductProperties(productId string) (map[string]model.Property, bool) {
	return d.getProductProperties(productId)
}

// GetProductPropertyByIdentifier 根据产品id与code获取属性信息
func (d *PluginService) GetProductPropertyByCode(productId, code string) (model.Property, bool) {
	return d.getProductPropertyByIdentifier(productId, code)
}

// GetProductEvents 根据产品id获取产品所有事件信息
func (d *PluginService) GetProductEvents(productId string) (map[string]model.Event, bool) {
	return d.getProductEvents(productId)
}

// GetProductEventByIdentifier 根据产品id与code获取事件信息
func (d *PluginService) GetProductEventByCode(productId, identifier string) (model.Event, bool) {
	return d.getProductEventByIdentifier(productId, identifier)
}

// GetProductServices 根据产品id获取产品所有服务信息
func (d *PluginService) GetProductServices(productId string) (map[string]model.Service, bool) {
	return d.getProductServices(productId)
}

// GetProductServiceByIdentifier 根据产品id与code获取服务信息
func (d *PluginService) GetProductServiceByIdentifier(productId, identifier string) (model.Service, bool) {
	return d.getProductServiceByIdentifier(productId, identifier)
}

// PropertyReport 物模型属性上报 如果data参数中的Sys.Ack设置为1，则该方法会同步阻塞等待云端返回结果。
func (d *PluginService) PropertyReport(deviceId string, data model.PropertyReport) (model.CommonResponse, error) {
	return d.propertyReport(deviceId, data)
}

// EventReport 物模型事件上报
func (d *PluginService) EventReport(deviceId string, data model.EventReport) (model.CommonResponse, error) {
	return d.eventReport(deviceId, data)
}

// BatchReport 设备批量上报属性和事件 如果data参数中的Sys.Ack设置为1，则该方法会同步阻塞等待云端返回结果。
// 如非必要，不建议设置Sys.Ack
// 废弃
// func (d *PluginService) BatchReport(deviceId string, data model.BatchReport) (model.CommonResponse, error) {
// 	return d.batchReport(deviceId, data)
// }

// PropertyDesiredGet 设备拉取属性期望值 如果data参数中的Sys.Ack设置为1，则该方法会同步阻塞等待云端返回结果。
// 废弃
//func (d *PluginService) PropertyDesiredGet(deviceId string, data model.PropertyDesiredGet) (model.PropertyDesiredGetResponse, error) {
//	return d.propertyDesiredGet(deviceId, data)
//}

// PropertyDesiredDelete 设备删除属性期望值 如果data参数中的Sys.Ack设置为1，则该方法会同步阻塞等待云端返回结果。
// 废弃
//func (d *PluginService) PropertyDesiredDelete(deviceId string, data model.PropertyDesiredDelete) (model.PropertyDesiredDeleteResponse, error) {
//	return d.propertyDesiredDelete(deviceId, data)
//}

// PropertySetResponse 设备属性下发响应
func (d *PluginService) PropertySetResponse(deviceId string, data model.PropertySetResponse) error {
	return d.propertySetResponse(deviceId, data)
}

// PropertyGetResponse 设备属性查询响应
func (d *PluginService) PropertyGetResponse(deviceId string, data model.PropertyGetResponse) error {
	return d.propertyGetResponse(deviceId, data)
}

// ServiceExecuteResponse 设备动作执行响应
func (d *PluginService) ServiceExecuteResponse(deviceId string, data model.ServiceExecuteResponse) error {
	return d.serviceExecuteResponse(deviceId, data)
}

// GetCustomStorage 根据key值获取驱动存储的自定义内容
func (d *PluginService) GetCustomStorage(keys []string) (map[string][]byte, error) {
	return d.getCustomStorage(keys)
}

// PutCustomStorage 存储驱动的自定义内容
func (d *PluginService) PutCustomStorage(kvs map[string][]byte) error {
	return d.putCustomStorage(kvs)
}

// DeleteCustomStorage 根据key值删除驱动存储的自定义内容
func (d *PluginService) DeleteCustomStorage(keys []string) error {
	return d.deleteCustomStorage(keys)
}

// GetAllCustomStorage 获取所有驱动存储的自定义内容
func (d *PluginService) GetAllCustomStorage() (map[string][]byte, error) {
	return d.getAllCustomStorage()
}

// GetAllCustomStorage 获取所有驱动存储的自定义内容
func (d *PluginService) AppSendCommand(deviceId, serviceId, data string) error {
	return d.appSendCommandRequest(deviceId, serviceId, data)
}
