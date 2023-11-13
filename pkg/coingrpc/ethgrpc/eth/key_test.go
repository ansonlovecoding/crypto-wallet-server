package eth_test

import (
	"Share-Wallet/pkg/coingrpc/ethgrpc/eth"
	"Share-Wallet/pkg/testutil"
	"log"
	"testing"
)

// TestGetPrivKey is test for GetPrivKey
// Note: keydir in config must be fullpath when testing
func TestGetPrivKey(t *testing.T) {
	et := testutil.GetETH()
	if et == nil {
		t.Error("et is nil!")
	}

	type args struct {
		addr string
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
				addr: "0x24b11b06de55b09cb1c2d667af4abf570ac29088",
			},
			want: want{false},
		},
		{
			name: "wrong address",
			args: args{
				addr: "0x5357135e0D3CbBD37cFCeb9F06257Bb133548Exx",
			},
			want: want{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("tt.args.addr:%s, eth.Password:%s", tt.args.addr, eth.Password)
			prikey, err := et.GetPrivKey(tt.args.addr, eth.Password)
			if (err == nil) == tt.want.isErr {
				t.Errorf("readPrivKey() = %v, want error = %v", err, tt.want.isErr)
				return
			}
			if err == nil && prikey == nil {
				t.Error("prikey is nil")
			}
		})
	}
}
