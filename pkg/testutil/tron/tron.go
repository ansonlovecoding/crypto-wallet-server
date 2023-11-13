package tron

import (
	"Share-Wallet/pkg/coingrpc/trongrpc"
	"Share-Wallet/pkg/wallet/coin"
	"Share-Wallet/pkg/wallet/config"
	"log"
	"path/filepath"
	"runtime"
)

var troner trongrpc.Troner

func GetTRON() trongrpc.Troner {
	if troner != nil {
		return troner
	}

	_, b, _, _ := runtime.Caller(0)
	// Root folder of this project
	Root := filepath.Join(filepath.Dir(b), "../../..")
	confPath := Root + "/config/wallet/tron.toml"
	log.Printf("config path: %s", confPath)

	conf, err := config.NewWallet(confPath, coin.TRX)
	if err != nil {
		log.Printf("fail to create config: %v", err)
		return nil
	}
	// TODO: if config should be overridden, here
	// client
	client, err := trongrpc.NewRPCClient(&conf.Tron)
	if err != nil {
		log.Printf("fail to create TRON rpc client: %v", err)
	}
	troner, err = trongrpc.NewTron(client, &conf.Tron, &conf.USDTTRC20, conf.CoinTypeCode)
	if err != nil {
		log.Printf("fail to create TRON instance: %v", err)
	}
	return troner
}
