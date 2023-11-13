package transfer

import (
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/http"
	db "Share-Wallet/pkg/db/local_db"
	sdkstruct "Share-Wallet/pkg/sdk_struct"
	"Share-Wallet/pkg/struct/wallet_api"
	"Share-Wallet/pkg/utils"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/params"
	"github.com/shopspring/decimal"
)

// Transfer struct
type Transfer struct {
	Db               *db.DataBase
	LoginUserID      string
	LoginTime        int64
	API              *http.PostAPI
	UserRandomSecret []byte
	SupportTokens    []*wallet_api.SupportTokenAddress
}

// Function receiver param
var TransferMgr *Transfer

func NewTransfer(dataBase *db.DataBase, loginUserID string, loginTime int32, p *http.PostAPI, userRandomSecret []byte, supportTokens []*wallet_api.SupportTokenAddress) *Transfer {

	return &Transfer{Db: dataBase,
		LoginUserID:      loginUserID,
		LoginTime:        int64(loginTime),
		API:              p,
		UserRandomSecret: userRandomSecret,
		SupportTokens:    supportTokens,
	}
}

func (t *Transfer) Transferfn(coinType int, fromAddress, toAddress string, secret string, amount, gasPrice float64) (string, error) {

	wallet, err := t.Db.GetLocalWalletByUserID(t.LoginUserID, coinType)
	if err != nil {
		log.Println(fmt.Errorf("error in GetLocalWalletByUserID() %w", err))
		return "", errors.New("It is failed in getting wallet information")
	}

	key := fmt.Sprintf("%s%s", secret, t.UserRandomSecret)
	pkey, err := utils.DecryptAES(wallet.WalletImportFormat, key)
	if err != nil {
		log.Println(fmt.Errorf("error in DecryptAES() %w", err))
		return "", errors.New("It is failed in getting the private key!")
	}

	//deal transaction
	var txHash string
	if coinType == constant.ETHCoin || coinType == constant.USDTERC20 {
		if !utils.IsValidAddress(fromAddress) || utils.IsZeroAddress(fromAddress) {
			log.Println(fmt.Errorf("The sending address is incorrect!"))
			return "", errors.Wrap(constant.ErrSendingAddressIncorrect, fmt.Sprintf("%v", constant.ErrSendingAddressIncorrect.ErrCode))
		}
		if !utils.IsValidAddress(toAddress) || utils.IsZeroAddress(toAddress) {
			log.Println(fmt.Errorf("The receiving address is incorrect!"))
			return "", errors.Wrap(constant.ErrReceiverAddressIncorrect, fmt.Sprintf("%v", constant.ErrReceiverAddressIncorrect.ErrCode))
		}

		return t.dealETHTransaction(coinType, fromAddress, toAddress, pkey, amount, gasPrice)
	} else if coinType == constant.TRX || coinType == constant.USDTTRC20 {
		if !utils.IsValidTRXAddress(fromAddress) {
			log.Println(fmt.Errorf("The sending address is incorrect!"))
			return "", errors.Wrap(constant.ErrSendingAddressIncorrect, fmt.Sprintf("%v", constant.ErrSendingAddressIncorrect.ErrCode))
		}
		if !utils.IsValidTRXAddress(toAddress) {
			log.Println(fmt.Errorf("The receiving address is incorrect!"))
			return "", errors.Wrap(constant.ErrReceiverAddressIncorrect, fmt.Sprintf("%v", constant.ErrReceiverAddressIncorrect.ErrCode))
		}
		return t.dealTronTransaction(coinType, fromAddress, toAddress, pkey, amount)
	}

	return txHash, nil
}

