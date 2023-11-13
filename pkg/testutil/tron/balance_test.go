package tron

import (
	"testing"
)

func TestBalanceAt(t *testing.T) {
	tron := GetTRON()
	balance, err := tron.BalanceAt("TTy7o4hXwuiztVe24EesKAB8haMcE5Keyo")
	if err != nil {
		t.Log("tron.BalanceAt failed", err.Error())
		return
	}
	t.Log("tron.BalanceAt success", balance)
}

func TestGetTokenBalance(t *testing.T) {
	tron := GetTRON()
	contractAddr := tron.GetUSDTTRC20ContractAddress()
	t.Log("contractAddr", contractAddr)
	balance, err := tron.GetTokenBalance("TTy7o4hXwuiztVe24EesKAB8haMcE5Keyo", contractAddr)
	if err != nil {
		t.Log("tron.GetTokenBalance failed", err.Error())
		return
	}
	t.Log("tron.GetTokenBalance success", balance)
}
