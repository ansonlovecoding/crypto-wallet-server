package eth_test

import (
	"Share-Wallet/pkg/testutil"
	"testing"
)

// TestGetTotalBalance is test for GetTotalBalance
func TestGetTotalBalance(t *testing.T) {
	// t.SkipNow()
	et := testutil.GetETH()

	type args struct {
		addrs []string
	}
	type want struct {
		total uint64
		isErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{[]string{
				"0x36587c80f8652875bcb4bb85de44409ef9a35245",
				"0x24b11b06de55b09cb1c2d667af4abf570ac29088",
				"0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b",
			}},
			want: want{100, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total, userAmounts := et.GetTotalBalance(tt.args.addrs)
			t.Log(total)
			t.Log(userAmounts)
		})
	}
	// et.Close()
}