func (t *Transfer) GetBalance(coinType int, publicAddress string) string {
	operationID := utils.OperationIDGenerator()

	//coinType: 1 BTC, 2 ETH, 3 USDT-ERC20, 4 TRX, 5 USDT-TRC20
	req := sdkstruct.GetBalance{
		CoinType:    coinType,
		Address:     publicAddress,
		OperationID: operationID,
	}

	resp, err := t.API.PostWalletAPI(constant.GetAccountBalanceURL, req, constant.APITimeout)
	if err != nil {
		log.Println(fmt.Errorf("error in PostWalletAPI() %w", err))
		return ""
	}

	type Balance struct {
		Balance string `json:"balance"`
	}
	type Obj struct {
		Code   int     `json:"code"`
		ErrMsg string  `json:"err_msg"`
		Data   Balance `json:"data"`
	}
	var obj Obj
	err = utils.JsonStringToStruct(string(resp), &obj)
	if err != nil {
		return ""
	}

	if obj.Data.Balance == "" {
		return "0"
	}
	if coinType == constant.ETHCoin {
		log.Println("obj.Data.Balance", obj.Data.Balance)
		balanceInt := new(big.Int)
		balanceInt, _ = balanceInt.SetString(obj.Data.Balance, 10)
		log.Println("balanceInt", balanceInt)
		bal := Wei2Eth_str(balanceInt)
		log.Println("bal", bal)
		balStr := utils.Float64WithoutRound(bal, 8)
		log.Println("balStr", balStr)
		return balStr
	} else {
		return obj.Data.Balance
	}

}
func (t *Transfer) GetGasPrice(coinType int) string {
	operationID := utils.OperationIDGenerator()
	req := sdkstruct.GetGasPrice{
		CoinType:    coinType,
		IsEstimated: false,
		OperationID: operationID, // Remove after new requestID injection from middleware
	}
	resp, err := t.API.PostWalletAPI(constant.GetGasPriceURL, req, constant.APITimeout)
	if err != nil {
		log.Println(fmt.Errorf("error in PostWalletAPI() %w", err))
		return ""
	}
	type Balance struct {
		GasPrice *big.Int `json:"gas_price"`
	}
	type Obj struct {
		Code   int     `json:"code"`
		ErrMsg string  `json:"err_msg"`
		Data   Balance `json:"data"`
	}
	var obj Obj
	err = utils.JsonStringToStruct(string(resp), &obj)
	if err != nil {
		return ""
	}

	var formattedGasprice string
	switch coinType {
	case constant.ETHCoin, constant.USDTERC20:
		log.Println("obj.Data.GasPrice", obj.Data.GasPrice)
		if obj.Data.GasPrice != nil {
			gasPrice := Wei2Eth_str(obj.Data.GasPrice)
			formattedGasprice = strings.TrimRight(gasPrice, "0")
		}

	default:
		if obj.Data.GasPrice != nil {

			formattedGasprice = obj.Data.GasPrice.String()
		}
	}

	return formattedGasprice

}

// no use, deprecated
/*
func (t *Transfer) GetTransaction(coinType int, transactionHash string) string {

	// To do: Refactor
	req := sdkstruct.GetTransaction{
		CoinType:        coinType,
		TransactionHash: transactionHash,
	}

	resp, err := t.API.PostWalletAPI(constant.GetTransactionURL, req, constant.APITimeout)
	if err != nil {
		fmt.Println(err)
		log.Println(fmt.Errorf("error in PostWalletAPI() %w", err))
		return ""
	}

	type Transaction struct {
		BlockHash        string   `json:"blockHash"`
		BlockNumber      int64    `json:"blockNumber"`
		From             string   `json:"from"`
		Gas              int64    `json:"gas"`
		GasPrice         *big.Int `json:"gasPrice"`
		Hash             string   `json:"hash"`
		Input            string   `json:"-"`
		Nonce            int64    `json:"-"`
		To               string   `json:"to"`
		TransactionIndex int64    `json:"-"`
		Value            *big.Int `json:"value"`
		V                int64    `json:"-"`
		R                string   `json:"-"`
		S                string   `json:"-"`
		ValueETH         float64  `json:"value_conv"`
		GasPriceETH      float64  `json:"gas_price_conv"`
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
	err = utils.JsonStringToStruct(string(resp), &obj)
	if err != nil {
		return ""
	}

	amt := Wei2Eth_str(obj.Data.Trans.Value)
	amtFloat, _ := strconv.ParseFloat(amt, 64)
	obj.Data.Trans.ValueETH = amtFloat

	gasPrice := Wei2Eth_str(obj.Data.Trans.GasPrice)
	gasPriceFloat, _ := strconv.ParseFloat(gasPrice, 64)
	obj.Data.Trans.ValueETH = gasPriceFloat

	// confirmationdata := t.GetConfirmation(coinType, transactionHash)
	// err = utils.JsonStringToStruct(string(confirmationdata), &obj.Data.Confirmation)
	// if err != nil {
	// 	return ""
	// }
	return utils.StructToJsonString(obj.Data)

}

*/

