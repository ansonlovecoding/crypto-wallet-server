package eth

import (
	"Share-Wallet/internal/common"
	"Share-Wallet/pkg/coingrpc"
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	db2 "Share-Wallet/pkg/db"
	wdb "Share-Wallet/pkg/db/mysql"
	"Share-Wallet/pkg/db/mysql/mysql_model"
	walletdb "Share-Wallet/pkg/db/mysql/mysql_model"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/eth"
	"Share-Wallet/pkg/utils"
	"context"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"

	uuid2 "github.com/google/uuid"

	"github.com/pkg/errors"

	"github.com/shopspring/decimal"

	"google.golang.org/grpc"
)

type ethRPCServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewEthRPCServer(port int) *ethRPCServer {
	return &ethRPCServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.EthRPC,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *ethRPCServer) Run() {
	log.NewInfo("0", "ethRPCServer start ")
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
	eth.RegisterEthServer(srv, s)
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
	log.NewInfo("0", "message cms rpc success")
}

func (s *ethRPCServer) TestEthRPC(_ context.Context, req *eth.CommonReq) (*eth.CommonResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TestEthRPC!!", req.String())
	resp := &eth.CommonResp{
		ErrCode: 0,
		ErrMsg:  "Test Success!",
	}
	return resp, nil
}

func (s *ethRPCServer) GetEthBalanceRPC(_ context.Context, req *eth.GetBalanceReq) (*eth.GetBalanceResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetEthBalanceRPC!!", req.String())

	if req.Address == "" {
		resp := &eth.GetBalanceResp{
			ErrCode: constant.ErrArgs.ErrCode,
			ErrMsg:  "address is nil!",
		}
		return resp, nil
	}

	ethClient, err := coingrpc.GetETHInstance()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetETHInstance failed!", err.Error())
		resp := &eth.GetBalanceResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  err.Error(),
		}
		return resp, nil
	}

	var balance *big.Int
	if req.CoinType == constant.ETHCoin {
		//ETH balance
		balance, err = ethClient.BalanceAt(req.Address)
	} else {
		//Token balance
		balance, err = ethClient.GetTokenBalance(req.Address)
	}

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ethClient get balance failed!", err.Error())
		resp := &eth.GetBalanceResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  err.Error(),
		}
		return resp, nil
	}
	if balance == nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "balance is nil!")
		resp := &eth.GetBalanceResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  "balance is nil",
		}
		return resp, nil
	}

	resp := &eth.GetBalanceResp{
		ErrCode: 0,
		ErrMsg:  "get balance success!",
		Balance: balance.String(),
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "balance:", balance.String())
	return resp, nil
}

func (s *ethRPCServer) GetEthGasPriceRPC(_ context.Context, req *eth.GetGasPriceReq) (*eth.GetGasPriceRes, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetEthGasPrice!!", req.String())

	ethClient, err := coingrpc.GetETHInstance()
	if err != nil {
		resp := &eth.GetGasPriceRes{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  err.Error(),
		}
		return resp, nil
	}

	ethGasPrice, err := ethClient.GasPrice()
	if err != nil {
		resp := &eth.GetGasPriceRes{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  err.Error(),
		}
		return resp, nil
	}

	if ethGasPrice == nil {
		resp := &eth.GetGasPriceRes{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  "GasPrice is nil",
		}
		return resp, nil
	}

	resp := &eth.GetGasPriceRes{
		ErrCode:  0,
		ErrMsg:   "create account success!",
		GasPrice: ethGasPrice.Int64(),
	}

	return resp, nil
}

