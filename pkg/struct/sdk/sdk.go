package sdk

import "Share-Wallet/pkg/struct/wallet_api"

type CoinAddress struct {
	Name            string `json:"name"`
	CoinType        uint8  `json:"coin_type"`
	Address         string `json:"address"`
	ContractAddress string `json:"contract_address"`
}

type CheckBalanceAndNonceData struct {
	Nonce                    uint64 `json:"nonce"`
	ChainID                  string `json:"chain_id"`
	USDTERC20ContractAddress string `json:"usdterc_20_contract_address"`
}

type CheckBalanceAndNonceResp struct {
	Code   int                      `json:"code"`
	ErrMsg string                   `json:"err_msg"`
	Data   CheckBalanceAndNonceData `json:"data"`
}

type CreateTronTransactionResp struct {
	Code   int                                  `json:"code"`
	ErrMsg string                               `json:"err_msg"`
	Data   wallet_api.CreateTransactionResponse `json:"data"`
}
