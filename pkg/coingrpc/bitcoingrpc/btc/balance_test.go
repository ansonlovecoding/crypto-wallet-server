package btc_test

import (
	"Share-Wallet/pkg/testutil"
	"Share-Wallet/pkg/wallet/account"
	"testing"
)

// TestListAccounts is test for ListAccounts
func TestGetBalance(t *testing.T) {
	// t.SkipNow()
	bc := testutil.GetBTC()

	// GetBalance
	if res, err := bc.GetBalance(); err != nil {
		t.Errorf("fail to call GetBalance(): %v", err)
	} else {
		t.Log(res)
	}

	// bc.Close()
}

// TestGetBalanceByAccount is test for GetBalanceByAccount
func TestGetBalanceByAccount(t *testing.T) {
	// t.SkipNow()
	bc := testutil.GetBTC()

	type args struct {
		account account.AccountType
	}
	type want struct {
		isErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{account.AccountTypeClient},
			want: want{false},
		},
		{
			name: "happy path",
			args: args{account.AccountTypeDeposit},
			want: want{false},
		},
		{
			name: "happy path",
			args: args{account.AccountTypePayment},
			want: want{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := bc.GetBalanceByAccount(tt.args.account, bc.ConfirmationBlock()); (err != nil) != tt.want.isErr {
				t.Errorf("GetBalanceByAccount() = %v, isErr %v", err, tt.want.isErr)
			} else {
				t.Log(got)
			}
		})
	}

	// bc.Close()
}
