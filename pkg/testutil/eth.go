package testutil

import (
	"Share-Wallet/pkg/coingrpc/ethgrpc"
	"Share-Wallet/pkg/wallet/coin"
	"Share-Wallet/pkg/wallet/config"
	"log"
	"path/filepath"
	"runtime"
)

var et ethgrpc.Ethereumer

// GetETH returns eth instance
// FIXME: hard coded
func GetETH() ethgrpc.Ethereumer {
	if et != nil {
		return et
	}

	_, b, _, _ := runtime.Caller(0)
	// Root folder of this project
	Root := filepath.Join(filepath.Dir(b), "../..")
	confPath := Root + "/config/wallet/eth.toml"
	log.Printf("config path: %s", confPath)

	conf, err := config.NewWallet(confPath, coin.ETH)
	if err != nil {
		log.Printf("fail to create config: %v", err)
	}
	// TODO: if config should be overridden, here
	// client
	client, err := ethgrpc.NewRPCClient(&conf.Ethereum)
	if err != nil {
		log.Printf("fail to create ethereum rpc client: %v", err)
	}
	et, err = ethgrpc.NewEthereum(client, &conf.Ethereum, &conf.USDTERC20, conf.CoinTypeCode)
	if err != nil {
		log.Printf("fail to create eth instance: %v", err)
	}
	return et
}
