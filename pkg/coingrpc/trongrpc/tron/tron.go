package tron

import (
	"Share-Wallet/pkg/wallet/coin"
	"Share-Wallet/pkg/wallet/config"
	context "context"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
)

// Tron includes client to call JSON-RPC
type Tron struct {
	ctx          context.Context
	rpcClient    *client.GrpcClient
	coinTypeCode coin.CoinTypeCode
	conf         *config.Tron
	confTrc20    *config.TRC20
}

func NewTron(ctx context.Context, rpcClient *client.GrpcClient, coinTypeCode coin.CoinTypeCode, conf *config.Tron, confTrc20 *config.TRC20) *Tron {
	return &Tron{
		ctx:          ctx,
		rpcClient:    rpcClient,
		coinTypeCode: coinTypeCode,
		conf:         conf,
		confTrc20:    confTrc20,
	}
}
