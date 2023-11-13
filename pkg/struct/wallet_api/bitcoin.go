package wallet_api

type GetBlockChainInfoRequest struct {
	OperationID string `json:"operationID" binding:"required"`
}

type GetBlockChainInfoResponse struct {
	Data string `json:"data"`
}