func (s *ethRPCServer) TransferRPC(ctx context.Context, req *eth.PostTransferReq) (*eth.TransferRPCResponse, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TransferRPC!!!", req.String())

	var (
		resp = &eth.TransferRPCResponse{}
	)

	//check the account
	_, err := mysql_model.GetAccountInformationByMerchantUidAndUid(req.FromMerchantUID, req.FromAccountUID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAccountInformationByMerchantUidAndUid failed", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrTransactMerchantIncorrect)
	}

	var receiverAccountUID string
	receiveAccount, err := mysql_model.GetAccountInformationByPublicAddress(req.ToAddress, req.CoinType)
	if err == nil {
		receiverAccountUID = receiveAccount.UUID
	}

	ethClient, err := coingrpc.GetETHInstance()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetETHInstance failed", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrEthInstance)
	}

	newAmount, _ := decimal.NewFromString(req.Amount)
	newFee, _ := decimal.NewFromString(req.Fee)
	gasPrice, _ := decimal.NewFromString(req.GasPrice)

	// SendSignedRawTransaction
	txHash, err := ethClient.SendSignedRawTransaction(req.TxHash)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SendSignedRawTransaction failed", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrTransactFailed)
	}

	//insert eth detail
	uuid, _ := uuid2.NewUUID()
	timeNow := time.Now()
	ethDetail := wdb.EthDetailTX{
		UUID:            uuid.String(),
		SenderAccount:   req.FromAccountUID,
		SenderAddress:   req.FromAddress,
		ReceiverAccount: receiverAccountUID,
		ReceiverAddress: req.ToAddress,
		Amount:          newAmount,
		Fee:             newFee,
		GasLimit:        req.GasLimit,
		Nonce:           req.Nounce,
		SentHashTX:      txHash,
		SentUpdatedAt:   &timeNow,
		ConfirmTime:     nil,
		Status:          constant.TxStatusPending,
		GasPrice:        gasPrice,
		CoinType:        utils.GetCoinName(uint8(req.CoinType)),
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Completed Stage 3!!!", req.String())
	err = walletdb.AddEthTransaction(&ethDetail)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "AddEthTransaction failed", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrDB)
	}

	var ethAmount, usdtAmount, txFeeInETH float64
	if req.CoinType == constant.ETHCoin {
		ethAmount = Wei2Eth_float(newAmount.BigInt())
		txFeeInETH = Wei2Eth_float(newFee.BigInt())

		ethFundLog := wdb.FundsLog{
			UUID:            ethDetail.UUID,
			UID:             req.FromAccountUID,
			MerchantUid:     req.FromMerchantUID,
			Txid:            ethDetail.SentHashTX,
			TransactionType: constant.TransactionTypeSendString,
			CoinType:        utils.GetCoinName(uint8(req.CoinType)),
			UserAddress:     ethDetail.SenderAddress,
			OppositeAddress: ethDetail.ReceiverAddress,
			AmountOfCoins:   decimal.NewFromFloat(ethAmount),
			CreationTime:    time.Now().Unix(),
			NetworkFee:      decimal.NewFromFloat(txFeeInETH),
			GasLimit:        ethDetail.GasLimit,
			GasPrice:        ethDetail.GasPrice,
			State:           constant.FundlogPending,
		}
		err = walletdb.CreateFundLog(ethFundLog)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreateFundLog failed", err.Error())
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = constant.ErrDB.ErrMsg
			return nil, utils.WrapErrorWithCode(constant.ErrDB)
		}

	} else {
		usdtAmount = utils.ConvertBigUSDT2Float(ethDetail.Amount.BigInt())
		txFeeInETH = Wei2Eth_float(newFee.BigInt())

		usdtFundLog := wdb.FundsLog{
			UUID:            ethDetail.UUID,
			UID:             req.FromAccountUID,
			MerchantUid:     req.FromMerchantUID,
			Txid:            ethDetail.SentHashTX,
			TransactionType: constant.TransactionTypeSendString,
			CoinType:        utils.GetCoinName(uint8(req.CoinType)),
			UserAddress:     ethDetail.SenderAddress,
			OppositeAddress: ethDetail.ReceiverAddress,
			AmountOfCoins:   decimal.NewFromFloat(usdtAmount),
			CreationTime:    time.Now().Unix(),
			NetworkFee:      decimal.NewFromFloat(txFeeInETH),
			GasLimit:        ethDetail.GasLimit,
			GasPrice:        ethDetail.GasPrice,
			State:           constant.FundlogPending,
		}
		err = walletdb.CreateFundLog(usdtFundLog)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreateFundLog failed", err.Error())
			return nil, utils.WrapErrorWithCode(constant.ErrDB)
		}
	}

	//insert to redis queue
	err = db2.DB.RedisDB.InsertUnconfirmedOrder(ethDetail.SentHashTX, uint8(req.CoinType))
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "InsertUnconfirmedOrder failed", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrDB)
	}
	//update confirm_status
	err = mysql_model.UpdateConfirmStatus(ethDetail.SentHashTX, uint8(req.CoinType), constant.ConfirmationStatusPending)
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "UpdateConfirmStatus error", err.Error())
		return nil, utils.WrapErrorWithCode(constant.ErrDB)
	}

	resp.Transaction = &eth.EthTransactionDetail{
		UUID:            ethDetail.UUID,
		SenderAccount:   ethDetail.SenderAccount,
		SenderAddress:   ethDetail.SenderAddress,
		ReceiverAccount: ethDetail.ReceiverAccount,
		ReceiverAddress: ethDetail.ReceiverAddress,
		Amount:          ethDetail.Amount.String(),
		Fee:             ethDetail.Fee.String(),
		GasLimit:        ethDetail.GasLimit,
		Nonce:           ethDetail.Nonce,
		SentHashTX:      ethDetail.SentHashTX,
		SentUpdatedAt:   uint64(ethDetail.SentUpdatedAt.Unix()),
		GasPrice:        ethDetail.GasPrice.String(),
		Status:          int32(ethDetail.Status),
	}
	return resp, nil
}

