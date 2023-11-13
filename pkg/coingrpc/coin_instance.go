package coingrpc

import (
	"Share-Wallet/pkg/coingrpc/bitcoingrpc"
	"Share-Wallet/pkg/coingrpc/ethgrpc"
	clientTron "Share-Wallet/pkg/coingrpc/trongrpc"
	"Share-Wallet/pkg/wallet/coin"
	"Share-Wallet/pkg/wallet/config"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	eth  ethgrpc.Ethereumer
	btc  bitcoingrpc.Bitcoiner
	tron clientTron.Troner
)

// GetETH returns eth instance
// FIXME: hard coded
func GetETHInstance() (ethgrpc.Ethereumer, error) {
	if eth != nil {
		return eth, nil
	}

	cfgName := os.Getenv("ETH_TOML_PATH")
	if len(cfgName) == 0 {
		_, b, _, _ := runtime.Caller(0)
		// Root folder of this project
		Root := filepath.Join(filepath.Dir(b), "../..")
		cfgName = Root + "/config/wallet/eth.toml"
	}

	log.Printf("config path: %s", cfgName)

	conf, err := config.NewWallet(cfgName, coin.ETH)
	if err != nil {
		log.Printf("fail to create config: %v", err)
		return nil, err
	}
	// TODO: if config should be overridden, here
	// client
	client, err := ethgrpc.NewRPCClient(&conf.Ethereum)
	if err != nil {
		log.Printf("fail to create ethereum rpc client: %v", err)
		return nil, err
	}
	eth, err = ethgrpc.NewEthereum(client, &conf.Ethereum, &conf.USDTERC20, conf.CoinTypeCode)
	if err != nil {
		log.Printf("fail to create eth instance: %v", err)
		return nil, err
	}
	return eth, nil
}

func GetETHWebsocketInstance() (ethgrpc.Ethereumer, error) {
	if eth != nil {
		return eth, nil
	}

	cfgName := os.Getenv("ETH_TOML_PATH")
	if len(cfgName) == 0 {
		_, b, _, _ := runtime.Caller(0)
		// Root folder of this project
		Root := filepath.Join(filepath.Dir(b), "../..")
		cfgName = Root + "/config/wallet/eth.toml"
	}

	log.Printf("config path: %s", cfgName)

	conf, err := config.NewWallet(cfgName, coin.ETH)
	if err != nil {
		log.Printf("fail to create config: %v", err)
		return nil, err
	}
	// TODO: if config should be overridden, here
	// client
	client, err := ethgrpc.NewWsClient(&conf.Ethereum)
	if err != nil {
		log.Printf("fail to create ethereum rpc client: %v", err)
		return nil, err
	}
	eth, err = ethgrpc.NewEthereum(client, &conf.Ethereum, &conf.USDTERC20, conf.CoinTypeCode)
	if err != nil {
		log.Printf("fail to create eth instance: %v", err)
		return nil, err
	}
	return eth, nil
}

func GetBTCInstance() (bitcoingrpc.Bitcoiner, error) {
	if btc != nil {
		return btc, nil
	}

	cfgName := os.Getenv("BTC_TOML_PATH")
	if len(cfgName) == 0 {
		_, b, _, _ := runtime.Caller(0)
		// Root folder of this project
		Root := filepath.Join(filepath.Dir(b), "../..")
		cfgName = Root + "/config/wallet/btc.toml"
	}

	log.Printf("config path: %s", cfgName)

	conf, err := config.NewWallet(cfgName, coin.BTC)
	if err != nil {
		log.Printf("fail to create config: %v", err)
		return nil, err
	}
	// TODO: if config should be overridden, here
	// client
	client, err := bitcoingrpc.NewRPCClient(&conf.Bitcoin)
	if err != nil {
		log.Printf("fail to create ethereum rpc client: %v", err)
		return nil, err
	}
	btc, err = bitcoingrpc.NewBitcoin(client, &conf.Bitcoin, conf.CoinTypeCode)
	if err != nil {
		log.Printf("fail to create eth instance: %v", err)
		return nil, err
	}
	return btc, nil
}

func GetTronInstance() (clientTron.Troner, error) {
	if tron != nil {
		return tron, nil
	}
	cfgName := os.Getenv("TRON_TOML_PATH")
	if len(cfgName) == 0 {
		_, b, _, _ := runtime.Caller(0)
		// Root folder of this project
		Root := filepath.Join(filepath.Dir(b), "../..")
		cfgName = Root + "/config/wallet/tron.toml"
	}

	log.Printf("config path: %s", cfgName)

	conf, err := config.NewWallet(cfgName, coin.TRX)
	if err != nil {
		log.Printf("fail to create config: %v", err)
		return nil, err
	}

	// client
	client, err := clientTron.NewRPCClient(&conf.Tron)
	if err != nil {
		log.Printf("fail to create tron rpc client: %v", err)
		return nil, err
	}
	tron, err = clientTron.NewTron(client, &conf.Tron, &conf.USDTTRC20, coin.TRX)
	if err != nil {
		log.Printf("fail to create eth instance: %v", err)
		return nil, err
	}
	return tron, nil
}
