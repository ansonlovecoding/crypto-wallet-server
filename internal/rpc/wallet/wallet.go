package wallet

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/http"
	"Share-Wallet/pkg/common/log"
	db "Share-Wallet/pkg/db/mysql"
	walletdb "Share-Wallet/pkg/db/mysql/mysql_model"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/wallet"
	walletStruct "Share-Wallet/pkg/struct/wallet_api"
	"Share-Wallet/pkg/utils"
	"Share-Wallet/pkg/wallet/coin"
	config2 "Share-Wallet/pkg/wallet/config"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"google.golang.org/grpc"
)

type walletRPCServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewWalletRPCServer(port int) *walletRPCServer {
	return &walletRPCServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.WalletRPC,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *walletRPCServer) Run() {
	log.NewInfo("0", "walletRPCServer start ")
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
	wallet.RegisterWalletServer(srv, s)
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

func (s *walletRPCServer) TestWalletRPC(_ context.Context, req *wallet.CommonReq) (*wallet.CommonResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TestWalletRPC!!", req.String())
	resp := &wallet.CommonResp{
		ErrCode: 0,
		ErrMsg:  "Test Success!",
	}
	return resp, nil
}

func (s *walletRPCServer) GetSupportTokenAddressesRPC(_ context.Context, req *wallet.GetSupportTokenAddressesReq) (*wallet.GetSupportTokenAddressesResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetSupportTokenAddressesRPC!!", req.String())
	var addressList []*wallet.SupportTokenAddress
	//ETH Token
	cfgETH := os.Getenv("ETH_TOML_PATH")
	if len(cfgETH) == 0 {
		_, b, _, _ := runtime.Caller(0)
		// Root folder of this project
		Root := filepath.Join(filepath.Dir(b), "../..")
		cfgETH = Root + "/config/wallet/eth.toml"
	}

	confETH, err := config2.NewWallet(cfgETH, coin.ETH)
	if err != nil {
		return nil, err
	}

	//USDT-ERC20
	usdterc20 := wallet.SupportTokenAddress{
		BelongCoin:      constant.ETHCoin,
		CoinType:        constant.USDTERC20,
		ContractAddress: confETH.USDTERC20.ContractAddress,
	}
	addressList = append(addressList, &usdterc20)

	//TRON Token
	cfgTron := os.Getenv("TRON_TOML_PATH")
	if len(cfgTron) == 0 {
		_, b, _, _ := runtime.Caller(0)
		// Root folder of this project
		Root := filepath.Join(filepath.Dir(b), "../..")
		cfgTron = Root + "/config/wallet/tron.toml"
	}

	confTron, err := config2.NewWallet(cfgTron, coin.TRX)
	if err != nil {
		return nil, err
	}

	//USDT-TRC20
	usdttrc20 := wallet.SupportTokenAddress{
		BelongCoin:      constant.TRX,
		CoinType:        constant.USDTTRC20,
		ContractAddress: confTron.USDTTRC20.ContractAddress,
	}
	addressList = append(addressList, &usdttrc20)

	resp := &wallet.GetSupportTokenAddressesResp{
		ErrCode:     0,
		ErrMsg:      "request Success!",
		AddressList: addressList,
	}
	return resp, nil
}

func (s *walletRPCServer) CreateAccountInformation(c context.Context, req *wallet.CreateAccountInfoReq) (*wallet.CreateAccountInfoResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "create account information RPC")
	resp := &wallet.CreateAccountInfoResp{}

	//checking if the UUID(public key of wallet) was taken by another user
	accountByUid, err := walletdb.GetAccountInformationByUUID(req.Uid)

	//if this wallet was taken, not allow to create account information
	if err == nil && accountByUid != nil && accountByUid.MerchantUid != "" {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "the UUID(public key of wallet) was taken by another user")
		//check if the user is trying to recover his wallet, update login info
		if req.MerchantId == accountByUid.MerchantUid {
			account := db.AccountInformation{
				LastLoginIp:       req.LastLoginIp,
				LastLoginRegion:   req.LastLoginRegion,
				LastLoginTerminal: req.LastLoginTerminal,
				LastLoginTime:     time.Now().Unix(),
			}
			err := walletdb.UpdateAccountLoginInfo(account, req.MerchantId, req.Uid)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateAccountInformation failed", err.Error())
				return resp, http.WrapError(constant.ErrDB)
			}
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "User recovered his own wallet")
		}
		return resp, http.WrapError(constant.ErrNotYourWallet)
		/*
			//update account info
			account := db.AccountInformation{
				ID:                accountByUid.ID,
				UUID:              req.Uid,
				MerchantUid:       req.MerchantId,
				AccountSource:     accountByUid.AccountSource,
				BtcPublicAddress:  req.BtcPublicAddress,
				EthPublicAddress:  req.EthPublicAddress,
				ErcPublicAddress:  req.ErcPublicAddress,
				TrcPublicAddress:  req.TrcPublicAddress,
				TrxPublicAddress:  req.TrxPublicAddress,
				BtcBalance:        req.BtcBalance,
				EthBalance:        req.EthBalance,
				ErcBalance:        req.ErcBalance,
				TrcBalance:        req.TrcBalance,
				TrxBalance:        req.TrxBalance,
				LastLoginIp:       req.LastLoginIp,
				LastLoginRegion:   req.LastLoginRegion,
				LastLoginTerminal: req.LastLoginTerminal,
				LastLoginTime:     req.LastLoginTime,
			}
			log.NewInfo(req.OperationID, "Updating account information", "previous merchantUID", accountByUid.MerchantUid, "new merchantUID", req.MerchantId, req.Uid)
			err := walletdb.UpdateAccountInfoWithAccount(&account)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateAccountInfoWithAccount failed", err.Error())
				return resp, http.WrapError(constant.ErrDB)
			}

		*/

	} else {
		//checking if the MerchantUID(userID from the merchant) was taken by another user
		accountByMerchantUID, err := walletdb.GetAccountInformationByMerchantUid(req.MerchantId)
		//if it was taken, cannot create again
		if accountByMerchantUID != nil && accountByMerchantUID.MerchantUid != "" {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "the MerchantUID(userID from the merchant) was taken by another user")
			return resp, http.WrapError(constant.ErrYouBoundAlready)
			/*
				//update account info
				account := db.AccountInformation{
					ID:          accountByMerchantUID.ID,
					MerchantUid: accountByMerchantUID.MerchantUid + "_old",
				}
				log.NewInfo(req.OperationID, "Updating account information", "previous merchantUID", accountByMerchantUID.MerchantUid, "new merchantUID", account.MerchantUid, req.Uid)
				err := walletdb.UpdateAccountInfoWithAccount(&account)
				if err != nil {
					log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateAccountInfoWithAccount failed", err.Error())
					return resp, http.WrapError(constant.ErrDB)
				}

			*/
		}

		//otherwise create the account information
		account := db.AccountInformation{
			UUID:                  req.Uid,
			MerchantUid:           req.MerchantId,
			AccountSource:         req.AccountSource,
			BtcPublicAddress:      req.BtcPublicAddress,
			EthPublicAddress:      req.EthPublicAddress,
			ErcPublicAddress:      req.ErcPublicAddress,
			TrcPublicAddress:      req.TrcPublicAddress,
			TrxPublicAddress:      req.TrxPublicAddress,
			BtcBalance:            req.BtcBalance,
			EthBalance:            req.EthBalance,
			ErcBalance:            req.ErcBalance,
			TrcBalance:            req.TrcBalance,
			TrxBalance:            req.TrxBalance,
			CreationLoginIp:       req.CreationLoginIp,
			CreationLoginRegion:   req.CreationLoginRegion,
			CreationLoginTerminal: req.CreationLoginTerminal,
			CreationTime:          req.CreationLoginTime,
			LastLoginIp:           req.LastLoginIp,
			LastLoginRegion:       req.LastLoginRegion,
			LastLoginTerminal:     req.LastLoginTerminal,
			LastLoginTime:         req.LastLoginTime,
		}
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "creating new account with new uid and new merchant uid")
		err = walletdb.CreateAccountInformation(account)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreateAccountInformation failed", err.Error())
			return resp, http.WrapError(constant.ErrDB)
		}
	}

	resp.Uuid = req.Uid
	return resp, nil
}
func (s *walletRPCServer) UpdateAccountInformation(c context.Context, req *wallet.UpdateAccountInfoReq) (*wallet.CommonResp, error) {
	resp := &wallet.CommonResp{}
	account := db.AccountInformation{
		LastLoginIp:       req.LastLoginIp,
		LastLoginRegion:   req.LastLoginRegion,
		LastLoginTerminal: req.LastLoginTerminal,
		LastLoginTime:     time.Now().Unix(),
	}
	err := walletdb.UpdateAccountLoginInfo(account, req.MerchantId, req.Uid)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateAccountInformation failed", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	return resp, nil
}

