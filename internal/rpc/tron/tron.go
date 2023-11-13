package tron

import (
	common2 "Share-Wallet/internal/common"
	"Share-Wallet/pkg/coingrpc"
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	db2 "Share-Wallet/pkg/db"
	wdb "Share-Wallet/pkg/db/mysql"
	"Share-Wallet/pkg/db/mysql/mysql_model"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/tron"
	"Share-Wallet/pkg/utils"
	"context"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"

	api2 "github.com/fbsobreira/gotron-sdk/pkg/proto/api"

	"github.com/fbsobreira/gotron-sdk/pkg/common"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"

	uuid2 "github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"google.golang.org/grpc"
)

type tronRPCServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewTronRPCServer(port int) *tronRPCServer {
	return &tronRPCServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.TronRPC,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *tronRPCServer) Run() {
	log.NewInfo("0", "tronRPCServer start ")
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)

	// listener network
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + s.rpcRegisterName)
	}
	log.NewInfo("0", "listen network success, ", address, listener)
	defer listener.Close()
	// grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()

	// Service registers with etcd
	tron.RegisterTronServer(srv, s)
	rpcRegisterIP := ""
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "message tron rpc success")
}

func (s *tronRPCServer) GetTronBalanceRPC(_ context.Context, req *tron.GetBalanceReq) (*tron.GetBalanceResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetTronBalanceRPC!!", req.String())

	if req.Address == "" {
		resp := &tron.GetBalanceResp{
			ErrCode: constant.ErrArgs.ErrCode,
			ErrMsg:  "address is nil!",
		}
		return resp, nil
	}

	tronClient, err := coingrpc.GetTronInstance()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetETHInstance failed!", err.Error())
		resp := &tron.GetBalanceResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  err.Error(),
		}
		return resp, nil
	}

	var balance *big.Int
	if req.CoinType == constant.TRX {
		//TRX balance
		balance, err = tronClient.BalanceAt(req.Address)
	} else {
		//Token balance
		contractAddr := tronClient.GetUSDTTRC20ContractAddress()
		balance, err = tronClient.GetTokenBalance(req.Address, contractAddr)
	}

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "tronClient get balance failed!", err.Error())
		resp := &tron.GetBalanceResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  err.Error(),
		}
		return resp, nil
	}
	if balance == nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "balance is nil!")
		resp := &tron.GetBalanceResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  "balance is nil",
		}
		return resp, nil
	}

	resp := &tron.GetBalanceResp{
		ErrCode: 0,
		ErrMsg:  "get balance success!",
		Balance: balance.String(),
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "balance:", balance.String())
	return resp, nil
}

