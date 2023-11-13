package tron

import (
	"math/big"
)

func (t Tron) BalanceAt(address string) (*big.Int, error) {
	account, err := t.rpcClient.GetAccount(address)
	if err != nil {
		return nil, err
	}
	return big.NewInt(account.Balance), nil
}
