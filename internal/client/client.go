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

package client

import (
	"context"
	"errors"
	"time"

	"github.com/ytuox/elink-sdk-go/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	pb_app "github.com/ytuox/elink-plugin-proto/app"
	pb_common "github.com/ytuox/elink-plugin-proto/common"
	pb_device "github.com/ytuox/elink-plugin-proto/device"
	pb_product "github.com/ytuox/elink-plugin-proto/product"
	pb_storage "github.com/ytuox/elink-plugin-proto/storage"
	pb_thingmodel "github.com/ytuox/elink-plugin-proto/thingmodel"
)

type ResourceClient struct {
	address string
	Conn    *grpc.ClientConn
	pb_common.CommonClient
	pb_device.RPCDeviceClient
	pb_product.RPCProductClient
	pb_app.RPCAppClient

	pb_storage.StorageClient
	pb_thingmodel.RPCThingModelClient
}

var connParams = grpc.ConnectParams{
	Backoff: backoff.Config{
		BaseDelay:  time.Second * 1.0,
		Multiplier: 1.0,
		Jitter:     0,
		MaxDelay:   10 * time.Second,
	},
	MinConnectTimeout: time.Second * 3,
}

var keep = keepalive.ClientParameters{
	Time:                10 * time.Second,
	Timeout:             3 * time.Second,
	PermitWithoutStream: true,
}

func dial(address string) (*grpc.ClientConn, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithKeepaliveParams(keep), grpc.WithConnectParams(connParams))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewCoreClient(cfg config.AdapterRPC) (*ResourceClient, error) {

	if cfg.Address == "" {
		return nil, errors.New("required address")
	}

	conn, err := dial(cfg.Address)
	if err != nil {
		return nil, err
	}

	rc := &ResourceClient{
		address:             cfg.Address,
		Conn:                conn,
		CommonClient:        pb_common.NewCommonClient(conn),
		StorageClient:       pb_storage.NewStorageClient(conn),
		RPCAppClient:        pb_app.NewRPCAppClient(conn),
		RPCDeviceClient:     pb_device.NewRPCDeviceClient(conn),
		RPCProductClient:    pb_product.NewRPCProductClient(conn),
		RPCThingModelClient: pb_thingmodel.NewRPCThingModelClient(conn),
	}
	return rc, nil
}

func (c *ResourceClient) Close() error {
	return c.Conn.Close()
}
