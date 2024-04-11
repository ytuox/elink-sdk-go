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
	"context"
	"errors"

	"time"

	"github.com/ytuox/elink-sdk-go/common"
	"github.com/ytuox/elink-sdk-go/interfaces"
	"github.com/ytuox/elink-sdk-go/internal/cache"
	"github.com/ytuox/elink-sdk-go/internal/client"
	"github.com/ytuox/elink-sdk-go/internal/config"
	"github.com/ytuox/elink-sdk-go/internal/logger"
	"github.com/ytuox/elink-sdk-go/internal/server"
	"github.com/ytuox/elink-sdk-go/internal/snowflake"
	"github.com/ytuox/elink-sdk-go/model"
	"github.com/ytuox/elink-sdk-go/util"

	pb_app "github.com/ytuox/elink-plugin-proto/app"
	pb_common "github.com/ytuox/elink-plugin-proto/common"
	pb_device "github.com/ytuox/elink-plugin-proto/device"
	pb_storage "github.com/ytuox/elink-plugin-proto/storage"
	pb_thingmodel "github.com/ytuox/elink-plugin-proto/thingmodel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PluginService struct {
	ctx          context.Context
	direction    common.DataDirection
	cfg          config.AdapterCfg
	logger       logger.Logger
	deviceCache  cache.DeviceProvider
	productCache cache.ProductProvider
	plugin       interfaces.Plugin
	rpcClient    *client.ResourceClient
	rpcServer    *server.RPCService
	baseMessage  common.BaseMessage
	node         *snowflake.Worker
}

func NewPluginService(ctx context.Context, conf string, direction common.DataDirection) (*PluginService, error) {
	var cfg config.AdapterCfg

	util.StringDecoder(conf, &cfg)

	log := logger.NewLogger(cfg.Logger.Path, cfg.Logger.Level, cfg.AdapterId)

	// Start rpc client
	coreClient, err := client.NewCoreClient(cfg.AdapterRPC)
	if err != nil {
		log.Errorf("new resource client error: %s", err)
		return nil, err
	}

	// Snowflake node
	node, err := snowflake.NewWorker(1)
	if err != nil {
		log.Errorf("new msg id generator error: %s", err)
		return nil, err
	}

	pluginService := &PluginService{
		ctx:       ctx,
		rpcClient: coreClient,
		logger:    log,
		cfg:       cfg,
		node:      node,
		direction: direction,
	}

	if err = pluginService.buildRpcBaseMessage(); err != nil {
		log.Error("buildRpcBaseMessage error:", err)
		return nil, err
	}

	if err = pluginService.initCache(); err != nil {
		log.Error("initCache error:", err)
		return nil, err
	}

	return pluginService, nil
}

func (d *PluginService) buildRpcBaseMessage() error {
	var baseMessage common.BaseMessage
	baseMessage.PluginId = d.cfg.GetAdapterId()
	baseMessage.Direction = d.direction
	d.baseMessage = baseMessage
	return nil
}

func (d *PluginService) initCache() error {
	// Sync device
	if deviceCache, err := cache.InitDeviceCache(d.baseMessage, d.rpcClient, d.logger); err != nil {
		d.logger.Errorf("sync device error: %s", err.Error())
		return err
	} else {
		d.deviceCache = deviceCache
	}

	// Sync product
	if productCache, err := cache.InitProductCache(d.baseMessage, d.rpcClient, d.logger); err != nil {
		d.logger.Errorf("sync tsl error: %s", err.Error())
		return err
	} else {
		d.productCache = productCache
	}
	return nil
}

func (d *PluginService) start(plugin interfaces.Plugin) error {

	if plugin == nil {
		return errors.New("plugin unimplemented")
	}
	d.plugin = plugin

	var err error

	// rpc server
	d.rpcServer, err = server.NewRPCService(d.ctx, d.cfg.PluginRPC, d.deviceCache, d.productCache, d.plugin, d.rpcClient, d.logger)
	if err != nil {
		return err
	}

	err = d.rpcServer.Start()
	if err != nil {
		return err
	}

	return nil
}

