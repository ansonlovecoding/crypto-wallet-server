package wallet

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	http2 "Share-Wallet/pkg/common/http"
	"Share-Wallet/pkg/common/log"
	sql "Share-Wallet/pkg/db/mysql/mysql_model"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/eth"
	"Share-Wallet/pkg/proto/tron"
	"Share-Wallet/pkg/proto/wallet"
	adminStruct "Share-Wallet/pkg/struct/admin_api"
	walletStruct "Share-Wallet/pkg/struct/wallet_api"
	"Share-Wallet/pkg/utils"
	"context"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/ip2location/ip2location-go/v9"
)

// TestWalletApi godoc
// @Summary      Testing wallet server
// @Description  Testing wallet server
// @Tags         Test
// @Accept       json
// @Produce      json
// @Param        req body wallet_api.TestRequest true "operationID is only for tracking"
// @Success      200  {object}  wallet_api.TestResponse
// @Failure      400  {object}  wallet_api.TestResponse
// @Router       /wallet/test [post]
func TestWalletApi(c *gin.Context) {
	var (
		req   walletStruct.TestRequest
		resp  walletStruct.TestResponse
		reqPb wallet.CommonReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TestWalletApi!!")

	reqPb.OperationID = req.OperationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := wallet.NewWalletClient(etcdConn)
	_, err := client.TestWalletRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	resp.Name = req.Name
	http2.RespHttp200(c, constant.OK, resp)
}

// GetSupportTokenAddresses godoc
// @Summary      Get token address list that we are support
// @Description  Get token address list that we are support
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Param        req body wallet_api.TestRequest true "operationID is only for tracking"
// @Success      200  {object}  wallet_api.GetSupportTokenAddressesRequest
// @Failure      400  {object}  wallet_api.GetSupportTokenAddressesResponse
// @Router       /eth/get_support_token_addresses [post]
func GetSupportTokenAddresses(c *gin.Context) {
	var (
		req   walletStruct.GetSupportTokenAddressesRequest
		resp  walletStruct.GetSupportTokenAddressesResponse
		reqPb wallet.GetSupportTokenAddressesReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := wallet.NewWalletClient(etcdConn)
	respPb, err := client.GetSupportTokenAddressesRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}

	var addressList []*walletStruct.SupportTokenAddress
	for _, address := range respPb.AddressList {
		var newAddress walletStruct.SupportTokenAddress
		err := utils.CopyStructFields(&newAddress, &address)
		if err != nil {
			http2.RespHttp200(c, err, nil)
		}
		addressList = append(addressList, &newAddress)
	}
	resp.AddressList = addressList
	http2.RespHttp200(c, constant.OK, resp)
}

func CreateAccountInformation(c *gin.Context) {
	var (
		req   adminStruct.CreateAccountInfoRequest
		reqPb wallet.CreateAccountInfoReq
		resp  adminStruct.CreateAccountInfoResponse
	)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "account info")
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	// account, err := sql.GetAccountInformationByMerchantUid(req.MerchantUid)
	// if account.ID != 0 {
	// 	http2.RespHttp200(c, err, nil)
	// 	return
	// }
	ip := c.ClientIP()
	//db, err := ip2location.OpenDB("/data/wallet/server/IP2LOCATION-LITE-DB5.BIN")
	absPath, _ := filepath.Abs("../IP2LOCATION-LITE-DB5.BIN")
	db, err := ip2location.OpenDB(absPath)
	if err != nil {
		log.NewError(req.OperationID, "Failed to get locations")
		http2.RespHttp200(c, err, "Failed to get locations")
		return
	}
	result, err := db.Get_all(ip)
	if err != nil {
		log.NewError(req.OperationID, "Failed location info")
		return
	}
	platform := "Android"
	if req.LastLoginTerminal == "Ios" {
		platform = "Ios"

	}
	reqPb.Uid = req.Uid
	reqPb.MerchantId = req.MerchantUid
	reqPb.AccountSource = req.AccountSource
	reqPb.BtcBalance = req.BtcBalance
	reqPb.EthBalance = req.EthBalance
	reqPb.ErcBalance = req.ErcBalance
	reqPb.TrcBalance = req.TrcBalance
	reqPb.TrxBalance = req.TrxBalance
	reqPb.BtcPublicAddress = req.BtcPublicAddress
	reqPb.EthPublicAddress = req.EthPublicAddress
	reqPb.ErcPublicAddress = req.ErcPublicAddress
	reqPb.TrcPublicAddress = req.TrcPublicAddress
	reqPb.TrxPublicAddress = req.TrxPublicAddress
	reqPb.CreationLoginIp = c.ClientIP()
	reqPb.CreationLoginRegion = result.City
	reqPb.CreationLoginTerminal = platform
	reqPb.CreationLoginTime = req.CreationLoginTime
	reqPb.LastLoginIp = c.ClientIP()
	reqPb.LastLoginRegion = result.City
	reqPb.LastLoginTerminal = platform
	reqPb.LastLoginTime = req.LastLoginTime

	client := wallet.NewWalletClient(etcdConn)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "create account info rpc")
	respPb, err := client.CreateAccountInformation(context.Background(), &reqPb)
	if err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "create account info FAILED", err)
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	resp.Uid = respPb.Uuid
	http2.RespHttp200(c, constant.OK, resp)
}