func (s *walletRPCServer) GetFundsLog(c context.Context, req *wallet.GetFundsLogReq) (*wallet.GetFundsLogResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req.String())
	resp := &wallet.GetFundsLogResp{}
	filters := make(map[string]string)
	if req.From != "" {
		filters["to"] = req.To
		filters["from"] = req.From
	}
	if req.TransactionType != "" {
		filters["transaction_type"] = req.TransactionType
	}
	if req.UserAddress != "" {
		filters["user_address"] = req.UserAddress
	}
	if req.OppositeAddress != "" {
		filters["opposite_address"] = req.OppositeAddress
	}
	if req.CoinsType != "" {
		filters["coins_type"] = req.CoinsType
	}
	if req.State != "" {
		filters["state"] = req.State
	}
	if req.Txid != "" {
		filters["txid"] = req.Txid
	}
	if req.MerchantUid != "" {
		filters["merchant_uid"] = req.MerchantUid
	}

	fundsLog, count, err := walletdb.GetRecentRecords(filters, req.Pagination.Page, req.Pagination.PageSize, req.Uid)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFundsLog", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	currencies, _ := walletdb.GetCoinStatuses()

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "fundslog length:", len(fundsLog))
	if len(fundsLog) > 0 {
		for _, v := range fundsLog {
			//get the coin id
			coinId := utils.GetCoinType(v.CoinType) - 1
			coinsAmount, _ := v.AmountOfCoins.Float64()
			networkFee, _ := v.NetworkFee.Float64()
			totalCoins := v.AmountOfCoins.Add(v.NetworkFee)
			totalCoinsAmount, _ := totalCoins.Float64()
			totalUsdCoins := currencies[coinId].Usd * coinsAmount
			totalYuanCoins := currencies[coinId].Yuan * coinsAmount
			totalEuroCoins := currencies[coinId].Euro * coinsAmount
			usdFee := currencies[coinId].Usd * networkFee
			yuanFee := currencies[coinId].Yuan * networkFee
			euroFee := currencies[coinId].Euro * networkFee

			if v.CoinType == "USDT-ERC20" {
				coinId = utils.GetCoinType("ETH") - 1
				ethRate := currencies[coinId]

				usdFee = ethRate.Usd * networkFee
				yuanFee = ethRate.Yuan * networkFee
				euroFee = ethRate.Euro * networkFee
			}

			totalUsdTransfered := totalUsdCoins + usdFee
			totalYuanTransfered := totalYuanCoins + yuanFee
			totalEuroTransfered := totalEuroCoins + euroFee
			// creationTime := time.Unix(v.CreationTime, 0)
			// confirmationTime := time.Unix(v.ConfirmationTime, 0)
			//state 0 failed, 1 success, 2 pending
			state := utils.GetFundLogStateToString(v.State)

			fundLog := &wallet.FundsLog{
				ID:                      v.ID,
				Txid:                    v.Txid,
				Uid:                     v.UID,
				MerchantUid:             v.MerchantUid,
				TransactionType:         v.TransactionType,
				UserAddress:             v.UserAddress,
				OppositeAddress:         v.OppositeAddress,
				CoinType:                v.CoinType,
				AmountOfCoins:           utils.RoundFloat(coinsAmount, 8),
				UsdAmount:               utils.RoundFloat(totalUsdCoins, 8),
				YuanAmount:              utils.RoundFloat(totalYuanCoins, 8),
				EuroAmount:              utils.RoundFloat(totalEuroCoins, 8),
				NetworkFee:              utils.RoundFloat(networkFee, 8),
				UsdNetworkFee:           utils.RoundFloat(usdFee, 8),
				YuanNetworkFee:          utils.RoundFloat(yuanFee, 8),
				EuroNetworkFee:          utils.RoundFloat(euroFee, 8),
				TotalCoinsTransfered:    utils.RoundFloat(totalCoinsAmount, 8),
				TotalUsdTransfered:      utils.RoundFloat(totalUsdTransfered, 8),
				TotalYuanTransfered:     utils.RoundFloat(totalYuanTransfered, 8),
				TotalEuroTransfered:     utils.RoundFloat(totalEuroTransfered, 8),
				CreationTime:            v.CreationTime,
				State:                   state,
				ConfirmationTime:        v.ConfirmationTime,
				GasUsed:                 v.GasUsed,
				GasPrice:                uint64(v.GasPrice.IntPart()),
				GasLimit:                uint64(v.GasLimit),
				ConfirmationBlockNumber: v.ConfirmationBlockNumber,
			}
			resp.FundLog = append(resp.FundLog, fundLog)
		}
	}
	resp.TotalFundLogs = count
	return resp, nil
}

