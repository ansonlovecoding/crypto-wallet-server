package db

import (
	"time"

	"github.com/shopspring/decimal"
)

// UserAccounts: store the user wallet address
type UserAccounts struct {
	ID            int    `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	MerchantUid   string `gorm:"column:merchant_uid;type:varchar(255);index:merchant_uid;NOT NULL" json:"merchant_uid"`
	PublicKey     string `gorm:"column:public_key;type:varchar(255);index:public_key;NOT NULL" json:"public_key"`
	WalletAddress string `gorm:"column:wallet_address;type:varchar(255);index:wallet_address;NOT NULL" json:"wallet_address"`
	Source        int    `gorm:"column:source;type:tinyint(1);default:1;comment:1 created by our platform, 2 from other platform;NOT NULL" json:"source"`
	LoginTime     int64  `gorm:"column:login_time;type:bigint(11)" json:"login_time"`
	LoginIp       string `gorm:"column:login_ip;type:varchar(255)" json:"login_ip"`
	LoginArea     string `gorm:"column:login_area;type:varchar(255)" json:"login_area"`
	LoginDevice   int    `gorm:"column:login_device;type:tinyint(1);comment:1 ios, 2 android" json:"login_device"`
	CreateUser    string `gorm:"column:create_user;type:varchar(64)" json:"create_user"`
	CreateTime    int64  `gorm:"column:create_time;type:bigint(11)" json:"create_time"`
	UpdateUser    string `gorm:"column:update_user;type:varchar(64)" json:"update_user"`
	UpdateTime    int64  `gorm:"column:update_time;type:bigint(11)" json:"update_time"`
	DeleteUser    string `gorm:"column:delete_user;type:varchar(64)" json:"delete_user"`
	DeleteTime    int64  `gorm:"column:delete_time;type:bigint(11);default:0;NOT NULL" json:"delete_time"`
}

func (UserAccounts) TableName() string {
	return "w_user_accounts"
}

// UserBalances: store the user balance for each currency
type UserBalances struct {
	ID         int     `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	AccountId  int     `gorm:"column:account_id;type:bigint(11);index:account_id;NOT NULL" json:"account_id"`
	CurrencyId int     `gorm:"column:currency_id;type:bigint(11);index:currency_id;NOT NULL" json:"currency_id"`
	Address    string  `gorm:"column:address;type:varchar(255);index:address;NOT NULL" json:"address"`
	Amount     float64 `gorm:"column:amount;type:decimal(28,8);default:0;NOT NULL" json:"amount"`
	CreateUser string  `gorm:"column:create_user;type:varchar(64)" json:"create_user"`
	CreateTime int64   `gorm:"column:create_time;type:bigint(11)" json:"create_time"`
	UpdateUser string  `gorm:"column:update_user;type:varchar(64)" json:"update_user"`
	UpdateTime int64   `gorm:"column:update_time;type:bigint(11)" json:"update_time"`
	DeleteUser string  `gorm:"column:delete_user;type:varchar(64)" json:"delete_user"`
	DeleteTime int64   `gorm:"column:delete_time;type:bigint(11);default:0;NOT NULL" json:"delete_time"`
}

func (UserBalances) TableName() string {
	return "w_user_balances"
}

// LocalCurrencyRates: store the rate for how many usdt in one local currency, like 1 usd = 1 usdt
type LocalCurrencyRates struct {
	ID         int     `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	Name       string  `gorm:"column:name;type:varchar(64);NOT NULL" json:"name"`
	Rate       float64 `gorm:"column:rate;type:decimal(28,8)" json:"rate"`
	CreateUser string  `gorm:"column:create_user;type:varchar(64)" json:"create_user"`
	CreateTime int64   `gorm:"column:create_time;type:bigint(11)" json:"create_time"`
	UpdateUser string  `gorm:"column:update_user;type:varchar(64)" json:"update_user"`
	UpdateTime int64   `gorm:"column:update_time;type:bigint(11)" json:"update_time"`
	DeleteUser string  `gorm:"column:delete_user;type:varchar(64)" json:"delete_user"`
	DeleteTime int64   `gorm:"column:delete_time;type:bigint(11);default:0;NOT NULL" json:"delete_time"`
}

func (LocalCurrencyRates) TableName() string {
	return "w_local_currency_rates"
}

// WalletCurrency: store the data for digtal currencies, include the rate (1 eth = 1160usdt) and status(if it was closed, it will not show in App)
type WalletCurrency struct {
	ID         int     `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	Name       string  `gorm:"column:name;type:varchar(64);NOT NULL" json:"name"`
	Rate       float64 `gorm:"column:rate;type:decimal(28,8);" json:"rate"`
	Status     int     `gorm:"column:status;type:tinyint(1);default:1;comment:1 opened, 2 closed;NOT NULL" json:"status"`
	CreateUser string  `gorm:"column:create_user;type:varchar(64)" json:"create_user"`
	CreateTime int64   `gorm:"column:create_time;type:bigint(11)" json:"create_time"`
	UpdateUser string  `gorm:"column:update_user;type:varchar(64)" json:"update_user"`
	UpdateTime int64   `gorm:"column:update_time;type:bigint(11)" json:"update_time"`
	DeleteUser string  `gorm:"column:delete_user;type:varchar(64)" json:"delete_user"`
	DeleteTime int64   `gorm:"column:delete_time;type:bigint(11);default:0;NOT NULL" json:"delete_time"`
}

