package eth_test

import (
	"Share-Wallet/pkg/testutil"
	"math/big"
	"testing"
)

func TestERC20_GetTokenBalance(t *testing.T) {
	erc20 := testutil.GetETH()
	addr := "0x048Caa04B0976aA80F8a18616d0f6c13b27D4E5b"
	balance, err := erc20.GetTokenBalance(addr)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("balance", balance)
}

func TestERC20_EstimateContractGas(t *testing.T) {
	erc20 := testutil.GetETH()
	fromAddr := "0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b"
	toAddr := "0x24b11b06de55b09cb1c2d667af4abf570ac29088"
	amount := big.NewInt(1000)
	data := erc20.CreateTransferData(toAddr, amount)
	gas, err := erc20.EstimateContractGas(data, fromAddr)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("gas", gas)
}

func TestERC20_CreateRawTransaction(t *testing.T) {
	erc20 := testutil.GetETH()
	fromAddr := "0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b"
	toAddr := "0x24b11b06de55b09cb1c2d667af4abf570ac29088"
	amount := 1000
	gasPrice := 23
	gasLimit := 70000
	rawTx, _, _, _, err := erc20.CreateTokenRawTransaction(fromAddr, toAddr, big.NewInt(int64(amount)), 1, big.NewInt(int64(gasPrice)), big.NewInt(int64(gasLimit)))
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("rawTx hash", rawTx.Hash)
}

func TestEthereum_SendSignedRawTransaction(t *testing.T) {
	erc20 := testutil.GetETH()

	fromAddr := "0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b"
	privateKey := "6f57398ac37f7e20dd38abb39bfc2a7a6a085320907e21fa288e635a09568788"
	toAddr := "0x24b11b06de55b09cb1c2d667af4abf570ac29088"
	amount := 1000
	gasPrice := 23
	gasLimit := 70000
	rawTx, _, _, _, err := erc20.CreateTokenRawTransaction(fromAddr, toAddr, big.NewInt(int64(amount)), 1, big.NewInt(int64(gasPrice)), big.NewInt(int64(gasLimit)))
	if err != nil {
		t.Log(err.Error())
		return
	}

	signedRawTx, err := erc20.SignOnRawTransaction(rawTx, privateKey)
	if err != nil {
		t.Log(err.Error())
		return
	}

	txHash, err := erc20.SendSignedRawTransaction(signedRawTx.TxHex)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("txHash", txHash)
}