func (s *tronRPCServer) CreateTransactionRPC(ctx context.Context, req *tron.CreateTransactionReq) (*tron.CreateTransactionResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "CreateTransactionRPC!!!", req.String())

	var (
		resp = &tron.CreateTransactionResp{}
	)

	var err error
	tronClient, err := coingrpc.GetTronInstance()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTronInstance failed", err.Error())
		resp.ErrCode = constant.ErrEthInstance.ErrCode
		resp.ErrMsg = constant.ErrEthInstance.ErrMsg
		return resp, fmt.Errorf("GetETHInstance failed %w", err)
	}

	newAmount, _ := decimal.NewFromString(req.Amount)
	var tx *api2.TransactionExtention
	if req.CoinType == constant.TRX {
		tx, err = tronClient.CreateTransaction(req.FromAddress, req.ToAddress, newAmount.BigInt())
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreateTransaction failed", err.Error())
			resp.ErrCode = constant.ErrCreateTronTransaction.ErrCode
			resp.ErrMsg = constant.ErrCreateTronTransaction.ErrMsg
			return resp, fmt.Errorf("CreateTransaction failed %w", err)
		}
	} else {
		contractAddress := tronClient.GetUSDTTRC20ContractAddress()
		fee, err := tronClient.EstimateFee(req.FromAddress, req.ToAddress, contractAddress, newAmount.BigInt())
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "EstimateFee failed", err.Error())
			return nil, utils.WrapErrorWithCode(constant.ErrCreateTronTransaction)
		}

		tx, err = tronClient.CreateTokenTransaction(req.FromAddress, req.ToAddress, contractAddress, newAmount.BigInt(), fee)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreateTokenTransaction failed", err.Error())
			resp.ErrCode = constant.ErrCreateTronTransaction.ErrCode
			resp.ErrMsg = constant.ErrCreateTronTransaction.ErrMsg
			return resp, fmt.Errorf("CreateTokenTransaction failed %w", err)
		}
	}

	//txByte, err := proto.Marshal(tx)
	//if err != nil {
	//	log.NewError(req.OperationID, utils.GetSelfFuncName(), "proto.Marshal failed", err.Error())
	//	resp.ErrCode = constant.ErrArgs.ErrCode
	//	resp.ErrMsg = "proto.Marshal failed"
	//	return resp, fmt.Errorf("proto.Marshal failed %w", err)
	//}

	resp.ErrCode = constant.OK.Code()
	resp.ErrMsg = constant.OK.ErrMsg
	txID := common.BytesToHexString(tx.Txid)
	resp.TxID = strings.TrimPrefix(txID, "0x")

	rawTxByte, err := proto.Marshal(tx.Transaction)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "proto.Marshal(tx.Transaction.RawData) failed", err.Error())
		resp.ErrCode = constant.ErrCreateTronTransaction.ErrCode
		resp.ErrMsg = constant.ErrCreateTronTransaction.ErrMsg
		return resp, fmt.Errorf("CreateTransaction failed %w", err)
	}
	resp.RawTXData = common.Bytes2Hex(rawTxByte)

	return resp, nil
}