func (t *Transfer) GetConfirmation(coinType int, transactionHash string) string {

	req := sdkstruct.GetTransaction{
		CoinType:        coinType,
		TransactionHash: transactionHash,
	}

	var getEthConfirmResp wallet_api.GetETHConfirmationResponseObj
	var getTronConfirmResp wallet_api.GetTronConfirmationResponseObj
	var resp []byte
	var err error
	if coinType == constant.ETHCoin || coinType == constant.USDTERC20 {
		resp, err = t.API.PostWalletAPI(constant.GetETHTransactionConfirmationURL, req, constant.APITimeout)
		if err != nil {
			log.Println(fmt.Errorf("error in PostWalletAPI() %w", err))
			return ""
		}
		err = utils.JsonStringToStruct(string(resp), &getEthConfirmResp)
		if err != nil {
			return ""
		}
		return utils.StructToJsonString(getEthConfirmResp.Block)
	} else if coinType == constant.TRX || coinType == constant.USDTTRC20 {
		resp, err = t.API.PostWalletAPI(constant.GetTronTransactionConfirmationURL, req, constant.APITimeout)
		if err != nil {
			log.Println(fmt.Errorf("error in PostWalletAPI() %w", err))
			return ""
		}
		err = utils.JsonStringToStruct(string(resp), &getTronConfirmResp)
		if err != nil {
			return ""
		}
		return utils.StructToJsonString(getTronConfirmResp.Block)
	}
	return ""

}
func (t *Transfer) GetPublicAddress(coinType int) string {
	userID := t.LoginUserID
	publicAddress, err := t.Db.GetPublicAddressByUserID(userID, coinType)
	if err != nil {
		log.Println(fmt.Errorf("error in GetUserByUserID() %w", err))
		return ""
	}
	return utils.StructToJsonString(publicAddress)

}

