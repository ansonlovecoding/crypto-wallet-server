package tron

import (
	"Share-Wallet/pkg/utils"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
)

func (t Tron) GetTokenBalance(fromAddr, contractAddr string) (*big.Int, error) {
	triggerData := fmt.Sprintf("[{\"address\":\"%s\"}]", fromAddr)
	//log.Println("triggerData", triggerData)
	tx, err := t.rpcClient.TriggerConstantContract(fromAddr, contractAddr, "balanceOf(address)", triggerData)
	if err != nil {
		return nil, err
	}
	//log.Println("tx.ConstantResult", tx.ConstantResult)
	if tx.ConstantResult == nil {
		return nil, errors.New("result is nil")
	}
	resultByte := tx.ConstantResult[0]
	resultHex := common.Bytes2Hex(resultByte)
	amountDecimal, err := utils.HexToDecimal(resultHex)
	if err != nil {
		return nil, err
	}
	return amountDecimal.BigInt(), nil
}

func (t Tron) CreateTokenTransaction(fromAddress, toAddress, contractAddress string, amount *big.Int, feeLimit *big.Int) (*api.TransactionExtention, error) {
	tx, err := t.rpcClient.TRC20Send(fromAddress, toAddress, contractAddress, amount, feeLimit.Int64())
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetUSDTTRC20ContractAddress returns the contract address of USDT-TRC20
func (t Tron) GetUSDTTRC20ContractAddress() string {
	return t.confTrc20.ContractAddress
}

func (t Tron) GetUSDTTRC20AbiJSON() string {
	return t.confTrc20.AbiJson
}

// GetEstimateFee
func (t Tron) EstimateFee(fromAddress, toAddress, contractAddress string, amount *big.Int) (*big.Int, error) {
	triggerData := fmt.Sprintf("[{\"address\":\"%s\"},{\"uint256\":\"%s\"}]", toAddress, amount.String())
	tx, err := t.rpcClient.TriggerConstantContract(fromAddress, contractAddress, "transfer(address,uint256)", triggerData)
	if err != nil {
		return nil, err
	}
	fee := big.NewInt(tx.EnergyUsed * 420)
	return fee, nil
}
