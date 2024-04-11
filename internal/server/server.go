package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"time"

	pb_app "github.com/ytuox/elink-plugin-proto/app"
	pb_app_callback "github.com/ytuox/elink-plugin-proto/appcallback"
	pb_common "github.com/ytuox/elink-plugin-proto/common"
	pb_device_callback "github.com/ytuox/elink-plugin-proto/devicecallback"
	pb_product_callback "github.com/ytuox/elink-plugin-proto/productcallback"
	pb_thingmodel "github.com/ytuox/elink-plugin-proto/thingmodel"
	"github.com/ytuox/elink-sdk-go/common"
	"github.com/ytuox/elink-sdk-go/interfaces"
	"github.com/ytuox/elink-sdk-go/internal/cache"
	"github.com/ytuox/elink-sdk-go/internal/client"
	"github.com/ytuox/elink-sdk-go/internal/config"
	"github.com/ytuox/elink-sdk-go/internal/logger"
	"github.com/ytuox/elink-sdk-go/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RPCService struct {
	pb_app_callback.UnimplementedAppCallBackServiceServer
	pb_product_callback.UnimplementedProductCallBackServiceServer
	pb_device_callback.UnimplementedDeviceCallBackServiceServer
	pb_thingmodel.UnimplementedRPCThingModelServer

	*CommonRPCServer
	ctx             context.Context
	lis             net.Listener
	rpcs            *grpc.Server
	deviceProvider  cache.DeviceProvider
	productProvider cache.ProductProvider
	pluginProvider  interfaces.Plugin
	logger          logger.Logger
	cli             *client.ResourceClient
	isRunning       bool
}

func (server *RPCService) AppStatusCallback(ctx context.Context,
	request *pb_app_callback.AppStatusCallbackRequest) (*emptypb.Empty, error) {
	var notifyType common.PluginNotifyType
	if request.GetStatus() == pb_app.AppStatus_Stop {
		notifyType = common.PluginStopNotify
	} else if request.GetStatus() == pb_app.AppStatus_Start {
		notifyType = common.PluginStartNotify
	}
	if err := server.pluginProvider.PluginNotify(ctx, notifyType, request.GetAppName()); err != nil {
		return new(emptypb.Empty), status.Errorf(codes.Internal, err.Error())
	}
	return new(emptypb.Empty), nil
}

func (server *RPCService) CreateDeviceCallback(ctx context.Context, request *pb_device_callback.CreateDeviceCallbackRequest) (*emptypb.Empty, error) {
	server.logger.Info("CreateDeviceCallback:", request.String())

	dev := model.TransformDeviceModel(request.GetData())
	server.deviceProvider.Add(dev)
	if err := server.pluginProvider.DeviceNotify(ctx, common.DeviceAddNotify, dev.Id, dev); err != nil {
		return new(emptypb.Empty), status.Errorf(codes.Internal, err.Error())
	}
	return new(emptypb.Empty), nil
}

func (server *RPCService) UpdateDeviceCallback(ctx context.Context, request *pb_device_callback.UpdateDeviceCallbackRequest) (*emptypb.Empty, error) {
	server.logger.Info("UpdateDeviceCallback:", request.String())

	deviceId := request.GetData().GetId()
	if len(deviceId) == 0 {
		return new(emptypb.Empty), fmt.Errorf("")
	}
	dev, ok := server.deviceProvider.SearchById(deviceId)
	if !ok {
		return new(emptypb.Empty), status.Errorf(codes.NotFound, "failed to find device %s", deviceId)
	}
	model.UpdateDeviceModelFieldsFromProto(&dev, request.Data)
	server.deviceProvider.Update(dev)

	if err := server.pluginProvider.DeviceNotify(ctx, common.DeviceUpdateNotify, dev.Id, dev); err != nil {
		return new(emptypb.Empty), status.Errorf(codes.Internal, err.Error())
	}

	return new(emptypb.Empty), nil
}

func (server *RPCService) DeleteDeviceCallback(ctx context.Context, request *pb_device_callback.DeleteDeviceCallbackRequest) (*emptypb.Empty, error) {
	server.logger.Info("DeleteDeviceCallback:", request.String())
	id := request.GetDeviceId()
	dev, ok := server.deviceProvider.SearchById(id)
	if !ok {
		server.logger.Errorf("failed to find device %s", id)
		return new(emptypb.Empty), status.Errorf(codes.NotFound, "failed to find device %s", id)
	}
	server.deviceProvider.RemoveById(id)
	if err := server.pluginProvider.DeviceNotify(ctx, common.DeviceDeleteNotify, dev.Id, model.Device{}); err != nil {
		return new(emptypb.Empty), status.Errorf(codes.Internal, err.Error())
	}
	return new(emptypb.Empty), nil
}