func (t *Transfer) GetTransactionList(coinType int, publicAddress string, transactionType int, page int, pageSize int, orderBy string, transactionHash string) string {
	operationID := utils.OperationIDGenerator()
	userID := t.LoginUserID

	//publicAddress := "0x2ee5280e641e7c3533b5b513a00f12e275a71242"
	req := sdkstruct.GetTransactionList{
		CoinType:        coinType,
		UserID:          userID,
		Address:         publicAddress,
		TransactionType: transactionType,
		OperationID:     operationID,
		Page:            page,
		PageSize:        pageSize,
		OrderBy:         orderBy,
		TransactionHash: transactionHash,
	}
	resp, err := t.API.PostWalletAPI(constant.GetTransactionListURL, req, constant.APITimeout)
	if err != nil {
		log.Println(fmt.Errorf("error in PostWalletAPI() %w", err))
		return ""
	}
	type TransactionInfo struct {
		TransactionID          uint64                    `json:"transactionID"`
		UUID                   string                    `json:"uuID"`
		CurrentTransactionType int32                     `json:"current_transaction_type"`
		SenderAddress          string                    `json:"sender_address"`
		Sender                 *sdkstruct.GetAddressBook `json:"sender_account"`
		Receiver               *sdkstruct.GetAddressBook `json:"receiver_account"`
		ReceiverAddress        string                    `json:"receiver_address"`
		Amount                 string                    `json:"amount"`
		Fee                    string                    `json:"fee"`
		ConfirmationTime       uint64                    `json:"confirm_time"`
		TransactionHash        string                    `json:"transaction_hash"`
		SentTime               uint64                    `json:"sent_time"`
		Status                 int8                      `json:"status"`
		AmountFloat            float64                   `json:"amount_conv"`
		FeeFloat               float64                   `json:"fee_conv"`
		GasPriceFloat          decimal.Decimal           `json:"gas_price_conv"`
		IsSend                 bool                      `json:"is_send"`
		GasUsed                float64                   `json:"gas_used"`
		GasLimit               float64                   `json:"gas_limit"`
		GasPrice               string                    `json:"gas_price"`
		ConfirmBlockNumber     string                    `json:"confirm_block_number"`
	}
	type Data struct {
		OperationID    string            `json:"operationID"`
		Transaction    []TransactionInfo `json:"transaction"`
		TransactionNum uint64            `json:"tran_nums"`
		Page           uint64            `json:"page"`
		PageSize       uint64            `json:"page_size"`
	}
	type Obj struct {
		Code   int    `json:"code"`
		ErrMsg string `json:"err_msg"`
		Data   Data   `json:"data"`
	}
	var obj Obj

	utils.JsonStringToStruct(string(resp), &obj)
	for i, v := range obj.Data.Transaction {
		amountDecimal, _ := decimal.NewFromString(v.Amount)
		// Amount converted for ETH
		if coinType == constant.ETHCoin {

			amt := Wei2Eth_str(amountDecimal.BigInt())
			amtFloat, _ := strconv.ParseFloat(amt, 64)
			s := fmt.Sprintf("%.8f", amtFloat)

			amtFloatRounded, _ := strconv.ParseFloat(s, 64)
			obj.Data.Transaction[i].AmountFloat = amtFloatRounded

		} else if coinType == constant.TRX {
			//convert amount to trx
			amt := utils.SunToTrx(amountDecimal.BigInt())
			obj.Data.Transaction[i].AmountFloat, _ = amt.Float64()

		} else if coinType == constant.USDTERC20 {
			// Convert Amount when other coins are introduced. USDT-ERC20 need to devide by 1000000
			obj.Data.Transaction[i].AmountFloat = utils.ConvertBigUSDT2Float(amountDecimal.BigInt())
		} else if coinType == constant.USDTTRC20 {
			// Convert Amount when other coins are introduced. USDT-TRC20 need to devide by 1000000
			obj.Data.Transaction[i].AmountFloat = utils.ConvertBigUSDT2Float(amountDecimal.BigInt())
		} else {
			obj.Data.Transaction[i].AmountFloat, _ = amountDecimal.Float64()
		}

		feeDecimal, _ := decimal.NewFromString(v.Fee)
		if coinType == constant.ETHCoin || coinType == constant.USDTERC20 {
			fee := Wei2Eth_str(feeDecimal.BigInt())
			feeFloat, _ := strconv.ParseFloat(fee, 64)
			s2 := fmt.Sprintf("%.15f", feeFloat)

			FeeFloatRounded, _ := strconv.ParseFloat(s2, 64)
			obj.Data.Transaction[i].FeeFloat = FeeFloatRounded

			gasDecimal, _ := decimal.NewFromString(v.GasPrice)
			gasPrice := Wei2Eth_str(gasDecimal.BigInt())
			n, _ := decimal.NewFromString(gasPrice)
			obj.Data.Transaction[i].GasPriceFloat = n
		} else if coinType == constant.TRX || coinType == constant.USDTTRC20 {
			feeTrx := utils.SunToTrx(feeDecimal.BigInt())
			obj.Data.Transaction[i].FeeFloat, _ = feeTrx.Float64()
			obj.Data.Transaction[i].GasPriceFloat = decimal.NewFromInt(0)
		}

		if v.SenderAddress == publicAddress {
			obj.Data.Transaction[i].IsSend = true
		}
		var senderObj, receiverObj *sdkstruct.GetAddressBook
		sender, _ := t.Db.GetLocalUserAddressBookbyAddress(userID, coinType, v.SenderAddress)
		if sender.Name != "" {
			var send sdkstruct.GetAddressBook
			send.Address = v.SenderAddress
			send.Name = sender.Name
			senderObj = &send
		}
		receiver, _ := t.Db.GetLocalUserAddressBookbyAddress(userID, coinType, v.ReceiverAddress)
		if receiver.Name != "" {
			var receive sdkstruct.GetAddressBook
			receive.Address = v.ReceiverAddress
			receive.Name = receiver.Name
			receiverObj = &receive
		}
		obj.Data.Transaction[i].Sender = senderObj
		obj.Data.Transaction[i].Receiver = receiverObj
		/*
			type GetConfirmationResponse struct {
				BlockNumber      uint64 `json:"block_num"`
				ConfirmationTime uint64 `json:"confirm_time"`
				Status           int8   `json:"status"`
				GasUsed          uint64 `json:"gas_used"`
			}
			var confirm GetConfirmationResponse
			respConfirmations := t.GetConfirmation(coinType, v.TransactionHash)
			utils.JsonStringToStruct(string(respConfirmations), &confirm)
			obj.Data.Transaction[i].Status = confirm.Status
			obj.Data.Transaction[i].ConfirmBlockNumber = uint64(confirm.BlockNumber)

		*/
		//}
	}
	if len(obj.Data.Transaction) == 0 {
		// To do: Move this logic to server side
		tran := []TransactionInfo{}
		obj.Data.Transaction = tran
	}

	return utils.StructToJsonString(obj.Data)
}

