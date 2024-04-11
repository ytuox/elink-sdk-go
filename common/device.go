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

package common

import pb_device "github.com/ytuox/elink-plugin-proto/device"

type DeviceConnectStatus string

const (
	Online  DeviceConnectStatus = "online"  //在线
	Offline DeviceConnectStatus = "offline" //离线
)

type DeviceStatus uint

const (
	DeviceUnActive DeviceStatus = 0 //未激活
	DeviceOffline  DeviceStatus = 1 //离线
	DeviceOnline   DeviceStatus = 2 //在线
	DeviceDisable  DeviceStatus = 3 //禁用
)

func TransformRpcDeviceStatusToModel(deviceStatus pb_device.DeviceStatus) DeviceStatus {
	switch deviceStatus {
	case pb_device.DeviceStatus_Online:
		return DeviceOnline
	case pb_device.DeviceStatus_Offline:
		return DeviceOffline
	case pb_device.DeviceStatus_Disable:
		return DeviceDisable
	default:
		return DeviceUnActive
	}
}
