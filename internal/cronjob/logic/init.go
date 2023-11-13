package logic

import (
	"Share-Wallet/pkg/coingrpc"
	"Share-Wallet/pkg/coingrpc/ethgrpc"
	"Share-Wallet/pkg/coingrpc/trongrpc"
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	log2 "Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/db"
	dbModel "Share-Wallet/pkg/db/mysql"
	"Share-Wallet/pkg/proto/tron"
	"encoding/hex"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"

	"github.com/ethereum/go-ethereum/common"
	uuid2 "github.com/google/uuid"

	"github.com/panjf2000/ants"

	"Share-Wallet/pkg/db/mysql/mysql_model"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/admin"
	"Share-Wallet/pkg/proto/eth"
	"Share-Wallet/pkg/proto/wallet"
	"Share-Wallet/pkg/utils"
	"context"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/core/types"
)

type CronJob struct {
	FetchingUnconfirmedLock sync.Mutex
	CheckConfirmationLock   sync.Mutex
	FetchingOrderNumber     int
	MaxCheckingTimes        int
	OrderLockDuration       int
	TronNowBlockNum         int64
	TronWg                  sync.WaitGroup
}

func Run() {
	log.Println("the cronjob is running!")
	log2.NewInfo("", utils.GetSelfFuncName(), "the cronjob is running!")

	var cronJob CronJob
	cronJob.FetchingOrderNumber = 10
	cronJob.MaxCheckingTimes = 10
	cronJob.OrderLockDuration = constant.OrderLockDuration

	c := cron.New()

	//fetching ETH unconfirmed orders
	c.AddFunc("*/30 * * * * ?", func() {
		log.Println("fetching unconfirmed order", time.Now().String())
		cronJob.startFetchingUnconfirmedOrders(constant.ETHCoin)
	})
	//checking ETH confirmation
	c.AddFunc("*/5 * * * * ?", func() {
		log.Println("checking confirmation", time.Now().String())
		cronJob.startCheckingConfirmation(constant.ETHCoin)
	})

	//fetching USDT-ERC20 unconfirmed orders
	c.AddFunc("*/30 * * * * ?", func() {
		log.Println("fetching unconfirmed order", time.Now().String())
		cronJob.startFetchingUnconfirmedOrders(constant.USDTERC20)
	})
	//checking USDT-ERC20 confirmation
	c.AddFunc("*/5 * * * * ?", func() {
		log.Println("checking confirmation", time.Now().String())
		cronJob.startCheckingConfirmation(constant.USDTERC20)
	})

	//update account balance
	c.AddFunc("*/15 * * * * ?", func() {
		log.Println("Updating account balances", time.Now().String())
		cronJob.updateAccountsBalance()
	})
	//updating coin rates
	c.AddFunc("*/30 * * * * ?", func() {
		log.Println("Updating coin rates", time.Now().String())
		cronJob.updateCoinRates()
	})

	// Getting Tron new block
	c.AddFunc("*/3 * * * * ?", func() {
		log.Println("Getting Tron new block", time.Now().String())
		cronJob.GettingTronNewBlock()
	})
	//fetching TRX unconfirmed orders
	c.AddFunc("*/30 * * * * ?", func() {
		log.Println("fetching TRX unconfirmed order", time.Now().String())
		cronJob.startFetchingUnconfirmedOrders(constant.TRX)
	})
	//checking TRX confirmation
	c.AddFunc("*/5 * * * * ?", func() {
		log.Println("checking TRX confirmation", time.Now().String())
		cronJob.startCheckingConfirmation(constant.TRX)
	})

	//fetching Trc20 unconfirmed orders
	c.AddFunc("*/30 * * * * ?", func() {
		log.Println("fetching Trc20 unconfirmed order", time.Now().String())
		cronJob.startFetchingUnconfirmedOrders(constant.USDTTRC20)
	})
	//checking Trc20 confirmation
	c.AddFunc("*/5 * * * * ?", func() {
		log.Println("checking Trc20 confirmation", time.Now().String())
		cronJob.startCheckingConfirmation(constant.USDTTRC20)
	})
	c.Start()

	go StartGettingETHNewBlock()
}

