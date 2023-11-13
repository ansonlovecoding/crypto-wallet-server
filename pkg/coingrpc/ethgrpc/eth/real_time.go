package eth

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func (e *Ethereum) SubscribeNewHead(ch chan *types.Header) (*rpc.ClientSubscription, error) {
	return e.rpcClient.Subscribe(e.ctx, "eth", ch, "newHeads")
}
