package eth_test

import (
	"Share-Wallet/pkg/coingrpc/ethgrpc/eth"
	"Share-Wallet/pkg/testutil"
	"math/big"
	"testing"
	"time"
)

// TestCreateRawTransaction is test for CreateRawTransaction
func TestCreateRawTransaction(t *testing.T) {
	// t.SkipNow()
	et := testutil.GetETH()

	type args struct {
		senderAddr   string
		receiverAddr string
		amount       int64
		gasPrice     int64
		gasLimit     int64
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
			name: "happy path, send all",
			args: args{
				senderAddr:   "0xe52307Deb1a7dC3985D2873b45AE23b91D57a36d",
				receiverAddr: "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:       0,
			},
			want: want{false},
		},
		{
			name: "happy path, send specific amount",
			args: args{
				senderAddr:   "0xe52307Deb1a7dC3985D2873b45AE23b91D57a36d",
				receiverAddr: "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:       40000000000000000, // 0.04 Ether
			},
			want: want{false},
		},
		{
			name: "sender balance is insufficient",
			args: args{
				senderAddr:   "0xe52307Deb1a7dC3985D2873b45AE23b91D57a36d",
				receiverAddr: "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:       250000000000000000, // 0.25 Ether
			},
			want: want{true},
		},
		{
			name: "sender doesn't have amount",
			args: args{
				senderAddr:   "0x0Dd4d77D8b3bf210974332d1E16275bbEDdbF1CE",
				receiverAddr: "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:       0,
			},
			want: want{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawTx, txDetail, _, _, err := et.CreateRawTransaction(tt.args.senderAddr, tt.args.receiverAddr, big.NewInt(tt.args.amount), 0, big.NewInt(tt.args.gasPrice), big.NewInt(tt.args.gasLimit))
			if (err == nil) == tt.want.isErr {
				t.Errorf("CreateRawTransaction() = %v, want error = %v", err, tt.want.isErr)
				return
			}
			if rawTx != nil {
				t.Log(rawTx)
				// grok.Value(rawTx)
			}
			if txDetail != nil {
				t.Log(txDetail)
				// grok.Value(txDetail)
			}
		})
	}
	// et.Close()
}

// TestSignAndSendRawTransaction is test for SignOnRawTransaction and SendSignedRawTransaction
func TestSignAndSendRawTransaction(t *testing.T) {
	// t.SkipNow()
	et := testutil.GetETH()

	type args struct {
		senderAddr   string
		receiverAddr string
		amount       int64
		password     string
		gasPrice     int64
		gasLimit     int64
	}
	type want struct {
		isSignErr bool
		isSendErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "wrong password",
			args: args{
				senderAddr:   "0x24b11b06de55b09cb1c2d667af4abf570ac29088",
				receiverAddr: "0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b",
				amount:       0,
				password:     "foobar",
			},
			want: want{true, false},
		},
		{
			name: "happy path",
			args: args{
				senderAddr:   "0x24b11b06de55b09cb1c2d667af4abf570ac29088",
				receiverAddr: "0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b",
				amount:       0,
				password:     eth.Password,
			},
			want: want{false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create raw transaction
			rawTx, _, _, _, err := et.CreateRawTransaction(tt.args.senderAddr, tt.args.receiverAddr, big.NewInt(tt.args.amount), 0, big.NewInt(tt.args.gasPrice), big.NewInt(tt.args.gasLimit))
			if err != nil {
				t.Fatal(err)
			}
			// sign on raw transaction
			signedTx, err := et.SignOnRawTransaction(rawTx, tt.args.password)
			if (err == nil) == tt.want.isSignErr {
				t.Errorf("SignOnRawTransaction() = %v, want error = %v", err, tt.want.isSignErr)
				return
			}
			if err != nil {
				return
			}
			if signedTx != nil {
				t.Log(signedTx)
			}
			// send signed transaction
			txHash, err := et.SendSignedRawTransaction(signedTx.TxHex)
			if (err == nil) == tt.want.isSendErr {
				t.Errorf("SendSignedRawTransaction() = %v, want error = %v", err, tt.want.isSignErr)
				return
			}
			if txHash != "" {
				t.Log(txHash)
				// check transaction
				time.Sleep(3 * time.Second)
				res, err := et.GetTransactionByHash(txHash)
				if err != nil {
					t.Fatal(err)
				}
				t.Log(res)

				// check balance
				balance, err := et.GetBalance(tt.args.receiverAddr, eth.QuantityTagPending)
				if err != nil {
					t.Fatal(err)
				}
				if balance.Uint64() == 0 {
					t.Error("balance must be NOT zero")
				}

				// check confirmation
				confirmNum, err := et.GetConfirmation(txHash)
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("confirmation is %d", confirmNum)
			}
		})
	}
	// et.Close()
}

func TestEthereum_GetConfirmation(t *testing.T) {
	et := testutil.GetETH()
	blockNum, err := et.GetConfirmation("0xd9c111e3c1047b4285383bc732d483cbc69d5a40c395e9e861cd3f63b8e842cc")
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("blockNum", blockNum)
}