// Getting 100 unconfirmed orders which check_times below 10 from mysql and store the txid to redis list
func (job *CronJob) startFetchingUnconfirmedOrders(coinType uint8) {
	job.FetchingUnconfirmedLock.Lock()
	defer job.FetchingUnconfirmedLock.Unlock()
	unconfirmedTransactionList, err := mysql_model.GetUnconfirmedOrders(job.FetchingOrderNumber, job.MaxCheckingTimes, coinType)
	if err != nil {
		log2.NewError("", utils.GetSelfFuncName(), err.Error())
		return
	}
	for _, tx := range unconfirmedTransactionList {
		err = db.DB.RedisDB.LockOrder(tx.SentHashTX, coinType, job.OrderLockDuration)
		if err != nil {
			break
		}
		//insert to redis queue
		err = db.DB.RedisDB.InsertUnconfirmedOrder(tx.SentHashTX, coinType)
		if err != nil {
			log2.NewError("", utils.GetSelfFuncName(), "InsertUnconfirmedOrder error", err.Error())
			db.DB.RedisDB.UnLockOrder(tx.SentHashTX, coinType)
			break
		}
		//update confirm_status
		err = mysql_model.UpdateConfirmStatus(tx.SentHashTX, coinType, constant.ConfirmationStatusPending)
		if err != nil {
			log2.NewError("", utils.GetSelfFuncName(), "UpdateConfirmStatus error", err.Error())
		}
		db.DB.RedisDB.UnLockOrder(tx.SentHashTX, coinType)
	}
}

// Getting 100 records from redis each time, and confirm the orders
func (job *CronJob) startCheckingConfirmation(coinType uint8) {
	job.CheckConfirmationLock.Lock()
	defer job.CheckConfirmationLock.Unlock()
	txidList, err := db.DB.RedisDB.ReadUnconfirmedOrder(coinType, 100)
	if err != nil {
		//log2.NewError("", utils.GetSelfFuncName(), err.Error())
		return
	}

	switch coinType {
	case constant.BTCCoin:
	case constant.ETHCoin, constant.USDTERC20:
		operationID := utils.OperationIDGenerator()
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, operationID)
		if etcdConn == nil {
			log2.NewError(operationID, "getcdv3.GetConn == nil")
			return
		}
		client := eth.NewEthClient(etcdConn)

		var wg sync.WaitGroup
		for _, s := range txidList {
			wg.Add(1)
			go func(txid string, coin uint8) {
				var reqPb eth.GetEthConfirmationReq
				reqPb.CoinType = uint32(coin)

				//txid format: txID+"|"+messageID
				strList := strings.Split(txid, "|")
				messageID := strList[1]
				realTxID := strList[0]
				reqPb.OperationID = operationID
				reqPb.TransactionHash = realTxID
				reqPb.MessageID = messageID

				respPb, err := client.GetConfirmationRPC(context.Background(), &reqPb)
				wg.Done()
				if err != nil {
					log.Println(operationID, utils.GetSelfFuncName(), "GetConfirmationRPC failed in cronjob", err.Error())
					log2.NewError(operationID, utils.GetSelfFuncName(), "GetConfirmationRPC failed in cronjob", err.Error())
				} else {
					log.Println(operationID, utils.GetSelfFuncName(), "GetConfirmationRPC resp:", respPb.Status, respPb.BlockNum, respPb.ConfirmTime)
					log2.NewInfo(operationID, utils.GetSelfFuncName(), "GetConfirmationRPC resp:", respPb.Status, respPb.BlockNum, respPb.ConfirmTime)
				}

			}(s, coinType)
		}
		wg.Wait()

	case constant.TRX, constant.USDTTRC20:
		operationID := utils.OperationIDGenerator()
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.TronRPC, operationID)
		if etcdConn == nil {
			log2.NewError(operationID, "getcdv3.GetConn == nil")
			return
		}
		client := tron.NewTronClient(etcdConn)

		var wg sync.WaitGroup
		for _, s := range txidList {
			wg.Add(1)
			go func(txid string, coin uint8) {
				var reqPb tron.GetTronConfirmationReq
				reqPb.CoinType = uint32(coin)

				//txid format: txID+"|"+messageID
				strList := strings.Split(txid, "|")
				messageID := strList[1]
				realTxID := strList[0]
				reqPb.OperationID = operationID
				reqPb.TransactionHash = realTxID
				reqPb.MessageID = messageID

				respPb, err := client.GetConfirmationRPC(context.Background(), &reqPb)
				wg.Done()
				if err != nil {
					log.Println(operationID, utils.GetSelfFuncName(), "Get tron confirmation failed in cronjob", err.Error(), "txid", realTxID)
					log2.NewError(operationID, utils.GetSelfFuncName(), "Get tron confirmation failed in cronjob", err.Error(), "txid", realTxID)
				} else {
					log.Println(operationID, utils.GetSelfFuncName(), "GetConfirmationRPC resp:", respPb.Status, respPb.BlockNum, respPb.ConfirmTime)
					log2.NewInfo(operationID, utils.GetSelfFuncName(), "GetConfirmationRPC resp:", respPb.Status, respPb.BlockNum, respPb.ConfirmTime)
				}

			}(s, coinType)
		}
		wg.Wait()

	}

}

