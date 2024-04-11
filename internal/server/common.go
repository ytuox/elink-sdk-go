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

package server

import (
	"time"

	h "github.com/ytuox/elink-sdk-go"
	"github.com/ytuox/elink-sdk-go/common"
	"github.com/ytuox/elink-sdk-go/interfaces"

	"context"
	"encoding/json"
	"strings"

	pb_common "github.com/ytuox/elink-plugin-proto/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommonRPCServer struct {
	pb_common.UnimplementedCommonServer
	pluginProvider interfaces.Plugin
}

func NewCommonRPCServer(pluginProvider interfaces.Plugin) *CommonRPCServer {
	return &CommonRPCServer{
		pluginProvider: pluginProvider,
	}
}

// Ping tests whether the service is working
func (crs *CommonRPCServer) Ping(context.Context, *emptypb.Empty) (*pb_common.Pong, error) {
	return &pb_common.Pong{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// Version obtains version information from the target service.
func (crs *CommonRPCServer) Version(context.Context, *emptypb.Empty) (*pb_common.VersionResponse, error) {
	return &pb_common.VersionResponse{
		Version: h.SdkVersion,
	}, nil
}

func (crs *CommonRPCServer) PluginNotify(ctx context.Context, request *pb_common.PluginNotifyRequest) (*emptypb.Empty, error) {
	var notifyType common.PluginNotifyType
	if request.GetStatus() == pb_common.PluginStatus_Stop {
		notifyType = common.PluginStopNotify
	} else if request.GetStatus() == pb_common.PluginStatus_Start {
		notifyType = common.PluginStartNotify
	}
	if err := crs.pluginProvider.PluginNotify(ctx, notifyType, request.GetName()); err != nil {
		return new(emptypb.Empty), status.Errorf(codes.Internal, err.Error())
	}
	return new(emptypb.Empty), nil
}

func decoder(payload string, v interface{}) error {
	d := json.NewDecoder(strings.NewReader(payload))
	d.UseNumber()
	return d.Decode(v)
}