func (s *tronRPCServer) TransferRPC(ctx context.Context, req *tron.PostTransferReq) (*tron.TransferRPCResponse, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TransferRPC!!!", req.String())

	var resp tron.TransferRPCResponse

	//check the account
	_, err := mysql_model.GetAccountInformationByMerchantUidAndUid(req.FromMerchantUID, req.FromAccountUID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAccountInformationByMerchantUidAndUid failed", err.Error(), "resp.ErrCode", resp.ErrCode, "resp.ErrMsg", resp.ErrMsg)
		return nil, utils.WrapErrorWithCode(constant.ErrTransactMerchantIncorrect)
	}

	var receiverAccountUID string
	receiveAccount, err := mysql_model.GetAccountInformationByPublicAddress(req.ToAddress, req.CoinType)
	if err == nil {
		receiverAccountUID = receiveAccount.UUID
	}

	tronClient, err := coingrpc.GetTronInstance()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTronInstance failed", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrTronInstance)
	}

	newAmount, _ := decimal.NewFromString(req.Amount)

	//check the raw transaction
	if len(req.TxDataStr) == 0 {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "req.TxDataStr is empty", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrInfo{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: "req.TxDataStr is empty"})
	}

	// SendSignedRawTransaction
	var rawTx core.Transaction
	rawTxByte, err := common.HexStringToBytes(req.TxDataStr)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "common.HexStringToBytes(req.TxDataStr) failed", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrInfo{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: "common.HexStringToBytes failed"})
	}

	err = proto.Unmarshal(rawTxByte, &rawTx)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "proto.Unmarshal failed", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrInfo{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: "proto.Unmarshal failed"})
	}

	err = tronClient.SendTransaction(&rawTx)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SendTransaction failed", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrTransactFailed)
	}

	//insert eth detail
	txHash := req.TxID

	uuid, _ := uuid2.NewUUID()
	energyUsed, _ := decimal.NewFromString(req.EnergyUsed)
	timeNow := time.Now()
	tronDetail := wdb.TronDetailTX{
		UUID:            uuid.String(),
		SenderAccount:   req.FromAccountUID,
		SenderAddress:   req.FromAddress,
		ReceiverAccount: receiverAccountUID,
		ReceiverAddress: req.ToAddress,
		Amount:          newAmount,
		EnergyUsed:      energyUsed,
		SentHashTX:      txHash,
		SentUpdatedAt:   &timeNow,
		ConfirmTime:     nil,
		Status:          constant.TxStatusPending,
		CoinType:        utils.GetCoinName(uint8(req.CoinType)),
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Completed Stage 3!!!", req.String())
	err = mysql_model.AddTronTransaction(&tronDetail)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "AddTronTransaction failed", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrDB)
	}

	var trxAmount, usdtAmount float64
	if req.CoinType == constant.TRX {
		trxAmount, _ = utils.SunToTrx(newAmount.BigInt()).Float64()

		trxFundLog := wdb.FundsLog{
			UUID:            tronDetail.UUID,
			UID:             req.FromAccountUID,
			MerchantUid:     req.FromMerchantUID,
			Txid:            tronDetail.SentHashTX,
			TransactionType: constant.TransactionTypeSendString,
			CoinType:        utils.GetCoinName(uint8(req.CoinType)),
			UserAddress:     tronDetail.SenderAddress,
			OppositeAddress: tronDetail.ReceiverAddress,
			AmountOfCoins:   decimal.NewFromFloat(trxAmount),
			CreationTime:    time.Now().Unix(),
			State:           constant.FundlogPending,
		}
		err = mysql_model.CreateFundLog(trxFundLog)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreateFundLog failed", err.Error())
			return nil, utils.WrapErrorWithCode(constant.ErrDB)
		}

	} else {
		usdtAmount = utils.ConvertBigUSDT2Float(tronDetail.Amount.BigInt())

		usdtFundLog := wdb.FundsLog{
			UUID:            tronDetail.UUID,
			UID:             req.FromAccountUID,
			MerchantUid:     req.FromMerchantUID,
			Txid:            tronDetail.SentHashTX,
			TransactionType: constant.TransactionTypeSendString,
			CoinType:        utils.GetCoinName(uint8(req.CoinType)),
			UserAddress:     tronDetail.SenderAddress,
			OppositeAddress: tronDetail.ReceiverAddress,
			AmountOfCoins:   decimal.NewFromFloat(usdtAmount),
			CreationTime:    time.Now().Unix(),
			State:           constant.FundlogPending,
		}
		err = mysql_model.CreateFundLog(usdtFundLog)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreateFundLog failed", err.Error())
			return nil, utils.WrapErrorWithCode(constant.ErrDB)
		}
	}

	//insert to redis queue
	err = db2.DB.RedisDB.InsertUnconfirmedOrder(tronDetail.SentHashTX, uint8(req.CoinType))
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "InsertUnconfirmedOrder failed", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrDB)
	}
	//update confirm_status
	err = mysql_model.UpdateConfirmStatus(tronDetail.SentHashTX, uint8(req.CoinType), constant.ConfirmationStatusPending)
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "UpdateConfirmStatus error", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrDB)
	}

	//resp.Transaction = &eth.EthTransactionDetail{
	//	UUID:            ethDetail.UUID,
	//	SenderAccount:   ethDetail.SenderAccount,
	//	SenderAddress:   ethDetail.SenderAddress,
	//	ReceiverAccount: ethDetail.ReceiverAccount,
	//	ReceiverAddress: ethDetail.ReceiverAddress,
	//	Amount:          ethDetail.Amount.String(),
	//	Fee:             ethDetail.Fee.String(),
	//	GasLimit:        ethDetail.GasLimit,
	//	Nonce:           ethDetail.Nonce,
	//	SentHashTX:      ethDetail.SentHashTX,
	//	SentUpdatedAt:   uint64(ethDetail.SentUpdatedAt.Unix()),
	//	GasPrice:        ethDetail.GasPrice.String(),
	//	Status:          int32(ethDetail.Status),
	//}
	return &resp, nil
}