// subscribe from the eth node websocket
func StartGettingETHNewBlock() {
	log.Println("StartGettingETHNewBlock", time.Now().String())
	operationID := utils.OperationIDGenerator()
	eth, err := coingrpc.GetETHWebsocketInstance()
	if err != nil {
		log.Println(operationID, utils.GetSelfFuncName(), "GetETHWebsocketInstance failed", err.Error())
		log2.NewError(operationID, utils.GetSelfFuncName(), "GetETHWebsocketInstance failed", err.Error())
		return
	}

	headerCh := make(chan *types.Header, 1000)
	txCh := make(chan *dbModel.EthDetailTX, 1000)
	sub, err := eth.SubscribeNewHead(headerCh)
	if err != nil {
		log.Println(operationID, utils.GetSelfFuncName(), "client.Subscribe failed", err.Error())
		log2.NewError(operationID, utils.GetSelfFuncName(), "client.Subscribe failed", err.Error())
		return
	}

	abiObj, err := abi.JSON(strings.NewReader(eth.GetUSDTERC20AbiJSON()))
	if err != nil {
		log.Println(operationID, utils.GetSelfFuncName(), "abi.JSON failed", err.Error())
		log2.NewError(operationID, utils.GetSelfFuncName(), "abi.JSON failed", err.Error())
		return
	}

	for {
		select {
		case err := <-sub.Err():
			log2.NewError(operationID, utils.GetSelfFuncName(), "sub error", err.Error())
			sub, err = eth.SubscribeNewHead(headerCh)
			if err != nil {
				log.Println(operationID, utils.GetSelfFuncName(), "client.Subscribe reconnection failed", err.Error())
				log2.NewError(operationID, utils.GetSelfFuncName(), "client.Subscribe reconnection failed", err.Error())
				return
			}
		case header := <-headerCh:
			log.Println("new header, block number:", header.Number)
			//log2.NewInfo(operationID, "new header, block number:", header.Number)

			blockInfo, err := eth.GetBlockByNumber(header.Number)
			// log.Println("blockInfo", blockInfo.Hash, blockInfo.Transactions)
			if err != nil {
				log.Println(operationID, utils.GetSelfFuncName(), "eth.GetBlockByNumber failed", err.Error())
				log2.NewError(operationID, utils.GetSelfFuncName(), "eth.GetBlockByNumber failed", err.Error())
				continue
			}

			if len(blockInfo.Transactions) > 0 {
				txNum := len(blockInfo.Transactions)
				poolNum := 1000
				if txNum < 1000 {
					poolNum = poolNum
				}
				workerPool, _ := ants.NewPool(poolNum)
				var wg sync.WaitGroup
				wg.Add(txNum)
				for _, txHash := range blockInfo.Transactions {
					workerPool.Submit(dealETHTransaction(operationID, txHash, abiObj, txCh, eth, &wg))
				}

				wg.Wait()
				workerPool.Release()
			}

		case tx := <-txCh:
			log.Println("new transaction hash:", tx.SentHashTX)
			if tx != nil {
				err := mysql_model.AddEthTransaction(tx)
				if err != nil {
					log2.NewError(operationID, "Failed to insert transaction", err.Error())
				}
			}
		}
	}
}

