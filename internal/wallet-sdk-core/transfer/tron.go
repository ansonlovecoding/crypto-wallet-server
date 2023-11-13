package transfer

import (
	"Share-Wallet/pkg/coingrpc/trongrpc/tron"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/struct/sdk"
	"Share-Wallet/pkg/struct/wallet_api"
	"Share-Wallet/pkg/utils"
	"fmt"
	"log"
	"math/big"

	"github.com/fbsobreira/gotron-sdk/pkg/common"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"

	"google.golang.org/protobuf/proto"

	"github.com/pkg/errors"
)

func (t *Transfer) dealTronTransaction(coinType int, fromAddress, toAddress string, pkey string, amount float64) (string, error) {
	operationID := utils.OperationIDGenerator()

	localUser, err := t.Db.GetUserByUserID(t.LoginUserID)
	if err != nil {
		log.Println(fmt.Errorf("error in GetUserByUserID() %w", err))
		return "", errors.New("It is failed in getting user information!")
	}
	fromAccountUID := localUser.PublicKey
	fromMerchantUID := localUser.UserID

	var sunAmount *big.Int
	if coinType == constant.TRX {
		sunAmount = utils.TrxToSun(amount)
	} else if coinType == constant.USDTTRC20 {
		sunAmount = big.NewInt(int64(amount * 1000000))
	}

	//create raw transaction
	var txRaw core.Transaction
	var createTxResp *wallet_api.CreateTransactionResponse
	var tronClient tron.Tron

	createTxResp, err = t.createTronTransaction(operationID, coinType, fromAddress, toAddress, sunAmount.String())
	if err != nil {
		log.Println(operationID, utils.GetSelfFuncName(), "createTronTransaction failed", err.Error())
		return "", err
	}
	if createTxResp == nil {
		log.Println(operationID, utils.GetSelfFuncName(), "createTronTransaction failed", err.Error())
		return "", errors.New("createTronTransaction failed, return nil response")
	}
	byteResp, err := common.Hex2Bytes(createTxResp.RawTxData)
	if err != nil {
		log.Println(operationID, utils.GetSelfFuncName(), "common.Hex2Bytes(createTxResp.RawTxData) failed", err.Error())
		return "", err
	}
	if byteResp == nil {
		log.Println(operationID, utils.GetSelfFuncName(), "byteResp is nil", err.Error())
		return "", errors.New("convert string to bytes failed")
	}
	//err = utils.JsonStringToStruct(createTxResp.RawTxData, txRaw)
	err = proto.Unmarshal(byteResp, &txRaw)
	if err != nil {
		log.Println(operationID, utils.GetSelfFuncName(), "proto.Unmarshal failed", err.Error(), "createTxResp.RawTxData", createTxResp.RawTxData)
		return "", err
	}

	//sign raw transaction
	// SignTransactionLocal sign the transaction with private key.
	signedTx, err := tronClient.SignTransactionLocal(&txRaw, pkey)
	if err != nil {
		log.Println(operationID, utils.GetSelfFuncName(), "SignTransactionLocal failed", err.Error())
		return "", fmt.Errorf("SignTransactionLocal failed %w", err)
	}

	orderID := createTxResp.TxID
	//send raw transaction
	txByte, err := proto.Marshal(signedTx)
	txByteStr := common.BytesToHexString(txByte)

	req2 := wallet_api.TransferTronRequest{
		OperationID:     operationID,
		CoinType:        uint32(coinType),
		FromAccountUID:  fromAccountUID,
		FromMerchantUID: fromMerchantUID,
		FromAddress:     fromAddress,
		ToAddress:       toAddress,
		Amount:          sunAmount.String(),
		TxID:            orderID,
		TxData:          txByteStr,
	}

	resp2, err := t.API.PostWalletAPI(constant.TronTransferAccountURL, req2, constant.APITimeout)
	if err != nil {
		log.Println(fmt.Errorf("error in PostWalletAPI() %w", err))
		return "", err
	}
	log.Println("string(resp2)", utils.Bytes2String(resp2))
	type Obj struct {
		Code   int    `json:"code"`
		ErrMsg string `json:"err_msg"`
	}

	var obj Obj
	err = utils.JsonStringToStruct(string(resp2), &obj)
	if err != nil {
		log.Println(fmt.Errorf("transfer error in JsonStringToStruct() %w", err, string(resp2)))
		return "", err
	}

	if obj.ErrMsg != "" {
		//log.Println("check error", obj.Code, errors.Wrap(constant.ErrInfo{ErrCode: int32(obj.Code), ErrMsg: obj.ErrMsg}, fmt.Sprintf("%v", obj.Code)).Error())
		return orderID, errors.New(obj.ErrMsg)
	}

	return orderID, nil
}

func (t *Transfer) createTronTransaction(operationID string, coinType int, fromAddress, toAddress, amount string) (*wallet_api.CreateTransactionResponse, error) {
	req2 := wallet_api.CreateTransactionRequest{
		OperationID: operationID,
		CoinType:    uint32(coinType),
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
	}

	resp2, err := t.API.PostWalletAPI(constant.CreateTronTransactionURL, req2, constant.APITimeout)
	if err != nil {
		log.Println(fmt.Errorf("error in PostWalletAPI() %w", err))
		return nil, errors.Wrap(constant.ErrCreateTronTransaction, fmt.Sprintf("%v", constant.ErrCreateTronTransaction.ErrCode))
	}

	var obj sdk.CreateTronTransactionResp
	err = utils.JsonStringToStruct(string(resp2), &obj)
	if err != nil {
		log.Println(fmt.Errorf("create transaction error in JsonStringToStruct() %w", err, string(resp2)))
		return nil, err
	}

	if obj.ErrMsg != "" {
		return nil, errors.New(obj.ErrMsg)
	}

	return &obj.Data, nil
}
