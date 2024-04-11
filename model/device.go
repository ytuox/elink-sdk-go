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

package model

import (
	pb_device "github.com/ytuox/elink-plugin-proto/device"
	"github.com/ytuox/elink-sdk-go/common"
)

type (
	Device struct {
		Id          string
		Name        string
		ProductId   string
		DeviceSn    string
		Description string
		Status      common.DeviceStatus
		Secret      string
		External    map[string]string
		AppId       string
		Location    string
		CreatedBy   string
	}
)

type (
	AddDevice struct {
		Name        string
		ProductId   string
		DeviceSn    string
		Description string
		AppId       string
		Location    string
		CreatedBy   string
	}
)

func NewAddDevice(name, productId, deviceSn, description string) AddDevice {
	return AddDevice{
		Name:        name,
		ProductId:   productId,
		DeviceSn:    deviceSn,
		Description: description,
	}
}

func TransformDeviceModel(dev *pb_device.Device) Device {
	var d Device
	d.Id = dev.GetId()
	d.Name = dev.GetName()
	d.ProductId = dev.GetProductId()
	d.Description = dev.GetDescription()
	d.DeviceSn = dev.GetDeviceSn()
	d.Status = common.TransformRpcDeviceStatusToModel(dev.GetStatus())
	d.Secret = dev.GetSecret()
	d.AppId = dev.GetAppId()
	d.Location = dev.GetLocation()
	d.CreatedBy = dev.GetCreatedBy()
	return d
}

func UpdateDeviceModelFieldsFromProto(dev *Device, patch *pb_device.Device) {
	if patch.GetName() != "" {
		dev.Name = patch.GetName()
	}
	if patch.GetProductId() != "" {
		dev.ProductId = patch.GetProductId()
	}
	if patch.GetDescription() != "" {
		dev.Description = patch.GetDescription()
	}
	if patch.GetDescription() != "" {
		dev.Description = patch.GetDescription()
	}
	if patch.GetStatus().String() != "" {
		dev.Status = common.TransformRpcDeviceStatusToModel(patch.GetStatus())
	}
	if patch.GetDeviceSn() != "" {
		dev.DeviceSn = patch.GetDeviceSn()
	}
	if patch.GetAppId() != "" {
		dev.AppId = patch.GetAppId()
	}
	if patch.GetLocation() != "" {
		dev.Location = patch.GetLocation()
	}
	if patch.GetCreatedBy() != "" {
		dev.CreatedBy = patch.GetCreatedBy()
	}
}
