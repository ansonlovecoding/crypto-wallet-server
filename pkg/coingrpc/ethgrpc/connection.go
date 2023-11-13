package ethgrpc

import (
	"Share-Wallet/pkg/coingrpc/ethgrpc/contract"
	"Share-Wallet/pkg/coingrpc/ethgrpc/eth"
	"Share-Wallet/pkg/wallet/coin"
	"Share-Wallet/pkg/wallet/config"
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

// NewRPCClient try to connect Ethereum node RPC Server to create client instance
func NewRPCClient(conf *config.Ethereum) (*ethrpc.Client, error) {
	url := conf.Host
	if conf.IPCPath != "" {
		log.Println("IPC connection")
		url = conf.IPCPath
	}

	rpcClient, err := ethrpc.Dial(url)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.Dial()")
	}
	return rpcClient, nil
}

// NewWsClient try to connect Ethereum node Websocket Server to create client instance
func NewWsClient(conf *config.Ethereum) (*ethrpc.Client, error) {
	wsClient, err := ethrpc.Dial(conf.WSPath)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.Dial()")
	}
	return wsClient, nil
}

// NewEthereum creates ethereum instance according to coinType
func NewEthereum(rpcClient *ethrpc.Client, conf *config.Ethereum, confERC *config.ERC20, coinTypeCode coin.CoinTypeCode) (Ethereumer, error) {
	client := ethclient.NewClient(rpcClient)

	tokenClient, err := contract.NewContractToken(confERC.ContractAddress, client)
	if err != nil {
		log.Println("tokenClient", tokenClient)
		fmt.Println("err instance", err)
	}

	erc20coin := eth.NewERC20(
		tokenClient,
		confERC.Symbol,
		confERC.Name,
		confERC.ContractAddress,
		confERC.AbiJson,
		confERC.Decimals,
	)

	eth, err := eth.NewEthereum(
		context.Background(),
		client,
		rpcClient,
		erc20coin,
		coinTypeCode,
		conf,
	)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.NewEthereum()")
	}
	return eth, err
}