func (WalletCurrency) TableName() string {
	return "w_wallet_currency"
}

// Statistics: store the statistics for each currency and every date
type Statistics struct {
	ID            int     `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	CurrencyId    int     `gorm:"column:currency_id;type:bigint(11);index:currency_id;NOT NULL" json:"currency_id"`
	Type          int     `gorm:"column:type;type:tinyint(1);comment:1 received, 2 sent;NOT NULL" json:"type"`
	Amount        float64 `gorm:"column:amount;type:decimal(28,8);default:0;NOT NULL" json:"amount"`
	StatisticDate int64   `gorm:"column:statistic_date;type:date;index:statistic_date;NOT NULL" json:"statistic_date"`
	CreateUser    string  `gorm:"column:create_user;type:varchar(64)" json:"create_user"`
	CreateTime    int64   `gorm:"column:create_time;type:bigint(11)" json:"create_time"`
	UpdateUser    string  `gorm:"column:update_user;type:varchar(64)" json:"update_user"`
	UpdateTime    int64   `gorm:"column:update_time;type:bigint(11)" json:"update_time"`
	DeleteUser    string  `gorm:"column:delete_user;type:varchar(64)" json:"delete_user"`
	DeleteTime    int64   `gorm:"column:delete_time;type:bigint(11);default:0;NOT NULL" json:"delete_time"`
}

func (Statistics) TableName() string {
	return "w_statistics"
}

// AdminUser: store the users for management platform
type AdminUser struct {
	ID                int    `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	UserName          string `gorm:"column:user_name;type:varchar(64);index:user_name" json:"user_name"`
	NickName          string `gorm:"column:nick_name;type:varchar(64);NOT NULL" json:"nick_name"`
	Password          string `gorm:"column:password;type:varchar(255);NOT NULL" json:"password"`
	Salt              string `gorm:"column:salt;type:varchar(32);NOT NULL" json:"salt"`
	Google2fSecretKey string `gorm:"column:google_2f_secret_key;type:varchar(32)" json:"google_2f_secret_key"`
	GoogleStatus      int    `gorm:"column:google_status;type:tinyint(1);default:1;comment:1 opened, 2 closed;NOT NULL" json:"google_status"`
	Status            int    `gorm:"column:status;type:tinyint(1);default:1;comment:1 opened, 2 banned;NOT NULL" json:"status"`
	RoleId            int    `gorm:"column:role_id;type:bigint(11)" json:"role_id"`
	LoginIp           string `gorm:"column:login_ip;type:varchar(255)" json:"login_ip"`
	LoginTime         int64  `gorm:"column:login_time;type:bigint(11)" json:"login_time"`
	Remarks           string `gorm:"column:remarks;type:varchar(255)" json:"remarks"`
	CreateUser        string `gorm:"column:create_user;type:varchar(64)" json:"create_user"`
	CreateTime        int64  `gorm:"column:create_time;type:bigint(11)" json:"create_time"`
	UpdateUser        string `gorm:"column:update_user;type:varchar(64)" json:"update_user"`
	UpdateTime        int64  `gorm:"column:update_time;type:bigint(11)" json:"update_time"`
	DeleteUser        string `gorm:"column:delete_user;type:varchar(64)" json:"delete_user"`
	DeleteTime        int64  `gorm:"column:delete_time;type:bigint(11);default:0;NOT NULL" json:"delete_time"`
	TwoFactorEnabled  bool   `gorm:"column:two_factor_enabled;default:0"`
	User2FAuthEnable  bool   `gorm:"column:user_two_factor_control_status;default:1"`
}