func Wei2Eth_str(amount *big.Int) string {
	compact_amount := big.NewInt(0)
	reminder := big.NewInt(0)
	divisor := big.NewInt(1e18)
	compact_amount.QuoRem(amount, divisor, reminder)
	return fmt.Sprintf("%v.%018s", compact_amount.String(), reminder.String())
}

// Wei   = 1
// GWei  = 1e9  (Giga)
// Ether = 1e18

// Wei :  1000000000000000000
// GWei:  1000000000
// Ether: 1
func FromWei(v int64) *big.Int {
	return big.NewInt(v * params.Wei)
}

func (t *Transfer) GetTransactionFee(coinType int64, gasPrice float64) string {
	// Eth Only
	gasLimit := 0.0
	switch coinType {
	case constant.ETHCoin:
		gasLimit = 21000
	case constant.USDTERC20:
		gasLimit = 70000
	}
	transactionFee := gasLimit * gasPrice
	s := fmt.Sprintf("%.10f", transactionFee)
	formmatedTransactionFee := strings.TrimRight(s, "0")
	return formmatedTransactionFee
}

func (t *Transfer) GetTransactionDetails(coinType int, publicAddress, transactionHash string) string {
	operationID := utils.OperationIDGenerator()
	userID := t.LoginUserID

	req := sdkstruct.GetTransaction{
		CoinType:        coinType,
		OperationID:     operationID,
		TransactionHash: transactionHash,
	}
	resp, err := t.API.PostWalletAPI(constant.GetTransactionURL, req, constant.APITimeout)
	if err != nil {
		log.Println(fmt.Errorf("error in GetTransactionURL() %w", err))
		return ""
	}

	var obj sdkstruct.GetTransactionObj

	err = utils.JsonStringToStruct(string(resp), &obj)
	if err != nil {
		log.Println(fmt.Errorf("error in GetTransactionURL() %w", err))
		return ""
	}

	log.Println("check transaction", obj.Code, obj.ErrMsg, obj.Data.Transaction.TransactionHash)
	// var confirm GetConfirmationResponse
	var details sdkstruct.GetTransactionDetailResponse
	if obj.Code == int(constant.OK.ErrCode) {
		transaction := obj.Data.Transaction
		amountDecimal, _ := decimal.NewFromString(transaction.Amount)
		var amtFloatRounded float64
		// Amount converted for ETH
		if coinType == constant.ETHCoin {
			amt := Wei2Eth_str(amountDecimal.BigInt())
			amtFloat, _ := strconv.ParseFloat(amt, 64)
			s := fmt.Sprintf("%.8f", amtFloat)

			amtFloatRounded, _ = strconv.ParseFloat(s, 64)

		} else if coinType == constant.TRX {
			//convert amount to trx
			amt := utils.SunToTrx(amountDecimal.BigInt())
			amtFloatRounded, _ = amt.Float64()

		} else if coinType == constant.USDTERC20 {
			// Convert Amount when other coins are introduced. USDT-ERC20 need to devide by 1000000
			amtFloatRounded = utils.ConvertBigUSDT2Float(amountDecimal.BigInt())
		} else if coinType == constant.USDTTRC20 {
			// Convert Amount when other coins are introduced. USDT-TRC20 need to devide by 1000000
			amtFloatRounded = utils.ConvertBigUSDT2Float(amountDecimal.BigInt())
		} else {
			amtFloatRounded, _ = amountDecimal.Float64()
		}

		feeDecimal, _ := decimal.NewFromString(transaction.Fee)
		var feeFloat float64
		var n decimal.Decimal
		if coinType == constant.ETHCoin || coinType == constant.USDTERC20 {
			fee := Wei2Eth_str(feeDecimal.BigInt())
			feeFloat, _ = strconv.ParseFloat(fee, 64)
			s2 := fmt.Sprintf("%.15f", feeFloat)

			FeeFloatRounded, _ := strconv.ParseFloat(s2, 64)
			feeFloat = FeeFloatRounded

			gasDecimal, _ := decimal.NewFromString(transaction.GasPrice)
			gasPrice := Wei2Eth_str(gasDecimal.BigInt())
			n, _ = decimal.NewFromString(gasPrice)
		} else if coinType == constant.TRX || coinType == constant.USDTTRC20 {
			feeTrx := utils.SunToTrx(feeDecimal.BigInt())
			feeFloat, _ = feeTrx.Float64()
		}

		var senderObj, receiverObj *sdkstruct.GetAddressBook
		sender, _ := t.Db.GetLocalUserAddressBookbyAddress(userID, coinType, transaction.SenderAddress)
		if sender.Name != "" {
			var send sdkstruct.GetAddressBook
			send.Address = transaction.SenderAddress
			send.Name = sender.Name
			send.CoinType = int(sender.CoinType)
			senderObj = &send
		}
		receiver, _ := t.Db.GetLocalUserAddressBookbyAddress(userID, coinType, transaction.ReceiverAddress)
		if receiver.Name != "" {
			var receive sdkstruct.GetAddressBook
			receive.Address = transaction.ReceiverAddress
			receive.Name = receiver.Name
			receive.CoinType = int(receiver.CoinType)
			receiverObj = &receive
		}
		if strings.ToLower(transaction.SenderAddress) == strings.ToLower(publicAddress) {
			details.IsSend = true
		} else {
			details.IsSend = false
		}
		details.TransactionID = transaction.TransactionID
		details.UUID = transaction.UUID
		details.ConfirmationTime = transaction.ConfirmationTime
		details.SentTime = transaction.SentTime
		details.TransactionHash = transaction.TransactionHash
		details.Status = transaction.Status
		details.Amount = transaction.Amount
		details.Fee = transaction.Fee
		details.AmountFloat = amtFloatRounded
		details.FeeFloat = feeFloat
		details.GasPriceFloat = n
		details.Sender = senderObj
		details.SenderAddress = transaction.SenderAddress
		details.Receiver = receiverObj
		details.ReceiverAddress = transaction.ReceiverAddress
		details.ConfirmationBlockNumber = transaction.ConfirmBlockNumber
		details.GasUsed = transaction.GasUsed
		details.GasLimit = transaction.GasLimit
	}

	return utils.StructToJsonString(details)
}