func (d *PluginService) stop() error {
	return d.rpcServer.Stop()
}

func (d *PluginService) propertySetResponse(cid string, data model.PropertySetResponse) error {
	msg, err := common.TransformToProtoMsg(cid, common.PropertySetResponse, data, d.baseMessage)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err = d.rpcClient.ThingModelMsgUp(ctx, msg); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	return nil
}

func (d *PluginService) propertyGetResponse(cid string, data model.PropertyGetResponse) error {
	msg, err := common.TransformToProtoMsg(cid, common.PropertyGetResponse, data, d.baseMessage)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err = d.rpcClient.ThingModelMsgUp(ctx, msg); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	return nil
}

func (d *PluginService) serviceExecuteResponse(cid string, data model.ServiceExecuteResponse) error {
	msg, err := common.TransformToProtoMsg(cid, common.ServiceExecuteResponse, data, d.baseMessage)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err = d.rpcClient.ThingModelMsgUp(ctx, msg); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	return nil
}

func (d *PluginService) propertyReport(cid string, data model.PropertyReport) (model.CommonResponse, error) {
	msg, err := common.TransformToProtoMsg(cid, common.PropertyReport, data, d.baseMessage)
	if err != nil {
		return model.CommonResponse{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	thingModelResp := new(pb_common.CommonResponse)
	if thingModelResp, err = d.rpcClient.ThingModelMsgUp(ctx, msg); err != nil {
		return model.CommonResponse{}, errors.New(status.Convert(err).Message())
	}

	return model.NewCommonResponse(thingModelResp), nil
}

func (d *PluginService) eventReport(cid string, data model.EventReport) (model.CommonResponse, error) {
	msgId := d.node.GetId().String()
	data.MsgId = msgId
	msg, err := common.TransformToProtoMsg(cid, common.EventReport, data, d.baseMessage)
	if err != nil {
		return model.CommonResponse{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	thingModelResp := new(pb_common.CommonResponse)
	if thingModelResp, err = d.rpcClient.ThingModelMsgUp(ctx, msg); err != nil {
		return model.CommonResponse{}, errors.New(status.Convert(err).Message())
	}

	return model.NewCommonResponse(thingModelResp), nil
}

func (d *PluginService) batchReport(cid string, data model.BatchReport) (model.CommonResponse, error) {
	msgId := d.node.GetId().String()
	data.MsgId = msgId
	msg, err := common.TransformToProtoMsg(cid, common.BatchReport, data, d.baseMessage)
	if err != nil {
		return model.CommonResponse{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	thingModelResp := new(pb_common.CommonResponse)
	if thingModelResp, err = d.rpcClient.ThingModelMsgUp(ctx, msg); err != nil {
		return model.CommonResponse{}, errors.New(status.Convert(err).Message())
	}

	return model.NewCommonResponse(thingModelResp), nil
}

func (d *PluginService) propertyDesiredGet(deviceId string, data model.PropertyDesiredGet) (model.PropertyDesiredGetResponse, error) {
	msgId := d.node.GetId().String()
	data.MsgId = msgId
	msg, err := common.TransformToProtoMsg(deviceId, common.PropertyDesiredGet, data, d.baseMessage)
	if err != nil {
		return model.PropertyDesiredGetResponse{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	thingModelResp := new(pb_common.CommonResponse)
	if thingModelResp, err = d.rpcClient.ThingModelMsgUp(ctx, msg); err != nil {
		return model.PropertyDesiredGetResponse{}, errors.New(status.Convert(err).Message())
	}
	d.logger.Info(thingModelResp)
	return model.PropertyDesiredGetResponse{}, nil
}

func (d *PluginService) propertyDesiredDelete(deviceId string, data model.PropertyDesiredDelete) (model.CommonResponse, error) {
	msgId := d.node.GetId().String()
	data.MsgId = msgId
	msg, err := common.TransformToProtoMsg(deviceId, common.PropertyDesiredDelete, data, d.baseMessage)
	if err != nil {
		return model.CommonResponse{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	thingModelResp := new(pb_common.CommonResponse)
	if thingModelResp, err = d.rpcClient.ThingModelMsgUp(ctx, msg); err != nil {
		return model.CommonResponse{}, errors.New(status.Convert(err).Message())
	}

	return model.NewCommonResponse(thingModelResp), nil
}

func (d *PluginService) connectIotPlatform(deviceId string) error {
	var (
		err  error
		resp *pb_device.ConnectIotPlatformResponse
	)
	if len(deviceId) == 0 {
		return errors.New("required device id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req := pb_device.ConnectIotPlatformRequest{
		BaseRequest: d.baseMessage.BuildBaseRequest(),
		DeviceId:    deviceId,
	}
	if resp, err = d.rpcClient.ConnectIotPlatform(ctx, &req); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	if resp != nil {
		if !resp.BaseResponse.Success {
			return errors.New(resp.BaseResponse.Message)
		}
		if resp.Data.Status == pb_device.ConnectStatus_ONLINE {
			return nil
		} else if resp.Data.Status == pb_device.ConnectStatus_OFFLINE {

		}
	}
	return errors.New("unKnow error")
}

func (d *PluginService) disconnectIotPlatform(deviceId string) error {
	var (
		err  error
		resp *pb_device.DisconnectIotPlatformResponse
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req := pb_device.DisconnectIotPlatformRequest{
		BaseRequest: d.baseMessage.BuildBaseRequest(),
		DeviceId:    deviceId,
	}
	if resp, err = d.rpcClient.DisconnectIotPlatform(ctx, &req); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	if resp != nil {
		if !resp.BaseResponse.Success {
			return errors.New(resp.BaseResponse.Message)
		}
		if resp.Data.Status == pb_device.ConnectStatus_ONLINE {

		} else if resp.Data.Status == pb_device.ConnectStatus_OFFLINE {
			return nil
		}
	}
	return errors.New("unKnow error")
}

func (d *PluginService) getConnectStatus(deviceId string) (common.DeviceConnectStatus, error) {
	var (
		err  error
		resp *pb_device.GetDeviceConnectStatusResponse
	)

	if len(deviceId) == 0 {
		return "", errors.New("required device cid")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req := pb_device.GetDeviceConnectStatusRequest{
		BaseRequest: d.baseMessage.BuildBaseRequest(),
		DeviceId:    deviceId,
	}
	if resp, err = d.rpcClient.GetDeviceConnectStatus(ctx, &req); err != nil {
		return "", errors.New(status.Convert(err).Message())
	}
	if resp != nil {
		if !resp.BaseResponse.Success {
			return "", errors.New(resp.BaseResponse.Message)
		}
		if resp.Data.Status == pb_device.ConnectStatus_ONLINE {
			return common.Online, nil
		} else if resp.Data.Status == pb_device.ConnectStatus_OFFLINE {
			return common.Offline, nil
		}
	}
	return "", errors.New("unKnow error")
}

func (d *PluginService) getDeviceList() []model.Device {
	var devices []model.Device
	for _, v := range d.deviceCache.All() {
		devices = append(devices, v)
	}
	return devices
}

func (d *PluginService) getDeviceListByAppId(id string) []model.Device {
	var devices []model.Device
	for _, v := range d.deviceCache.All() {
		if v.AppId == id {
			devices = append(devices, v)
		}
	}
	return devices
}

func (d *PluginService) getDeviceListByUserId(userId string) ([]model.Device, error) {
	var devices []model.Device

	c, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	resp, err := d.rpcClient.RPCDeviceClient.QueryDeviceListByUserId(c, &pb_device.QueryDeviceListByUserIdRequest{
		BaseRequest: d.baseMessage.BuildBaseRequest(),
		UserId:      userId,
	})

	if err != nil {
		return nil, err
	}

	if !resp.BaseResponse.Success {
		return nil, errors.New(resp.BaseResponse.Message)
	}

	if resp.Data != nil {
		for _, device := range resp.Data.Devices {
			devices = append(devices, model.TransformDeviceModel(device))
		}
	}
	return devices, nil
}

func (d *PluginService) getDeviceById(deviceId string) (model.Device, bool) {
	device, ok := d.deviceCache.SearchById(deviceId)
	if !ok {
		return model.Device{}, false
	}
	return device, true
}

func (d *PluginService) getDevicePropertyShadow(deviceId, identifier string) ([]model.PropertyShadowData, error) {

	c, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := d.rpcClient.RPCThingModelClient.QueryThingModelShadow(c, &pb_thingmodel.QueryThingModelShadowRequest{
		BaseRequest:   d.baseMessage.BuildBaseRequest(),
		DeviceId:      deviceId,
		OperationType: pb_thingmodel.OperationType_PROPERTY_REPORT,
		Identifier:    identifier,
	})

	if err != nil {
		return nil, err
	}

	if !resp.BaseResponse.Success {
		return nil, errors.New(resp.BaseResponse.Message)
	}

	var data []model.PropertyShadowData
	if err := util.StringDecoder(resp.GetData(), &data); err != nil {
		d.logger.Errorf("decode data error: %s", err)
		return nil, status.Errorf(codes.Internal, "decode data error: %s", err)
	}
	return data, nil
}

func (d *PluginService) getDeviceServiceShadow(deviceId, identifier string) ([]model.ServiceShadowData, error) {

	c, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := d.rpcClient.RPCThingModelClient.QueryThingModelShadow(c, &pb_thingmodel.QueryThingModelShadowRequest{
		BaseRequest:   d.baseMessage.BuildBaseRequest(),
		DeviceId:      deviceId,
		OperationType: pb_thingmodel.OperationType_SERVICE_EXECUTE,
		Identifier:    identifier,
	})

	if err != nil {
		return nil, err
	}

	if !resp.BaseResponse.Success {
		return nil, errors.New(resp.BaseResponse.Message)
	}

	var data []model.ServiceShadowData
	if err := util.StringDecoder(resp.GetData(), &data); err != nil {
		d.logger.Errorf("decode data error: %s", err)
		return nil, status.Errorf(codes.Internal, "decode data error: %s", err)
	}
	return data, nil
}

func (d *PluginService) createDevice(addDevice model.AddDevice) (model.Device, error) {

	if addDevice.ProductId == "" || addDevice.Name == "" || addDevice.DeviceSn == "" {
		return model.Device{}, errors.New("param failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	reqDevice := new(pb_device.AddDevice)
	reqDevice.Name = addDevice.Name
	reqDevice.ProductId = addDevice.ProductId
	reqDevice.DeviceSn = addDevice.DeviceSn
	reqDevice.Description = addDevice.Description
	req := pb_device.CreateDeviceRequest{
		BaseRequest: d.baseMessage.BuildBaseRequest(),
		Device:      reqDevice,
	}
	resp, err := d.rpcClient.CreateDevice(ctx, &req)
	if err != nil {
		return model.Device{}, errors.New(status.Convert(err).Message())
	}

	var deviceInfo model.Device
	if resp != nil {
		if resp.GetBaseResponse().GetSuccess() {
			deviceInfo.Id = resp.Data.Devices.Id
			deviceInfo.Name = resp.Data.Devices.Name
			deviceInfo.ProductId = resp.Data.Devices.ProductId
			deviceInfo.DeviceSn = resp.Data.Devices.DeviceSn
			deviceInfo.Secret = resp.Data.Devices.Secret
			deviceInfo.Status = common.TransformRpcDeviceStatusToModel(resp.Data.Devices.Status)
			d.deviceCache.Add(deviceInfo)
			return deviceInfo, nil
		} else {
			return deviceInfo, errors.New(resp.GetBaseResponse().GetMessage())
		}
	}
	return deviceInfo, errors.New("unKnow error")
}

func (d *PluginService) getProductProperties(productId string) (map[string]model.Property, bool) {
	return d.productCache.GetProductProperties(productId)
}

func (d *PluginService) getProductPropertyByIdentifier(productId, identifier string) (model.Property, bool) {
	return d.productCache.GetPropertySpecByIdentifier(productId, identifier)
}

func (d *PluginService) getProductEvents(productId string) (map[string]model.Event, bool) {
	return d.productCache.GetProductEvents(productId)
}

func (d *PluginService) getProductEventByIdentifier(productId, identifier string) (model.Event, bool) {
	return d.productCache.GetEventSpecByIdentifier(productId, identifier)
}

func (d *PluginService) getProductServices(productId string) (map[string]model.Service, bool) {
	return d.productCache.GetProductServices(productId)
}

func (d *PluginService) getProductServiceByIdentifier(productId, identifier string) (model.Service, bool) {
	return d.productCache.GetServiceSpecByIdentifier(productId, identifier)
}

func (d *PluginService) getCustomStorage(keys []string) (map[string][]byte, error) {
	if len(keys) <= 0 {
		return nil, errors.New("required keys")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var req = pb_storage.GetReq{
		PluginId: d.cfg.GetAdapterId(),
		Keys:     keys,
	}

	if resp, err := d.rpcClient.StorageClient.Get(ctx, &req); err != nil {
		return nil, errors.New(status.Convert(err).Message())
	} else {
		kvs := make(map[string][]byte, len(resp.GetKvs()))
		for _, value := range resp.GetKvs() {
			kvs[value.GetKey()] = value.GetValue()
		}
		return kvs, nil
	}
}

func (d *PluginService) putCustomStorage(kvs map[string][]byte) error {
	if len(kvs) <= 0 {
		return errors.New("required key value")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var kv []*pb_storage.KV
	for k, v := range kvs {
		kv = append(kv, &pb_storage.KV{
			Key:   k,
			Value: v,
		})
	}
	var req = pb_storage.PutReq{
		PluginId: d.cfg.GetAdapterId(),
		Data:     kv,
	}

	if _, err := d.rpcClient.StorageClient.Put(ctx, &req); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	return nil

}

func (d *PluginService) deleteCustomStorage(keys []string) error {
	if len(keys) <= 0 {
		return errors.New("required keys")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var req = pb_storage.DeleteReq{
		PluginId: d.cfg.GetAdapterId(),
		Keys:     keys,
	}

	if _, err := d.rpcClient.StorageClient.Delete(ctx, &req); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	return nil
}

func (d *PluginService) getAllCustomStorage() (map[string][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if resp, err := d.rpcClient.StorageClient.All(ctx, &pb_storage.AllReq{
		PluginId: d.cfg.GetAdapterId(),
	}); err != nil {
		return nil, errors.New(status.Convert(err).Message())
	} else {
		kvs := make(map[string][]byte, len(resp.Kvs))
		for _, v := range resp.GetKvs() {
			kvs[v.GetKey()] = v.GetValue()
		}
		return kvs, nil
	}
}

func (d *PluginService) appSendCommandRequest(deviceId, serviceId, data string) error {

	// msgId := d.node.GetId().String()
	// data.MsgId = msgId
	// msg, err := pb_common.TransformToProtoMsg(cid, pb_common.PropertyReport, data, d.baseMessage)
	// if err != nil {
	// 	return  err
	// }
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	reportPlatformInfoRequest := pb_app.AppSendCommandRequest{
		DeviceId:  deviceId,
		ServiceId: serviceId,
		Data:      data,
	}
	pluginReportPlatformResp, err := d.rpcClient.RPCAppClient.SendCommand(ctx, &reportPlatformInfoRequest)
	if err != nil {
		return err
	}
	if !pluginReportPlatformResp.BaseResponse.Success {
		return errors.New(pluginReportPlatformResp.BaseResponse.Message)
	}

	return nil
}
