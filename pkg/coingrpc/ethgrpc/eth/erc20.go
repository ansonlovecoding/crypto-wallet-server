package eth

import (
	"Share-Wallet/pkg/coingrpc/ethgrpc/contract"
	"Share-Wallet/pkg/coingrpc/ethgrpc/ethtx"
	"Share-Wallet/pkg/common/constant"
	db "Share-Wallet/pkg/db/mysql"
	"Share-Wallet/pkg/utils"
	"context"
	"fmt"
	"log"
	"math/big"

	"go.uber.org/zap"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/sha3"
)

type ERC20 struct {
	tokenClient     *contract.Token
	symbol          string
	name            string
	contractAddress string
	abiJSON         string
	decimals        int
}

func NewERC20(
	tokenClient *contract.Token,
	symbol string,
	name string,
	contractAddress string,
	abiJSON string,
	decimals int,
) *ERC20 {
	return &ERC20{
		tokenClient:     tokenClient,
		symbol:          symbol,
		name:            name,
		contractAddress: contractAddress,
		abiJSON:         abiJSON,
		decimals:        decimals,
	}
}

func (e *Ethereum) GetTokenBalance(hexAddr string) (*big.Int, error) {
	balance, err := e.erc20.tokenClient.BalanceOf(nil, common.HexToAddress(hexAddr))
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call e.contract.BalanceOf(%s)", hexAddr)
	}
	return balance, nil
}

func (e *Ethereum) CreateTokenRawTransaction(fromAddr, toAddr string, amount *big.Int, additionalNonce int, gasPrice *big.Int, gasLimit *big.Int) (*ethtx.RawTx, *db.EthDetailTX, *big.Int, *big.Int, error) {
	// validation check
	if e.ValidateAddr(fromAddr) != nil || e.ValidateAddr(toAddr) != nil {
		return nil, nil, nil, nil, errors.New("address validation error")
	}

	balance, err := e.GetTokenBalance(fromAddr)
	if err != nil {
		return nil, nil, nil, nil, errors.New("fail to get balance")
	}

	if balance.Cmp(amount) == -1 {
		return nil, nil, nil, nil, errors.New("your balance is no enough to transact")
	}

	ethBalance, err := e.GetBalance(fromAddr, QuantityTagLatest)
	if err != nil {
		return nil, nil, nil, nil, errors.New("fail to get eth balance")
	}
	log.Println("ethBalance", zap.Int64("ethBalance", ethBalance.Int64()))
	if ethBalance.Uint64() == 0 {
		return nil, nil, nil, nil, errors.New("eth balance is needed to send usdt")
	}

	if amount.Uint64() == 0 {
		return nil, nil, nil, nil, errors.New("amount can not be zero")
	}

	data := e.CreateTransferData(toAddr, amount)

	//gasLimit2, err := e.EstimateContractGas(data, fromAddr)
	//if err != nil {
	//	return nil, nil, nil, nil, errors.Wrap(err, "fail to call estimateGas(data)")
	//}
	//
	//gasPrice2, err := e.ethClient.SuggestGasPrice(context.Background())
	//if err != nil {
	//	return nil, nil, nil, nil, errors.Wrap(err, "fail to call client.SuggestGasPrice()")
	//}

	// txFee := gasPrice * estimatedGas
	txFeeWei := new(big.Int).Mul(gasPrice, gasLimit)
	if ethBalance.Cmp(txFeeWei) == -1 {
		// networkFee, _ := e.ConvertWeiToEther(txFeeWei).Float64()
		// errorMsg := fmt.Sprintf("The amount of ETH couldn’t less than transaction fee, around %f ETH", networkFee)
		return nil, nil, nil, nil, errors.Wrap(constant.ErrEthBalanceLessThanFee, fmt.Sprintf("%v", constant.ErrEthBalanceLessThanFee.ErrCode))
	}

	//log2.NewInfo("", utils.GetSelfFuncName(), "gasLimit", gasLimit, "gasPrice", gasPrice, "gasLimit2", gasLimit2, "gasPrice2", gasPrice2, "txFeeWei", txFeeWei)
	// nonce
	nonce, err := e.GetNonce(fromAddr, additionalNonce)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "fail to call e.getNonce()")
	}

	// create transaction
	contractAddr := common.HexToAddress(e.erc20.contractAddress)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &contractAddr,
		Value:    new(big.Int), // value must be 0 for ERC-20
		Gas:      gasLimit.Uint64(),
		GasPrice: gasPrice,
		Data:     data,
	})
	// From here, same as CreateRawTransaction() in ethgrop/eth/transaction.go
	txHash := tx.Hash().Hex()
	rawTxHex, err := ethtx.EncodeTx(tx)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "fail to call encodeTx()")
	}

	// generate UUID to trace transaction because unsignedTx is not unique
	uid := uuid.NewV4().String()

	// create insert data for　eth_detail_tx
	txDetailItem := &db.EthDetailTX{
		UUID:            uid,
		SenderAccount:   "",
		SenderAddress:   fromAddr,
		ReceiverAccount: "",
		ReceiverAddress: toAddr,
		Amount:          decimal.NewFromBigInt(amount, 0),
		Fee:             decimal.NewFromBigInt(txFeeWei, 0),
		GasLimit:        gasLimit.Uint64(),
		Nonce:           nonce,
		CoinType:        utils.GetCoinName(constant.USDTERC20),
	}

	// RawTx
	rawtx := &ethtx.RawTx{
		From:  fromAddr,
		To:    toAddr,
		Value: *amount,
		Nonce: nonce,
		TxHex: *rawTxHex,
		Hash:  txHash,
	}

	return rawtx, txDetailItem, balance, txFeeWei, nil
}

