package tron

import (
	"Share-Wallet/pkg/utils"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/big"
	"strings"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/pkg/errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

func (t Tron) CreateTransaction(fromAddress, toAddress string, amount *big.Int) (*api.TransactionExtention, error) {
	tx, err := t.rpcClient.Transfer(fromAddress, toAddress, amount.Int64())
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (t Tron) SignTransactionLocal(tx *core.Transaction, privateKey string) (*core.Transaction, error) {
	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		log.Println("proto.Marshal failed", err.Error())
		return nil, err
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	// btcec.PrivKeyFromBytes only returns a secret key and public key
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		log.Println("hex.DecodeString failed", err.Error())
		return nil, err
	}

	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	signature, err := crypto.Sign(hash, sk.ToECDSA())
	if err != nil {
		log.Println("crypto.Sign failed", err.Error())
		return nil, err
	}

	tx.Signature = append(tx.Signature, signature)
	return tx, nil
}

func (t Tron) SendTransaction(tx *core.Transaction) error {
	result, err := t.rpcClient.Broadcast(tx)
	if err != nil {
		return err
	}

	if result.Code != api.Return_SUCCESS {
		msg := utils.Bytes2String(result.Message)
		return errors.New(msg)
	}

	log.Println("SendTransaction result", result.Message, "result.Message string", utils.Bytes2String(result.Message), "result.Result", result.Result, "result.Code", result.Code.String())
	return nil
}

func (t Tron) GetTransactionInfo(hashTx string) (*core.TransactionInfo, error) {
	return t.rpcClient.GetTransactionInfoByID(hashTx)
}

func (t Tron) GetNowBlockNum() (int64, error) {
	nowBlock, err := t.rpcClient.GetNowBlock()
	if err != nil {
		return 0, err
	}

	if nowBlock.BlockHeader != nil && nowBlock.BlockHeader.RawData != nil {
		return nowBlock.BlockHeader.RawData.Number, nil
	} else {
		return 0, err
	}
}

func (t Tron) GetBlockByNum(num int64) (*api.BlockExtention, error) {
	return t.rpcClient.GetBlockByNum(num)
}