func (s *walletRPCServer) GetCoinStatuses(c context.Context, req *wallet.GetCoinStatusesReq) (*wallet.GetCoinStatusesResp, error) {
	resp := &wallet.GetCoinStatusesResp{}
	coins, err := walletdb.GetCoinStatuses()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetCoinStatuses failed", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	for _, v := range coins {
		currency := &wallet.Currency{
			ID:             int32(v.ID),
			CoinType:       v.Coin,
			LastEditedTime: v.UpdateTime,
			Editor:         v.UpdateUser,
			State:          int32(v.State),
		}
		resp.Currencies = append(resp.Currencies, currency)
	}

	return resp, nil
}

func (s *walletRPCServer) GetCoinRatio(c context.Context, req *wallet.GetCoinRatioReq) (*wallet.GetCoinRatioResp, error) {
	resp := &wallet.GetCoinRatioResp{}
	coins, err := walletdb.GetCoinStatuses()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetCoinStatuses failed", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	for _, v := range coins {
		coin := &wallet.Coin{
			ID:       int32(v.ID),
			CoinType: v.Coin,
			Usd:      v.Usd,
			Yuan:     v.Yuan,
			Euro:     v.Euro,
		}
		resp.Coins = append(resp.Coins, coin)
	}

	return resp, nil
}
func (s *walletRPCServer) GetUserWallet(c context.Context, req *wallet.GetUserWalletReq) (*wallet.GetUserWalletResp, error) {
	resp := &wallet.GetUserWalletResp{}
	wallet, err := walletdb.GetAccountInformationByMerchantUid(req.UserId)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAccountInformationByMerchantUid failed", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	if wallet.UUID != "" {
		resp.HasWallet = true
	} else {
		resp.HasWallet = false
	}
	var address string
	if req.CoinType != 0 {
		switch req.CoinType {
		case constant.BTCCoin:
			address = wallet.BtcPublicAddress
		case constant.ETHCoin:
			address = wallet.EthPublicAddress
		case constant.USDTERC20:
			address = wallet.EthPublicAddress
		case constant.TRX:
			address = wallet.TrxPublicAddress
		case constant.USDTTRC20:
			address = wallet.TrxPublicAddress
		}
	}
	resp.Address = address
	return resp, nil
}
func (s *walletRPCServer) UpdateCoinRates(c context.Context, req *wallet.UpdateCoinRatesReq) (*wallet.CommonResp, error) {
	resp := &wallet.CommonResp{}
	var btcResult, ethResult, usdtResult, trxResult walletStruct.Obj
	var trxResultTemp walletStruct.TempRateObj

	//fetch coin rates using coinbase
	btcResp, _ := http.Get(constant.GetBtc)
	ethResp, _ := http.Get(constant.GetEth)
	usdtResp, _ := http.Get(constant.GetUsdt)
	trxResp, _ := http.Get(constant.GetTrx)
	trxRespTemp, _ := http.Get(constant.GetTrxTemp)

	//unmarshall objects to structs
	err := utils.JsonStringToStruct(string(btcResp), &btcResult)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "btcResp failed", err.Error())
		return resp, err
	}
	err = utils.JsonStringToStruct(string(ethResp), &ethResult)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ethResp failed", err.Error())
		return resp, err
	}
	err = utils.JsonStringToStruct(string(usdtResp), &usdtResult)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "usdtResp failed", err.Error())
		return resp, err
	}
	err = utils.JsonStringToStruct(string(trxResp), &trxResult)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "trxResult failed", err.Error())
		return resp, err
	}
	err = utils.JsonStringToStruct(string(trxRespTemp), &trxResultTemp)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "trxResult failed", err.Error())
		return resp, err
	}
	log.NewInfo(req.OperationID, string(trxRespTemp), trxResultTemp)
	log.NewInfo(req.OperationID, "values", trxResultTemp.Tron.Usd, trxResultTemp.Tron.Cny, trxResultTemp.Tron.Eur)

	//convert string values from the response to float and form the objects
	floatEthUsd, _ := strconv.ParseFloat(ethResult.Data.Rates.USD, 64)
	floatEthEur, _ := strconv.ParseFloat(ethResult.Data.Rates.EUR, 64)
	floatEthCny, _ := strconv.ParseFloat(ethResult.Data.Rates.CNY, 64)
	eth := &db.CoinCurrencyValues{
		Coin: "ETH",
		Usd:  floatEthUsd,
		Yuan: floatEthCny,
		Euro: floatEthEur,
	}
	floatErcUsd, _ := strconv.ParseFloat(usdtResult.Data.Rates.USD, 64)
	floatErcEur, _ := strconv.ParseFloat(usdtResult.Data.Rates.EUR, 64)
	floatErcCny, _ := strconv.ParseFloat(usdtResult.Data.Rates.CNY, 64)
	erc := &db.CoinCurrencyValues{
		Coin: "USDT-ERC20",
		Usd:  floatErcUsd,
		Yuan: floatErcCny,
		Euro: floatErcEur,
	}

	floatBtcUsd, _ := strconv.ParseFloat(btcResult.Data.Rates.USD, 64)
	floatBtcEur, _ := strconv.ParseFloat(btcResult.Data.Rates.EUR, 64)
	floatBtcCny, _ := strconv.ParseFloat(btcResult.Data.Rates.CNY, 64)
	btc := &db.CoinCurrencyValues{
		Coin: "BTC",
		Usd:  floatBtcUsd,
		Yuan: floatBtcCny,
		Euro: floatBtcEur,
	}
	floatTrxUsd := trxResultTemp.Tron.Usd
	floatTrxEur := trxResultTemp.Tron.Eur
	floatTrxCny := trxResultTemp.Tron.Cny
	trx := &db.CoinCurrencyValues{
		Coin: "TRX",
		Usd:  floatTrxUsd,
		Yuan: floatTrxCny,
		Euro: floatTrxEur,
	}
	floatTrcUsd, _ := strconv.ParseFloat(usdtResult.Data.Rates.USD, 64)
	floatTrcEur, _ := strconv.ParseFloat(usdtResult.Data.Rates.EUR, 64)
	floatTrcCny, _ := strconv.ParseFloat(usdtResult.Data.Rates.CNY, 64)
	trc := &db.CoinCurrencyValues{
		Coin: "USDT-TRC20",
		Usd:  floatTrcUsd,
		Yuan: floatTrcCny,
		Euro: floatTrcEur,
	}
	coinsMap := make(map[string]*db.CoinCurrencyValues)
	coinsMap["btc"] = btc
	coinsMap["eth"] = eth
	coinsMap["erc"] = erc
	coinsMap["trx"] = trx
	coinsMap["trc"] = trc
	coinTypes := []string{utils.GetCoinName(constant.BTCCoin),
		utils.GetCoinName(constant.ETHCoin),
		utils.GetCoinName(constant.USDTERC20),
		utils.GetCoinName(constant.TRX),
		utils.GetCoinName(constant.USDTTRC20)}
	//Bulk update coin rates in the db
	err = walletdb.BulkUpdateCoinRates(coinsMap, coinTypes)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateCoinRates eth failed", err.Error())
		return resp, err
	}
	return resp, nil
}

