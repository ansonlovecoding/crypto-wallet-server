package btc_test

import (
	"Share-Wallet/pkg/testutil"
	"Share-Wallet/pkg/wallet/account"
	"Share-Wallet/pkg/wallet/key"
	"testing"
)

// TestSetLabel is test for SetLabel
// Note: this test will contaminate wallet.dat
func TestSetLabel(t *testing.T) {
	// t.SkipNow()
	bc := testutil.GetBTC()

	// create address
	_, pubKey, err := key.GenerateWIF(bc.GetChainConf())
	if err != nil {
		t.Fatalf("fail to call GenerateWIF(): %v", err)
	}
	// import address
	// this scan address
	// err = bc.ImportAddressWithoutReScan(pubKey)
	err = bc.ImportAddressWithoutReScan(pubKey)
	if err != nil {
		t.Fatalf("fail to call ImportAddress(): %v", err)
	}

	// set label
	err = bc.SetLabel(pubKey, account.AccountTypeTest.String())
	if err != nil {
		t.Errorf("fail to call SetLabel(): %v", err)
	}

	// check addr
	res, err := bc.GetAddressInfo(pubKey)
	if err != nil {
		t.Fatalf("fail to call GetAddressInfo(): %v", err)
	}
	t.Log(res.Labels)
}