func (t *Transfer) GetRecentTransactions(page int, pageSize int) string {
	operationID := utils.OperationIDGenerator()
	userID := t.LoginUserID

	//publicAddress := "0x2ee5280e641e7c3533b5b513a00f12e275a71242"
	localUser, _ := t.Db.GetUserByUserID(userID)
	req := sdkstruct.GetRecentRecordsRequest{
		OperationID: operationID,
		UID:         localUser.PublicKey,
		Page:        page,
		PageSize:    pageSize,
	}

	resp, err := t.API.PostWalletAPI(constant.GetRecentRecordsURL, req, constant.APITimeout)
	if err != nil {
		log.Println(fmt.Errorf("error in PostWalletAPI() %w", err))
		return ""
	}
	log.Println("GetRecentTransactions resp:", string(resp))
	type FundLog struct {
		ID                      int64   `json:"id"`
		Txid                    string  `json:"txid"`
		Uid                     int64   `json:"uid"`
		MerchantUid             string  `json:"merchant_uid"`
		TransactionType         string  `json:"transaction_type"`
		UserAddress             string  `json:"user_address"`
		UserAddressName         string  `json:"user_address_name"`
		OppositeAddress         string  `json:"opposite_address"`
		OppositeAddressName     string  `json:"opposite_address_name"`
		CoinType                string  `json:"coin_type"`
		AmountOfCoins           float64 `json:"amount_of_coins"`
		UsdAmount               float64 `json:"usd_amount"`
		YuanAmount              float64 `json:"yuan_amount"`
		EuroAmount              float64 `json:"euro_amount"`
		NetworkFee              float64 `json:"network_fee"`
		UsdNetworkFee           float64 `json:"usd_network_fee"`
		YuanNetworkFee          float64 `json:"yuan_network_fee"`
		EuroNetworkFee          float64 `json:"euro_network_fee"`
		TotalCoinsTransfered    float64 `json:"total_coins_transfered"`
		TotalUsdTransfered      float64 `json:"total_usd_transfered"`
		TotalYuanTransfered     float64 `json:"total_yuan_transfered"`
		TotalEuroTransfered     float64 `json:"total_euro_transfered"`
		CreationTime            int64   `json:"creation_time"`
		State                   string  `json:"state"`
		ConfirmationTime        int64   `json:"confirmation_time"`
		GasLimit                uint32  `json:"gas_limit"`
		GasPrice                uint64  `json:"gas_price"`
		GasUsed                 uint64  `json:"gas_used"`
		ConfirmationBlockNumber uint64  `json:"confirm_block_number"`
	}
	type Data struct {
		FundsLog []*FundLog `json:"funds_log"`
		TotalNum uint64     `json:"total_num"`
		Page     uint64     `json:"page"`
		PageSize uint64     `json:"page_size"`
	}
	type Obj struct {
		Code   int    `json:"code"`
		ErrMsg string `json:"err_msg"`
		Data   Data   `json:"data"`
	}
	var obj Obj

	utils.JsonStringToStruct(string(resp), &obj)
	for i, v := range obj.Data.FundsLog {
		sender, _ := t.Db.GetLocalUserAddressBookbyAddress(userID, utils.GetCoinType(v.CoinType), v.UserAddress)
		if sender.Name != "" {
			obj.Data.FundsLog[i].UserAddressName = sender.Name
		}
		receiver, _ := t.Db.GetLocalUserAddressBookbyAddress(userID, utils.GetCoinType(v.CoinType), v.OppositeAddress)
		if receiver.Name != "" {
			obj.Data.FundsLog[i].UserAddressName = receiver.Name
		}
	}
	if len(obj.Data.FundsLog) == 0 {
		// To do: Move this logic to server side
		tran := []*FundLog{}
		obj.Data.FundsLog = tran
	}

	return utils.StructToJsonString(obj.Data)
}