func (s *walletRPCServer) GetTransactionDetailRPC(ctx context.Context, req *wallet.GetTransactionDetailReq) (*wallet.GetTransactionDetailRes, error) {

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetTransaction!!!", req.String())

	var (
		resp = &wallet.GetTransactionDetailRes{}
	)

	if req.CoinType == constant.ETHCoin || req.CoinType == constant.USDTERC20 {
		txDetail, err := walletdb.GetETHTxDetailByTxid(req.TransactionHash)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetETHTxDetailByTxid failed", err.Error())
			return nil, err
		}

		var sentUpdateUnix, confirmTimeUnix uint64
		if txDetail.SentUpdatedAt != nil {
			sentUpdateUnix = uint64(txDetail.SentUpdatedAt.Unix())
		}
		if txDetail.ConfirmTime != nil {
			confirmTimeUnix = uint64(txDetail.ConfirmTime.Unix())
		}
		resp.TransactionDetails = &wallet.TransactionDetail{
			UUID:               txDetail.UUID,
			SenderAccount:      txDetail.SenderAccount,
			SenderAddress:      txDetail.SenderAddress,
			ReceiverAccount:    txDetail.ReceiverAccount,
			ReceiverAddress:    txDetail.ReceiverAddress,
			Amount:             txDetail.Amount.String(),
			Fee:                txDetail.Fee.String(),
			GasLimit:           txDetail.GasLimit,
			Nonce:              txDetail.Nonce,
			SentHashTX:         txDetail.SentHashTX,
			SentUpdatedAt:      sentUpdateUnix,
			Status:             int32(txDetail.Status),
			ConfirmTime:        confirmTimeUnix,
			GasPrice:           txDetail.GasPrice.String(),
			GasUsed:            txDetail.GasUsed,
			ConfirmBlockNumber: txDetail.ConfirmationBlockNumber,
		}
	} else if req.CoinType == constant.TRX || req.CoinType == constant.USDTTRC20 {
		txDetail, err := walletdb.GetTronTxDetailByTxid(req.TransactionHash)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTronTxDetailByTxid failed", err.Error())
			return nil, err
		}

		var sentUpdateUnix, confirmTimeUnix uint64
		if txDetail.SentUpdatedAt != nil {
			sentUpdateUnix = uint64(txDetail.SentUpdatedAt.Unix())
		}
		if txDetail.ConfirmTime != nil {
			confirmTimeUnix = uint64(txDetail.ConfirmTime.Unix())
		}
		resp.TransactionDetails = &wallet.TransactionDetail{
			UUID:               txDetail.UUID,
			SenderAccount:      txDetail.SenderAccount,
			SenderAddress:      txDetail.SenderAddress,
			ReceiverAccount:    txDetail.ReceiverAccount,
			ReceiverAddress:    txDetail.ReceiverAddress,
			Amount:             txDetail.Amount.String(),
			Fee:                txDetail.Fee.String(),
			GasLimit:           0,
			Nonce:              0,
			SentHashTX:         txDetail.SentHashTX,
			SentUpdatedAt:      sentUpdateUnix,
			Status:             int32(txDetail.Status),
			ConfirmTime:        confirmTimeUnix,
			GasPrice:           "",
			GasUsed:            0,
			ConfirmBlockNumber: txDetail.ConfirmationBlockNumber,
		}
	}

	return resp, nil
}