func (job *CronJob) updateAccountsBalance() {
	var (
		ethBalance float64
		ercBalance float64
		trxBalance float64
		trcBalance float64
	)
	operationID := utils.OperationIDGenerator()
	addressList, err := db.DB.RedisDB.ReadUpdateAddresses(25)
	if err != nil {
		//log2.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		return
	}

	adminEtcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, operationID)
	if adminEtcdConn == nil {
		errMsg := operationID + "getcdv3.GetConn == nil"
		log2.NewError(operationID, errMsg)
		return
	}

	for _, address := range addressList {
		var reqETHPb eth.GetBalanceReq
		var reqTronPb tron.GetBalanceReq

		//address format: address+"|"+messageID
		strList := strings.Split(address, "|")
		messageID := strList[1]
		publicAddress := strList[0]

		// account, err := mysql_model.GetAccountInformationByPublicAddress(publicAddress, uint32(coinType))
		// if err != nil {
		// 	log2.NewInfo(operationID, utils.GetSelfFuncName(), err.Error())
		// 	return
		// }
		log2.NewInfo(operationID, utils.GetSelfFuncName(), "Fetching account(s) by public address")
		accounts, err := mysql_model.GetAccountInformationByPublicAddressTemp(publicAddress)
		if err != nil || len(accounts) == 0 {
			log2.NewError(operationID, utils.GetSelfFuncName(), "Failed to get accounts")
			db.DB.RedisDB.RemoveAddressFromList(messageID)
			log2.NewInfo(operationID, utils.GetSelfFuncName(), "removed message id", messageID)
			continue
		}
		//ETH
		ethEtcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqETHPb.OperationID)
		if ethEtcdConn == nil {
			errMsg := operationID + "getcdv3.GetConn == nil"
			log2.NewError(operationID, errMsg)
			return
		}
		ethClient := eth.NewEthClient(ethEtcdConn)
		//TRON
		tronEtcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.TronRPC, reqTronPb.OperationID)
		if tronEtcdConn == nil {
			errMsg := operationID + "getcdv3.GetConn == nil"
			log2.NewError(operationID, errMsg)
			return
		}
		tronClient := tron.NewTronClient(tronEtcdConn)
		for _, account := range accounts {
			log.Println("checking balance", account.MerchantUid)
			log2.NewInfo(operationID, utils.GetSelfFuncName(), "Account merchant_uid: ", account.MerchantUid)
			if account.EthPublicAddress != "" {
				reqETHPb.CoinType = uint32(constant.ETHCoin)
				reqETHPb.Address = account.EthPublicAddress

				ethResp, err := ethClient.GetEthBalanceRPC(context.Background(), &reqETHPb)
				if err != nil {
					log2.NewError(operationID, "GetEthBalanceRPC failed", err.Error())
					return
				}
				balanceStr := ethResp.Balance
				if balanceStr == "" {
					balanceStr = "0"
				}
				balanceInt := new(big.Int)
				balanceInt, _ = balanceInt.SetString(balanceStr, 10)
				bal := utils.Wei2Eth_str(balanceInt)
				balFloat, err := strconv.ParseFloat(bal, 64)
				if err != nil {
					return
				}
				ethBalance = utils.RoundFloat(balFloat, 8)
			}
			if account.ErcPublicAddress != "" {
				log2.NewInfo(operationID, utils.GetSelfFuncName(), "goto eth rpc!")
				reqETHPb.CoinType = uint32(constant.USDTERC20)
				reqETHPb.Address = account.ErcPublicAddress
				ercResp, err := ethClient.GetEthBalanceRPC(context.Background(), &reqETHPb)
				if err != nil {
					log2.NewError(operationID, "GetEthBalanceRPC failed", err.Error())
					return
				}
				balanceStr := ercResp.Balance
				if balanceStr == "" {
					balanceStr = "0"
				}
				tmpBalance, _ := strconv.ParseFloat(balanceStr, 64)
				tmpBalance = (tmpBalance / 1000000)
				usdtBalance, _ := utils.FormatFloat(tmpBalance, 6)
				ercBalance = usdtBalance
			}
			if account.TrxPublicAddress != "" {
				reqTronPb.CoinType = uint32(constant.TRX)
				reqTronPb.Address = account.TrxPublicAddress

				tronResp, err := tronClient.GetTronBalanceRPC(context.Background(), &reqTronPb)
				if err != nil {
					log2.NewError(operationID, "GetTronBalanceRPC failed", err.Error())
					return
				}
				balanceStr := tronResp.Balance
				if balanceStr == "" {
					balanceStr = "0"
				}
				balanceInt := new(big.Int)
				balanceInt, _ = balanceInt.SetString(balanceStr, 10)
				bal := utils.SunToTrx(balanceInt)
				balFloat, _ := bal.Float64()
				trxBalance = utils.RoundFloat(balFloat, 6)
			}
			if account.TrcPublicAddress != "" {
				log2.NewInfo(operationID, utils.GetSelfFuncName(), "goto tron rpc!")
				reqTronPb.CoinType = uint32(constant.USDTTRC20)
				reqTronPb.Address = account.TrcPublicAddress
				trcResp, err := tronClient.GetTronBalanceRPC(context.Background(), &reqTronPb)
				if err != nil {
					log2.NewError(operationID, "GetTronBalanceRPC failed", err.Error())
					return
				}
				balanceStr := trcResp.Balance
				if balanceStr == "" {
					balanceStr = "0"
				}
				tmpBalance, _ := strconv.ParseFloat(balanceStr, 64)
				tmpBalance = (tmpBalance / 1000000)
				usdtBalance, _ := utils.FormatFloat(tmpBalance, 6)
				trcBalance = usdtBalance
			}

			var updateReqPb admin.UpdateAccountBalanceReq
			updateReqPb.MerchantUid = account.MerchantUid
			updateReqPb.Uuid = account.UUID
			updateReqPb.EthBalance = ethBalance
			updateReqPb.Erc20Balance = ercBalance
			updateReqPb.TrxBalance = trxBalance
			updateReqPb.Trc20Balance = trcBalance
			updateReqPb.MessageID = messageID
			log2.NewInfo(operationID, "Updating account balances: ", ethBalance, ercBalance)
			client := admin.NewAdminClient(adminEtcdConn)
			_, er := client.UpdateAccountBalance(context.Background(), &updateReqPb)
			if er != nil {
				log2.NewError(operationID, "UpdateAccountBalance failed", err.Error())
				return
			}
		}
	}
}

