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

package config

type (
	PluginRPC struct {
		Address  string
		UseTLS   bool
		CertFile string
		KeyFile  string
	}

	AdapterRPC struct {
		Address  string
		UseTLS   bool
		CertFile string
		KeyFile  string
	}

	LogConfig struct {
		Path  string
		Level string
	}

	AdapterCfg struct {
		Connect     string
		AdapterId   string
		AdapterRPC  AdapterRPC
		PluginRPC   PluginRPC
		PluginParam string
		Logger      LogConfig
	}
)

func (d *AdapterCfg) GetAdapterId() string {
	return d.AdapterId
}
