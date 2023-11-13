package db

import (
	"Share-Wallet/pkg/db/local_db/model_struct"
	"Share-Wallet/pkg/struct/sdk"
	"Share-Wallet/pkg/utils"
	"log"
	"time"

	"gorm.io/gorm/clause"
)

func (d *DataBase) InitLocalWalletCoinType() error {
	coins := []model_struct.LocalWalletType{
		{CoinType: 1, CoinName: "BTC", Description: "Bitcoin", Status: 1, CreateTime: time.Now()},
		{CoinType: 2, CoinName: "ETH", Description: "Ethereum", Status: 1, CreateTime: time.Now()},
		{CoinType: 3, CoinName: "USDT-ERC20", Description: "Thther USDT", Status: 1, CreateTime: time.Now()},
		{CoinType: 4, CoinName: "TRX", Description: "Tron", Status: 1, CreateTime: time.Now()},
		{CoinType: 5, CoinName: "USDT-TRC20", Description: "Thther USDT", Status: 1, CreateTime: time.Now()},
	}
	return utils.Wrap(d.conn.Clauses(clause.OnConflict{DoNothing: true}).Create(coins).Error, "InitLocalWalletCoinType() failed")
}
func (d *DataBase) GetLocalWalletCoinType() ([]*model_struct.LocalWalletType, error) {
	var localWallet []model_struct.LocalWalletType
	err := utils.Wrap(d.conn.Order("id asc").Where("status = ?", 1).Find(&localWallet).Error, "GetLocalWalletCoinType() failed")
	if err != nil {
		return nil, utils.Wrap(err, "GetLocalWalletCoinType() failed")
	}
	var transfer []*model_struct.LocalWalletType
	for _, v := range localWallet {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, nil
}
func (d *DataBase) GetPublicAddressByUserID(userID string, coinType int) ([]*sdk.CoinAddress, error) {
	var localWallet []*model_struct.LocalWallet
	var err error
	if coinType > 0 {
		err = utils.Wrap(d.conn.Where("coin_type = ? AND user_id = ?", coinType, userID).Debug().Find(&localWallet).Error, "GetPublicAddressByUserID() failed")
	} else {
		err = utils.Wrap(d.conn.Where("user_id = ?", userID).Debug().Find(&localWallet).Order("coin_type ASC").Error, "GetPublicAddressByUserID() failed")
	}

	if err != nil {
		return nil, utils.Wrap(err, "GetPublicAddressByUserID() failed")
	}

	log.Println("localWallet", localWallet)
	var coinAddressList []*sdk.CoinAddress
	for _, wallet := range localWallet {
		coinAddress := &sdk.CoinAddress{
			utils.GetCoinName(wallet.CoinType),
			wallet.CoinType,
			wallet.P2PKHAddress,
			wallet.ContractAddress,
		}
		coinAddressList = append(coinAddressList, coinAddress)
	}
	return coinAddressList, nil
}
