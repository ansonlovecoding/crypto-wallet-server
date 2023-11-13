package wallet_api

type CreateTransactionRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	CoinType    uint32 `json:"coin_type" binding:"required"`
	FromAddress string `json:"from_address" binding:"required"`
	ToAddress   string `json:"to_address" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
}

type CreateTransactionResponse struct {
	TxID      string `json:"txID" binding:"required"`
	RawTxData string `json:"raw_tx_data" binding:"required"`
}

type TransferTronRequest struct {
	OperationID     string `json:"operationID" binding:"required"`
	CoinType        uint32 `json:"coin_type" binding:"required"`
	FromAccountUID  string `json:"from_account_uid" binding:"required"`
	FromMerchantUID string `json:"from_merchant_uid" binding:"required"`
	FromAddress     string `json:"from_address" binding:"required"`
	ToAddress       string `json:"to_address" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
	TxID            string `json:"tx_id" binding:"required"`
	TxData          string `json:"tx_data" binding:"required"`
	EnergyUsed      string `json:"energy_used"`
	EnergyPenalty   string `json:"energy_penalty"`
}

type GetTronConfirmationRequest struct {
	OperationID           string `json:"operationID"`
	CoinType              int    `json:"coin_type" binding:"required"`
	TransactionSignedHash string `json:"transaction_hash" binding:"required"`
}

type GetTronConfirmationResponse struct {
	BlockNum         string `json:"block_num"`
	ConfirmationTime uint64 `json:"confirm_time"`
	Status           int8   `json:"status"`
	NetFee           string `json:"net_fee"`
	EnergyUsage      string `json:"energy_usage"`
}

type GetTronConfirmationResponseObj struct {
	Code   int                         `json:"code"`
	ErrMsg string                      `json:"err_msg"`
	Block  GetTronConfirmationResponse `json:"data"`
}
