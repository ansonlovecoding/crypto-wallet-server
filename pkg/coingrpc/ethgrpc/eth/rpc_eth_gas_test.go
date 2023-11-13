package eth_test

import (
	"Share-Wallet/pkg/testutil"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

// TestGasPrice is test for GasPrice
func TestGasPrice(t *testing.T) {
	// t.SkipNow()
	et := testutil.GetETH()

	price, err := et.GasPrice()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("gasPrice:", price)
}

// TestEstimateGas is test for EstimateGas
func TestEstimateGas(t *testing.T) {
	et := testutil.GetETH()

	toAddr := common.HexToAddress("0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b")
	amount := big.NewInt(100000)
	var msg = &ethereum.CallMsg{
		From:  common.HexToAddress("0x24b11b06de55b09cb1c2d667af4abf570ac29088"),
		To:    &toAddr,
		GasPrice: new(big.Int).SetInt64(1000000000),
		Value: amount,
	}
	gas, err := et.EstimateGas(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("gas:",gas)

}