func (server *RPCService) CreateProductCallback(ctx context.Context, request *pb_product_callback.CreateProductCallbackRequest) (*emptypb.Empty, error) {
	server.logger.Info("CreateProductCallback:", request.String())
	product := model.TransformProductModel(request.GetData())
	server.productProvider.Add(product)
	if err := server.pluginProvider.ProductNotify(ctx, common.ProductAddNotify, product.Id, product); err != nil {
		return new(emptypb.Empty), status.Errorf(codes.Internal, err.Error())
	}
	return new(emptypb.Empty), nil
}

func (server *RPCService) UpdateProductCallback(ctx context.Context, request *pb_product_callback.UpdateProductCallbackRequest) (*emptypb.Empty, error) {
	server.logger.Info("UpdateProductCallback:", request.String())
	productId := request.GetData().GetId()
	if len(productId) == 0 {
		return new(emptypb.Empty), fmt.Errorf("")
	}
	_, ok := server.productProvider.SearchById(productId)
	if !ok {
		return new(emptypb.Empty), status.Errorf(codes.NotFound, "failed to find product %s", productId)
	}
	product := model.TransformProductModel(request.GetData())

	server.productProvider.Update(product)

	if err := server.pluginProvider.ProductNotify(ctx, common.ProductUpdateNotify, productId, product); err != nil {
		return new(emptypb.Empty), status.Errorf(codes.Internal, err.Error())
	}

	return new(emptypb.Empty), nil
}

func (server *RPCService) DeleteProductCallback(ctx context.Context, request *pb_product_callback.DeleteProductCallbackRequest) (*emptypb.Empty, error) {
	server.logger.Info("DeleteProductCallback:", request.String())
	productId := request.GetProductId()
	product, ok := server.productProvider.SearchById(productId)
	if !ok {
		server.logger.Errorf("failed to find product %s", productId)
		return new(emptypb.Empty), status.Errorf(codes.NotFound, "failed to find device %s", productId)
	}
	server.deviceProvider.RemoveById(productId)
	if err := server.pluginProvider.ProductNotify(ctx, common.ProductDeleteNotify, product.Id, model.Product{}); err != nil {
		return new(emptypb.Empty), status.Errorf(codes.Internal, err.Error())
	}
	return new(emptypb.Empty), nil
}

