package push

type PushDetailContent struct {
	CoinType        int    `json:"coin_type"`
	PublicAddress   string `json:"public_address"`
	TransactionHash string `json:"transaction_hash"`
}
