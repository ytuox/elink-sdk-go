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
	pb_product "github.com/ytuox/elink-plugin-proto/product"
	"github.com/ytuox/elink-sdk-go/common"
)

type (
	Product struct {
		Id               string
		Name             string
		Description      string
		NodeType         common.ProductNodeType
		DataFormat       string
		NetType          common.ProductNetType
		ProtocolType     string
		CategoryKey      string
		ManufacturerName string
		ModelNumber      string
		Properties       []Property //属性
		Events           []Event    //事件
		Services         []Service  //服务
	}

	Service struct {
		ProductId  string
		Identifier string
		Name       string
		Desc       string
		Required   bool
		CallType   string
		Input      []InputOutput
		Output     []InputOutput
	}

	Define struct {
		Type  string
		Specs string
	}

	Event struct {
		ProductId  string
		Identifier string
		Name       string
		Desc       string
		Type       string
		Required   bool
		Params     []InputOutput
	}

	Property struct {
		ProductId  string
		Identifier string
		Name       string
		Desc       string
		Mode       string
		Required   bool
		Define     Define
	}

	InputOutput struct {
		Identifier string
		Name       string
		Define     Define
	}
)

func TransformProductModel(p *pb_product.Product) Product {
	return Product{
		Id:               p.GetId(),
		Name:             p.GetName(),
		Description:      p.GetDescription(),
		NodeType:         common.TransformRpcNodeTypeToModel(p.NodeType),
		NetType:          common.TransformRpcNetTypeToModel(p.NetType),
		CategoryKey:      p.GetCategoryKey(),
		ProtocolType:     p.GetProtocolType(),
		ManufacturerName: p.GetManufacturerName(),
		ModelNumber:      p.GetModelNumber(),
		Properties:       propertyModels(p.GetProperties()),
		Events:           eventModels(p.GetEvents()),
		Services:         serviceModels(p.GetServices()),
	}
}

func propertyModels(p []*pb_product.Properties) []Property {
	rets := make([]Property, 0, len(p))
	for i := range p {
		rets = append(rets, Property{
			ProductId:  p[i].GetProductId(),
			Name:       p[i].GetName(),
			Identifier: p[i].GetIdentifier(),
			Desc:       p[i].GetDesc(),
			Required:   p[i].GetRequired(),
			Mode:       p[i].GetMode(),
			Define:     TransformDefineModel(p[i].GetDefine()),
		})
	}
	return rets
}

func eventModels(e []*pb_product.Events) []Event {
	rets := make([]Event, 0, len(e))
	for i := range e {
		rets = append(rets, Event{
			ProductId:  e[i].GetProductId(),
			Name:       e[i].GetName(),
			Identifier: e[i].GetIdentifier(),
			Required:   e[i].GetRequired(),
			Type:       e[i].GetType(),
			Desc:       e[i].GetDesc(),
			Params:     TransformOutputData(e[i].GetParams()),
		})
	}
	return rets
}

func serviceModels(as []*pb_product.Services) []Service {
	rets := make([]Service, 0, len(as))
	for i := range as {
		rets = append(rets, Service{
			ProductId:  as[i].GetProductId(),
			Name:       as[i].GetName(),
			Identifier: as[i].GetIdentifier(),
			Required:   as[i].GetRequired(),
			CallType:   as[i].CallType,
			Desc:       as[i].GetDesc(),
			Input:      TransformInputData(as[i].GetInput()),
			Output:     TransformOutputData(as[i].GetOutput()),
		})
	}
	return rets
}

func TransformInputData(params []*pb_product.InputOutput) []InputOutput {
	rets := make([]InputOutput, 0, len(params))
	for i := range params {
		rets = append(rets, InputOutput{
			Identifier: params[i].GetIdentifier(),
			Name:       params[i].GetName(),
			Define:     TransformDefineModel(params[i].GetDefine()),
		})
	}
	return rets
}

func TransformOutputData(params []*pb_product.InputOutput) []InputOutput {
	rets := make([]InputOutput, 0, len(params))
	for i := range params {
		rets = append(rets, InputOutput{
			Identifier: params[i].GetIdentifier(),
			Name:       params[i].GetName(),
			Define:     TransformDefineModel(params[i].GetDefine()),
		})
	}
	return rets
}

func TransformDefineModel(spec *pb_product.Define) Define {
	return Define{
		Type:  spec.GetType(),
		Specs: spec.GetSpecs(),
	}
}
