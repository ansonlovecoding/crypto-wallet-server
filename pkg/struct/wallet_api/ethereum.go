package wallet_api

import (
	"Share-Wallet/pkg/struct/common"

	"github.com/shopspring/decimal"
)

type GetEthBalanceRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	CoinType    int    `json:"coin_type" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

type GetEthBalanceResponse struct {
	Balance string `json:"balance"`
}
type GetGasPriceRequest struct {
	OperationID string `json:"operationID" binding:"required"`
}
type GetGasPriceResponse struct {
	GasPrice int64 `json:"gas_price"`
}

type PostTransferRequest struct {
	OperationID     string `json:"operationID" binding:"required"`
	CoinType        uint32 `json:"coin_type" binding:"required"`
	FromAccountUID  string `json:"from_account_uid" binding:"required"`
	FromMerchantUID string `json:"from_merchant_uid" binding:"required"`
	FromAddress     string `json:"from_address" binding:"required"`
	ToAddress       string `json:"to_address" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
	Fee             string `json:"fee" binding:"omitempty"`
	GasLimit        uint64 `json:"gas_limit" binding:"omitempty"`
	Nonce           uint64 `json:"nonce"`
	GasPrice        string `json:"gas_price" binding:"required"`
	TxHash          string `json:"tx_hash" binding:"required"`
}

type PostTransferResponse struct {
	OperationID          string      `json:"operationID"`
	EthTransactionDetail Transaction `json:"transaction"`
}
type Transaction struct {
	TransactionID          uint64          `json:"transactionID"`
	UUID                   string          `json:"uuID"`
	CurrentTransactionType int32           `json:"current_transaction_type"`
	SenderAccount          string          `json:"sender_account"`
	SenderAddress          string          `json:"sender_address"`
	ReceiverAccount        string          `json:"receiver_account"`
	ReceiverAddress        string          `json:"receiver_address"`
	Amount                 decimal.Decimal `json:"amount"`
	Fee                    decimal.Decimal `json:"fee"`
	GasLimit               uint64          `json:"gas_limit"`
	Nonce                  uint64          `json:"nonce"`
	UnsignedHexTX          string          `json:"unsigned_hex_tx"`
	SignedHexTX            string          `json:"signed_hex_tx"`
	SentHashTX             string          `json:"sent_hash_tx"`
	UnsignedUpdatedAt      uint64          `json:"unsigned_updated_at"`
	SentUpdatedAt          uint64          `json:"sent_updated_at"`
}

type TransactionDetail struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              int64  `json:"gas"`
	GasPrice         int64  `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            int64  `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex int64  `json:"transactionIndex"`
	Value            string `json:"value"`
	V                int64  `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

type GetEthConfirmationRequest struct {
	OperationID           string `json:"operationID"`
	CoinType              int    `json:"coin_type" binding:"required"`
	TransactionSignedHash string `json:"transaction_hash" binding:"required"`
}

type GetEthConfirmationResponse struct {
	BlockNum         string `json:"block_num"`
	ConfirmationTime uint64 `json:"confirm_time"`
	Status           int8   `json:"status"`
	GasUsed          string `json:"gas_used"`
}

type GetTransactionDetailRequest struct {
	OperationID           string `json:"operationID"`
	CoinType              int    `json:"coin_type"`
	TransactionSignedHash string `json:"transaction_hash" binding:"required"`
}

type GetTransactionDetailResponse struct {
	OperationID       string          `json:"operationID"`
	TransactionDetail TransactionInfo `json:"transaction"`
}
type GetTransactionListRequest struct {
	OperationID      string `json:"operationID"`
	TransactionType  int    `json:"transaction_type" binding:"omitempty"`  // 1-Send 2- Receive
	TransactionState int    `json:"transaction_state" binding:"omitempty"` // 0-all, 1-pending, 2-success, 3-failed, 4-exclude pending
	PublicAddress    string `json:"address" binding:"required"`
	CoinType         int    `json:"coin_type" binding:"omitempty"`
	TransactionHash  string `json:"transaction_hash" binding:"omitempty"` // transaction_hash
	OrderBy          string `json:"order_by" binding:"omitempty,oneof=create_time:asc create_time:desc"`
	PageSize         int    `json:"page_size" binding:"omitempty,min=1,max=9223372036854775807" swaggertype:"integer"`
	Page             int    `json:"page" binding:"omitempty,min=-1,max=9223372036854775807" swaggertype:"integer"`
}

type GetTransactionListResponse struct {
	OperationID       string            `json:"operationID"`
	TransactionDetail []TransactionInfo `json:"transaction"`
	TranNums          int64             `json:"tran_nums"`
	common.ResponsePagination
}

type TransactionInfo struct {
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
type PostTransferRequestV2 struct {
	OperationID string `json:"operationID" binding:"required"`
	RawHex      string `json:"raw_hex" binding:"required"`
}
type PostTransferResponseV2 struct {
	OperationID     string `json:"operationID"`
	TransactionHash string `json:"transaction_hash"`
}
type PostEstimatedGasPriceRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	RawHex      string `json:"raw_hex" binding:"required"`
}
type PostEstimatedGasPriceResponse struct {
	OperationID     string `json:"operationID"`
	TransactionHash string `json:"transaction_hash"`
}

type CheckBalanceAndNonceRequest struct {
	OperationID    string `json:"operationID" binding:"required"`
	CoinType       uint32 `json:"coin_type" binding:"required"`
	FromAddress    string `json:"from_address" binding:"required"`
	TransactAmount string `json:"transact_amount" binding:"required"`
	GasPrice       string `json:"gas_price"`
	GasLimit       string `json:"gas_limit"`
}

type CheckBalanceAndNonceResponse struct {
	Nonce                    uint64 `json:"nonce"`
	ChainID                  string `json:"chain_id"`
	USDTERC20ContractAddress string `json:"usdterc_20_contract_address"`
}

type GetETHConfirmationResponseObj struct {
	Code   int                        `json:"code"`
	ErrMsg string                     `json:"err_msg"`
	Block  GetEthConfirmationResponse `json:"data"`
}
