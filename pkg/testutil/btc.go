package testutil

import (
	btcgrp "Share-Wallet/pkg/coingrpc/bitcoingrpc"
	"Share-Wallet/pkg/wallet/coin"
	"Share-Wallet/pkg/wallet/config"
	"log"
	"path/filepath"
	"runtime"
)

var bc btcgrp.Bitcoiner

// GetBTC returns btc instance
// FIXME: hard coded
func GetBTC() btcgrp.Bitcoiner {
	if bc != nil {
		return bc
	}

	_, b, _, _ := runtime.Caller(0)
	// Root folder of this project
	Root := filepath.Join(filepath.Dir(b), "../..")
	confPath := Root + "/config/wallet/btc.toml"
	log.Printf("config path: %s", confPath)

	conf, err := config.NewWallet(confPath, coin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	// TODO: if config should be overridden, here

	// client
	client, err := btcgrp.NewRPCClient(&conf.Bitcoin)
	if err != nil {
		log.Fatalf("fail to create bitcoin core client: %v", err)
	}
	bc, err = btcgrp.NewBitcoin(client, &conf.Bitcoin, conf.CoinTypeCode)
	if err != nil {
		log.Fatalf("fail to create btc instance: %v", err)
	}
	return bc
}
