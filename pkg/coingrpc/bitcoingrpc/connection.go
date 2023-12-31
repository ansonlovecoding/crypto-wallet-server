package bitcoingrpc

import (
	"Share-Wallet/pkg/coingrpc/bitcoingrpc/bch"
	"Share-Wallet/pkg/coingrpc/bitcoingrpc/btc"
	"Share-Wallet/pkg/wallet/coin"
	"Share-Wallet/pkg/wallet/config"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/pkg/errors"
)

// NewRPCClient try to connect bitcoin core RPCserver to create client instance
// using HTTP POST mode
func NewRPCClient(conf *config.Bitcoin) (*rpcclient.Client, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         conf.Host,
		User:         conf.User,
		Pass:         conf.Pass,
		HTTPPostMode: conf.PostMode,   // Bitcoin core only supports HTTP POST mode
		DisableTLS:   conf.DisableTLS, // Bitcoin core does not provide TLS by default
	}

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, errors.Errorf("rpcclient.New() error: %s", err)
	}
	return client, err
}

// NewBitcoin creates bitcoin/bitcoin cash instance according to coinType
func NewBitcoin(client *rpcclient.Client, conf *config.Bitcoin, coinTypeCode coin.CoinTypeCode) (Bitcoiner, error) {
	switch coinTypeCode {
	case coin.BTC:
		bit, err := btc.NewBitcoin(client, coinTypeCode, conf)
		if err != nil {
			return nil, errors.Wrap(err, "fail to call btc.NewBitcoin()")
		}

		return bit, err
	case coin.BCH:
		// BCH
		bitc, err := bch.NewBitcoinCash(client, coinTypeCode, conf)
		if err != nil {
			return nil, errors.Wrap(err, "fail to call bch.NewBitcoinCash()")
		}

		return bitc, err
	}
	return nil, errors.Errorf("coinType %s is not defined", coinTypeCode.String())
}