func (AdminUser) TableName() string {
	return "w_admin_user"
}

// AdminRole: store the roles for admin users
type AdminRole struct {
	ID          int    `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	RoleName    string `gorm:"column:role_name;type:varchar(64);NOT NULL" json:"role_name"`
	Description string `gorm:"column:description;type:varchar(255)" json:"description"`
	Remarks     string `gorm:"column:remarks;type:varchar(255)" json:"remarks"`
	Status      int    `gorm:"column:status;type:tinyint(1);default:1;comment:1 opened, 2 banned;NOT NULL" json:"status"`
	MemberNum   int    `gorm:"column:member_num;type:bigint(11);default:0;NOT NULL" json:"member_num"`
	ActionIds   string `gorm:"column:action_ids;type:varchar(255)" json:"action_ids"`
	CreateUser  string `gorm:"column:create_user;type:varchar(64)" json:"create_user"`
	CreateTime  int64  `gorm:"column:create_time;type:bigint(11)" json:"create_time"`
	UpdateUser  string `gorm:"column:update_user;type:varchar(64)" json:"update_user"`
	UpdateTime  int64  `gorm:"column:update_time;type:bigint(11)" json:"update_time"`
	DeleteUser  string `gorm:"column:delete_user;type:varchar(64)" json:"delete_user"`
	DeleteTime  int64  `gorm:"column:delete_time;type:bigint(11);default:0;NOT NULL" json:"delete_time"`
}

func (AdminRole) TableName() string {
	return "w_admin_role"
}

// AdminActions: store the page url and api url for each action
type AdminActions struct {
	ID         int    `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	ActionName string `gorm:"column:action_name;type:varchar(64);NOT NULL" json:"action_name"`
	Remarks    string `gorm:"column:remarks;type:varchar(255)" json:"remarks"`
	Status     int    `gorm:"column:status;type:tinyint(1);default:1;comment:1 opened, 2 banned;NOT NULL" json:"status"`
	Pid        int    `gorm:"column:pid;type:bigint(11);default:0;comment:parent action id;NOT NULL" json:"pid"`
	ApiUrl     string `gorm:"column:api_url;type:varchar(255)" json:"api_url"`
	PageUrl    string `gorm:"column:page_url;type:varchar(255)" json:"page_url"`
	CreateUser string `gorm:"column:create_user;type:varchar(64)" json:"create_user"`
	CreateTime int64  `gorm:"column:create_time;type:bigint(11)" json:"create_time"`
	UpdateUser string `gorm:"column:update_user;type:varchar(64)" json:"update_user"`
	UpdateTime int64  `gorm:"column:update_time;type:bigint(11)" json:"update_time"`
	DeleteUser string `gorm:"column:delete_user;type:varchar(64)" json:"delete_user"`
	DeleteTime int64  `gorm:"column:delete_time;type:bigint(11);default:0;NOT NULL" json:"delete_time"`
}

func (AdminActions) TableName() string {
	return "w_admin_actions"
}