func (s *tronRPCServer) GetConfirmationRPC(ctx context.Context, req *tron.GetTronConfirmationReq) (*tron.GetTronConfirmationRes, error) {

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetConfirmation!!!", req.String())
	var (
		resp = &tron.GetTronConfirmationRes{}
	)

	//check the lock
	err := db2.DB.RedisDB.LockOrder(req.TransactionHash, uint8(req.CoinType), constant.OrderLockDuration)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "You can't do it, because this order is under confirming now!", err.Error())
		return nil, errors.New("You can't do it, because this order is under confirming now!")
	}
	defer func() {
		db2.DB.RedisDB.UnLockOrder(req.TransactionHash, uint8(req.CoinType))
	}()

	tronClient, err := coingrpc.GetTronInstance()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTronInstance failed!", err.Error())
		return nil, err
	}

	//increase the check_times
	err = mysql_model.IncreaseCheckTimes(req.TransactionHash, uint8(req.CoinType))
	if err != nil {
		if req.MessageID != "" {
			db2.DB.RedisDB.RemoveUnconfirmedOrder(req.MessageID, uint8(req.CoinType))
		}
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "IncreaseCheckTimes failed!", err.Error())
		return nil, err
	}

	//check the eth transaction detail
	tronDetail, err := mysql_model.GetTronTxDetailByTxid(req.TransactionHash)
	if err != nil {
		if req.MessageID != "" {
			db2.DB.RedisDB.RemoveUnconfirmedOrder(req.MessageID, uint8(req.CoinType))
		}
		return nil, err
	}
	if tronDetail.ConfirmStatus == constant.ConfirmationStatusCompleted {
		if req.MessageID != "" {
			db2.DB.RedisDB.RemoveUnconfirmedOrder(req.MessageID, uint8(req.CoinType))
		}
		return nil, errors.New("The order was confirmed before")
	}

	//if the transaction is pending, it will return nil
	transactionInfo, err := tronClient.GetTransactionInfo(req.TransactionHash)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTransactionReceipt failed!", err.Error(), req.TransactionHash)
		return nil, err
	}

	var txStatus int32
	var fundLogState int32
	if transactionInfo.Result == core.TransactionInfo_SUCESS {
		//success
		txStatus = constant.TxStatusSuccess
		fundLogState = constant.FundlogSuccess
	} else if transactionInfo.Result == core.TransactionInfo_FAILED {
		//failed
		txStatus = constant.TxStatusFailed
		fundLogState = constant.FundlogFailed
	} else {
		txStatus = constant.TxStatusPending
		fundLogState = constant.FundlogPending
	}

	blockerNumber := transactionInfo.BlockNumber
	netFee := transactionInfo.Fee
	confirmTime := transactionInfo.BlockTimeStamp / 1000 //convert to second
	energyUsage := transactionInfo.Receipt.EnergyUsageTotal
	netUsage := transactionInfo.Receipt.NetUsage

	resp.BlockNum = decimal.NewFromInt(blockerNumber).String()
	resp.ConfirmTime = uint64(confirmTime)
	resp.NetFee = decimal.NewFromInt(netFee).String()
	resp.Status = txStatus
	resp.EnergyUsage = decimal.NewFromInt(energyUsage).String()

	blockDecimal, _ := decimal.NewFromString(resp.BlockNum)

	//convert the transaction amount
	var transactAmount, transactFee string
	if req.CoinType == constant.TRX {
		newAmount := utils.SunToTrx(tronDetail.Amount.BigInt())
		transactAmount = newAmount.String()
		txFeeInTrx := utils.SunToTrx(big.NewInt(netFee))
		transactFee = txFeeInTrx.String()
	} else {
		newAmount := utils.ConvertBigUSDT2Float(tronDetail.Amount.BigInt())
		transactAmount = strconv.FormatFloat(newAmount, 'f', 6, 64)
		txFeeInTrx := utils.SunToTrx(big.NewInt(netFee))
		transactFee = txFeeInTrx.String()
	}

	amountOfCoinsStr, _ := decimal.NewFromString(transactAmount)
	networkFeeStr, _ := decimal.NewFromString(transactFee)

	//check the sender information
	var senderMerchantUid string
	senderAccount, _ := mysql_model.GetAccountInformationByPublicAddress(tronDetail.SenderAddress, req.CoinType)
	if senderAccount != nil {
		senderMerchantUid = senderAccount.MerchantUid
	}

	//check the sending fund log, if not exist that means the transaction is from the other wallet application
	sentLog, _ := mysql_model.GetFundLog(uint8(req.CoinType), req.TransactionHash, constant.TransactionTypeSendString)
	if (sentLog == nil || sentLog.Txid == "") && senderMerchantUid != "" {
		//the sender is our user, but he is doing transaction in other wallet, we will also create the fund log for him
		sentLog = &wdb.FundsLog{
			UUID:                    tronDetail.UUID,
			UID:                     tronDetail.SenderAccount,
			MerchantUid:             senderMerchantUid,
			Txid:                    tronDetail.SentHashTX,
			TransactionType:         constant.TransactionTypeSendString,
			UserAddress:             tronDetail.SenderAddress,
			OppositeAddress:         tronDetail.ReceiverAddress,
			CoinType:                utils.GetCoinName(uint8(req.CoinType)),
			AmountOfCoins:           amountOfCoinsStr,
			NetworkFee:              networkFeeStr,
			GasUsed:                 uint64(netFee),
			ConfirmationBlockNumber: resp.BlockNum,
			ConfirmationTime:        confirmTime,
			State:                   int8(fundLogState),
			CreationTime:            time.Now().Unix(),
		}
	}

	//check the receiveAccount, if not exist that means the receiver is from the other wallet application
	var receiveFundLog wdb.FundsLog
	receiveAccount1, err := mysql_model.GetAccountInformationByPublicAddress(tronDetail.ReceiverAddress, req.CoinType)
	if err == nil && receiveAccount1.UUID != "" {
		//create the fund log for receiver
		receiveFundLog.UUID = tronDetail.UUID
		receiveFundLog.UID = receiveAccount1.UUID
		receiveFundLog.MerchantUid = receiveAccount1.MerchantUid
		receiveFundLog.Txid = tronDetail.SentHashTX
		receiveFundLog.TransactionType = constant.TransactionTypeReceiveString
		receiveFundLog.UserAddress = tronDetail.ReceiverAddress
		receiveFundLog.OppositeAddress = tronDetail.SenderAddress
		receiveFundLog.CoinType = utils.GetCoinName(uint8(req.CoinType))
		receiveFundLog.AmountOfCoins = amountOfCoinsStr
		receiveFundLog.NetworkFee = networkFeeStr
		receiveFundLog.GasUsed = uint64(netFee)
		receiveFundLog.ConfirmationBlockNumber = resp.BlockNum
		receiveFundLog.ConfirmationTime = confirmTime
		receiveFundLog.State = int8(fundLogState)
		receiveFundLog.CreationTime = time.Now().Unix()
	}
	log.NewInfo(req.OperationID, "Network fee is: ", networkFeeStr)
	err = mysql_model.UpdateTxConfirmation(int(txStatus), int(req.CoinType), req.TransactionHash, uint64(confirmTime), big.NewInt(netFee), nil, big.NewInt(energyUsage), big.NewInt(netUsage), blockDecimal.BigInt(), sentLog, &receiveFundLog, networkFeeStr)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateTxConfirmation failed!", err.Error())

		//revert confirm status to waiting
		err = mysql_model.UpdateConfirmStatus(req.TransactionHash, uint8(req.CoinType), constant.ConfirmationStatusWaitting)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateConfirmStatus failed!", err.Error())
			return nil, err
		}

		return nil, errors.New("UpdateTxConfirmation failed")
	}
	resp.Status = txStatus

	if req.MessageID != "" {
		db2.DB.RedisDB.RemoveUnconfirmedOrder(req.MessageID, uint8(req.CoinType))
	}

	//record the addresses to update the account information
	if sentLog != nil && sentLog.Txid != "" {
		err = db2.DB.RedisDB.InsertAddressToUpdate(tronDetail.SenderAddress)
		if err != nil {
			log.NewError(req.OperationID, "the sender address is from other wallet application!")
		}
	}
	if receiveAccount1 != nil {
		err = db2.DB.RedisDB.InsertAddressToUpdate(tronDetail.ReceiverAddress)
		if err != nil {
			log.NewError(req.OperationID, "Failed to insert transaction addresses to redis stream")
		}
	}

	//push notifications
	go common2.PushMsg(req.OperationID, senderMerchantUid, receiveAccount1.MerchantUid, tronDetail.SenderAddress, tronDetail.ReceiverAddress, tronDetail.SentHashTX, transactAmount, uint32(txStatus), req.CoinType)

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetConfirmation success!")
	return resp, nil
}
