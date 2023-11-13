package wallet

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	http2 "Share-Wallet/pkg/common/http"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/eth"
	walletStruct "Share-Wallet/pkg/struct/wallet_api"
	"Share-Wallet/pkg/utils"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// TestEthApi godoc
// @Summary      Testing eth-rpc server
// @Description  Testing eth-rpc server
// @Tags         Test
// @Accept       json
// @Produce      json
// @Param        req body wallet_api.TestRequest true "operationID is only for tracking"
// @Success      200  {object}  wallet_api.TestResponse
// @Failure      400  {object}  wallet_api.TestResponse
// @Router       /wallet/test_eth [post]
func TestEthApi(c *gin.Context) {
	var (
		req   walletStruct.TestRequest
		resp  walletStruct.TestResponse
		reqPb eth.CommonReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TestWalletApi!!")

	reqPb.OperationID = req.OperationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := eth.NewEthClient(etcdConn)
	_, err := client.TestEthRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	resp.Name = req.Name
	http2.RespHttp200(c, constant.OK, resp)
}

// GetETHBalance godoc
// @Summary      Get Balance for eth address
// @Description  Get Balance for eth address, coinType: 1 BTC, 2 ETH, 3 USDT-ERC20, 4 TRX, 5 USDT-TRC20
// @Tags         ETH
// @Accept       json
// @Produce      json
// @Param        req body wallet_api.GetEthBalanceRequest true "operationID is only for tracking"
// @Success      200  {object}  wallet_api.GetEthBalanceResponse
// @Failure      400  {object}  wallet_api.GetEthBalanceResponse
// @Router       /wallet/test_eth [post]
func GetETHBalance(c *gin.Context) {
	var (
		req   walletStruct.GetEthBalanceRequest
		resp  walletStruct.GetEthBalanceResponse
		reqPb eth.GetBalanceReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.Address = req.Address
	reqPb.CoinType = uint32(req.CoinType)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := eth.NewEthClient(etcdConn)
	respPb, err := client.GetEthBalanceRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}

	log.NewInfo(req.OperationID, "balance before convert:", respPb.Balance)
	if req.CoinType == constant.ETHCoin {
		resp.Balance = respPb.Balance
	} else {
		resp.Balance = utils.ConvertUSDTBalance(respPb.Balance)
	}

	log.NewInfo(req.OperationID, "balance after convert:", resp.Balance)
	http2.RespHttp200(c, constant.OK, resp)
}

// Get ETH gas price
func GetGasPrice(c *gin.Context) {
	var (
		req   walletStruct.GetGasPriceRequest
		resp  walletStruct.GetGasPriceResponse
		reqPb eth.GetGasPriceReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := eth.NewEthClient(etcdConn)
	respPb, err := client.GetEthGasPriceRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}

	log.NewInfo(req.OperationID, "gasprice:", respPb.GasPrice)
	resp.GasPrice = respPb.GasPrice
	http2.RespHttp200(c, constant.OK, resp)
}

// ETH transfer
func Transfer(c *gin.Context) {
	var (
		req   walletStruct.PostTransferRequest
		resp  walletStruct.PostTransferResponse
		reqPb eth.PostTransferReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.CoinType = req.CoinType
	reqPb.FromAccountUID = req.FromAccountUID
	reqPb.FromMerchantUID = req.FromMerchantUID
	reqPb.FromAddress = req.FromAddress
	reqPb.ToAddress = req.ToAddress
	reqPb.Amount = req.Amount
	reqPb.Fee = req.Fee
	reqPb.GasLimit = req.GasLimit
	reqPb.Nounce = req.Nonce
	reqPb.GasPrice = req.GasPrice
	reqPb.TxHash = req.TxHash

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := eth.NewEthClient(etcdConn)
	respPb, err := client.TransferRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}

	// log.NewInfo(req.OperationID, "gasprice:", respPb.GasPrice)
	amountDecimal, _ := decimal.NewFromString(respPb.Transaction.Amount)
	feeDecimal, _ := decimal.NewFromString(respPb.Transaction.Fee)
	resp.EthTransactionDetail = walletStruct.Transaction{
		UUID:            respPb.Transaction.UUID,
		SenderAccount:   respPb.Transaction.SenderAccount,
		SenderAddress:   respPb.Transaction.SenderAddress,
		ReceiverAccount: respPb.Transaction.ReceiverAccount,
		ReceiverAddress: respPb.Transaction.ReceiverAddress,
		Amount:          amountDecimal,
		Fee:             feeDecimal,
		GasLimit:        respPb.Transaction.GasLimit,
		Nonce:           respPb.Transaction.Nonce,
		SentHashTX:      respPb.Transaction.SentHashTX,
		SentUpdatedAt:   respPb.Transaction.SentUpdatedAt,
	}
	http2.RespHttp200(c, constant.OK, resp)
}

func GetConfirmationTransaction(c *gin.Context) {
	var (
		req   walletStruct.GetEthConfirmationRequest
		resp  walletStruct.GetEthConfirmationResponse
		reqPb eth.GetEthConfirmationReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.CoinType = uint32(req.CoinType)
	reqPb.TransactionHash = req.TransactionSignedHash
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := eth.NewEthClient(etcdConn)
	respPb, err := client.GetConfirmationRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}

	log.NewInfo(req.OperationID, "BlockNumber:", respPb.BlockNum)
	resp.BlockNum = respPb.BlockNum
	resp.ConfirmationTime = respPb.ConfirmTime
	resp.GasUsed = respPb.GasUsed
	resp.Status = int8(respPb.Status)
	http2.RespHttp200(c, constant.OK, resp)
}

// ETH transfer
func TransferV2(c *gin.Context) {
	var (
		req   walletStruct.PostTransferRequestV2
		resp  walletStruct.PostTransferResponseV2
		reqPb eth.PostTransferReq2
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.Rawhex = req.RawHex

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := eth.NewEthClient(etcdConn)
	respPb, err := client.TransferRPCV2(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}

	// log.NewInfo(req.OperationID, "gasprice:", respPb.GasPrice)
	// resp.EthTransactionDetail = walletStruct.Transaction{
	// 	TransactionID:          respPb.Transaction.TransactionID,
	// 	UUID:                   respPb.Transaction.UUID,
	// 	CurrentTransactionType: respPb.Transaction.CurrentTransactionType,
	// 	SenderAccount:          respPb.Transaction.SenderAccount,
	// 	SenderAddress:          respPb.Transaction.SenderAddress,
	// 	ReceiverAccount:        respPb.Transaction.ReceiverAccount,
	// 	ReceiverAddress:        respPb.Transaction.ReceiverAddress,
	// 	Amount:                 respPb.Transaction.Amount,
	// 	Fee:                    respPb.Transaction.Fee,
	// 	GasLimit:               respPb.Transaction.GasLimit,
	// 	Nonce:                  respPb.Transaction.Nonce,
	// 	UnsignedHexTX:          respPb.Transaction.UnsignedHexTX,
	// 	SignedHexTX:            respPb.Transaction.SignedHexTX,
	// 	SentHashTX:             respPb.Transaction.SentHashTX,
	// 	UnsignedUpdatedAt:      respPb.Transaction.UnsignedUpdatedAt,
	// 	SentUpdatedAt:          respPb.Transaction.SentUpdatedAt,
	// }
	resp.TransactionHash = respPb.TranHash
	http2.RespHttp200(c, constant.OK, resp)
}

//Check balance and nonce
func CheckBalanceAndNonce(c *gin.Context) {
	var (
		req   walletStruct.CheckBalanceAndNonceRequest
		resp  walletStruct.CheckBalanceAndNonceResponse
		reqPb eth.CheckBalanceAndGetNonceReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.TransactAmount = req.TransactAmount
	reqPb.FromAddress = req.FromAddress
	reqPb.CoinType = req.CoinType
	reqPb.GasLimit = req.GasLimit
	reqPb.GasPrice = req.GasPrice

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := eth.NewEthClient(etcdConn)
	respPb, err := client.CheckBalanceAndGetNonceRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	if respPb.ErrCode != constant.OK.ErrCode {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "CheckBalanceAndGetNonceRPC failed", respPb.ErrCode, respPb.ErrMsg)
		http2.RespHttp200(c, constant.ErrInfo{ErrCode: respPb.ErrCode, ErrMsg: respPb.ErrMsg}, nil)
		return
	}
	resp.Nonce = respPb.Nonce
	resp.ChainID = respPb.ChainID
	resp.USDTERC20ContractAddress = respPb.USDTERC20ContractAddress
	http2.RespHttp200(c, constant.OK, resp)
}
