package trongrpc

import (
	"math/big"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
)

type Troner interface {
	// client
	BalanceAt(hexAddr string) (*big.Int, error)
	// transaction
	CreateTransaction(fromAddress, toAddress string, amount *big.Int) (*api.TransactionExtention, error)
	SignTransactionLocal(tx *core.Transaction, privateKey string) (*core.Transaction, error)
	SendTransaction(tx *core.Transaction) error
	GetTransactionInfo(hashTx string) (*core.TransactionInfo, error)
	GetNowBlockNum() (int64, error)
	GetBlockByNum(num int64) (*api.BlockExtention, error)
	//ABI
	GetUSDTTRC20ContractAddress() string
	GetTokenBalance(fromAddr, contractAddr string) (*big.Int, error)
	CreateTokenTransaction(fromAddress, toAddress, contractAddress string, amount *big.Int, feeLimit *big.Int) (*api.TransactionExtention, error)
	GetUSDTTRC20AbiJSON() string
	EstimateFee(fromAddress, toAddress, contractAddress string, amount *big.Int) (*big.Int, error)
}
