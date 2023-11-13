package eth_test

import (
	"Share-Wallet/pkg/testutil"
	"testing"
)

// TestBalanceAt is test for BalanceAt
func TestBalanceAt(t *testing.T) {
	// t.SkipNow()
	et := testutil.GetETH()

	type args struct {
		addr string
	}
	type want struct {
		balance uint64
		isErr   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{"0x24b11b06de55b09cb1c2d667af4abf570ac29088"},
			want: want{100, false},
		},
		{
			name: "happy path",
			args: args{"0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b"},
			want: want{100, false},
		},
		{
			name: "address is random string",
			args: args{"0xe933a3318C3f5D94c2A3B2BEAEF772F67a45314d"},
			want: want{100, false},
		},
		{
			name: "address has no 0x",
			args: args{"e933a3318C3f5D94c2A3B2BEAEF772F67a45311c"},
			want: want{100, false},
		},
		{
			name: "address is btc address",
			args: args{"2N4TcHSCteXwiF2dj8SQijj3w2HieR4x6r5"},
			want: want{100, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balance, err := et.BalanceAt(tt.args.addr)
			if (err == nil) == tt.want.isErr {
				t.Errorf("BalanceAt() = %v, want error = %v", err, tt.want.isErr)
			}
			if balance != nil {
				t.Log(balance)
			}
		})
	}
	// et.Close()
}
