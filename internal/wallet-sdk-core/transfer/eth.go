package transfer

import (
	"Share-Wallet/pkg/coingrpc/ethgrpc/eth"
	"Share-Wallet/pkg/coingrpc/ethgrpc/ethtx"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/struct/sdk"
	"Share-Wallet/pkg/struct/wallet_api"
	"Share-Wallet/pkg/utils"
	"fmt"
	"log"
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func FromFloatEther(v float64) *big.Int {
	weiAmount := utils.ToWei(v, 18)
	return weiAmount
}

func (t *Transfer) dealETHTransaction(coinType int, fromAddress, toAddress string, pkey string, amount, gasPrice float64) (string, error) {
	operationID := utils.OperationIDGenerator()

	localUser, err := t.Db.GetUserByUserID(t.LoginUserID)
	if err != nil {
		log.Println(fmt.Errorf("error in GetUserByUserID() %w", err))
		return "", errors.New("It is failed in getting user information!")
	}
	fromAccountUID := localUser.PublicKey
	fromMerchantUID := localUser.UserID

	var weiAmount *big.Int
	var gasLimit *big.Int
	if coinType == constant.ETHCoin {
		weiAmount = FromFloatEther(amount)
		gasLimit = big.NewInt(constant.ETHGasLimit)
	} else if coinType == constant.USDTERC20 {
		weiAmount = big.NewInt(int64(amount * 1000000))
		gasLimit = big.NewInt(constant.USDTERC20GasLimit)
	}

	//checking balance and return nonce
	weiGasPrice := FromFloatEther(gasPrice)
	req := wallet_api.CheckBalanceAndNonceRequest{
		OperationID:    operationID,
		CoinType:       uint32(coinType),
		FromAddress:    fromAddress,
		TransactAmount: weiAmount.String(),
		GasPrice:       weiGasPrice.String(),
		GasLimit:       gasLimit.String(),
	}
	resp, err := t.API.PostWalletAPI(constant.CheckBalanceAndNonceURL, req, constant.APITimeout)
	if err != nil {
		log.Println(fmt.Errorf("error in CheckBalanceAndNonce() %w", err))
		return "", errors.Wrap(constant.ErrCreateTronTransaction, fmt.Sprintf("%v", constant.ErrCreateTronTransaction.ErrCode))
	}
	var respObj sdk.CheckBalanceAndNonceResp
	err = utils.JsonStringToStruct(string(resp), &respObj)
	if err != nil {
		log.Println(fmt.Errorf("CheckBalanceAndNonceResp error in JsonStringToStruct() %w", err, string(resp)))
		return "", err
	}
	if int32(respObj.Code) != constant.OK.ErrCode {
		log.Println(fmt.Errorf("error in CheckBalanceAndNonce %w", respObj))
		return "", errors.New(respObj.ErrMsg)
	}
	nonce := respObj.Data.Nonce
	chainID := respObj.Data.ChainID
	usdterc20Contract := respObj.Data.USDTERC20ContractAddress

	//create raw transaction
	var rawTx *ethtx.RawTx
	var ethClient eth.Ethereum
	if req.CoinType == constant.ETHCoin {
		rawTx, err = ethClient.CreateRawTransactionLocal(fromAddress, toAddress, weiAmount, nonce, weiGasPrice)
	} else {
		rawTx, err = ethClient.CreateTokenRawTransactionLocal(fromAddress, toAddress, weiAmount, nonce, weiGasPrice, gasLimit, usdterc20Contract)
	}
	if err != nil {
		log.Println(req.OperationID, utils.GetSelfFuncName(), "CreateRawTransaction failed", err.Error())
		return "", err
	}

	//sign raw transaction
	// SignOnRawTransaction sign the transaction with private key and password.
	chainIDDecimal, _ := decimal.NewFromString(chainID)
	rawTx, err = ethClient.SignOnRawTransactionV2(rawTx, pkey, chainIDDecimal.BigInt())
	if err != nil {
		log.Println(req.OperationID, utils.GetSelfFuncName(), "SignOnRawTransactionV2 failed", err.Error())
		return "", fmt.Errorf("SignOnRawTransaction failed %w", err)
	}

	//send raw transaction
	txFee := new(big.Int)
	txFee = txFee.Mul(gasLimit, weiGasPrice)
	req2 := wallet_api.PostTransferRequest{
		OperationID:     operationID,
		CoinType:        uint32(coinType),
		FromAccountUID:  fromAccountUID,
		FromMerchantUID: fromMerchantUID,
		FromAddress:     fromAddress,
		ToAddress:       toAddress,
		Amount:          weiAmount.String(),
		Fee:             txFee.String(),
		GasLimit:        gasLimit.Uint64(),
		Nonce:           nonce,
		GasPrice:        weiGasPrice.String(),
		TxHash:          rawTx.TxHex,
	}

	log.Println("Amount", decimal.NewFromBigInt(weiAmount, 0), "gasFee", decimal.NewFromBigInt(weiGasPrice, 0), "txfee", txFee.String())
	resp2, err := t.API.PostWalletAPI(constant.ETHTransferAccountURL, req2, constant.APITimeout)
	if err != nil {
		log.Println(fmt.Errorf("error in PostWalletAPI() %w", err))
		return "", err
	}
	type Transaction struct {
		TransactionID          uint64 `json:"transactionID"`
		UUID                   string `json:"uuID"`
		CurrentTransactionType int32  `json:"current_transaction_type"`
		SenderAccount          string `json:"sender_account"`
		SenderAddress          string `json:"sender_address"`
		ReceiverAccount        string `json:"receiver_account"`
		ReceiverAddress        string `json:"receiver_address"`
		Amount                 string `json:"amount"`
		Fee                    string `json:"fee"`
		GasLimit               uint64 `json:"gas_limit"`
		Nonce                  uint64 `json:"nonce"`
		UnsignedHexTX          string `json:"unsigned_hex_tx"`
		SignedHexTX            string `json:"signed_hex_tx"`
		SentHashTX             string `json:"sent_hash_tx"`
		UnsignedUpdatedAt      uint64 `json:"unsigned_updated_at"`
		SentUpdatedAt          uint64 `json:"sent_updated_at"`
	}

	type Data struct {
		OperationID string      `json:"operationID"`
		Trans       Transaction `json:"transaction"`
	}

	type Obj struct {
		Code   int    `json:"code"`
		ErrMsg string `json:"err_msg"`
		Data   Data   `json:"data"`
	}

	var obj Obj
	err = utils.JsonStringToStruct(string(resp2), &obj)
	if err != nil {
		log.Println(fmt.Errorf("transfer error in JsonStringToStruct() %w", err, string(resp2)))
		return "", err
	}

	if obj.ErrMsg != "" {
		return obj.Data.Trans.SentHashTX, errors.New(obj.ErrMsg)
	}

	return obj.Data.Trans.SentHashTX, nil
}