func (s *walletRPCServer) GetTransactionListRPC(ctx context.Context, req *wallet.GetTransactionListReq) (*wallet.GetTransactionListRes, error) {

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetTransactionListRPC!!!", req.String())

	var (
		resp = wallet.GetTransactionListRes{}
	)

	whereConditionMap := map[string]interface{}{}

	if req.CoinType != 0 {
		whereConditionMap["coin_type"] = utils.GetCoinName(uint8(req.CoinType))
	}
	if req.TransactionHash != "" {
		whereConditionMap["sent_hash_tx"] = req.TransactionHash
	}

	if req.TransactionType != constant.TransactionTypeAll {
		switch req.TransactionType {
		case constant.TransactionTypeSend:
			whereConditionMap["sender_address"] = req.PublicAddress
		case constant.TransactionTypeReceive:
			whereConditionMap["receiver_address"] = req.PublicAddress
		default:
			//no allow to search without address
			whereConditionMap["sender_address"] = ""
			whereConditionMap["receiver_address"] = ""
		}
	} else {
		//search both sender_address and receiver_address
		whereConditionMap["sender_address"] = req.PublicAddress
		whereConditionMap["receiver_address"] = req.PublicAddress
	}

	if req.TransactionState != 0 {
		switch req.TransactionState - 1 {
		case constant.TxStatusPending:
			whereConditionMap["status"] = constant.TxStatusPending
		case constant.TxStatusSuccess:
			whereConditionMap["status"] = constant.TxStatusSuccess
		case constant.TxStatusFailed:
			whereConditionMap["status"] = constant.TxStatusFailed
		case constant.TxStatusExcludePending:
			whereConditionMap["status"] = constant.TxStatusExcludePending
		}
	}

	var transCount int64
	var err error
	if req.CoinType == constant.ETHCoin || req.CoinType == constant.USDTERC20 {
		transCount, err = walletdb.GetETHTransactionListCount(whereConditionMap)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetETHTransactionListCount failed", err.Error())
			return &resp, http.WrapError(constant.ErrDB)
		}
		if transCount == 0 {
			return &resp, nil
		}

		transactionList, err := walletdb.GetETHTransactionList(whereConditionMap, req.Pagination.Page, req.Pagination.PageSize, req.OrderBy)
		if err != nil {
			return &resp, fmt.Errorf("GetTransactionList failed %w", err)
		}
		for _, v := range transactionList {
			var confirmTime uint64
			if v.ConfirmTime != nil {
				confirmTime = uint64(v.ConfirmTime.Unix())
			}
			trans := &wallet.TransactionDetail{
				UUID:               v.UUID,
				SenderAccount:      v.SenderAccount,
				SenderAddress:      v.SenderAddress,
				ReceiverAccount:    v.ReceiverAccount,
				ReceiverAddress:    v.ReceiverAddress,
				Amount:             v.Amount.String(),
				Fee:                v.Fee.String(),
				GasLimit:           uint64(v.GasLimit),
				Nonce:              v.Nonce,
				SentHashTX:         v.SentHashTX,
				SentUpdatedAt:      uint64(v.SentUpdatedAt.Unix()),
				GasPrice:           v.GasPrice.String(),
				GasUsed:            v.GasUsed,
				ConfirmTime:        confirmTime,
				ConfirmBlockNumber: v.ConfirmationBlockNumber,
				Status:             int32(v.Status),
			}
			resp.TransactionDetails = append(resp.TransactionDetails, trans)
		}
	} else if req.CoinType == constant.TRX || req.CoinType == constant.USDTTRC20 {
		transCount, err = walletdb.GetTronTransactionListCount(whereConditionMap)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTronTransactionListCount failed", err.Error())
			return &resp, http.WrapError(constant.ErrDB)
		}
		if transCount == 0 {
			return &resp, nil
		}

		transactionList, err := walletdb.GetTronTransactionList(whereConditionMap, req.Pagination.Page, req.Pagination.PageSize, req.OrderBy)
		if err != nil {
			return &resp, fmt.Errorf("GetTransactionList failed %w", err)
		}
		for _, v := range transactionList {
			var confirmTime uint64
			if v.ConfirmTime != nil {
				confirmTime = uint64(v.ConfirmTime.Unix())
			}
			trans := &wallet.TransactionDetail{
				UUID:               v.UUID,
				SenderAccount:      v.SenderAccount,
				SenderAddress:      v.SenderAddress,
				ReceiverAccount:    v.ReceiverAccount,
				ReceiverAddress:    v.ReceiverAddress,
				Amount:             v.Amount.String(),
				Fee:                v.Fee.String(),
				GasLimit:           0,
				Nonce:              0,
				SentHashTX:         v.SentHashTX,
				SentUpdatedAt:      uint64(v.SentUpdatedAt.Unix()),
				GasPrice:           "",
				GasUsed:            0,
				ConfirmTime:        confirmTime,
				ConfirmBlockNumber: v.ConfirmationBlockNumber,
				Status:             int32(v.Status),
			}
			resp.TransactionDetails = append(resp.TransactionDetails, trans)
		}
	}

	resp.TotalTrans = transCount
	return &resp, nil
}