func (s *ethRPCServer) GetConfirmationRPC(ctx context.Context, req *eth.GetEthConfirmationReq) (*eth.GetEthConfirmationRes, error) {

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetConfirmation!!!", req.String())
	var (
		resp = &eth.GetEthConfirmationRes{}
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

	ethClient, err := coingrpc.GetETHInstance()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetETHInstance failed!", err.Error())
		return nil, err
	}

	//increase the check_times
	err = walletdb.IncreaseCheckTimes(req.TransactionHash, uint8(req.CoinType))
	if err != nil {
		if req.MessageID != "" {
			db2.DB.RedisDB.RemoveUnconfirmedOrder(req.MessageID, uint8(req.CoinType))
		}
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "IncreaseCheckTimes failed!", err.Error())
		return nil, err
	}

	//check the eth transaction detail
	ethDetail, err := walletdb.GetETHTxDetailByTxid(req.TransactionHash)
	if err != nil {
		if req.MessageID != "" {
			db2.DB.RedisDB.RemoveUnconfirmedOrder(req.MessageID, uint8(req.CoinType))
		}
		return nil, err
	}
	if ethDetail.ConfirmStatus == constant.ConfirmationStatusCompleted {
		if req.MessageID != "" {
			db2.DB.RedisDB.RemoveUnconfirmedOrder(req.MessageID, uint8(req.CoinType))
		}
		return nil, errors.New("The order was confirmed before")
	}

	//if the transaction is pending, it will return nil
	transactionReceipt, err := ethClient.GetTransactionReceipt(req.TransactionHash)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTransactionReceipt failed!", err.Error())
		return nil, err
	}

	var txStatus int32
	var fundLogState int32
	if transactionReceipt.Status == 1 {
		//success
		txStatus = constant.TxStatusSuccess
		fundLogState = constant.FundlogSuccess
	} else if transactionReceipt.Status == 0 {
		//failed
		txStatus = constant.TxStatusFailed
		fundLogState = constant.FundlogFailed
	} else {
		txStatus = constant.TxStatusPending
		fundLogState = constant.FundlogPending
	}

	gasUsed := transactionReceipt.GasUsed
	confirmTime := uint64(time.Now().Unix())
	resp.ConfirmTime = confirmTime
	resp.BlockNum = decimal.NewFromBigInt(transactionReceipt.BlockNumber, 0).String()
	resp.GasUsed = decimal.NewFromBigInt(gasUsed, 0).String()
	blockDecimal, _ := decimal.NewFromString(resp.BlockNum)

	//convert the transaction amount
	var transactAmount, transactFee string
	if req.CoinType == constant.ETHCoin {
		newAmount := utils.ConvertWeiToEther(ethDetail.Amount.BigInt())
		transactAmount = newAmount.String()
		txFeeInETH := utils.ConvertWeiToEther(ethDetail.Fee.BigInt())
		transactFee = txFeeInETH.String()
	} else {
		newAmount := utils.ConvertBigUSDT2Float(ethDetail.Amount.BigInt())
		transactAmount = strconv.FormatFloat(newAmount, 'f', 6, 64)
		txFeeInETH := utils.ConvertWeiToEther(ethDetail.Fee.BigInt())
		transactFee = txFeeInETH.String()
	}

	amountOfCoinsStr, _ := decimal.NewFromString(transactAmount)
	networkFeeStr, _ := decimal.NewFromString(transactFee)

	//check the sender information
	var senderMerchantUid string
	senderAccount, _ := walletdb.GetAccountInformationByPublicAddress(ethDetail.SenderAddress, req.CoinType)
	if senderAccount != nil {
		senderMerchantUid = senderAccount.MerchantUid
	}

	//check the sending fund log, if not exist that means the transaction is from the other wallet application
	sentLog, _ := walletdb.GetFundLog(uint8(req.CoinType), req.TransactionHash, constant.TransactionTypeSendString)
	if (sentLog == nil || sentLog.Txid == "") && senderMerchantUid != "" {
		//the sender is our user, but he is doing transaction in other wallet, we will also create the fund log for him
		sentLog = &wdb.FundsLog{
			UUID:                    ethDetail.UUID,
			UID:                     ethDetail.SenderAccount,
			MerchantUid:             senderMerchantUid,
			Txid:                    ethDetail.SentHashTX,
			TransactionType:         constant.TransactionTypeSendString,
			UserAddress:             ethDetail.SenderAddress,
			OppositeAddress:         ethDetail.ReceiverAddress,
			CoinType:                utils.GetCoinName(uint8(req.CoinType)),
			AmountOfCoins:           amountOfCoinsStr,
			NetworkFee:              networkFeeStr,
			GasLimit:                ethDetail.GasLimit,
			GasPrice:                ethDetail.GasPrice,
			GasUsed:                 gasUsed.Uint64(),
			ConfirmationBlockNumber: resp.BlockNum,
			ConfirmationTime:        int64(confirmTime),
			State:                   int8(fundLogState),
			CreationTime:            time.Now().Unix(),
		}
	}

	//check the receiveAccount, if not exist that means the receiver is from the other wallet application
	var receiveFundLog wdb.FundsLog
	receiveAccount1, err := walletdb.GetAccountInformationByPublicAddress(ethDetail.ReceiverAddress, req.CoinType)
	if err == nil && receiveAccount1.UUID != "" {
		//create the fund log for receiver
		receiveFundLog.UUID = ethDetail.UUID
		receiveFundLog.UID = receiveAccount1.UUID
		receiveFundLog.MerchantUid = receiveAccount1.MerchantUid
		receiveFundLog.Txid = ethDetail.SentHashTX
		receiveFundLog.TransactionType = constant.TransactionTypeReceiveString
		receiveFundLog.UserAddress = ethDetail.ReceiverAddress
		receiveFundLog.OppositeAddress = ethDetail.SenderAddress
		receiveFundLog.CoinType = utils.GetCoinName(uint8(req.CoinType))
		receiveFundLog.AmountOfCoins = amountOfCoinsStr
		receiveFundLog.NetworkFee = networkFeeStr
		receiveFundLog.GasLimit = ethDetail.GasLimit
		receiveFundLog.GasPrice = ethDetail.GasPrice
		receiveFundLog.GasUsed = gasUsed.Uint64()
		receiveFundLog.ConfirmationBlockNumber = resp.BlockNum
		receiveFundLog.ConfirmationTime = int64(confirmTime)
		receiveFundLog.State = int8(fundLogState)
		receiveFundLog.CreationTime = time.Now().Unix()
	}

	err = walletdb.UpdateTxConfirmation(int(txStatus), int(req.CoinType), req.TransactionHash, confirmTime, gasUsed, transactionReceipt.EffectiveGasPrice, nil, nil, blockDecimal.BigInt(), sentLog, &receiveFundLog, networkFeeStr)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateTxConfirmation failed!", err.Error())

		//revert confirm status to waiting
		err = walletdb.UpdateConfirmStatus(req.TransactionHash, uint8(req.CoinType), constant.ConfirmationStatusWaitting)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateConfirmStatus failed!", err.Error())
			return nil, err
		}

		return nil, fmt.Errorf("UpdateTxConfirmation failed %w", err)
	}
	resp.Status = txStatus

	if req.MessageID != "" {
		db2.DB.RedisDB.RemoveUnconfirmedOrder(req.MessageID, uint8(req.CoinType))
	}

	//record the addresses to update the account information
	if sentLog != nil && sentLog.Txid != "" {
		err = db2.DB.RedisDB.InsertAddressToUpdate(ethDetail.SenderAddress)
		if err != nil {
			log.NewError(req.OperationID, "the sender address is from other wallet application!")
		}
	}
	if receiveAccount1 != nil {
		err = db2.DB.RedisDB.InsertAddressToUpdate(ethDetail.ReceiverAddress)
		if err != nil {
			log.NewError(req.OperationID, "Failed to insert transaction addresses to redis stream")
		}
	}

	//push notifications
	go common.PushMsg(req.OperationID, senderMerchantUid, receiveAccount1.MerchantUid, ethDetail.SenderAddress, ethDetail.ReceiverAddress, ethDetail.SentHashTX, transactAmount, uint32(txStatus), req.CoinType)

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetConfirmation success!")
	return resp, nil
}

