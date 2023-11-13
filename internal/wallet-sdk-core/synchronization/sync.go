package synchronization

import (
	"Share-Wallet/pkg/common/constant"
	http "Share-Wallet/pkg/common/http"
	db "Share-Wallet/pkg/db/local_db"
	"Share-Wallet/pkg/db/local_db/model_struct"
	walletStruct "Share-Wallet/pkg/struct/wallet_api"
	"Share-Wallet/pkg/utils"
	"fmt"
)

type SyncMgr struct {
	Db               *db.DataBase
	LoginUserID      string
	LoginTime        int64
	API              *http.PostAPI
	UserRandomSecret []byte
}

var Sg *SyncMgr

func NewWalletMgr(dataBase *db.DataBase, loginUserID string, loginTime int32, userRandomSecret []byte, p *http.PostAPI) (w *SyncMgr) {
	return &SyncMgr{Db: dataBase, LoginUserID: loginUserID, LoginTime: int64(loginTime), UserRandomSecret: userRandomSecret, API: p}
}

type AccountInfo struct {
	OperationID           string  `json:"operationID" binding:"omitempty" swaggertype:"string"`
	Uid                   string  `json:"uid" binding:"required" swaggertype:"string"`
	MerchantUid           string  `json:"merchant_uid" binding:"required" swaggertype:"string"`
	BtcPublicAddress      string  `json:"btc_public_address" binding:"omitempty" swaggertype:"string"`
	EthPublicAddress      string  `json:"eth_public_address" binding:"omitempty" swaggertype:"string"`
	TrxPublicAddress      string  `json:"trx_public_address" binding:"omitempty" swaggertype:"string"`
	ErcPublicAddress      string  `json:"erc_public_address" binding:"omitempty" swaggertype:"string"`
	TrcPublicAddress      string  `json:"trc_public_address" binding:"omitempty" swaggertype:"string"`
	BtcBalance            float64 `json:"btc_balance" swaggertype:"number"`
	EthBalance            float64 `json:"eth_balance" swaggertype:"number"`
	TrxBalance            float64 `json:"trx_balance" swaggertype:"number"`
	ErcBalance            float64 `json:"erc_balance" swaggertype:"number"`
	TrcBalance            float64 `json:"trc_balance" swaggertype:"number"`
	AccountSource         string  `json:"account_source" swaggertype:"string"`
	CreationLoginIp       string  `json:"creation_login_ip" swaggertype:"string"`
	CreationLoginRegion   string  `json:"creation_login_region" swaggertype:"string"`
	CreationLoginTerminal string  `json:"creation_login_terminal" swaggertype:"string"`
	CreationLoginTime     int64   `json:"creation_login_time" swaggertype:"integer"`
	LastLoginIp           string  `json:"last_login_ip" swaggertype:"string"`
	LastLoginRegion       string  `json:"last_login_region" swaggertype:"string"`
	LastLoginTerminal     string  `json:"last_login_terminal" swaggertype:"string"`
	LastLoginTime         int64   `json:"last_login_time" swaggertype:"integer"`
}