func (server *RPCService) ThingModelMsgDown(ctx context.Context, request *pb_thingmodel.ThingModelMsgDownRequest) (*emptypb.Empty, error) {
	deviceId := request.GetDeviceId()
	device, ok := server.deviceProvider.SearchById(deviceId)
	if !ok {
		server.logger.Errorf("can't find cid: %s in local cache", deviceId)
		return new(emptypb.Empty), status.Errorf(codes.NotFound, "can't find cid: %s in local cache", deviceId)
	}
	switch request.GetOperationType() {
	case pb_thingmodel.OperationType_PROPERTY_SET:
		var req model.PropertySet
		if err := decoder(request.GetData(), &req); err != nil {
			server.logger.Errorf("decode data error: %s", err)
			return new(emptypb.Empty), status.Errorf(codes.Internal, "decode data error: %s", err)
		}
		req.Spec = make(map[string]model.Property, len(req.Data))
		for k := range req.Data {
			if ps, ok := server.productProvider.GetPropertySpecByIdentifier(device.ProductId, k); !ok {
				server.logger.Warnf("can't find property(%s) spec in product(%s)", k, device.ProductId)
				continue
			} else {
				req.Spec[k] = ps
			}
		}
		err := server.pluginProvider.HandlePropertySet(ctx, deviceId, req)
		if err != nil {
			server.logger.Errorf("handlePropertySet error: %s", err)
			return new(emptypb.Empty), status.Errorf(codes.Unknown, err.Error())
		}
	case pb_thingmodel.OperationType_PROPERTY_GET:
		var req model.PropertyGet
		if err := decoder(request.GetData(), &req); err != nil {
			server.logger.Errorf("decode data error: %s", err)
			return new(emptypb.Empty), status.Errorf(codes.Internal, "decode data error: %s", err)
		}

		req.Spec = make(map[string]model.Property, len(req.Data))
		for _, k := range req.Data {
			if ps, ok := server.productProvider.GetPropertySpecByIdentifier(device.ProductId, k); !ok {
				server.logger.Warnf("can't find property(%s) spec in product(%s)", k, device.ProductId)
				continue
			} else {
				req.Spec[k] = ps
			}
		}
		err := server.pluginProvider.HandlePropertyGet(ctx, deviceId, req)
		if err != nil {
			server.logger.Errorf("handlePropertyGet error: %s", err)
			return new(emptypb.Empty), status.Errorf(codes.Unknown, err.Error())
		}
	case pb_thingmodel.OperationType_SERVICE_EXECUTE:
		var req model.ServiceExecuteRequest
		if err := decoder(request.GetData(), &req); err != nil {
			server.logger.Errorf("decode data error: %s", err)
			return new(emptypb.Empty), status.Errorf(codes.Internal, "decode data error: %s", err)
		}

		if action, ok := server.productProvider.GetServiceSpecByIdentifier(device.ProductId, req.Data.ServiceId); !ok {
			server.logger.Warnf("can't find action(%s) spec in product(%s)", req.Data.ServiceId, device.ProductId)
		} else {
			req.Spec = action
		}

		err := server.pluginProvider.HandleServiceExecute(ctx, deviceId, req)
		if err != nil {
			server.logger.Errorf("handleActionExecute error: %s", err)
		}
	case pb_thingmodel.OperationType_CUSTOM_MQTT_PUBLISH:
		//server.customMqttMessage.CustomMqttMessage("", request.Data)
	default:
		return new(emptypb.Empty), status.Errorf(codes.InvalidArgument, "unsupported operation type")
	}
	return new(emptypb.Empty), nil
}

func NewRPCService(ctx context.Context, cfg config.PluginRPC, dc cache.DeviceProvider, pc cache.ProductProvider,
	pluginProvider interfaces.Plugin, cli *client.ResourceClient, logger logger.Logger) (*RPCService, error) {

	if cfg.Address == "" {
		logger.Error("required rpc address")
		return nil, errors.New("required rpc address")
	}

	lis, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
		return nil, err
	}

	rpcs := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	}), grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     30 * time.Second,
		MaxConnectionAge:      30 * time.Second,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  5 * time.Second,
		Timeout:               3 * time.Second,
	}), grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if e := recover(); e != nil {
				logger.Errorf("%s", debug.Stack())
				err = fmt.Errorf("panic:%v", e)
			}
		}()
		reply, err := handler(ctx, req)
		return reply, err
	}))

	reflection.Register(rpcs)

	return &RPCService{
		CommonRPCServer: NewCommonRPCServer(pluginProvider),
		ctx:             ctx,
		lis:             lis,
		rpcs:            rpcs,
		deviceProvider:  dc,
		productProvider: pc,
		pluginProvider:  pluginProvider,
		cli:             cli,
		logger:          logger,
	}, nil
}

func (s *RPCService) Start() error {
	if s.isRunning {
		return errors.New("the grpc server is running")
	}

	// register method
	pb_common.RegisterCommonServer(s.rpcs, s)
	pb_app_callback.RegisterAppCallBackServiceServer(s.rpcs, s)
	pb_product_callback.RegisterProductCallBackServiceServer(s.rpcs, s)
	pb_device_callback.RegisterDeviceCallBackServiceServer(s.rpcs, s)
	pb_thingmodel.RegisterRPCThingModelServer(s.rpcs, s)

	s.logger.Infof("Server starting ( %s )", s.lis.Addr().String())
	s.logger.Info("Server start success")

	defer func() {
		s.isRunning = false
	}()

	s.isRunning = true
	err := s.rpcs.Serve(s.lis)
	if err != nil {
		s.logger.Errorf("Server failed: %v", err)
	} else {
		s.logger.Info("Server stopped")
	}

	return err
}

func (s *RPCService) Stop() error {
	if s == nil {
		return errors.New("server not start")
	}

	s.logger.Info("Server shutting down")
	_ = s.cli.Conn.Close()
	s.rpcs.Stop()
	s.logger.Info("Server shut down")
	return nil
}
