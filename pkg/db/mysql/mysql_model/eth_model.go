package mysql_model

import (
	"Share-Wallet/pkg/common/constant"
	db "Share-Wallet/pkg/db"
	dbModel "Share-Wallet/pkg/db/mysql"
	"Share-Wallet/pkg/utils"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

const (
	EthTransactionTable  = "w_eth_detail_tx"
	TronTransactionTable = "w_tron_detail_tx"
)

func AddEthTransaction(transaction *dbModel.EthDetailTX) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	result := dbConn.Table(EthTransactionTable).Create(&transaction)
	return result.Error
}

func UpdateAfterTxSent(uuid string, txType int8, signedHex, sentHashTx string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	nowTime := time.Now()
	t := dbConn.Table(EthTransactionTable).
		Where("uuid=?", uuid).
		Select("SignedHexTX", "SentHashTX", "SentUpdatedAt", "CurrentTXType").Updates(
		dbModel.EthDetailTX{
			SentHashTX:    sentHashTx,
			SentUpdatedAt: &nowTime,
		})

	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateAfterTxSent failed")
}

func GetETHTxDetailByTxid(txid string) (*dbModel.EthDetailTX, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var txDetail dbModel.EthDetailTX
	err = dbConn.Table(EthTransactionTable).Where("sent_hash_tx", txid).Take(&txDetail).Error
	return &txDetail, err
}

func GetETHTransactionList(WhereClauseMap map[string]interface{}, page int32, pageSize int32, orderBy string) ([]dbModel.EthDetailTX, error) {
	var (
		ethTransactionList []dbModel.EthDetailTX
	)
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	dbQuery := dbConn.Table(EthTransactionTable)
	if err != nil {
		return nil, err
	}
	if page > 0 && pageSize > 0 {
		dbQuery = dbQuery.Limit(int(pageSize)).Offset(int(pageSize * (page - 1)))
	}

	sortMap := map[string]string{}
	var orderByClause string
	sortMap["create_time"] = "sent_updated_at"

	if orderBy != "" {
		direction := "DESC"
		sort := strings.Split(orderBy, ":")
		if len(sort) == 2 {
			if sort[1] == "asc" {
				direction = "ASC"
			}
			col, ok := sortMap[sort[0]]
			if ok {
				orderByClause = fmt.Sprintf("%s %s ", col, direction)
			}
		}
	}
	if orderByClause != "" {
		dbQuery = dbQuery.Order(orderByClause)
	}

	if txstatus, ok := WhereClauseMap["transaction_state"]; ok {
		txState := txstatus.(int)
		if txState == constant.TxStatusExcludePending {
			dbQuery = dbQuery.Where("status <> ? ", constant.TxStatusPending)
		} else {
			dbQuery = dbQuery.Where("status = ? ", txState)
		}
	}

	var senderAddress, receiverAddress string
	if address, ok := WhereClauseMap["sender_address"]; ok {
		senderAddress = address.(string)
	}
	if address, ok := WhereClauseMap["receiver_address"]; ok {
		receiverAddress = address.(string)
	}
	if senderAddress != "" && receiverAddress != "" {
		dbQuery = dbQuery.Where("sender_address = ? OR (receiver_address = ? AND confirm_block_number != 0)", senderAddress, receiverAddress)
	} else if senderAddress != "" {
		dbQuery = dbQuery.Where("sender_address = ?", senderAddress)
	} else if receiverAddress != "" {
		dbQuery = dbQuery.Where("receiver_address = ? AND confirm_block_number != 0", receiverAddress)
	}

	if coinType, ok := WhereClauseMap["coin_type"]; ok {
		if coinType != "" && ok {
			dbQuery = dbQuery.Where("coin_type = ?", coinType)
		}
	}

	if hash, ok := WhereClauseMap["sent_hash_tx"]; ok {
		if hash != "" && ok {
			dbQuery = dbQuery.Where("sent_hash_tx = ?", hash)
		}
	}
	err = dbQuery.Debug().Find(&ethTransactionList).Error
	if err != nil {
		return nil, err
	}
	log.Println("ethTransactionList", ethTransactionList)
	fmt.Println("ethTransactionList", ethTransactionList)
	return ethTransactionList, nil
}
func GetETHTransactionListCount(WhereClauseMap map[string]interface{}) (int64, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	var count int64

	dbQuery := dbConn.Table(EthTransactionTable)

	if txstatus, ok := WhereClauseMap["transaction_state"]; ok {
		txState := txstatus.(int)
		if txState == constant.TxStatusExcludePending {
			dbQuery = dbQuery.Where("status <> ? ", constant.TxStatusPending)
		} else {
			dbQuery = dbQuery.Where("status = ? ", txState)
		}
	}

	var senderAddress, receiverAddress string
	if address, ok := WhereClauseMap["sender_address"]; ok {
		senderAddress = address.(string)
	}
	if address, ok := WhereClauseMap["receiver_address"]; ok {
		receiverAddress = address.(string)
	}
	if senderAddress != "" && receiverAddress != "" {
		dbQuery = dbQuery.Where("sender_address = ? OR (receiver_address = ? AND confirm_block_number != 0)", senderAddress, receiverAddress)
	} else if senderAddress != "" {
		dbQuery = dbQuery.Where("sender_address = ?", senderAddress)
	} else if receiverAddress != "" {
		dbQuery = dbQuery.Where("receiver_address = ? AND confirm_block_number != 0", receiverAddress)
	} else {
		dbQuery = dbQuery.Where("sender_address = '' AND receiver_address = ''")
	}

	if coinType, ok := WhereClauseMap["coin_type"]; ok {
		if coinType != "" && ok {
			dbQuery = dbQuery.Where("coin_type = ?", coinType)
		}
	}

	if hash, ok := WhereClauseMap["sent_hash_tx"]; ok {
		if hash != "" && ok {
			dbQuery = dbQuery.Where("sent_hash_tx = ?", hash)
		}
	}

	dbError := dbQuery.Debug().Count(&count).Error
	if dbError != nil {
		return 0, dbError
	}
	return count, nil
}

