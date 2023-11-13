package tron

import (
	"Share-Wallet/pkg/utils"
	"math/big"
	"testing"

	"github.com/fbsobreira/gotron-sdk/pkg/common"
)

func TestSendTransaction(t *testing.T) {
	tron := GetTRON()
	trx := 1000.0
	pKey := "7af46a1a965569f9d522b4ea39c10fecffb0ea3b962b98357b6dce00ca115abf"
	sunAmount := utils.TrxToSun(trx)
	txExtention, err := tron.CreateTransaction("TTy7o4hXwuiztVe24EesKAB8haMcE5Keyo", "TY1hnzqSDWRDSNaVETtpdCibLBdSLUN4uW", sunAmount)
	if err != nil {
		t.Error("CreateTransaction failed", err.Error())
		return
	}
	t.Log(txExtention)
	txID := common.BytesToHexString(txExtention.Txid)
	t.Log("txExtention.Txid", txID)
	txSign, err := tron.SignTransactionLocal(txExtention.Transaction, pKey)
	if err != nil {
		t.Error("SignTransactionLocal failed", err.Error())
		return
	}
	t.Log(txSign)
	err = tron.SendTransaction(txSign)
	if err != nil {
		t.Error("SendTransaction failed", err.Error())
	}
}

func TestTokenTransaction(t *testing.T) {
	tron := GetTRON()
	usdt := 1.0
	pKey := "7af46a1a965569f9d522b4ea39c10fecffb0ea3b962b98357b6dce00ca115abf"
	amount := utils.ConvertFloatUSDTToBigInt(usdt)
	t.Log("amount", amount)
	contractAddr := tron.GetUSDTTRC20ContractAddress()

	fromAddress := "TTy7o4hXwuiztVe24EesKAB8haMcE5Keyo"
	toAddress := "TMbsRYyymxA54JjU1H9kxvtcbK5ZiUabch"
	fee, err := tron.EstimateFee(fromAddress, toAddress, contractAddr, amount)
	if err != nil {
		t.Error("EstimateFee failed", err.Error())
		return
	}

	txExtention, err := tron.CreateTokenTransaction(fromAddress, toAddress, contractAddr, amount, fee)
	if err != nil {
		t.Error("CreateTokenTransaction failed", err.Error())
		return
	}
	t.Log(txExtention)
	txID := common.BytesToHexString(txExtention.Txid)
	t.Log("txExtention.Txid", txID)
	txSign, err := tron.SignTransactionLocal(txExtention.Transaction, pKey)
	if err != nil {
		t.Error("SignTransactionLocal failed", err.Error())
		return
	}
	t.Log(txSign)
	err = tron.SendTransaction(txSign)
	if err != nil {
		t.Error("SendTransaction failed", err.Error())
	}
}

func TestGetTransactionInfo(t *testing.T) {
	tron := GetTRON()
	tx, err := tron.GetTransactionInfo("1c7ee18f7cd202df96f673806333983b62cd6543fd133da052f7e97badca9e8c")
	if err != nil {
		t.Error("GetTransactionInfo failed", err.Error())
		return
	}

	t.Log("contract result", common.Bytes2Hex(tx.ResMessage))
	t.Log("transaction", tx.String(), tx.Result.String(), tx.Receipt.Result.String(), "EnergyUsageTotal:", tx.Receipt.EnergyUsageTotal, "fee:", tx.Fee)
}

func TestEstimateFee(t *testing.T) {
	tron := GetTRON()
	fee, err := tron.EstimateFee("TTy7o4hXwuiztVe24EesKAB8haMcE5Keyo", "TMbsRYyymxA54JjU1H9kxvtcbK5ZiUabch", "TG3XXyExBkPp9nzdajDZsozEu4BkaSJozs", big.NewInt(1000000))
	if err != nil {
		t.Error("EstimateFee failed", err.Error())
		return
	}

	t.Log("EstimateFee result", fee)
}
