package eth_test

import (
	"Share-Wallet/pkg/testutil"
	"github.com/bookerzzz/grok"
	"testing"
)

// TestGetTransactionByHash is test for GetTransactionByHash
func TestGetTransactionByHash(t *testing.T) {
	// t.SkipNow()
	et := testutil.GetETH()

	type args struct {
		txHash string
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
			args: args{
				txHash: "0xaeaa23d679d55ffa7db5a5da6825a1aca7a4fadd12144639be2024413f608cfd",
			},
			want: want{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// check GetTransactionByHash()
			res, err := et.GetTransactionByHash(tt.args.txHash)
			if err != nil {
				t.Fatal(err)
			}
			if res != nil {
				// t.Log(res)
				grok.Value(res)
			}

			// check GetTransactionReceipt()
			res2, err := et.GetTransactionReceipt(tt.args.txHash)
			if err != nil {
				t.Fatal(err)
			}
			if res2 != nil {
				// t.Log(res)
				grok.Value(res2)
			}
		})
	}
	// et.Close()
}