func UpdateTxConfirmation(status, coinType int, transactionHash string, confirmTime uint64, gasUsed, gasPrice, energyUsage, netUsage, confirmBlockNumber *big.Int, newSendFundLog *dbModel.FundsLog, receiveFundLog *dbModel.FundsLog, networkFee decimal.Decimal) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	//start transaction
	err = dbConn.Transaction(func(tx *gorm.DB) error {
		unixTimeUTC := time.Unix(int64(confirmTime), 0)
		//update transaction detail
		if coinType == constant.ETHCoin || coinType == constant.USDTERC20 {
			fee := new(big.Int).Mul(gasUsed, gasPrice)
			if rowEffected := tx.Set("gorm:query_option", "FOR UPDATE").Table(EthTransactionTable).
				Where("sent_hash_tx = ? and status != 1", transactionHash).Debug().
				Select("Status", "ConfirmTime", "GasUsed", "GasPrice", "Fee", "ConfirmationBlockNumber", "ConfirmStatus").Updates(
				dbModel.EthDetailTX{
					Status:                  int8(status),
					ConfirmTime:             &unixTimeUTC,
					GasUsed:                 gasUsed.Uint64(),
					GasPrice:                decimal.NewFromBigInt(gasPrice, 0),
					Fee:                     decimal.NewFromBigInt(fee, 0),
					ConfirmationBlockNumber: confirmBlockNumber.String(),
					ConfirmStatus:           constant.ConfirmationStatusCompleted,
				}).RowsAffected; rowEffected == 0 {
				return errors.New("eth detail no updated")
			}
		} else if coinType == constant.TRX || coinType == constant.USDTTRC20 {
			if rowEffected := tx.Set("gorm:query_option", "FOR UPDATE").Table(TronTransactionTable).
				Where("sent_hash_tx = ? and status != 1", transactionHash).Debug().
				Select("Status", "ConfirmTime", "Fee", "EnergyUsed", "NetUsed", "ConfirmationBlockNumber", "ConfirmStatus").Updates(
				dbModel.TronDetailTX{
					Status:                  int8(status),
					ConfirmTime:             &unixTimeUTC,
					Fee:                     decimal.NewFromBigInt(gasUsed, 0),
					EnergyUsed:              decimal.NewFromBigInt(energyUsage, 0),
					NetUsed:                 decimal.NewFromBigInt(netUsage, 0),
					ConfirmationBlockNumber: confirmBlockNumber.String(),
					ConfirmStatus:           constant.ConfirmationStatusCompleted,
				}).RowsAffected; rowEffected == 0 {
				return errors.New("tron detail no updated")
			}
		}

		//update fund log
		var gasPriceDecimal decimal.Decimal
		if gasPrice != nil {
			gasPriceDecimal = decimal.NewFromBigInt(gasPrice, 0)
		} else {
			gasPriceDecimal = decimal.NewFromInt(0)
		}

		if newSendFundLog.ID != 0 {
			var state int
			if status == constant.TxStatusFailed {
				state = constant.FundlogFailed
			} else if status == constant.TxStatusPending {
				state = constant.FundlogPending
			} else if status == constant.TxStatusSuccess {
				state = constant.FundlogSuccess
			}
			if coinType == constant.ETHCoin || coinType == constant.USDTERC20 {
				if err := tx.Table(FundsLogTableName).
					Where("txid = ? and coin_type = ? and transaction_type = ?", transactionHash, utils.GetCoinName(uint8(coinType)), constant.TransactionTypeSendString).Debug().
					Select("State", "ConfirmationTime", "GasUsed", "GasPrice", "ConfirmationBlockNumber").Updates(
					dbModel.FundsLog{
						State:                   int8(state),
						ConfirmationTime:        int64(confirmTime),
						GasUsed:                 gasUsed.Uint64(),
						GasPrice:                gasPriceDecimal,
						ConfirmationBlockNumber: confirmBlockNumber.String(),
					}).Error; err != nil {
					return err
				}
			}
			if coinType == constant.TRX || coinType == constant.USDTTRC20 {
				err = tx.Table(FundsLogTableName).
					Where("txid = ? and coin_type = ? and transaction_type = ?", transactionHash, utils.GetCoinName(uint8(coinType)), constant.TransactionTypeSendString).Debug().
					Select("State", "ConfirmationTime", "NetworkFee", "ConfirmationBlockNumber").Updates(
					dbModel.FundsLog{
						State:                   int8(state),
						ConfirmationTime:        int64(confirmTime),
						NetworkFee:              networkFee,
						ConfirmationBlockNumber: confirmBlockNumber.String(),
					}).Error
			}
		} else {
			//create the fund log for sender
			if newSendFundLog.Txid != "" {
				if err = tx.Create(&newSendFundLog).Error; err != nil {
					return err
				}
			}
		}

		//create the fund log for receiver
		if receiveFundLog != nil && receiveFundLog.UID != "" {
			if err = tx.Create(&receiveFundLog).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func GetFundLog(coinType uint8, transactionHash, transactionType string) (*dbModel.FundsLog, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}

	var fundLog dbModel.FundsLog
	if err = dbConn.Where("coin_type = ? and txid = ? and transaction_type = ?", utils.GetCoinName(coinType), transactionHash, transactionType).Find(&fundLog).Error; err != nil {
		return nil, err
	}
	return &fundLog, nil
}

func GetRecentRecords(filters map[string]string, page int32, pageSize int32, uid string) ([]dbModel.FundsLog, int64, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, 0, err
	}
	var fundsLog []dbModel.FundsLog
	dbQuery := dbConn.Table(FundsLogTableName)
	from, fok := filters["from"]
	to, tok := filters["to"]
	if tok && fok {
		dbQuery = dbQuery.Debug().Where("confirmation_time between ? and ?", from, to)
	}
	if transactionType, ok := filters["transaction_type"]; ok {
		if transactionType != "all" {
			dbQuery = dbQuery.Debug().Where("transaction_type=?", transactionType)
		}
	}
	if userAddress, ok := filters["user_address"]; ok {
		dbQuery = dbQuery.Debug().Where("user_address=?", userAddress)
	}
	if oppositeAddress, ok := filters["opposite_address"]; ok {
		dbQuery = dbQuery.Debug().Where("opposite_address=?", oppositeAddress)
	}

	// use coins type to decide which table to retreive the data from
	if coinsType, ok := filters["coins_type"]; ok {
		if coinsType != "all" {
			dbQuery = dbQuery.Debug().Where("coin_type=?", coinsType)
		}
	}
	if state, ok := filters["state"]; ok {
		if state != "all" {
			switch state {
			case "fail":
				dbQuery = dbQuery.Debug().Where("state=?", 0)
			case "success":
				dbQuery = dbQuery.Debug().Where("state=?", 1)
			}
		} else {
			dbQuery = dbQuery.Debug().Where("state != ?", 2)
		}
	}
	if txid, ok := filters["txid"]; ok {
		dbQuery = dbQuery.Debug().Where("txid=?", txid)
	}
	if merchantUid, ok := filters["merchant_uid"]; ok {
		dbQuery = dbQuery.Debug().Where("merchant_uid=?", merchantUid)
	}
	if uid != "" {
		dbQuery = dbQuery.Debug().Where("uid=?", uid)
	}
	var count int64
	dbQuery.Count(&count)
	if page > 0 && pageSize > 0 {
		dbQuery = dbQuery.Debug().Limit(int(pageSize)).Offset(int(pageSize * (page - 1)))
	}
	result := dbQuery.Debug().Where("confirmation_time != ?", 0).Order("confirmation_time DESC").Find(&fundsLog)
	err = result.Error
	return fundsLog, count, err
}

func GetUnconfirmedOrders(orderNum, maxCheckTimes int, coinType uint8) ([]*dbModel.EthDetailTX, error) {
	var (
		ethTransactionList []*dbModel.EthDetailTX
	)
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	var tableName string
	if coinType == constant.ETHCoin || coinType == constant.USDTERC20 {
		tableName = EthTransactionTable
	} else if coinType == constant.TRX || coinType == constant.USDTTRC20 {
		tableName = TronTransactionTable
	}
	dbQuery := dbConn.Table(tableName)
	if err != nil {
		return nil, err
	}
	if orderNum > 0 {
		dbQuery = dbQuery.Limit(orderNum)
	}
	if maxCheckTimes > 0 {
		dbQuery = dbQuery.Where("check_times < ?", maxCheckTimes)
	}

	err = dbQuery.Debug().Where("confirm_block_number = 0 and confirm_status = 0 and coin_type = ?", utils.GetCoinName(coinType)).Order("sent_updated_at ASC").Find(&ethTransactionList).Error
	if err != nil {
		return nil, err
	}
	log.Println("unconfirm TransactionList", ethTransactionList)
	return ethTransactionList, nil
}

func UpdateConfirmStatus(transactionHash string, coinType, confirmStatus uint8) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	var tableName string
	if coinType == constant.ETHCoin || coinType == constant.USDTERC20 {
		tableName = EthTransactionTable
	} else if coinType == constant.TRX || coinType == constant.USDTTRC20 {
		tableName = TronTransactionTable
	}
	err = dbConn.Table(tableName).
		Where("sent_hash_tx = ? AND coin_type = ? AND confirm_status < 2", transactionHash, utils.GetCoinName(coinType)).Debug().
		Select("ConfirmStatus").Updates(
		dbModel.EthDetailTX{
			ConfirmStatus: confirmStatus,
		}).Error
	return err
}

func IncreaseCheckTimes(transactionHash string, coinType uint8) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	var tableName string
	if coinType == constant.ETHCoin || coinType == constant.USDTERC20 {
		tableName = EthTransactionTable
	} else if coinType == constant.TRX || coinType == constant.USDTTRC20 {
		tableName = TronTransactionTable
	}
	err = dbConn.Set("gorm:query_option", "FOR UPDATE").Table(tableName).
		Where("sent_hash_tx = ? AND coin_type = ? AND check_times < ?", transactionHash, utils.GetCoinName(coinType), 10).Debug().
		UpdateColumn("check_times", gorm.Expr("check_times + ?", 1)).Error

	return err
}

// create eth_detail and fund_log for received transaction from blockchain
func AddEthTransactionFromBlock(transaction *dbModel.EthDetailTX, fundLog *dbModel.FundsLog) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	//start transaction
	err = dbConn.Transaction(func(tx *gorm.DB) error {
		err = tx.Table(EthTransactionTable).Create(&transaction).Error
		if err != nil {
			return err
		}

		err = tx.Table(FundsLogTableName).Create(&fundLog).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
