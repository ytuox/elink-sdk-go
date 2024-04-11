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

import "time"

type (
	// EventReport 设备向云端上报事件
	EventReport struct {
		CommonRequest `json:",inline"`
		Data          EventData `json:"data"`
	}
	EventData struct {
		Identifier string                 `json:"identifier"`
		Timestamp  int64                  `json:"timestamp"`
		Params     map[string]interface{} `json:"params"`
	}
)

func NewEventData(identifier string, params map[string]interface{}) EventData {
	return EventData{
		Identifier: identifier,
		Params:     params,
		Timestamp:  time.Now().UnixMilli(),
	}
}

func NewEventReport(data EventData) EventReport {
	return EventReport{
		CommonRequest: CommonRequest{
			Version: Version,
			//Time:    time.Now().UnixMilli(),
		},
		Data: data,
	}
}
