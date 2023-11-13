package admin_api

import (
	"Share-Wallet/pkg/struct/common"
)

type TestRequest struct {
	OperationID string `json:"operationID" binding:"required" swaggertype:"string"`
	Name        string `json:"name" binding:"required" swaggertype:"string"`
}

type TestResponse struct {
	Name string `json:"name" swaggertype:"string"`
}

// AdminLoginRequest struct
type AdminLoginRequest struct {
	AdminName string `json:"admin_name" binding:"required" swaggertype:"string"`
	Secret    string `json:"secret" binding:"required" swaggertype:"string"`
}

// AdminLoginResponse struct
type AdminLoginResponse struct {
	Token              string `json:"token"`
	GAuthEnabled       bool   `json:"gAuthEnabled"`
	GAuthSetupRequired bool   `json:"gAuthSetupRequired"`
	GAuthSetupProvUri  string `json:"gAuthSetupProvUri"`
	// UserName           string   `json:"user_name"`
	// Role               string   `json:"role"`
	// Permissions        []string `json:"permissions"`
	User User `json:"user"`
}
type User struct {
	UserName    string   `json:"user_name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

// AdminPasswordChangeResponse struct
type AdminPasswordChangeResponse struct {
	Token           string `json:"token"`
	PasswordUpdated bool   `json:"password_updated"`
}

// AdminPasswordChangeRequest
type AdminPasswordChangeRequest struct {
	Secret    string `json:"secret" binding:"required" swaggertype:"string"`
	NewSecret string `json:"new_secret" binding:"required" swaggertype:"string"`
	// TOTP      string `json:"totp" swaggertype:"string"`
}

// GetAdminUsersRequest struct
type GetAdminUsersRequest struct {
	OperationID string `form:"operationID" binding:"required" swaggertype:"string"`
	Name        string `form:"admin_name" swaggertype:"string"`
	OrderBy     string `form:"order_by" binding:"omitempty,oneof=create_time:asc create_time:desc"`
	common.RequestPagination
}

// GetAdminUsersResponse struct
type GetAdminUsersResponse struct {
	UserNums int64 `json:"user_nums"`
	common.ResponsePagination
	Users []*AdminUser `json:"users"`
}

// AdminUser struct
type AdminUser struct {
	Id               int32  `json:"id"`
	UserName         string `json:"user_name"`
	Role             string `json:"role"`
	LastLoginIP      string `json:"last_loginIP"`
	LastLoginTime    int64  `json:"last_login_time"`
	Remarks          string `json:"remarks"`
	TwoFactorEnabled bool   `json:"two_factor_enabled"`
	Status           int32  `json:"status"`
}

// AdminUserRole struct
type AdminUserRole struct {
	Id              int32  `json:"id"`
	RoleName        string `json:"role_name"`
	RoleDescription string `json:"description"`
	RoleNumber      int32  `json:"role_number"`
	CreateTime      int64  `json:"create_time"`
	CreateUser      string `json:"create_user"`
	UpdateTime      int64  `json:"update_time"`
	UpdateUser      string `json:"update_user"`
	Remarks         string `json:"remarks"`
	Status          int32  `json:"status"`
}

// GetAdminUserRoleRequest struct
type GetAdminUserRoleRequest struct {
	OperationID string `form:"operationID" binding:"required" swaggertype:"string"`
	Name        string `form:"admin_name" swaggertype:"string"`
	OrderBy     string `form:"order_by" binding:"omitempty,oneof=create_time:asc create_time:desc" swaggertype:"string"`
	common.RequestPagination
}

// GetAdminUserRoleResponse struct
type GetAdminUserRoleResponse struct {
	RoleNums int64 `json:"role_nums"`
	common.ResponsePagination
	Roles []*AdminUserRole `json:"roles"`
}

// PostAdminUserRequest struct
type PostAdminUserRequest struct {
	OperationID      string `json:"operationID" binding:"omitempty" swaggertype:"string"`
	UserName         string `json:"user_name" binding:"required" swaggertype:"string"`
	Secret           string `json:"secret" binding:"required" swaggertype:"integer"`
	Role             string `json:"role" binding:"required" swaggertype:"string"`
	Remarks          string `json:"remarks" binding:"required" swaggertype:"string"`
	TwoFactorEnabled bool   `json:"two_factor_enabled" swaggertype:"string"`
	Status           int32  `json:"status" binding:"required" swaggertype:"string"`
}

// PostAdminRoleRequest struct
type PostAdminRoleRequest struct { // add parent_id
	OperationID string `json:"operationID" binding:"omitempty" swaggertype:"string"`
	RoleName    string `json:"role_name" binding:"required" swaggertype:"string"`
	Description string `json:"description" binding:"required" swaggertype:"integer"`
	ActionIDs   string `json:"actionIDs" binding:"required" swaggertype:"integer"`
	Remarks     string `json:"remarks" binding:"required" swaggertype:"string"`
	Status      int32  `json:"status" binding:"omitempty" swaggertype:"string"`
	UserName    string `json:"user_name" binding:"required" swaggertype:"string"`
}

// DeleteAdminRequest struct
type DeleteAdminRequest struct {
	OperationID string `json:"operationID" binding:"omitempty" swaggertype:"string"`
	UserName    string `json:"user_name" binding:"required"`
	DeleteUser  string `json:"delete_user" binding:"omitempty" swaggertype:"string"`
}

// UpdateAdminReq struct
type UpdateAdminReq struct {
	OperationID      string `json:"operationID" binding:"omitempty" swaggertype:"string"`
	UserName         string `json:"user_name" binding:"required" swaggertype:"string"`
	Password         string `json:"password" binding:"omitempty" swaggertype:"string"`
	RoleName         string `json:"role_name" binding:"required" swaggertype:"string"`
	Remarks          string `json:"remarks" binding:"omitempty" swaggertype:"string"`
	Status           int32  `json:"status" binding:"required" swaggertype:"string"`
	TwoFactorEnabled bool   `json:"two_factor_enabled"`
}

//UpdateAdminResp struct
type UpdateAdminResp struct {
	Name         string `json:"user_name" `
	AdminUpdated bool   `json:"admin_updated"`
}

// UpdateAdminRoleRequest struct
type UpdateAdminRoleRequest struct {
	OperationID string `json:"operationID" binding:"omitempty" swaggertype:"string"`
	RoleName    string `json:"role_name" binding:"required" swaggertype:"string"`
	Description string `json:"description" binding:"omitempty" swaggertype:"string"`
	ActionIDs   string `json:"actionIDs" binding:"required" swaggertype:"integer"`
	Remarks     string `json:"remarks" binding:"required" swaggertype:"string"`
	UpdateUser  string `json:"update_user" binding:"omitempty" swaggertype:"string"`
}

//UpdateAdminRoleResponse struct
type UpdateAdminRoleResponse struct {
	Name             string `json:"user_name" `
	AdminRoleUpdated bool   `json:"admin_role_updated"`
}

//DeleteRoleRequest struct
type DeleteRoleRequest struct {
	OperationID string `json:"operationID" binding:"omitempty" swaggertype:"string"`
	RoleName    string `json:"role_name" binding:"required"`
	DeleteUser  string `json:"delete_user" binding:"omitempty" swaggertype:"string"`
}
type ParamsTOTPVerify struct {
	TOTP        string `json:"totp" binding:"required"`
	OperationID string `json:"operationID"`
}

type AdminActions struct {
	ActionName string `json:"action_name"`
	Pid        string `json:"pid"`
	Remarks    string `json:"remarks"`
	Status     int32  `json:"status"`
}
type GetAdminUserActionsRequest struct {
	OperationID string `json:"operationID" binding:"omitempty" swaggertype:"string"`
}
type Children struct {
	ID   int32  `json:"ID" swaggertype:"integer"`
	Name string `json:"Name" swaggertype:"string"`
}
type GetAdminUserActionsResponse struct {
	ID       int32      `json:"ID" swaggertype:"integer"`
	Name     string     `json:"Name" swaggertype:"string"`
	Children []Children `json:"Children" swaggertype:"array"`
}

type GetAdminUserRequest struct {
	// UserName    string `json:"user_name" swaggertype:"string"`
	OperationID string `json:"operationID" binding:"omitempty" swaggertype:"string"`
}
type GetAdminUserResponse struct {
	UserName    string  `json:"user_name" swaggertype:"string"`
	RoleName    string  `json:"role_name" swaggertype:"string"`
	Permissions []int64 `json:"permissions" swaggertype:"array"`
}

type GetAccountInformationRequest struct {
	OperationID    string `form:"operation_id" swaggertype:"string"`
	OrderBy        string `form:"order_by" binding:"omitempty"`
	Sort           string `form:"sort" binding:"omitempty,oneof=asc desc"`
	MerchantUid    string `form:"merchant_uid" binding:"omitempty"`
	Uid            string `form:"uid" binding:"omitempty" swaggertype:"integer"`
	AccountAddress string `form:"account_address" binding:"omitempty"`
	CoinsType      string `form:"coins_type" binding:"omitempty"`
	AccountSource  string `form:"account_source" binding:"omitempty"`
	Filter         string `form:"filter" binding:"omitempty,oneof=last_login_time create_time"`
	From           string `form:"from" binding:"omitempty"`
	To             string `form:"to" binding:"omitempty"`
	common.RequestPagination
}

type AccountAsset struct {
	UsdAmount  float64 `json:"usd_amount" swaggertype:"number"`
	YuanAmount float64 `json:"yuan_amount" swaggertype:"number"`
	EuroAmount float64 `json:"euro_amount" swaggertype:"number"`
}

type AccountAddresses struct {
	BtcPublicAddress string `json:"btc_public_address" binding:"omitempty" swaggertype:"string"`
	EthPublicAddress string `json:"eth_public_address" binding:"omitempty" swaggertype:"string"`
	TrxPublicAddress string `json:"trx_public_address" binding:"omitempty" swaggertype:"string"`
	ErcPublicAddress string `json:"erc_public_address" binding:"omitempty" swaggertype:"string"`
	TrcPublicAddress string `json:"trc_public_address" binding:"omitempty" swaggertype:"string"`
}

type LoginInformation struct {
	LoginIp       string `json:"login_ip" swaggertype:"string"`
	LoginRegion   string `json:"login_region" swaggertype:"string"`
	LoginTerminal string `json:"login_terminal" swaggertype:"string"`
	LoginTime     int64  `json:"login_time" swaggertype:"integer"`
}
type Coin struct {
	Balance     string `json:"balance" swaggertype:"string"`
	UsdBalance  string `json:"usd_balance" swaggertype:"string"`
	YuanBalance string `json:"yuan_balance" swaggertype:"string"`
	EuroBalance string `json:"euro_balance" swaggertype:"string"`
}
type AccountInformation struct {
	ID          int64    `json:"id" swaggertype:"integer"`
	Uid         string   `json:"uid" swaggertype:"string"`
	MerchantUid string   `json:"merchant_uid" swaggertype:"string"`
	CoinsType   []string `json:"coins_type" swaggertype:"array"`

	Addresses AccountAddresses `json:"addresses"`

	Btc Coin `json:"btc"`
	Eth Coin `json:"eth"`
	Erc Coin `json:"erc"`
	Trx Coin `json:"trx"`
	Trc Coin `json:"trc"`

	TotalBalance float64 `json:"total_balance" swaggertype:"number"`

	AccountAssets            AccountAsset     `json:"account_assets"`
	AccountSource            string           `json:"account_source" swaggertype:"string"`
	CreationLoginInformation LoginInformation `json:"creation_login_information"`
	LastLoginInformation     LoginInformation `json:"last_login_information"`
}

type GetAccountInformationResponse struct {
	TotalNum int64 `json:"total_num"`
	common.ResponsePagination
	Accounts    []*AccountInformation `json:"accounts"`
	TotalAssets AccountAsset          `json:"total_assets"`
	BtcTotal    Coin                  `json:"btc_total"`
	EthTotal    Coin                  `json:"eth_total"`
	ErcTotal    Coin                  `json:"erc_total"`
	TrcTotal    Coin                  `json:"trc_total"`
	TrxTotal    Coin                  `json:"trx_total"`
}

//GetFundsLogRequest struct
type GetFundsLogRequest struct {
	OperationID     string `form:"operation_id" swaggertype:"string"`
	Uid             string `form:"uid" swaggertype:"string"`
	MerchantUid     string `form:"merchant_uid" swaggertype:"string"`
	From            string `form:"from" binding:"omitempty"`
	To              string `form:"to" binding:"omitempty"`
	TransactionType string `form:"transaction_type" binding:"omitempty"`
	UserAddress     string `form:"user_address" binding:"omitempty"`
	OppositeAddress string `form:"opposite_address" binding:"omitempty"`
	CoinsType       string `form:"coins_type" binding:"required"`
	State           string `form:"state" binding:"omitempty,oneof=fail success all"`
	Txid            string `form:"txid" binding:"omitempty"`
	common.RequestPagination
}

//FundsLog struct
type FundsLog struct {
	ID              int64  `json:"id" swaggertype:"integer"`
	Txid            string `json:"txid" binding:"omitempty" swaggertype:"string"`
	Uid             string `json:"uid" swaggertype:"string"`
	MerchantUid     string `json:"merchant_uid" swaggertype:"string"`
	TransactionType string `json:"transaction_type" swaggertype:"string"`
	UserAddress     string `json:"user_address" swaggertype:"string"`
	OppositeAddress string `json:"opposite_address" swaggertype:"string"`
	// BalanceBefore        float64 `json:"balance_before" swaggertype:"number"`
	// BalanceAfter         float64 `json:"balance_after" swaggertype:"number"`
	CoinType             string `json:"coin_type" swaggertype:"string"`
	AmountOfCoins        string `json:"amount_of_coins" swaggertype:"string"`
	UsdAmount            string `json:"usd_amount" swaggertype:"string"`
	YuanAmount           string `json:"yuan_amount" swaggertype:"string"`
	EuroAmount           string `json:"euro_amount" swaggertype:"string"`
	NetworkFee           string `json:"network_fee" swaggertype:"string"`
	UsdNetworkFee        string `json:"usd_network_fee" swaggertype:"string"`
	YuanNetworkFee       string `json:"yuan_network_fee" swaggertype:"string"`
	EuroNetworkFee       string `json:"euro_network_fee" swaggertype:"string"`
	TotalCoinsTransfered string `json:"total_coins_transfered" swaggertype:"string"`
	TotalUsdTransfered   string `json:"total_usd_transfered" swaggertype:"string"`
	TotalYuanTransfered  string `json:"total_yuan_transfered" swaggertype:"string"`
	TotalEuroTransfered  string `json:"total_euro_transfered" swaggertype:"string"`
	CreationTime         int64  `json:"creation_time" swaggertype:"integer"`
	State                string `json:"state" swaggertype:"string"`
	ConfirmationTime     int64  `json:"confirmation_time" swaggertype:"integer"`
}

//GetFundsLogResponse struct
type GetFundsLogResponse struct {
	TotalNum int64 `json:"total_num"`
	common.ResponsePagination
	Funds []*FundsLog `json:"funds_log"`
}

type GetReceiveDetailsRequest struct {
	OperationID      string `form:"operation_id" swaggertype:"string"`
	From             string `form:"from" binding:"omitempty"`
	To               string `form:"to" binding:"omitempty"`
	Uid              string `form:"uid" swaggertype:"string"`
	MerchantUid      string `form:"merchant_uid" swaggertype:"string"`
	ReceivingAddress string `form:"receiving_address" binding:"omitempty"`
	CoinsType        string `form:"coins_type" binding:"omitempty"`
	DepositAddress   string `form:"deposit_address" binding:"omitempty"`
	Txid             string `form:"txid" binding:"omitempty"`
	common.RequestPagination
}

type ReceiveDetails struct {
	ID               int64  `json:"id" swaggertype:"integer"`
	Uid              string `json:"uid" swaggertype:"string"`
	MerchantUid      string `json:"merchant_uid" swaggertype:"string"`
	ReceivingAddress string `json:"receiving_address" swaggertype:"string"`
	CoinType         string `json:"coin_type" swaggertype:"string"`
	AmountOfReceived string `json:"amount_of_received" swaggertype:"string"`
	UsdAmount        string `json:"usd_amount" swaggertype:"string"`
	YuanAmount       string `json:"yuan_amount" swaggertype:"string"`
	EuroAmount       string `json:"euro_amount" swaggertype:"string"`
	DepositAddress   string `json:"deposit_address" swaggertype:"string"`
	Txid             string `json:"txid" binding:"omitempty"`
	CreationTime     int64  `json:"creation_time" swaggertype:"integer"`
}
type GetRecieveDetailsResponse struct {
	TotalNum                int64   `json:"total_num"`
	TotalAmountReceivedUsd  float64 `json:"total_amount_received_usd"`
	TotalAmountReceivedYuan float64 `json:"total_amount_received_yuan"`
	TotalAmountReceivedEuro float64 `json:"total_amount_received_euro"`
	GrandTotalUsd           float64 `json:"grand_total_usd"`
	GrandTotalYuan          float64 `json:"grand_total_yuan"`
	GrandTotalEuro          float64 `json:"grand_total_euro"`
	common.ResponsePagination
	ReceiveDetails []*ReceiveDetails `json:"receive_details"`
}

type GetTransferDetailsRequest struct {
	OperationID      string `form:"operation_id" swaggertype:"string"`
	OrderBy          string `form:"order_by" binding:"omitempty"`
	From             string `form:"from" binding:"omitempty"`
	To               string `form:"to" binding:"omitempty"`
	Uid              string `form:"uid" swaggertype:"string"`
	MerchantUid      string `form:"merchant_uid" swaggertype:"string"`
	TransferAddress  string `form:"transfer_address" binding:"omitempty"`
	CoinsType        string `form:"coins_type" binding:"omitempty"`
	State            string `form:"state" binding:"omitempty,oneof=fail pending success all"`
	ReceivingAddress string `form:"receiving_address" binding:"omitempty"`
	Txid             string `form:"txid" binding:"omitempty"`
	common.RequestPagination
}

type GetTransferDetailsResponse struct {
	TotalNum                  int64   `json:"total_num"`
	TotalAmountTransferedUsd  float64 `json:"total_amount_transfered_usd"`
	TotalAmountTransferedYuan float64 `json:"total_amount_transfered_yuan"`
	TotalAmountTransferedEuro float64 `json:"total_amount_transfered_euro"`
	TotalFeeAmountUsd         float64 `json:"total_fee_amount_usd"`
	TotalFeeAmountYuan        float64 `json:"total_fee_amount_Yuan"`
	TotalFeeAmountEuro        float64 `json:"total_fee_amount_euro"`
	GrandTotalUsd             float64 `json:"grand_total_usd"`
	GrandTotalYuan            float64 `json:"grand_total_yuan"`
	GrandTotalEuro            float64 `json:"grand_total_euro"`
	TotalTransferUsd          float64 `json:"total_transfer_usd"`
	TotalTransferYuan         float64 `json:"total_transfer_yuan"`
	TotalTransferEuro         float64 `json:"total_transfer_euro"`
	common.ResponsePagination
	TransferDetails []*FundsLog `json:"transfer_details"`
}

type ResetGoogleKeyRequest struct {
	OperationID string `json:"operationID"  swaggertype:"string"`
	UserName    string `json:"user_name" binding:"required" swaggertype:"string"`
}

type GetRoleActionsRequest struct {
	OperationID string `json:"operationID" swaggertype:"string"`
	RoleName    string `json:"role_name" binding:"required" swaggertype:"string"`
}

type GetRoleActionsResponse struct {
	Actions []int64 `json:"actions"`
}

type GetCurrenciesRequest struct {
	OperationID string `form:"operationID" swaggertype:"string"`
	common.RequestPagination
}

type Currency struct {
	ID             int32  `json:"id" swaggertype:"integer"`
	CoinType       string `json:"coin_type" swaggertype:"string"`
	LastEditedTime int64  `json:"last_edited_time" swaggertype:"integer"`
	Editor         string `json:"editor" swaggertype:"string"`
	State          int32  `json:"state" swaggertype:"integer"`
}
type GetCurrenciesResponse struct {
	Currencies []*Currency `json:"currencies"`
	common.ResponsePagination
	TotalNum int64 `json:"total_num"`
}

type UpdateCurrencyRequest struct {
	OperationID string `json:"operationID" swaggertype:"string"`
	CurrencyId  int32  `json:"currency_id" binding:"required" swaggertype:"integer"`
	State       int32  `json:"state" swaggertype:"integer"`
}

type CreateAccountInfoRequest struct {
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
type CreateAccountInfoResponse struct {
	Uid string `json:"uid" swaggertype:"string"`
}

type CreateFundLogRequest struct {
	OperationID     string `json:"operationID" binding:"omitempty" swaggertype:"string"`
	Uid             string `json:"uid" binding:"required" swaggertype:"string"`
	MerchantUid     string `json:"merchant_uid" binding:"required" swaggertype:"string"`
	Txid            string `json:"txid" binding:"required" swaggertype:"string"`
	TransactionType string `json:"transaction_type" binding:"omitempty" swaggertype:"string"`
	UserAddress     string `json:"user_address" binding:"required" swaggertype:"string"`
	OppositeAddress string `json:"opposite_address" binding:"required" swaggertype:"string"`
	// BalanceBefore   float64 `json:"balance_before" binding:"required" swaggertype:"number"`
	// BalanceAfter    float64 `json:"balance_after" binding:"omitempty" swaggertype:"number"`
	CoinType      string  `json:"coin_type" binding:"required" swaggertype:"string"`
	AmountOfCoins float64 `json:"amount_of_coins" binding:"required" swaggertype:"number"`
	NetworkFee    float64 `json:"network_fee" binding:"required" swaggertype:"number"`
	CreationTime  int64   `json:"creation_time" swaggertype:"integer"`
	// State            string  `json:"state" binding:"omitempty,oneof=fail success pending"`
	// ConfirmationTime int64 `json:"confirmation_time" swaggertype:"integer"`
}
type CreateFundLogResponse struct {
}

type UpdateFundLogRequest struct {
	OperationID      string  `json:"operationID" binding:"omitempty" swaggertype:"string"`
	Txid             string  `json:"txid" binding:"required" swaggertype:"string"`
	BalanceAfter     float64 `json:"balance_after" binding:"required" swaggertype:"number"`
	State            string  `json:"state" binding:"omitempty,oneof=fail success pending"`
	ConfirmationTime int64   `json:"confirmation_time" swaggertype:"integer"`
}

type GetOperationalReportRequest struct {
	OperationID string `form:"operationID" binding:"omitempty" swaggertype:"string"`
	From        string `form:"from" binding:"omitempty"`
	To          string `form:"to" binding:"omitempty"`
	common.RequestPagination
}

type Statistics struct {
	AccountAssets AccountAsset `json:"total_assets"`
	Btc           Coin         `json:"btc"`
	Eth           Coin         `json:"eth"`
	Erc           Coin         `json:"erc"`
	Trc           Coin         `json:"trc"`
	Trx           Coin         `json:"trx"`
}
type OperationalReport struct {
	Date            string     `json:"confirmation_time" swaggertype:"string "`
	NewUsers        int64      `json:"new_users" swaggertype:"integer"`
	TotalReceived   Statistics `json:"total_received"`
	TotalTransfered Statistics `json:"total_transfered"`
	NetworkFee      Statistics `json:"network_fee"`
}

type GetOperationalReportResponse struct {
	OperationalReports []OperationalReport `json:"operational_reports"`
	GrandTotals        OperationalReport   `json:"grand_totals"`
	TotalAssets        AccountAsset        `json:"total_assets"`
	TotalNum           int32               `json:"total_num"`
	TotalUsers         int64               `json:"total_users"`
	NewUsersToday      int64               `json:"new_users_today"`
	common.ResponsePagination
}

type ConfirmTransactionRequest struct {
	OperationID string `json:"operation_id" swaggertype:"string"`
	CoinsType   uint8  `json:"coin_type" binding:"omitempty"`
	TxHashId    string `json:"tx_hash_id" binding:"omitempty"`
}

type ConfirmTransactionResponse struct {
	TransferDetail *FundsLog `json:"transfer_detail"`
}

type UpdateAccountBalanceRequest struct {
	OperationID string `json:"operation_id" swaggertype:"string"`
	MerchantUid string `json:"merchant_uid" binding:"required"`
	Uuid        string `json:"uuid" binding:"required"`
}
type Operation struct {
	Date        string  `json:"date"`
	BtcTransfer float64 `json:"btc_transfer"`
	EthTransfer float64 `json:"eth_transfer"`
	ErcTransfer float64 `json:"erc_transfer"`
	TrcTransfer float64 `json:"trc_transfer"`
	TrxTransfer float64 `json:"trx_transfer"`
	BtcReceived float64 `json:"btc_received"`
	EthReceived float64 `json:"eth_received"`
	ErcReceived float64 `json:"erc_received"`
	TrcReceived float64 `json:"trc_received"`
	TrxReceived float64 `json:"trx_received"`
	BtcFee      float64 `json:"btc_fee"`
	EthFee      float64 `json:"eth_fee"`
	ErcFee      float64 `json:"erc_fee"`
	TrcFee      float64 `json:"trc_fee"`
	TrxFee      float64 `json:"trx_fee"`
}

type TotalCoins struct {
	Btc    float64 `json:"btc"`
	Eth    float64 `json:"eth"`
	Erc    float64 `json:"erc"`
	Trc    float64 `json:"trc"`
	Trx    float64 `json:"trx"`
	BtcFee float64 `json:"btc_fee"`
	EthFee float64 `json:"eth_fee"`
	ErcFee float64 `json:"erc_fee"`
	TrxFee float64 `json:"trx_fee"`
	TrcFee float64 `json:"trc_fee"`
}