func UpdateAccountInformation(c *gin.Context) {
	var (
		req   walletStruct.UpdateAccountInfoRequest
		reqPb wallet.UpdateAccountInfoReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	account, err := sql.GetAccountInformationByMerchantUidAndUid(req.MerchantUid, req.Uid)
	if account.ID == 0 {
		http2.RespHttp200(c, err, nil)
		return
	}
	ip := c.ClientIP()
	//db, err := ip2location.OpenDB("/data/wallet/server/IP2LOCATION-LITE-DB5.BIN")
	absPath, _ := filepath.Abs("../IP2LOCATION-LITE-DB5.BIN")
	db, err := ip2location.OpenDB(absPath)
	if err != nil {
		log.NewError(req.OperationID, "Failed to get locations")
		http2.RespHttp200(c, err, "Failed to get locations")
		return
	}
	result, err := db.Get_all(ip)
	if err != nil {
		log.NewError(req.OperationID, "Failed location info")
		return
	}
	platform := "Android"
	if req.LastLoginTerminal == "Ios" {
		platform = "Ios"

	}
	reqPb.Uid = req.Uid
	reqPb.MerchantId = req.MerchantUid
	reqPb.LastLoginIp = ip
	reqPb.LastLoginRegion = result.City
	reqPb.LastLoginTerminal = platform
	reqPb.LastLoginTime = time.Now().Unix()

	client := wallet.NewWalletClient(etcdConn)
	_, er := client.UpdateAccountInformation(context.Background(), &reqPb)
	if er != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	http2.RespHttp200(c, constant.OK, "updated")
}

// GetFundsLog godoc
// @Summary      GetFundsLog
// @Description  GetFundsLog
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Param        req query wallet_api.GetFundsLogRequest true
// @Success      200  {object}  wallet_api.GetFundsLogResponse
// @Failure      400  {object}  wallet_api.GetFundsLogResponse
// @Router       /wallet/funds-log [get]
func GetFundsLog(c *gin.Context) {
	var (
		req   walletStruct.GetFundsLogRequest
		reqPb wallet.GetFundsLogReq
		resp  walletStruct.GetFundsLogResponse
	)

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetFundsLog")

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.Uid = req.UID
	if req.Page == 0 {
		req.Page = constant.DefaultPageNumber
	}
	if req.PageSize == 0 && req.Page != -1 {
		req.PageSize = constant.DefaultPageSize
	}
	reqPb.Pagination = &wallet.RequestPagination{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := wallet.NewWalletClient(etcdConn)
	respPb, err := client.GetFundsLog(context.Background(), &reqPb)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), respPb)
	if err != nil {
		http2.RespHttp200(c, err, "error")
		return
	}
	if len(respPb.FundLog) > 0 {
		utils.CopyStructFields(&resp.Funds, respPb.FundLog)
	}
	if req.Page == -1 {
		resp.Page = constant.DefaultPageNumber
		resp.PageSize = 0
	} else {
		resp.Page = req.Page
		resp.PageSize = req.PageSize
	}
	resp.TotalNum = respPb.TotalFundLogs
	http2.RespHttp200(c, err, resp)
}