func dealETHTransaction(operationID, txHash string, abiObj abi.ABI, txCh chan *dbModel.EthDetailTX, eth ethgrpc.Ethereumer, wg *sync.WaitGroup) func() {
	return func() {
		log.Println("dealTransaction", txHash)
		defer wg.Done()
		transaction, err := eth.GetTransactionByHash(txHash)
		if err != nil {
			log.Println(operationID, utils.GetSelfFuncName(), "eth.GetTransactionByHash failed", err.Error(), txHash)
			log2.NewError(operationID, utils.GetSelfFuncName(), "eth.GetTransactionByHash failed", err.Error(), txHash)
			return
		}

		//log.Println(operationID, "transaction", transaction.Hash, transaction.From, transaction.To)
		//log2.NewInfo(operationID, "transaction", transaction.Hash, transaction.From, transaction.To)
		var toAddress string
		var amount *big.Int
		var coinType uint8
		//check if our user address
		if transaction.Input != "" && strings.ToLower(transaction.To) == strings.ToLower(eth.GetUSDTERC20ContractAddress()) {
			//this is a token transaction
			coinType = constant.USDTERC20
			//recover method
			decodedSig, err := hex.DecodeString(transaction.Input[2:10])
			if err != nil {
				log.Println(operationID, utils.GetSelfFuncName(), "hex.DecodeString failed", err.Error())
				log2.NewError(operationID, utils.GetSelfFuncName(), "hex.DecodeString failed", err.Error())
				return
			}
			method, err := abiObj.MethodById(decodedSig)
			if err != nil {
				log.Println(operationID, utils.GetSelfFuncName(), "abiObj.MethodById failed", err.Error())
				log2.NewError(operationID, utils.GetSelfFuncName(), "abiObj.MethodById failed", err.Error())
				return
			}
			log.Println("MethodById return", method.ID, method.Name)
			if method.Name != "transfer" {
				log.Println(operationID, utils.GetSelfFuncName(), "this transaction is not transfer type")
				log2.NewError(operationID, utils.GetSelfFuncName(), "his transaction is not transfer type")
				return
			}

			//decode txInput Payload
			decodeData, err := hex.DecodeString(transaction.Input[10:])
			if err != nil {
				log.Println(operationID, utils.GetSelfFuncName(), "hex.DecodeString failed", err.Error())
				log2.NewError(operationID, utils.GetSelfFuncName(), "hex.DecodeString failed", err.Error())
				return
			}

			inputData := map[string]interface{}{}
			//unpack method inputs
			err = method.Inputs.UnpackIntoMap(inputData, decodeData)
			if err != nil {
				log.Println(operationID, utils.GetSelfFuncName(), "method.Inputs.UnpackIntoMap failed", err.Error())
				log2.NewError(operationID, utils.GetSelfFuncName(), "method.Inputs.UnpackIntoMap failed", err.Error())
				return
			}
			log.Println("inputData return", inputData)
			toCommonAddress := inputData["_to"].(common.Address)
			toAddress = toCommonAddress.Hex()
			amount = inputData["_value"].(*big.Int)

		} else {
			coinType = constant.ETHCoin
			toAddress = transaction.To
			amount = transaction.Value
		}

		var senderAccountUID string
		sendAccount, err := mysql_model.GetAccountInformationByPublicAddress(transaction.From, uint32(coinType))
		if err == nil {
			senderAccountUID = sendAccount.UUID
		}

		var receiverAccountUID string
		receiverAccount, err := mysql_model.GetAccountInformationByPublicAddress(toAddress, uint32(coinType))
		if err == nil {
			receiverAccountUID = receiverAccount.UUID
		}

		if senderAccountUID == "" && receiverAccountUID == "" {
			return
		}

		//check if the transaction exist
		detail, err := mysql_model.GetETHTxDetailByTxid(transaction.Hash)
		if err != nil || detail.SentHashTX == "" {
			//push to another channel to deal the database persistence
			uuid, _ := uuid2.NewUUID()
			txFee := transaction.Gas * transaction.GasPrice
			timeNow := time.Now()
			ethDetail := dbModel.EthDetailTX{
				UUID:            uuid.String(),
				SenderAccount:   senderAccountUID,
				SenderAddress:   transaction.From,
				ReceiverAccount: receiverAccountUID,
				ReceiverAddress: toAddress,
				Amount:          decimal.NewFromBigInt(amount, 0),
				Fee:             decimal.NewFromInt(txFee),
				GasLimit:        uint64(transaction.Gas),
				Nonce:           0,
				SentHashTX:      transaction.Hash,
				SentUpdatedAt:   &timeNow,
				Status:          constant.TxStatusPending,
				GasPrice:        decimal.NewFromInt(transaction.GasPrice),
				GasUsed:         uint64(transaction.Gas),
				CoinType:        utils.GetCoinName(coinType),
				ConfirmStatus:   constant.ConfirmationStatusWaitting,
				ConfirmTime:     nil,
			}
			txCh <- &ethDetail
		}
		return
	}
}