// EthDetailTX is an object representing the database table.
type EthDetailTX struct {
	ID                      int64           `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	UUID                    string          `gorm:"column:uuid;type:varchar(255)" json:"uuid"`
	SenderAccount           string          `gorm:"column:sender_account;type:varchar(255)" json:"sender_account"`
	SenderAddress           string          `gorm:"column:sender_address;type:varchar(255)" json:"sender_address"`
	ReceiverAccount         string          `gorm:"column:receiver_account;type:varchar(255)" json:"receiver_account"`
	ReceiverAddress         string          `gorm:"column:receiver_address;type:varchar(255)" json:"receiver_address"`
	Amount                  decimal.Decimal `gorm:"column:amount;type:decimal(36,0)" json:"amount"`
	Fee                     decimal.Decimal `gorm:"column:fee;type:decimal(36,0)" json:"fee"`
	GasLimit                uint64          `gorm:"column:gas_limit;type:bigint(11)" json:"gas_limit"`
	Nonce                   uint64          `gorm:"column:nonce;type:bigint(11)" json:"nonce"`
	SentHashTX              string          `gorm:"column:sent_hash_tx;type:varchar(255);index:hash_unique,unique" json:"sent_hash_tx"`
	SentUpdatedAt           *time.Time      `gorm:"column:sent_updated_at" json:"sent_updated_at,omitempty"`
	ConfirmTime             *time.Time      `gorm:"column:confirm_time" json:"confirm_time,omitempty"`
	Status                  int8            `gorm:"column:status;comment:0 pending,1 success,2 failed" json:"status,omitempty"`
	GasPrice                decimal.Decimal `gorm:"column:gas_price;type:decimal(36,0)" json:"gas_price"`
	GasUsed                 uint64          `gorm:"column:gas_used;type:bigint(11)" json:"gas_used"`
	ConfirmationBlockNumber string          `gorm:"column:confirm_block_number;type:varchar(255)" json:"confirm_block_number"`
	CheckTimes              uint64          `gorm:"column:check_times;type:bigint(11);comment:record the times for checking confirmation;default:0" json:"check_times"`
	ConfirmStatus           uint8           `gorm:"column:confirm_status;type:tinyint(1);comment:status of confirmation process, 0 not yet checking, 1 pending, 2 completed;default:0" json:"confirm_status"`
	CoinType                string          `gorm:"column:coin_type;type:varchar(255);index:hash_unique,unique" json:"coin_type"`
}

func (EthDetailTX) TableName() string {
	return "w_eth_detail_tx"
}

// TronDetailTX is an object representing the database table.
type TronDetailTX struct {
	ID                      int64           `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	UUID                    string          `gorm:"column:uuid;type:varchar(255)" json:"uuid"`
	SenderAccount           string          `gorm:"column:sender_account;type:varchar(255)" json:"sender_account"`
	SenderAddress           string          `gorm:"column:sender_address;type:varchar(255)" json:"sender_address"`
	ReceiverAccount         string          `gorm:"column:receiver_account;type:varchar(255)" json:"receiver_account"`
	ReceiverAddress         string          `gorm:"column:receiver_address;type:varchar(255)" json:"receiver_address"`
	Amount                  decimal.Decimal `gorm:"column:amount;type:decimal(36,0)" json:"amount"`
	Fee                     decimal.Decimal `gorm:"column:fee;type:decimal(36,0)" json:"fee"`
	EnergyUsed              decimal.Decimal `gorm:"column:energy_used;type:decimal(36,0)" json:"energy_used"`
	NetUsed                 decimal.Decimal `gorm:"column:net_used;type:decimal(36,0)" json:"net_used"`
	SentHashTX              string          `gorm:"column:sent_hash_tx;type:varchar(255);index:hash_unique,unique" json:"sent_hash_tx"`
	SentUpdatedAt           *time.Time      `gorm:"column:sent_updated_at" json:"sent_updated_at,omitempty"`
	ConfirmTime             *time.Time      `gorm:"column:confirm_time" json:"confirm_time,omitempty"`
	Status                  int8            `gorm:"column:status;comment:0 pending,1 success,2 failed" json:"status,omitempty"`
	ConfirmationBlockNumber string          `gorm:"column:confirm_block_number;type:varchar(255)" json:"confirm_block_number"`
	CheckTimes              uint64          `gorm:"column:check_times;type:bigint(11);comment:record the times for checking confirmation;default:0" json:"check_times"`
	ConfirmStatus           uint8           `gorm:"column:confirm_status;type:tinyint(1);comment:status of confirmation process, 0 not yet checking, 1 pending, 2 completed;default:0" json:"confirm_status"`
	CoinType                string          `gorm:"column:coin_type;type:varchar(255);index:hash_unique,unique" json:"coin_type"`
}

func (TronDetailTX) TableName() string {
	return "w_tron_detail_tx"
}