func GetCoinStatuses(c *gin.Context) {
	var (
		req   walletStruct.GetCoinStatusesRequest
		reqPb wallet.GetCoinStatusesReq
		resp  walletStruct.GetCoinsStatusesResponse
	)
	if err := c.ShouldBind(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := wallet.NewWalletClient(etcdConn)
	respPb, er := client.GetCoinStatuses(context.Background(), &reqPb)
	if er != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", er.Error())
		http2.RespHttp200(c, er, nil)
		return
	}
	if len(respPb.Currencies) > 0 {
		utils.CopyStructFields(&resp.Currencies, respPb.Currencies)
	}
	http2.RespHttp200(c, constant.OK, resp)
}
func GetCoinRatio(c *gin.Context) {
	var (
		req   walletStruct.GetCoinRatioRequest
		reqPb wallet.GetCoinRatioReq
		resp  walletStruct.GetCoinRatioResponse
	)
	if err := c.ShouldBind(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := wallet.NewWalletClient(etcdConn)
	respPb, er := client.GetCoinRatio(context.Background(), &reqPb)
	if er != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", er.Error())
		http2.RespHttp200(c, er, nil)
		return
	}
	if len(respPb.Coins) > 0 {
		utils.CopyStructFields(&resp.Coins, respPb.Coins)
	}
	http2.RespHttp200(c, constant.OK, resp)
}
func GetUserWallet(c *gin.Context) {
	var (
		req   walletStruct.GetUserWalletRequest
		reqPb wallet.GetUserWalletReq
		resp  walletStruct.GetUserWalletResponse
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.UserId = req.UserId
	reqPb.CoinType = req.CoinType
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := wallet.NewWalletClient(etcdConn)
	respPb, er := client.GetUserWallet(context.Background(), &reqPb)
	if er != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", er.Error())
		http2.RespHttp200(c, er, nil)
		return
	}
	resp.HasWallet = respPb.HasWallet
	resp.Address = respPb.Address
	http2.RespHttp200(c, constant.OK, resp)
}
func UpdateCoinRates(c *gin.Context) {
	var (
		req   walletStruct.UpdateCoinRatesRequest
		reqPb wallet.UpdateCoinRatesReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := wallet.NewWalletClient(etcdConn)
	_, er := client.UpdateCoinRates(context.Background(), &reqPb)
	if er != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", er.Error())
		http2.RespHttp200(c, er, nil)
		return
	}
	http2.RespHttp200(c, constant.OK, "Updated")
}