func (s *SyncMgr) Synchronize(userID, publicKey string, platform int32) bool {
	fmt.Println("Synchronizing")
	user, err := s.Db.GetUserByUserID(userID)
	if err != nil {
		fmt.Println("GetUserByUserID error ", err)
		return false
	}
	type Data struct {
		Uid string `json:"uid"`
	}
	type Obj struct {
		Code   int    `json:"code"`
		ErrMsg string `json:"err_msg"`
		Data   Data   `json:"data"`
	}
	var obj Obj
	loginPlatform := ""
	switch platform {
	case 1:
		loginPlatform = "Android"
	case 2:
		loginPlatform = "Ios"
	}
	//check wallet id
	if user.WalletID == "" {
		addresses := make(map[string]string)
		wallets, err := s.Db.GetLocalUserAddresses(userID)
		if err != nil {
			fmt.Println("err", err)
			return false
		}
		for _, v := range wallets {
			cType := ""
			switch v.CoinType {
			case constant.BTCCoin:
				cType = "btc"
			case constant.ETHCoin:
				cType = "eth"
			case constant.USDTERC20:
				cType = "erc"
			case constant.USDTTRC20:
				cType = "trc"
			case constant.TRX:
				cType = "trx"
			}
			addresses[cType] = v.P2PKHAddress
		}
		accType := ""
		switch user.AccountType {
		case 1:
			accType = "share_wallet"
		case 2:
			accType = "imported"
		default:
			accType = "share_wallet"
		}
		request := &AccountInfo{
			Uid:                   publicKey,
			MerchantUid:           userID,
			AccountSource:         accType,
			BtcPublicAddress:      addresses["btc"],
			EthPublicAddress:      addresses["eth"],
			ErcPublicAddress:      addresses["erc"],
			TrcPublicAddress:      addresses["trc"],
			TrxPublicAddress:      addresses["trx"],
			BtcBalance:            0,
			EthBalance:            0,
			ErcBalance:            0,
			TrcBalance:            0,
			TrxBalance:            0,
			CreationLoginTime:     s.LoginTime,
			CreationLoginTerminal: loginPlatform,
			LastLoginTime:         s.LoginTime,
			LastLoginTerminal:     loginPlatform,
		}
		res, _ := s.API.PostWalletAPI(constant.CreateAccountInformationURL, request, constant.APITimeout)
		fmt.Println("result ", string(res))
		err = utils.JsonStringToStruct(string(res), &obj)
		if err != nil {
			fmt.Println("ersadadror", err)
			return false
		}
		// id := strconv.FormatInt(obj.Data.Uid, 10)
		updatedUser := model_struct.LocalUser{
			UserID:   userID,
			WalletID: publicKey,
		}
		err = s.Db.UpdateLocalUserWalletID(&updatedUser)
		if err != nil {
			fmt.Println("err", err, updatedUser)
			return false
		}
	} else {
		fmt.Println("WalletID", user.WalletID)
		// update login info
		request := &AccountInfo{
			Uid:               publicKey,
			MerchantUid:       userID,
			LastLoginTime:     Sg.LoginTime,
			LastLoginTerminal: loginPlatform,
		}
		res, err := s.API.PostWalletAPI(constant.UpdateLoginInformationURL, request, constant.APITimeout)
		if err != nil {
			fmt.Println("s.API.PostWalletAPI(constant.UpdateLoginInformationURL, request, constant.APITimeout) error", err)
			return false
		}
		fmt.Println("s.API.PostWalletAPI(constant.UpdateLoginInformationURL, request, constant.APITimeout) response", string(res))
		//err = utils.JsonStringToStruct(string(res), &obj)
		//fmt.Println("Data", obj.Data, err)

	}
	return true
}
func (s *SyncMgr) GetWallet(userID string, coinType uint32) (bool, string) {
	type req struct {
		UserId   string `json:"user_id"`
		CoinType uint32 `json:"coin_type"`
	}
	type Obj struct {
		Code   int                                 `json:"code"`
		ErrMsg string                              `json:"err_msg"`
		Data   *walletStruct.GetUserWalletResponse `json:"data"`
	}

	request := &req{
		UserId:   userID,
		CoinType: coinType,
	}
	var obj Obj
	res, err := s.API.PostWalletAPI(constant.GetUserWalletURL, request, constant.APITimeout)
	if err != nil {
		fmt.Println("GetUserWallet error ", err)
		return false, ""
	}
	err = utils.JsonStringToStruct(string(res), &obj)
	if err != nil {
		fmt.Println("GetUserWallet error ", err)
		return false, ""
	}
	if obj.ErrMsg != "" {
		return false, ""
	}
	return obj.Data.HasWallet, obj.Data.Address
}
func (s *SyncMgr) GetCoinStatuses() bool {
	type Currency struct {
		ID             int32  `json:"id"`
		CoinType       string `json:"coin_type"`
		LastEditedTime int64  `json:"last_edited_time"`
		Editor         string `json:"editor"`
		State          int32  `json:"state"`
	}
	type Data struct {
		Currencies []Currency `json:"currencies"`
	}
	type Obj struct {
		Code   int    `json:"code"`
		ErrMsg string `json:"err_msg"`
		Data   Data   `json:"data"`
	}
	var obj Obj
	res, err := http.Get(fmt.Sprintf("%s%s", s.API.BaseAddress, constant.GetCoinStatusesURL))
	if err != nil || len(res) == 0 {
		fmt.Println("failed to get coin statuses")
		return false
	}

	err = utils.JsonStringToStruct(string(res), &obj)
	if err != nil {
		fmt.Println("error", err)
	}
	if obj.Data.Currencies != nil {
		for _, v := range obj.Data.Currencies {
			localWallet := &model_struct.LocalWalletType{
				Status: uint8(v.State),
			}
			s.Db.UpdateLocalWalletTypes(localWallet, v.CoinType)
		}
		return true
	} else {
		return false
	}
}
func (s *SyncMgr) GetCoinRatio() string {
	type Coin struct {
		ID       int32   `json:"id"`
		CoinType string  `json:"coin_type"`
		Usd      float64 `json:"usd"`
		Yuan     float64 `json:"yuan"`
		Euro     float64 `json:"euro"`
	}
	type Data struct {
		Coins []Coin `json:"coins"`
	}
	type Obj struct {
		Code   int    `json:"code"`
		ErrMsg string `json:"err_msg"`
		Data   Data   `json:"data"`
	}
	var obj Obj
	res, err := http.Get(fmt.Sprintf("%s%s", s.API.BaseAddress, constant.GetCoinRatioURL))
	if err != nil {
		fmt.Println("failed to get coin ratio")
	}
	err = utils.JsonStringToStruct(string(res), &obj)
	if err != nil {
		fmt.Println("error", err)
	}
	// fmt.Println("Result", string(res), obj)
	return utils.StructToJsonString(obj.Data.Coins)
}
