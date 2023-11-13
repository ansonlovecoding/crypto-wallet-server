package mysql_model

import (
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/db"
	dbModel "Share-Wallet/pkg/db/mysql"
	"fmt"
	"log"
	"strings"
)

func AddTronTransaction(transaction *dbModel.TronDetailTX) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	result := dbConn.Create(&transaction)
	return result.Error
}

func GetTronTxDetailByTxid(txid string) (*dbModel.TronDetailTX, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var txDetail dbModel.TronDetailTX
	err = dbConn.Model(&txDetail).Where("sent_hash_tx", txid).Find(&txDetail).Error
	return &txDetail, err
}

func GetTronTransactionList(WhereClauseMap map[string]interface{}, page int32, pageSize int32, orderBy string) ([]dbModel.TronDetailTX, error) {
	var (
		tronTransactionList []dbModel.TronDetailTX
	)
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	dbQuery := dbConn.Table(TronTransactionTable)
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
	err = dbQuery.Debug().Find(&tronTransactionList).Error
	if err != nil {
		return nil, err
	}
	log.Println("ethTransactionList", tronTransactionList)
	fmt.Println("ethTransactionList", tronTransactionList)
	return tronTransactionList, nil
}

func GetTronTransactionListCount(WhereClauseMap map[string]interface{}) (int64, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	var count int64

	dbQuery := dbConn.Table(TronTransactionTable)

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
