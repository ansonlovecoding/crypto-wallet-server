package trongrpc

import (
	tron2 "Share-Wallet/pkg/coingrpc/trongrpc/tron"
	"Share-Wallet/pkg/wallet/coin"
	"Share-Wallet/pkg/wallet/config"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
)

// NewRPCClient try to connect TRON node RPC Server to create client instance
func NewRPCClient(conf *config.Tron) (*client.GrpcClient, error) {
	url := conf.Host
	rpcClient := client.NewGrpcClientWithTimeout(url, time.Duration(conf.Timeout)*time.Second)
	rpcClient.Start(grpc.WithTransportCredentials(insecure.NewCredentials()))
	return rpcClient, nil
}

// NewTron creates tron instance according to coinType
func NewTron(rpcClient *client.GrpcClient, conf *config.Tron, confTrc *config.TRC20, coinTypeCode coin.CoinTypeCode) (Troner, error) {
	tron := tron2.NewTron(context.Background(), rpcClient, coinTypeCode, conf, confTrc)
	return tron, nil
}