type AccountInformation struct {
	ID                    int64   `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	UUID                  string  `gorm:"column:uuid;type:varchar(255)" json:"uuid"`
	MerchantUid           string  `gorm:"column:merchant_uid;type:varchar(255)" json:"merchant_id"`
	BtcPublicAddress      string  `gorm:"column:btc_public_address;type:varchar(255)" json:"btc_public_address"`
	EthPublicAddress      string  `gorm:"column:eth_public_address;type:varchar(255)" json:"eth_public_address"`
	TrxPublicAddress      string  `gorm:"column:trx_public_address;type:varchar(255)" json:"trx_public_address"`
	ErcPublicAddress      string  `gorm:"column:erc_public_address;type:varchar(255)" json:"erc_public_address"`
	TrcPublicAddress      string  `gorm:"column:trc_public_address;type:varchar(255)" json:"trc_public_address"`
	BtcBalance            float64 `gorm:"column:btc_balance;type:decimal(28,8)" json:"btc_balance"`
	EthBalance            float64 `gorm:"column:eth_balance;type:decimal(28,8)" json:"eth_balance"`
	TrxBalance            float64 `gorm:"column:trx_balance;type:decimal(28,8)" json:"trx_balance"`
	ErcBalance            float64 `gorm:"column:erc_balance;type:decimal(28,8)" json:"erc_balance"`
	TrcBalance            float64 `gorm:"column:trc_balance;type:decimal(28,8)" json:"trc_balance"`
	AccountSource         string  `gorm:"column:account_source;type:varchar(255)" json:"account_source"`
	CreationLoginIp       string  `gorm:"column:creation_login_ip;type:varchar(255)" json:"creation_login_ip"`
	CreationLoginRegion   string  `gorm:"column:creation_login_region;type:varchar(255)" json:"creation_login_region"`
	CreationLoginTerminal string  `gorm:"column:creation_login_terminal;type:varchar(255)" json:"creation_login_terminal"`
	CreationTime          int64   `gorm:"column:creation_time;type:bigint(11)" json:"creation_time"`
	LastLoginIp           string  `gorm:"column:last_login_ip;type:varchar(255)" json:"last_login_ip"`
	LastLoginRegion       string  `gorm:"column:last_login_region;type:varchar(255)" json:"last_login_region"`
	LastLoginTime         int64   `gorm:"column:last_login_time;type:bigint(11)" json:"last_login_time"`
	LastLoginTerminal     string  `gorm:"column:last_login_terminal;type:varchar(255)" json:"last_login_terminal"`
}

func (AccountInformation) TableName() string {
	return "w_account_information"
}

type FundsLog struct {
	ID                      int64           `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	UUID                    string          `gorm:"column:uuid;type:varchar(255)" json:"uuid"`
	UID                     string          `gorm:"column:uid;type:varchar(255)" json:"uid"`
	MerchantUid             string          `gorm:"column:merchant_uid;type:varchar(255)" json:"merchant_id"`
	Txid                    string          `gorm:"column:txid;type:varchar(255);index:tx_unique,unique" json:"txid"`
	TransactionType         string          `gorm:"column:transaction_type;type:varchar(255);index:tx_unique,unique" json:"transaction_type"`
	UserAddress             string          `gorm:"column:user_address;type:varchar(255)" json:"user_address"`
	OppositeAddress         string          `gorm:"column:opposite_address;type:varchar(255)" json:"opposite_address"`
	CoinType                string          `gorm:"column:coin_type;type:varchar(255);index:tx_unique,unique" json:"coin_type"`
	AmountOfCoins           decimal.Decimal `gorm:"column:amount_of_coins;type:decimal(36,18)" json:"amount_of_coins"`
	NetworkFee              decimal.Decimal `gorm:"column:network_fee;type:decimal(36,18)" json:"network_fee"`
	GasLimit                uint64          `gorm:"column:gas_limit;type:bigint(11)" json:"gas_limit"`
	GasPrice                decimal.Decimal `gorm:"column:gas_price;type:decimal(36,18)" json:"gas_price"`
	GasUsed                 uint64          `gorm:"column:gas_used;type:bigint(11)" json:"gas_used"`
	ConfirmationBlockNumber string          `gorm:"column:confirm_block_number;type:varchar(255)" json:"confirm_block_number"`
	ConfirmationTime        int64           `gorm:"column:confirmation_time;type:bigint(11)" json:"confirmation_time"`
	State                   int8            `gorm:"column:state;type:int(8);comment:0 failed,1 success,2 pending" json:"state"`
	CreationTime            int64           `gorm:"column:creation_time;type:bigint(11)" json:"creation_time"`
}

func (FundsLog) TableName() string {
	return "w_funds_log"
}

type CoinCurrencyValues struct {
	ID         int64   `gorm:"column:id;primary_key;type:bigint(11);AUTO_INCREMENT" json:"id"`
	Coin       string  `gorm:"column:coin;type:varchar(255)" json:"coin"`
	Usd        float64 `gorm:"column:usd;type:decimal(28,8)" json:"usd"`
	Yuan       float64 `gorm:"column:yuan;type:decimal(28,8)" json:"yuan"`
	Euro       float64 `gorm:"column:euro;type:decimal(28,8)" json:"euro"`
	State      int8    `gorm:"column:state;type:int(8);default:1;" json:"state"`
	UpdateUser string  `gorm:"column:update_user;type:varchar(64)" json:"update_user"`
	UpdateTime int64   `gorm:"column:update_time;type:bigint(11)" json:"update_time"`
}

func (CoinCurrencyValues) TableName() string {
	return "w_coin_currency_value"
}