func (job *CronJob) updateCoinRates() {
	var reqPb wallet.UpdateCoinRatesReq
	operationID := utils.OperationIDGenerator()

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, operationID)
	if etcdConn == nil {
		errMsg := operationID + "getcdv3.GetConn == nil"
		log2.NewError(operationID, errMsg)
		return
	}
	client := wallet.NewWalletClient(etcdConn)
	reqPb.OperationID = operationID
	_, er := client.UpdateCoinRates(context.Background(), &reqPb)
	if er != nil {
		log2.NewError(operationID, utils.GetSelfFuncName(), "update coin rates rpc failed", er.Error())
		return
	}
}

// start getting tron new block
func (job *CronJob) GettingTronNewBlock() {
	job.TronWg.Wait()
	job.TronWg.Add(1)
	defer job.TronWg.Done()

	log.Println("StartGettingTronNewBlock", time.Now().String())
	operationID := utils.OperationIDGenerator()
	tron, err := coingrpc.GetTronInstance()
	if err != nil {
		log.Println(operationID, utils.GetSelfFuncName(), "GetTronInstance failed", err.Error())
		log2.NewError(operationID, utils.GetSelfFuncName(), "GetTronInstance failed", err.Error())
		return
	}

	nowBlockNum, err := tron.GetNowBlockNum()
	if err != nil {
		log.Println(operationID, utils.GetSelfFuncName(), "GetNowBlockNum failed", err.Error())
		log2.NewError(operationID, utils.GetSelfFuncName(), "GetNowBlockNum failed", err.Error())
		return
	}
	log.Println("tron now block number", nowBlockNum)
	log2.NewInfo(operationID, "tron now block number", nowBlockNum)
	if nowBlockNum < job.TronNowBlockNum {
		log.Println(operationID, utils.GetSelfFuncName(), "this blockNumber is already checked", nowBlockNum, job.TronNowBlockNum)
		log2.NewError(operationID, utils.GetSelfFuncName(), "this blockNumber is already checked", nowBlockNum, job.TronNowBlockNum)
		return
	}

	blockInfo, err := tron.GetBlockByNum(nowBlockNum)
	if err != nil {
		log.Println(operationID, utils.GetSelfFuncName(), "GetBlockByNum failed", err.Error())
		log2.NewError(operationID, utils.GetSelfFuncName(), "GetBlockByNum failed", err.Error())
		return
	}

	if len(blockInfo.Transactions) > 0 {
		var wg sync.WaitGroup
		txNum := len(blockInfo.Transactions)
		poolNum := 1000
		if txNum < 1000 {
			poolNum = poolNum
		}
		workerPool, _ := ants.NewPool(poolNum)

		wg.Add(txNum)
		for _, tx := range blockInfo.Transactions {
			workerPool.Submit(dealTronTransaction(operationID, tx, tron, &wg))
		}

		wg.Wait()
		workerPool.Release()
	}
}

