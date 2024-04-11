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

import pb_product "github.com/ytuox/elink-plugin-proto/product"

type ProductNodeType uint

const (
	NodeTypeUnKnow    ProductNodeType = 0 // unKnow
	NodeTypeDevice    ProductNodeType = 1 // device
	NodeTypeGateway   ProductNodeType = 2 // gateway
	NodeTypeSubDevice ProductNodeType = 3 // subDevice
)

func TransformRpcNodeTypeToModel(nodeType pb_product.ProductNodeType) ProductNodeType {
	switch nodeType {
	case pb_product.ProductNodeType_Device:
		return NodeTypeDevice
	case pb_product.ProductNodeType_Gateway:
		return NodeTypeGateway
	case pb_product.ProductNodeType_SubDevice:
		return NodeTypeSubDevice
	default:
		return NodeTypeUnKnow
	}
}

type ProductNetType string

const (
	ProductNetTypeWifi       ProductNetType = "WIFI"
	ProductNetTypeCellular   ProductNetType = "蜂窝"
	ProductNetTypeEthernet   ProductNetType = "以太网"
	ProductNetTypeBLE        ProductNetType = "BLE"
	ProductNetTypeLoRaWAN    ProductNetType = "LoRaWAN"
	ProductNetTypeSerialPort ProductNetType = "串口"
	ProductNetTypeOther      ProductNetType = "其他"
)

func TransformRpcNetTypeToModel(netType pb_product.ProductNetType) ProductNetType {
	switch netType {
	case pb_product.ProductNetType_Wifi:
		return ProductNetTypeWifi
	case pb_product.ProductNetType_Cellular:
		return ProductNetTypeCellular
	case pb_product.ProductNetType_Ethernet:
		return ProductNetTypeEthernet
	case pb_product.ProductNetType_BLE:
		return ProductNetTypeBLE
	case pb_product.ProductNetType_LoRaWAN:
		return ProductNetTypeLoRaWAN
	case pb_product.ProductNetType_SerialPort:
		return ProductNetTypeSerialPort
	default:
		return ProductNetTypeOther
	}
}
