package btc_test

import (
	"Share-Wallet/pkg/testutil"
	"testing"
)

// TestLogging is test for GetLogging
func TestLogging(t *testing.T) {
	// t.SkipNow()
	bc := testutil.GetBTC()

	// Logging
	if res, err := bc.Logging(); err != nil {
		t.Errorf("fail to call Logging(): %v", err)
	} else {
		t.Log(res)
	}

	// bc.Close()()()
}