func (s *ethRPCServer) TransferRPCV2(ctx context.Context, req *eth.PostTransferReq2) (*eth.TransferRPCResponse2, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TransferRPC!!!", req.String())

	var (
		resp = &eth.TransferRPCResponse2{}
	)

	ethClient, err := coingrpc.GetETHInstance()
	if err != nil {
		return resp, fmt.Errorf("GetETHInstance failed %w", err)
	}

	// // To do: Worker Pool implementation may require in later stage, if too many request.
	// rawTx, ethDetail, err := ethClient.CreateRawTransaction(req.FromAddress, req.ToAddress, uint64(req.Amount), int(req.Nounce))
	// if err != nil {
	// 	return resp, fmt.Errorf("CreateRawTransaction failed %w", err)
	// }

	// log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Completed Stage 1!!!", req.String())
	// ethDetail.UnsignedUpdatedAt = time.Now()
	// ethDetail.SentUpdatedAt = time.Now()

	// resp.Transaction = &eth.EthTransactionDetail{
	// 	TransactionID:          uint64(ethDetail.TXID),
	// 	UUID:                   ethDetail.UUID,
	// 	CurrentTransactionType: int32(ethDetail.CurrentTXType),
	// 	SenderAccount:          ethDetail.SenderAccount,
	// 	SenderAddress:          ethDetail.SenderAddress,
	// 	ReceiverAccount:        ethDetail.ReceiverAccount,
	// 	ReceiverAddress:        ethDetail.ReceiverAddress,
	// 	Amount:                 ethDetail.Amount,
	// 	Fee:                    ethDetail.Fee,
	// 	GasLimit:               uint64(ethDetail.GasLimit),
	// 	Nonce:                  ethDetail.Nonce,
	// 	UnsignedHexTX:          ethDetail.UnsignedHexTX,
	// 	SignedHexTX:            ethDetail.SignedHexTX,
	// 	SentHashTX:             ethDetail.SentHashTX,
	// 	UnsignedUpdatedAt:      uint64(ethDetail.UnsignedUpdatedAt.Unix()),
	// 	SentUpdatedAt:          uint64(ethDetail.SentUpdatedAt.Unix()),
	// }
	// Added temp flag for testing
	// if req.Flag == 1 {
	// 	// SignOnRawTransaction sign the transaction with private key and password.
	// 	rawTx, err = ethClient.SignOnRawTransaction(rawTx, req.Secret)
	// 	if err != nil {
	// 		return resp, fmt.Errorf("SignOnRawTransaction failed %w", err)
	// 	}
	// 	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Completed Stage 2!!!", req.String())
	// } else {
	// 	// SignOnRawTransaction sign the transaction with private key and password.
	// 	rawTx, err = ethClient.SignOnRawTransactionV2(rawTx, req.Pkey)
	// 	if err != nil {
	// 		return resp, fmt.Errorf("SignOnRawTransaction failed %w", err)
	// 	}
	// 	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Completed Stage 2!!!", req.String())
	// }

	// SendSignedRawTransaction
	txHash, err := ethClient.SendSignedRawTransaction(req.Rawhex)
	if err != nil {
		return resp, fmt.Errorf("SendSignedRawTransaction failed %w", err)
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Completed Stage 3!!!", req.String())
	// err = walletdb.AddEthTransaction(ethDetail)
	// if err != nil {
	// 	return resp, fmt.Errorf("AddEthTransaction failed %w", err)
	// }
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Completed Stage 4!!!", req.String())
	// resp.Transaction.SignedHexTX = rawTx.TxHex
	resp.TranHash = txHash
	return resp, nil
}
func Wei2Eth_float(amount *big.Int) float64 {
	compact_amount := big.NewInt(0)
	reminder := big.NewInt(0)
	divisor := big.NewInt(1e18)
	compact_amount.QuoRem(amount, divisor, reminder)
	x := fmt.Sprintf("%v.%018s", compact_amount.String(), reminder.String())
	floatAmount, _ := strconv.ParseFloat(x, 64)
	return floatAmount
}

func (s *ethRPCServer) CheckBalanceAndGetNonceRPC(_ context.Context, req *eth.CheckBalanceAndGetNonceReq) (*eth.CheckBalanceAndGetNonceResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetEthBalanceRPC!!", req.String())

	if req.FromAddress == "" {
		resp := &eth.CheckBalanceAndGetNonceResp{
			ErrCode: constant.ErrArgs.ErrCode,
			ErrMsg:  "address is nil!",
		}
		return resp, nil
	}

	ethClient, err := coingrpc.GetETHInstance()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetETHInstance failed!", err.Error())
		return nil, err
	}

	//transact amount
	transactDecimal, _ := decimal.NewFromString(req.TransactAmount)
	transactAmount := transactDecimal.BigInt()
	if transactAmount.Uint64() == 0 {
		resp := &eth.CheckBalanceAndGetNonceResp{
			ErrCode: constant.ErrTransactAmountZero.ErrCode,
			ErrMsg:  constant.ErrTransactAmountZero.ErrMsg,
		}
		return resp, nil
	}

	// txFee := gasPrice * estimatedGas
	gasPriceDecimal, _ := decimal.NewFromString(req.GasPrice)
	gasLimitDecimal, _ := decimal.NewFromString(req.GasLimit)
	txFee := new(big.Int).Mul(gasPriceDecimal.BigInt(), gasLimitDecimal.BigInt())

	//ETH balance
	ethBalance, err1 := ethClient.BalanceAt(req.FromAddress)
	if err1 != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ethClient get ETH balance failed!", err1.Error())
		return nil, err1
	}
	if ethBalance == nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ethBalance is nil!")
		resp := &eth.CheckBalanceAndGetNonceResp{
			ErrCode: constant.ErrEthBalanceNil.ErrCode,
			ErrMsg:  constant.ErrEthBalanceNil.ErrMsg,
		}
		return resp, nil
	}
	if ethBalance.Uint64() == 0 {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ethBalance is zero!")
		resp := &eth.CheckBalanceAndGetNonceResp{
			ErrCode: constant.ErrEthBalanceZero.ErrCode,
			ErrMsg:  constant.ErrEthBalanceZero.ErrMsg,
		}
		return resp, nil
	}

	//Token balance
	tokenBalance, err2 := ethClient.GetTokenBalance(req.FromAddress)
	if req.CoinType != constant.ETHCoin {
		//transact USDT-ERC20
		if err2 != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "ethClient get token balance failed!", err2.Error())
			return nil, err2
		}

		if tokenBalance == nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "tokenBalance is nil!")
			resp := &eth.CheckBalanceAndGetNonceResp{
				ErrCode: constant.ErrUSDTERC20BalanceNil.ErrCode,
				ErrMsg:  constant.ErrUSDTERC20BalanceNil.ErrMsg,
			}
			return resp, nil
		}

		if tokenBalance.Uint64() == 0 {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "tokenBalance is zero!")
			resp := &eth.CheckBalanceAndGetNonceResp{
				ErrCode: constant.ErrEthBalanceLessThanFee.ErrCode,
				ErrMsg:  constant.ErrEthBalanceLessThanFee.ErrMsg,
			}
			return resp, nil
		}

		if tokenBalance.Cmp(transactAmount) == -1 {
			resp := &eth.CheckBalanceAndGetNonceResp{
				ErrCode: constant.ErrUSDTERC20BalanceNotEnough.ErrCode,
				ErrMsg:  constant.ErrUSDTERC20BalanceNotEnough.ErrMsg,
			}
			return resp, nil
		}

		// txFee := gasPrice * estimatedGas
		if ethBalance.Cmp(txFee) == -1 {
			networkFee, _ := utils.ConvertWeiToEther(txFee).Float64()
			errorMsg := fmt.Sprintf("The amount of ETH couldn’t less than transaction fee, around %f ETH", networkFee)
			resp := &eth.CheckBalanceAndGetNonceResp{
				ErrCode: constant.ErrEthBalanceLessThanFee.ErrCode,
				ErrMsg:  errorMsg,
			}
			return resp, errors.Wrap(constant.ErrEthBalanceLessThanFee, fmt.Sprintf("%v", constant.ErrEthBalanceLessThanFee.ErrCode))
		}
	} else {
		//transact ETH

		if ethBalance.Cmp(transactAmount) == -1 {
			resp := &eth.CheckBalanceAndGetNonceResp{
				ErrCode: constant.ErrEthBalanceNotEnough.ErrCode,
				ErrMsg:  constant.ErrEthBalanceNotEnough.ErrMsg,
			}
			return resp, nil
		}

		totalAmount := new(big.Int)
		totalAmount = totalAmount.Add(transactAmount, txFee)
		if ethBalance.Cmp(totalAmount) == -1 {
			networkFee, _ := utils.ConvertWeiToEther(txFee).Float64()
			errorMsg := fmt.Sprintf("Your balance is no enough to pay the transaction fee, around %f ETH", networkFee)
			resp := &eth.CheckBalanceAndGetNonceResp{
				ErrCode: constant.ErrEthBalanceLessThanFee.ErrCode,
				ErrMsg:  errorMsg,
			}
			return resp, errors.Wrap(constant.ErrEthBalanceLessThanFee, fmt.Sprintf("%v", constant.ErrEthBalanceLessThanFee.ErrCode))
		}
	}

	// nonce
	nonce, err := ethClient.GetNonce(req.FromAddress, 0)
	if err != nil {
		return nil, err
	}

	//net id
	netID, err := ethClient.NetVersion()
	if err != nil {
		return nil, err
	}

	resp := &eth.CheckBalanceAndGetNonceResp{
		ErrCode:                  0,
		ErrMsg:                   "request success",
		Nonce:                    nonce,
		ChainID:                  strconv.Itoa(int(netID)),
		USDTERC20ContractAddress: ethClient.GetUSDTERC20ContractAddress(),
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "nonce:", nonce, "chainID", netID)
	return resp, nil
}
