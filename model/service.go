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

type (
	// 服务影子
	ServiceShadowData struct {
		Time      int64                  `json:"time"`
		ServiceId string                 `json:"serviceId"`
		Name      string                 `json:"name"`
		Input     map[string]interface{} `json:"input"`
		Output    map[string]interface{} `json:"output"`
	}

	ServiceDataIn struct {
		ServiceId string                 `json:"serviceId"`
		Input     map[string]interface{} `json:"input"`
	}

	ServiceDataOut struct {
		ServiceId string                 `json:"serviceId"`
		Output    map[string]interface{} `json:"output"`
	}

	ServiceExecuteRequest struct {
		CommonRequest `json:",inline"`
		Data          ServiceDataIn `json:"data"`
		Spec          Service       `json:"-"`
	}

	// ServiceExecuteResponse 执行设备动作响应
	ServiceExecuteResponse struct {
		MsgId string         `json:"msgId"`
		Data  ServiceDataOut `json:"data"`
	}
)

func NewServiceExecuteResponse(msgId string, data ServiceDataOut) ServiceExecuteResponse {
	return ServiceExecuteResponse{
		MsgId: msgId,
		Data:  data,
	}
}