func dealTronTransaction(operationID string, tx *api.TransactionExtention, tron trongrpc.Troner, wg *sync.WaitGroup) func() {
	return func() {
		log.Println("deal tron transaction", tx.Txid)
		defer wg.Done()

		txID := common.Bytes2Hex(tx.Txid)
		log.Println(operationID, "tron transaction", txID)
		log2.NewInfo(operationID, "tron transaction", txID)

		var fromAddress string
		var toAddress string
		var amount *big.Int
		var coinType uint8
		var contract *core.Transaction_Contract
		if tx.Transaction != nil && tx.Transaction.RawData != nil &&
			tx.Transaction.RawData.Contract != nil && len(tx.Transaction.RawData.Contract) > 0 {
			contract = tx.Transaction.RawData.Contract[0]
			if contract.Type == core.Transaction_Contract_TransferContract {
				log.Println("this is trx transaction")
				coinType = constant.TRX
				var contractParameter core.TransferContract
				err := contract.Parameter.UnmarshalTo(&contractParameter)
				if err != nil {
					log.Println(operationID, utils.GetSelfFuncName(), "contract.Parameter.UnmarshalTo failed", err.Error())
					log2.NewError(operationID, utils.GetSelfFuncName(), "contract.Parameter.UnmarshalTo failed", err.Error())
					return
				}

				fromAddress, err = utils.HexToBase58Address(common.Bytes2Hex(contractParameter.OwnerAddress))
				if err != nil {
					log.Println(operationID, utils.GetSelfFuncName(), "utils.HexToBase58Address failed", err.Error())
					log2.NewError(operationID, utils.GetSelfFuncName(), "utils.HexToBase58Address failed", err.Error())
					return
				}
				toAddress, err = utils.HexToBase58Address(common.Bytes2Hex(contractParameter.ToAddress))
				if err != nil {
					log.Println(operationID, utils.GetSelfFuncName(), "utils.HexToBase58Address failed", err.Error())
					log2.NewError(operationID, utils.GetSelfFuncName(), "utils.HexToBase58Address failed", err.Error())
					return
				}
				amount = big.NewInt(contractParameter.Amount)
				//log.Println(operationID, utils.GetSelfFuncName(), "fromAddress", fromAddress, "toAddress", toAddress, "amount", amount)
				//log2.NewError(operationID, utils.GetSelfFuncName(), "fromAddress", fromAddress, "toAddress", toAddress, "amount", amount)

			} else if contract.Type == core.Transaction_Contract_TriggerSmartContract {
				log.Println("this is token transaction")
				var contractParameter core.TriggerSmartContract
				err := contract.Parameter.UnmarshalTo(&contractParameter)
				if err != nil {
					log.Println(operationID, utils.GetSelfFuncName(), "contract.Parameter.UnmarshalTo failed", err.Error())
					log2.NewError(operationID, utils.GetSelfFuncName(), "contract.Parameter.UnmarshalTo failed", err.Error())
					return
				}
				fromAddress, err = utils.HexToBase58Address(common.Bytes2Hex(contractParameter.OwnerAddress))
				if err != nil {
					log.Println(operationID, utils.GetSelfFuncName(), "utils.HexToBase58Address failed", err.Error())
					log2.NewError(operationID, utils.GetSelfFuncName(), "utils.HexToBase58Address failed", err.Error())
					return
				}

				contractAddress, err := utils.HexToBase58Address(common.Bytes2Hex(contractParameter.ContractAddress))
				if err != nil {
					log.Println(operationID, utils.GetSelfFuncName(), "utils.HexToBase58Address failed", err.Error())
					log2.NewError(operationID, utils.GetSelfFuncName(), "utils.HexToBase58Address failed", err.Error())
					return
				}
				usdtContractAddress := tron.GetUSDTTRC20ContractAddress()
				//log.Println(operationID, utils.GetSelfFuncName(), "contractAddress", contractAddress, "usdtContractAddress", usdtContractAddress)
				//log2.NewError(operationID, utils.GetSelfFuncName(), "contractAddress", contractAddress, "usdtContractAddress", usdtContractAddress)
				if strings.ToLower(contractAddress) == strings.ToLower(usdtContractAddress) {
					coinType = constant.USDTTRC20
					if contractParameter.Data != nil {
						abiObj, err := abi.JSON(strings.NewReader(tron.GetUSDTTRC20AbiJSON()))
						if err != nil {
							log.Println(operationID, utils.GetSelfFuncName(), "abi.JSON failed", err.Error())
							log2.NewError(operationID, utils.GetSelfFuncName(), "abi.JSON failed", err.Error())
							return
						}

						method, err := abiObj.MethodById(contractParameter.Data[:4])
						if err != nil {
							log.Println(operationID, utils.GetSelfFuncName(), "abiObj.MethodById failed", err.Error())
							log2.NewError(operationID, utils.GetSelfFuncName(), "abiObj.MethodById failed", err.Error())
							return
						}
						log.Println("MethodById return", method.ID, method.Name)
						if method.Name != "transfer" {
							log.Println(operationID, utils.GetSelfFuncName(), "this transaction is not transfer type")
							log2.NewError(operationID, utils.GetSelfFuncName(), "his transaction is not transfer type")
							return
						}

						inputData := map[string]interface{}{}
						//unpack method inputs
						err = method.Inputs.UnpackIntoMap(inputData, contractParameter.Data[4:])
						if err != nil {
							log.Println(operationID, utils.GetSelfFuncName(), "method.Inputs.UnpackIntoMap failed", err.Error())
							log2.NewError(operationID, utils.GetSelfFuncName(), "method.Inputs.UnpackIntoMap failed", err.Error())
							return
						}
						log.Println("inputData return", inputData)
						toCommonAddress := inputData["_to"].(common.Address)
						ethHexAddr := toCommonAddress.Hex()
						ethHexAddr = strings.TrimPrefix(ethHexAddr, "0x")
						ethHexAddr = strings.TrimPrefix(ethHexAddr, "0X")
						tronHexAddr := "41" + ethHexAddr
						toAddress, err = utils.HexToBase58Address(tronHexAddr)
						if err != nil {
							log.Println(operationID, utils.GetSelfFuncName(), "utils.HexToBase58Address failed", err.Error())
							log2.NewError(operationID, utils.GetSelfFuncName(), "utils.HexToBase58Address failed", err.Error())
							return
						}
						amount = inputData["_value"].(*big.Int)

						//log.Println(operationID, utils.GetSelfFuncName(), "contractAddress", contractAddress, "fromAddress", fromAddress, "toAddress", toAddress, "amount", amount.String())
						//log2.NewInfo(operationID, utils.GetSelfFuncName(), "contractAddress", contractAddress, "fromAddress", fromAddress, "toAddress", toAddress, "amount", amount.String())

					} else {
						log.Println(operationID, utils.GetSelfFuncName(), "contractParameter.Data is nil")
						log2.NewError(operationID, utils.GetSelfFuncName(), "contractParameter.Data is nil")
						return
					}

				} else {
					log.Println(operationID, utils.GetSelfFuncName(), "contract address is not usdt-trc20 contract address")
					log2.NewError(operationID, utils.GetSelfFuncName(), "contract address is not usdt-trc20 contract address")
					return
				}
			} else {
				return
			}
		} else {
			log.Println(operationID, utils.GetSelfFuncName(), "tx.Transaction.RawData.Contract is nil")
			log2.NewError(operationID, utils.GetSelfFuncName(), "tx.Transaction.RawData.Contract is nil")
			return
		}

		//check if our user address
		var senderAccountUID string
		sendAccount, err := mysql_model.GetAccountInformationByPublicAddress(fromAddress, uint32(coinType))
		if err == nil {
			senderAccountUID = sendAccount.UUID
		}

		var receiverAccountUID string
		receiverAccount, err := mysql_model.GetAccountInformationByPublicAddress(toAddress, uint32(coinType))
		if err == nil {
			receiverAccountUID = receiverAccount.UUID
		}

		if senderAccountUID == "" && receiverAccountUID == "" {
			return
		}

		//check if the transaction exist
		detail, err := mysql_model.GetTronTxDetailByTxid(txID)
		if err != nil || detail.SentHashTX == "" {
			//push to another channel to deal the database persistence
			log2.NewInfo(operationID, utils.GetSelfFuncName(), "transaction not exist, insert it", txID)
			uuid, _ := uuid2.NewUUID()
			timeNow := time.Now()
			tronDetail := dbModel.TronDetailTX{
				UUID:            uuid.String(),
				SenderAccount:   senderAccountUID,
				SenderAddress:   fromAddress,
				ReceiverAccount: receiverAccountUID,
				ReceiverAddress: toAddress,
				Amount:          decimal.NewFromBigInt(amount, 0),
				Fee:             decimal.NewFromInt(0),
				SentHashTX:      txID,
				SentUpdatedAt:   &timeNow,
				Status:          constant.TxStatusPending,
				CoinType:        utils.GetCoinName(coinType),
				ConfirmStatus:   constant.ConfirmationStatusWaitting,
				ConfirmTime:     nil,
			}

			err := mysql_model.AddTronTransaction(&tronDetail)
			if err != nil {
				log2.NewError(operationID, "Failed to insert transaction", err.Error())
			}
		} else {
			log2.NewInfo(operationID, utils.GetSelfFuncName(), "transaction exist, no need to insert", txID)
		}
		return

	}
}