func (e *Ethereum) CreateTokenRawTransactionLocal(fromAddr, toAddr string, amount *big.Int, nonce uint64, gasPrice *big.Int, gasLimit *big.Int, contractAddress string) (*ethtx.RawTx, error) {
	// validation check
	if e.ValidateAddr(fromAddr) != nil || e.ValidateAddr(toAddr) != nil {
		return nil, errors.New("address validation error")
	}

	data := e.CreateTransferData(toAddr, amount)

	// create transaction
	contractAddr := common.HexToAddress(contractAddress)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &contractAddr,
		Value:    new(big.Int), // value must be 0 for ERC-20
		Gas:      gasLimit.Uint64(),
		GasPrice: gasPrice,
		Data:     data,
	})
	// From here, same as CreateRawTransaction() in ethgrop/eth/transaction.go
	txHash := tx.Hash().Hex()
	rawTxHex, err := ethtx.EncodeTx(tx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call encodeTx()")
	}

	// RawTx
	rawtx := &ethtx.RawTx{
		From:  fromAddr,
		To:    toAddr,
		Value: *amount,
		Nonce: nonce,
		TxHex: *rawTxHex,
		Hash:  txHash,
	}

	return rawtx, nil
}

func (e *Ethereum) CreateTransferData(toAddr string, amount *big.Int) []byte {
	// function signature as a byte slice
	transferFnSignature := []byte("transfer(address,uint256)")

	// methodID of function
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	// set parameter for account: to address
	paddedToAddr := common.LeftPadBytes(common.HexToAddress(toAddr).Bytes(), 32)
	// set parameter for amount
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	// create data
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedToAddr...)
	data = append(data, paddedAmount...)

	return data
}

func (e *Ethereum) EstimateContractGas(data []byte, fromAddr string) (uint64, error) {
	contractAddr := common.HexToAddress(e.erc20.contractAddress)
	masterAddr := common.HexToAddress(fromAddr)
	gasLimit, err := e.ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From: masterAddr,
		To:   &contractAddr,
		Data: data,
	})
	if err != nil {
		return 0, errors.Wrap(err, "fail to call client.EstimateGas()")
	}
	return gasLimit, nil
}
