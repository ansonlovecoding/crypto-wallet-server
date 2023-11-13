package config

import (
	"Share-Wallet/pkg/wallet/coin"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// NewWallet creates wallet config
func NewWallet(file string, coinTypeCode coin.CoinTypeCode) (*WalletRoot, error) {
	if file == "" {
		return nil, errors.New("config file should be passed")
	}

	var err error
	conf, err := loadWallet(file)
	if err != nil {
		return nil, err
	}

	// debug
	// debug.Debug(conf)

	// validate
	//if err = conf.validate(coinTypeCode); err != nil {
	//	return nil, err
	//}

	return conf, nil
}

// loadWallet load config file
func loadWallet(path string) (*WalletRoot, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read toml file. %s", path)
	}

	var config WalletRoot
	_, err = toml.Decode(string(d), &config)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call toml.Decode()")
	}

	return &config, nil
}

// validate config
//func (c *WalletRoot) validate(coinTypeCode coin.CoinTypeCode) error {
//	validate := validator.New()
//
//	switch coinTypeCode {
//	case coin.BTC, coin.BCH:
//		if err := validate.StructExcept(c, "Ethereum", "Ripple", "Tron"); err != nil {
//			return err
//		}
//	case coin.ETH, coin.USDTERC20:
//		if err := validate.StructExcept(c, "AddressType", "Bitcoin", "Ripple"); err != nil {
//			return err
//		}
//	case coin.TRX, coin.USDTTRC20:
//		if err := validate.StructExcept(c, "AddressType", "Bitcoin", "Ripple"); err != nil {
//			return err
//		}
//	case coin.XRP:
//		if err := validate.StructExcept(c, "AddressType", "Bitcoin", "Ethereum"); err != nil {
//			return err
//		}
//	default:
//	}
//
//	return nil
//}
