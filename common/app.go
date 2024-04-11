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

import pb_app "github.com/ytuox/elink-plugin-proto/app"

type CloudServiceInfo struct {
	AppId   string
	Address string
	AppName string
	Status  CloudServiceStatus
}

type CloudServiceStatus string

const (
	Start CloudServiceStatus = "start"
	Stop  CloudServiceStatus = "stop"
)

func (i CloudServiceStatus) TransformToRpcAppStatus() pb_app.AppStatus {
	switch i {
	case Stop:
		return pb_app.AppStatus_Stop
	case Start:
		return pb_app.AppStatus_Start
	default:
		return pb_app.AppStatus_Stop
	}
}
