package wallet_api

import "Share-Wallet/pkg/struct/common"

type TestRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	Name        string `json:"name" binding:"required"`
}

type TestResponse struct {
	Name string `json:"name"`
}

type GetSupportTokenAddressesRequest struct {
	OperationID string `json:"operationID" binding:"required"`
}

type SupportTokenAddress struct {
	BelongCoin      uint8  `json:"belong_coin" binding:"required"`
	CoinType        uint8  `json:"coin_type" binding:"required"`
	ContractAddress string `json:"contract_address" binding:"required"`
}

type GetSupportTokenAddressesResponse struct {
	AddressList []*SupportTokenAddress `json:"address_list"`
}

type GetSupportTokenAddressesBaseResp struct {
	Code   int32                            `json:"code"`
	ErrMsg string                           `json:"err_msg"`
	Data   GetSupportTokenAddressesResponse `json:"data"`
}

type UpdateAccountInfoRequest struct {
	OperationID       string `json:"operationID" binding:"omitempty" swaggertype:"string"`
	Uid               string `json:"uid" binding:"required" swaggertype:"string"`
	MerchantUid       string `json:"merchant_uid" binding:"required" swaggertype:"string"`
	LastLoginIp       string `json:"last_login_ip" swaggertype:"string"`
	LastLoginRegion   string `json:"last_login_region" swaggertype:"string"`
	LastLoginTerminal string `json:"last_login_terminal" swaggertype:"string"`
	LastLoginTime     int64  `json:"last_login_time" swaggertype:"integer"`
}

//GetFundsLogRequest struct
type GetFundsLogRequest struct {
	OperationID string `json:"operation_id" swaggertype:"string"`
	UID         string `json:"uid" swaggertype:"string"`
	PageSize    int    `json:"page_size" binding:"omitempty,min=1,max=9223372036854775807" swaggertype:"integer"`
	Page        int    `json:"page" binding:"omitempty,min=-1,max=9223372036854775807" swaggertype:"integer"`
}

//FundsLog struct
type FundsLog struct {
	ID                      int64   `json:"id" swaggertype:"integer"`
	Txid                    string  `json:"txid" binding:"omitempty" swaggertype:"string"`
	Uid                     int64   `json:"uid" swaggertype:"integer"`
	MerchantUid             string  `json:"merchant_uid" swaggertype:"string"`
	TransactionType         string  `json:"transaction_type" swaggertype:"string"`
	UserAddress             string  `json:"user_address" swaggertype:"string"`
	OppositeAddress         string  `json:"opposite_address" swaggertype:"string"`
	BalanceBefore           float64 `json:"balance_before" swaggertype:"number"`
	BalanceAfter            float64 `json:"balance_after" swaggertype:"number"`
	CoinType                string  `json:"coin_type" swaggertype:"string"`
	AmountOfCoins           float64 `json:"amount_of_coins" swaggertype:"number"`
	UsdAmount               float64 `json:"usd_amount" swaggertype:"number"`
	YuanAmount              float64 `json:"yuan_amount" swaggertype:"number"`
	EuroAmount              float64 `json:"euro_amount" swaggertype:"number"`
	NetworkFee              float64 `json:"network_fee" swaggertype:"number"`
	UsdNetworkFee           float64 `json:"usd_network_fee" swaggertype:"number"`
	YuanNetworkFee          float64 `json:"yuan_network_fee" swaggertype:"number"`
	EuroNetworkFee          float64 `json:"euro_network_fee" swaggertype:"number"`
	TotalCoinsTransfered    float64 `json:"total_coins_transfered" swaggertype:"number"`
	TotalUsdTransfered      float64 `json:"total_usd_transfered" swaggertype:"number"`
	TotalYuanTransfered     float64 `json:"total_yuan_transfered" swaggertype:"number"`
	TotalEuroTransfered     float64 `json:"total_euro_transfered" swaggertype:"number"`
	CreationTime            int64   `json:"creation_time" swaggertype:"integer"`
	State                   string  `json:"state" swaggertype:"string"`
	ConfirmationTime        int64   `json:"confirmation_time" swaggertype:"integer"`
	GasLimit                uint32  `json:"gas_limit" swaggertype:"integer"`
	GasPrice                uint64  `json:"gas_price" swaggertype:"integer"`
	GasUsed                 uint64  `json:"gas_used" swaggertype:"integer"`
	ConfirmationBlockNumber uint64  `json:"confirm_block_number" swaggertype:"integer"`
}

//GetFundsLogResponse struct
type GetFundsLogResponse struct {
	TotalNum int64 `json:"total_num"`
	common.ResponsePagination
	Funds []*FundsLog `json:"funds_log"`
}

type GetCoinStatusesRequest struct {
	OperationID string `json:"operationID" binding:"omitempty" swaggertype:"string"`
}

type Currency struct {
	ID             int32  `json:"id" swaggertype:"integer"`
	CoinType       string `json:"coin_type" swaggertype:"string"`
	LastEditedTime int64  `json:"last_edited_time" swaggertype:"integer"`
	Editor         string `json:"editor" swaggertype:"string"`
	State          int32  `json:"state" swaggertype:"integer"`
}
type GetCoinsStatusesResponse struct {
	Currencies []*Currency `json:"currencies"`
}

type GetCoinRatioRequest struct {
	OperationID string `json:"operationID" binding:"omitempty" swaggertype:"string"`
}

type Coin struct {
	ID       int32   `json:"id" swaggertype:"integer"`
	CoinType string  `json:"coin_type" swaggertype:"string"`
	Usd      float64 `json:"usd" swaggertype:"number"`
	Yuan     float64 `json:"yuan" swaggertype:"number"`
	Euro     float64 `json:"euro" swaggertype:"number"`
}
type GetCoinRatioResponse struct {
	Coins []*Coin `json:"coins"`
}
type GetUserWalletRequest struct {
	OperationID string `json:"operationID" binding:"omitempty" swaggertype:"string"`
	UserId      string `json:"user_id" binding:"required" swaggertype:"string"`
	CoinType    uint32 `json:"coin_type" binding:"omitempty" swaggertype:"integer"`
}
type GetUserWalletResponse struct {
	HasWallet bool   `json:"has_wallet" swaggertype:"boolean"`
	Address   string `json:"address" swaggertype:"string"`
}

type UpdateCoinRatesRequest struct {
	OperationID string `json:"operationID" binding:"omitempty" swaggertype:"string"`
}

type Rate struct {
	USD string `json:"USD"`
	CNY string `json:"CNY"`
	EUR string `json:"EUR"`
}
type Data struct {
	Currency string `json:"currency"`
	Rates    Rate   `json:"rates"`
}
type Obj struct {
	Data Data `json:"data"`
}

type GetAccountBalanceRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	CoinType    int    `json:"coin_type" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

type GetAccountBalanceResponse struct {
	Balance string `json:"balance"`
}

type TempCoinRate struct {
	Usd float64 `json:"usd"`
	Cny float64 `json:"cny"`
	Eur float64 `json:"eur"`
}

type TempRateObj struct {
	Tron TempCoinRate `json:"tron"`
}