// GetAccountBalance godoc
// @Summary      Get Balance for each coin address
// @Description  Get Balance for each coin address, coinType: 1 BTC, 2 ETH, 3 USDT-ERC20, 4 TRX, 5 USDT-TRC20
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Param        req body wallet_api.GetAccountBalanceRequest true "operationID is only for tracking"
// @Success      200  {object}  wallet_api.GetAccountBalanceResponse
// @Failure      400  {object}  wallet_api.GetAccountBalanceResponse
// @Router       /wallet/get_balance [post]
func GetAccountBalance(c *gin.Context) {
	var (
		req       walletStruct.GetAccountBalanceRequest
		resp      walletStruct.GetAccountBalanceResponse
		reqEthPb  eth.GetBalanceReq
		reqTronPb tron.GetBalanceReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	if req.CoinType == constant.ETHCoin || req.CoinType == constant.USDTERC20 {
		reqEthPb.OperationID = req.OperationID
		reqEthPb.Address = req.Address
		reqEthPb.CoinType = uint32(req.CoinType)
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqEthPb.OperationID)
		if etcdConn == nil {
			errMsg := reqEthPb.OperationID + "getcdv3.GetConn == nil"
			log.NewError(reqEthPb.OperationID, errMsg)
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
			return
		}
		client := eth.NewEthClient(etcdConn)
		respPb, err := client.GetEthBalanceRPC(c, &reqEthPb)
		if err != nil {
			log.NewError(reqEthPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
			http2.RespHttp200(c, err, nil)
			return
		}
		balanceStr := respPb.Balance
		if balanceStr == "" {
			balanceStr = "0"
		}
		log.NewInfo(req.OperationID, "balance before convert:", balanceStr)
		if req.CoinType == constant.ETHCoin {
			resp.Balance = balanceStr
		} else {
			resp.Balance = utils.ConvertUSDTBalance(balanceStr)
		}
	} else if req.CoinType == constant.TRX || req.CoinType == constant.USDTTRC20 {
		reqTronPb.OperationID = req.OperationID
		reqTronPb.Address = req.Address
		reqTronPb.CoinType = uint32(req.CoinType)
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.TronRPC, reqTronPb.OperationID)
		if etcdConn == nil {
			errMsg := reqTronPb.OperationID + "getcdv3.GetConn == nil"
			log.NewError(reqTronPb.OperationID, errMsg)
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
			return
		}
		client := tron.NewTronClient(etcdConn)
		respPb, err := client.GetTronBalanceRPC(c, &reqTronPb)
		if err != nil {
			log.NewError(reqTronPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
			http2.RespHttp200(c, err, nil)
			return
		}
		balanceStr := respPb.Balance
		if balanceStr == "" {
			balanceStr = "0"
		}
		log.NewInfo(req.OperationID, "balance before convert:", balanceStr)
		if req.CoinType == constant.TRX {
			var balance = balanceStr
			balanceDecimal, err := decimal.NewFromString(balance)
			if err != nil {
				balanceDecimal = decimal.NewFromInt(0)
			}
			trxDecimal := utils.SunToTrx(balanceDecimal.BigInt())
			resp.Balance = trxDecimal.String()
		} else {
			resp.Balance = utils.ConvertUSDTBalance(balanceStr)
		}
	}

	log.NewInfo(req.OperationID, "balance after convert:", resp.Balance)
	http2.RespHttp200(c, constant.OK, resp)
}

func GetTransaction(c *gin.Context) {
	var (
		req   walletStruct.GetTransactionDetailRequest
		resp  walletStruct.GetTransactionDetailResponse
		reqPb wallet.GetTransactionDetailReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.CoinType = int32(req.CoinType)
	reqPb.TransactionHash = req.TransactionSignedHash
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		http2.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	client := wallet.NewWalletClient(etcdConn)
	respPb, err := client.GetTransactionDetailRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}

	resp.TransactionDetail = walletStruct.TransactionInfo{
		UUID:               respPb.TransactionDetails.UUID,
		SenderAccount:      respPb.TransactionDetails.SenderAccount,
		SenderAddress:      respPb.TransactionDetails.SenderAddress,
		ReceiverAccount:    respPb.TransactionDetails.ReceiverAccount,
		ReceiverAddress:    respPb.TransactionDetails.ReceiverAddress,
		Amount:             respPb.TransactionDetails.Amount,
		Fee:                respPb.TransactionDetails.Fee,
		TransactionHash:    respPb.TransactionDetails.SentHashTX,
		ConfirmationTime:   respPb.TransactionDetails.ConfirmTime,
		SentTime:           respPb.TransactionDetails.SentUpdatedAt,
		Status:             int8(respPb.TransactionDetails.Status),
		GasUsed:            respPb.TransactionDetails.GasUsed,
		GasLimit:           respPb.TransactionDetails.GasLimit,
		GasPrice:           respPb.TransactionDetails.GasPrice,
		ConfirmBlockNumber: respPb.TransactionDetails.ConfirmBlockNumber,
	}
	// log.NewInfo(req.OperationID, "TransactionID:", respPb.TransactionDetail.UUID)
	http2.RespHttp200(c, constant.OK, resp)
}

func GetTransactionList(c *gin.Context) {
	var (
		req   walletStruct.GetTransactionListRequest
		resp  walletStruct.GetTransactionListResponse
		reqPb wallet.GetTransactionListReq
		err   error
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.PublicAddress = req.PublicAddress
	reqPb.TransactionType = int32(req.TransactionType)
	reqPb.TransactionState = int32(req.TransactionState)
	reqPb.TransactionHash = req.TransactionHash
	reqPb.CoinType = int32(req.CoinType)
	if req.Page == 0 {
		req.Page = constant.DefaultPageNumber
	}
	if req.PageSize == 0 && req.Page != -1 {
		req.PageSize = constant.DefaultPageSize
	}
	reqPb.Pagination = &wallet.RequestPagination{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	reqPb.OrderBy = req.OrderBy
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.WalletRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := wallet.NewWalletClient(etcdConn)
	respPb, err := client.GetTransactionListRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "respPb", respPb.String())
	var info []walletStruct.TransactionInfo

	for _, v := range respPb.TransactionDetails {
		info = append(info, walletStruct.TransactionInfo{
			UUID:               v.UUID,
			SenderAccount:      v.SenderAccount,
			SenderAddress:      v.SenderAddress,
			ReceiverAccount:    v.ReceiverAccount,
			ReceiverAddress:    v.ReceiverAddress,
			Amount:             v.Amount,
			Fee:                v.Fee,
			TransactionHash:    v.SentHashTX,
			ConfirmationTime:   v.ConfirmTime,
			SentTime:           v.SentUpdatedAt,
			Status:             int8(v.Status),
			GasUsed:            v.GasUsed,
			GasLimit:           v.GasLimit,
			GasPrice:           v.GasPrice,
			ConfirmBlockNumber: v.ConfirmBlockNumber,
		})
	}
	resp.TransactionDetail = info
	if req.Page == -1 {
		resp.Page = constant.DefaultPageNumber
		resp.PageSize = 0
	} else {
		resp.Page = int(req.Page)
		resp.PageSize = int(req.PageSize)
	}
	resp.TranNums = respPb.TotalTrans
	http2.RespHttp200(c, constant.OK, resp)
}
