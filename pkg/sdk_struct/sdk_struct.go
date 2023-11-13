package sdkstruct

import "github.com/shopspring/decimal"

type WalletCoinType struct {
	CoinType    int     `json:"coin_type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Balance     float64 `json:"balance"`
}

type GetBalance struct {
	CoinType    int    `json:"coin_type"`
	Address     string `json:"address"`
	OperationID string `json:"operationID"`
}
type GetGasPrice struct {
	OperationID string `json:"operationID"`
	CoinType    int    `json:"coin_type"`
	IsEstimated bool   `json:"is_estimated"`
}

type Transfer struct {
	OperationID string `json:"operationID"`
	CoinType    int    `json:"coin_type"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      string `json:"amount"`
	GasLimit    string `json:"gas_limit" binding:"omitempty"`
}
type GetTransaction struct {
	OperationID     string `json:"operationID"`
	CoinType        int    `json:"coin_type"`
	TransactionHash string `json:"transaction_hash"`
}

type SWConfig struct {
	Platform int32  `json:"platform"`
	ApiAddr  string `json:"api_addr"`
	DataDir  string `json:"data_dir"`
	LogLevel uint32 `json:"log_level"`
}

var SvrConf SWConfig

type GetTransactionList struct {
	CoinType        int    `json:"coin_type"`
	UserID          string `json:"user_id"`
	TransactionType int    `json:"transaction_type"`
	Address         string `json:"address"`
	OperationID     string `json:"operationID"`
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
	OrderBy         string `json:"order_by"`
	TransactionHash string `json:"transaction_hash"`
}
type GetAddressBook struct {
	Name     string `json:"name"`
	CoinType int    `json:"coin_type"`
	Address  string `json:"address"`
}
type TransactionListResponse struct {
	Code   int             `json:"code"`
	ErrMsg string          `json:"err_msg"`
	Data   TransactionData `json:"data"`
}
type TransactionData struct {
	OperationID string            `json:"operationID"`
	Transaction []TransactionInfo `json:"transaction"`
}
type TransactionInfo struct {
	TransactionID          uint64 `json:"transactionID"`
	UUID                   string `json:"uuID"`
	CurrentTransactionType int32  `json:"current_transaction_type"`
	SenderAccount          string `json:"sender_account"`
	SenderAddress          string `json:"sender_address"`
	ReceiverAccount        string `json:"receiver_account"`
	ReceiverAddress        string `json:"receiver_address"`
	Amount                 uint64 `json:"amount"`
	Fee                    uint64 `json:"fee"`
}

type TransactionDetail struct {
	TransactionID          uint64 `json:"transactionID"`
	UUID                   string `json:"uuID"`
	CurrentTransactionType int32  `json:"current_transaction_type"`
	SenderAccount          string `json:"sender_account"`
	SenderAddress          string `json:"sender_address"`
	ReceiverAccount        string `json:"receiver_account"`
	ReceiverAddress        string `json:"receiver_address"`
	Amount                 string `json:"amount"`
	Fee                    string `json:"fee"`
	TransactionHash        string `json:"transaction_hash"`
	ConfirmationTime       uint64 `json:"confirm_time"`
	SentTime               uint64 `json:"sent_time"`
	Status                 int8   `json:"status"`
	GasUsed                uint64 `json:"gas_used"`
	GasLimit               uint64 `json:"gas_limit"`
	GasPrice               string `json:"gas_price"`
	ConfirmBlockNumber     string `json:"confirm_block_number"`
}

type GetTransactionData struct {
	Transaction TransactionDetail `json:"transaction"`
}
type GetTransactionObj struct {
	Code   int                `json:"code"`
	ErrMsg string             `json:"err_msg"`
	Data   GetTransactionData `json:"data"`
}

type GetTransactionDetailResponse struct {
	TransactionID           uint64          `json:"transactionID"`
	UUID                    string          `json:"uuID"`
	CurrentTransactionType  int32           `json:"current_transaction_type"`
	SenderAccount           string          `json:"-"`
	SenderAddress           string          `json:"sender_address"`
	Sender                  *GetAddressBook `json:"sender_account"`
	Receiver                *GetAddressBook `json:"receiver_account"`
	ReceiverAccount         string          `json:"-"`
	ReceiverAddress         string          `json:"receiver_address"`
	Amount                  string          `json:"amount"`
	Fee                     string          `json:"fee"`
	TransactionHash         string          `json:"transaction_hash"`
	ConfirmationTime        uint64          `json:"confirm_time"`
	SentTime                uint64          `json:"sent_time"`
	Status                  int8            `json:"status"`
	AmountFloat             float64         `json:"amount_conv"`
	FeeFloat                float64         `json:"fee_conv"`
	IsSend                  bool            `json:"is_send"`
	GasPriceFloat           decimal.Decimal `json:"gas_price_conv"`
	GasUsed                 uint64          `json:"gas_used"`
	GasLimit                uint64          `json:"gas_limit"`
	GasPrice                string          `json:"gas_price"`
	ConfirmationBlockNumber string          `json:"confirm_block_number"`
}

type GetRecentRecordsRequest struct {
	OperationID string `json:"operationID"`
	UID         string `json:"uid"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
}
type GetRecentRecordsResponse struct {
	Code   int             `json:"code"`
	ErrMsg string          `json:"err_msg"`
	Data   TransactionData `json:"data"`
}
