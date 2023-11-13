package eth

import (
	"Share-Wallet/pkg/coingrpc/ethgrpc/ethtx"
	"Share-Wallet/pkg/common/constant"
	log2 "Share-Wallet/pkg/common/log"
	db "Share-Wallet/pkg/db/mysql"
	"Share-Wallet/pkg/utils"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// when creating multiple transaction from same address, nonce should be incremented
func (e *Ethereum) GetNonce(fromAddr string, additionalNonce int) (uint64, error) {
	// by calling GetTransactionCount()
	nonce, err := e.GetTransactionCount(fromAddr, QuantityTagPending)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call eth.GetTransactionCount()")
	}
	if additionalNonce != 0 {
		nonce = nonce.Add(nonce, new(big.Int).SetUint64(uint64(additionalNonce)))
	}
	log.Println("nonce",
		zap.Uint64("GetTransactionCount(fromAddr, QuantityTagPending)", nonce.Uint64()),
	)

	return nonce.Uint64(), nil
}

// How to calculate transaction fee?
// https://ethereum.stackexchange.com/questions/19665/how-to-calculate-transaction-fee
func (e *Ethereum) calculateFee(fromAddr, toAddr common.Address, balance, gasPrice, value *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	msg := &ethereum.CallMsg{
		From:     fromAddr,
		To:       &toAddr,
		Gas:      0,
		GasPrice: gasPrice,
		Value:    value,
		Data:     nil,
	}
	// gasLimit
	estimatedGas, err := e.EstimateGas(msg)
	if err != nil {
		return nil, nil, nil, err
	}
	// txFee := gasPrice * estimatedGas
	txFee := new(big.Int).Mul(gasPrice, estimatedGas)
	// newValue := value - txFee
	newValue := new(big.Int)
	if value.Uint64() == 0 {
		// receiver pays fee (deposit, transfer(pays all) action)
		newValue = newValue.Sub(balance, txFee)
	} else {
		// sender pays fee (payment, transfer(pays partially)
		newValue = value
		totalAmount := new(big.Int)
		totalAmount = totalAmount.Add(value, txFee)
		// if balance.Cmp(value) == -1 {
		if balance.Cmp(totalAmount) == -1 {
			//   -1 if x <  y
			//    0 if x == y
			//   +1 if x >  y
			return nil, nil, nil, errors.Errorf("balance`%s` is insufficient to send `%s`", balance.String(), value.String())
		}
	}

	return newValue, txFee, estimatedGas, nil
}

// CreateRawTransaction creates raw transaction for watch only wallet
// TODO: which QuantityTag should be used?
// - Creating offline/raw transactions with Go-Ethereum
//   https://medium.com/@akshay_111meher/creating-offline-raw-transactions-with-go-ethereum-8d6cc8174c5d
// Note: sender account owes fee
// - if sender sends 5ETH, receiver receives 5ETH
// - sender has to pay 5ETH + fee
func (e *Ethereum) CreateRawTransaction(fromAddr, toAddr string, amount *big.Int, additionalNonce int, gasPrice *big.Int, gasLimit *big.Int) (*ethtx.RawTx, *db.EthDetailTX, *big.Int, *big.Int, error) {
	// validation check
	if e.ValidateAddr(fromAddr) != nil || e.ValidateAddr(toAddr) != nil {
		return nil, nil, nil, nil, errors.New("address validation error")
	}
	log.Println("eth.CreateRawTransaction()",
		zap.String("fromAddr", fromAddr),
		zap.String("toAddr", toAddr),
		zap.String("amount", amount.String()),
	)

	// TODO: pending status should be included in target balance??
	// TODO: if block is still syncing, proper balance is not returned
	balance, err := e.GetBalance(fromAddr, QuantityTagLatest)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "fail to call eth.GetBalance()")
	}
	log.Println("balance", zap.Int64("balance", balance.Int64()))
	if balance.Uint64() == 0 {
		return nil, nil, nil, nil, errors.New("balance is needed to send eth")
	}
	if balance.Cmp(amount) == -1 {
		return nil, nil, nil, nil, errors.New("your balance is no enough to transact")
	}

	// nonce
	nonce, err := e.GetNonce(fromAddr, additionalNonce)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "fail to call eth.getNonce()")
	}

	// gasPrice
	// gasPrice, err := e.GasPrice()
	// if err != nil {
	// 	return nil, nil, errors.Wrap(err, "fail to call eth.GasPrice()")
	// }
	// fmt.Println("gasPrice", gasPrice)
	// log.Println("gas_price", zap.Int64("gas_price", gasPrice.Int64()))

	// fromAddr, toAddr common.Address, gasPrice, value *big.Int
	//newValue, txFee, estimatedGas, err := e.calculateFee(
	//	common.HexToAddress(fromAddr),
	//	common.HexToAddress(toAddr),
	//	balance,
	//	gasPrice,
	//	amount,
	//)
	//if err != nil {
	//	if strings.Contains(err.Error(), "gas required exceeds allowance") {
	//		return nil, nil, nil, nil, errors.New("Your balance is no enough to pay the transaction fee")
	//	}
	//	return nil, nil, nil, nil, err
	//}

	//log.Println("tx parameter",
	//	zap.Uint64("GasLimit", GasLimit),
	//	zap.Uint64("estimatedGas", estimatedGas.Uint64()),
	//	zap.Uint64("txFee", txFee.Uint64()))

	// txFee := gasPrice * estimatedGas
	txFee := new(big.Int).Mul(gasPrice, gasLimit)

	totalAmount := new(big.Int)
	totalAmount = totalAmount.Add(amount, txFee)
	if balance.Cmp(totalAmount) == -1 {
		// networkFee, _ := e.ConvertWeiToEther(txFee).Float64()
		// errorMsg := fmt.Sprintf("Your balance is no enough to pay the transaction fee, around %f ETH", networkFee)
		return nil, nil, nil, nil, errors.Wrap(constant.ErrEthBalanceLessThanFee, fmt.Sprintf("%v", constant.ErrEthBalanceLessThanFee.ErrCode))
	}

	log2.NewInfo("", utils.GetSelfFuncName(), "totalAmount", totalAmount, "txFee", txFee, "balance", balance)

	// create transaction
	tmpToAddr := common.HexToAddress(toAddr)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &tmpToAddr,
		Value:    amount,
		Gas:      GasLimit,
		GasPrice: gasPrice,
	})
	txHash := tx.Hash().Hex()
	rawTxHex, err := ethtx.EncodeTx(tx)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "fail to call encodeTx()")
	}
	// generate UUID to trace transaction because unsignedTx is not unique
	uid := uuid.NewV4().String()

	// create insert data for　eth_detail_tx
	txDetailItem := &db.EthDetailTX{
		UUID:            uid,
		SenderAccount:   "",
		SenderAddress:   fromAddr,
		ReceiverAccount: "",
		ReceiverAddress: toAddr,
		Amount:          decimal.NewFromBigInt(amount, 0),
		Fee:             decimal.NewFromBigInt(txFee, 0),
		GasLimit:        gasLimit.Uint64(),
		Nonce:           nonce,
		Status:          0,
		GasPrice:        decimal.NewFromBigInt(gasPrice, 0),
		CoinType:        utils.GetCoinName(constant.ETHCoin),
	}

	// RawTx
	rawtx := &ethtx.RawTx{
		From:  fromAddr,
		To:    toAddr,
		Value: *amount,
		Nonce: nonce,
		TxHex: *rawTxHex,
		Hash:  txHash,
	}
	return rawtx, txDetailItem, balance, txFee, nil
}

// SignOnRawTransaction signs on raw transaction
// - https://ethereum.stackexchange.com/questions/16472/signing-a-raw-transaction-in-go
// - Note: this requires private key on this machine, if node is working remotely, it would not work.
func (e *Ethereum) SignOnRawTransaction(rawTx *ethtx.RawTx, passphrase string) (*ethtx.RawTx, error) {
	txHex := rawTx.TxHex
	fromAddr := rawTx.From
	tx, err := ethtx.DecodeTx(txHex)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call decodeTx(txHex)")
	}

	// get private key
	key, err := e.GetPrivKey(fromAddr, passphrase)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call e.GetPrivKey()")
	}

	// chain id
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md
	chainID := big.NewInt(int64(e.netID))
	if chainID.Uint64() == 0 {
		return nil, errors.Errorf("chainID can't get from netID:  %d", e.netID)
	}

	log.Println("call types.SignTx",
		zap.Any("tx", tx),
		zap.Uint64("chainID", chainID.Uint64()),
		zap.Any("key.PrivateKey", key.PrivateKey),
	)

	// sign
	signedTX, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call types.SignTx()")
	}
	// TODO: baseFee *big.Int param is added in AsMessage method and maybe useful
	msg, err := signedTX.AsMessage(types.NewEIP155Signer(chainID), nil)
	if err != nil {
		return nil, errors.Wrap(err, "fail to cll signedTX.AsMessage()")
	}

	encodedTx, err := ethtx.EncodeTx(signedTX)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call encodeTx()")
	}

	resTx := &ethtx.RawTx{
		From:  msg.From().Hex(),
		To:    msg.To().Hex(),
		Value: *msg.Value(),
		Nonce: msg.Nonce(),
		TxHex: *encodedTx,
		Hash:  signedTX.Hash().Hex(),
	}

	return resTx, nil
}

// SendSignedRawTransaction sends signed raw transaction
// - SendRawTransaction in rpc_eth_tx.go
// - SendRawTx in client.go
func (e *Ethereum) SendSignedRawTransaction(signedTxHex string) (string, error) {
	decodedTx, err := ethtx.DecodeTx(signedTxHex)
	if err != nil {
		return "", errors.Wrap(err, "fail to call decodeTx(signedTxHex)")
	}

	txHash, err := e.SendRawTransactionWithTypesTx(decodedTx)
	if err != nil {
		return "", errors.Wrap(err, "fail to call SendRawTransactionWithTypesTx()")
	}

	return txHash, err
}

// GetConfirmation returns confirmation number
func (e *Ethereum) GetConfirmation(hashTx string) (*big.Int, error) {
	txInfo, err := e.GetTransactionByHash(hashTx)
	if err != nil {
		return nil, err
	}
	if txInfo.BlockNumber.Int64() == 0 {
		return nil, errors.New("block number can't retrieved")
	}
	currentBlockNum, err := e.BlockNumber()
	if err != nil {
		return nil, err
	}
	var confirmationBlockNum *big.Int
	confirmationBlockNum = confirmationBlockNum.Sub(currentBlockNum, txInfo.BlockNumber)
	confirmationBlockNum = confirmationBlockNum.Add(confirmationBlockNum, big.NewInt(1))

	return confirmationBlockNum, nil
}

// SignOnRawTransaction signs on raw transaction
// - https://ethereum.stackexchange.com/questions/16472/signing-a-raw-transaction-in-go
// - Note: this requires private key on this machine, if node is working remotely, it would not work.
func (e *Ethereum) SignOnRawTransactionV2(rawTx *ethtx.RawTx, pkey string, chainID *big.Int) (*ethtx.RawTx, error) {
	txHex := rawTx.TxHex
	log.Println("SignOnRawTransactionV2", txHex, pkey)
	tx, err := ethtx.DecodeTx(txHex)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call decodeTx(txHex)")
	}

	// chain id
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md
	//chainID := big.NewInt(int64(e.netID))
	//if chainID.Uint64() == 0 {
	//	return nil, errors.Errorf("chainID can't get from netID:  %d", 5)
	//}

	//log.Println("call types.SignTx",
	//	zap.Any("tx", tx),
	//	zap.Uint64("chainID", chainID.Uint64()),
	//	zap.Any("key.PrivateKey", pkey),
	//)

	// sign
	ECDSAkey, err := e.ToECDSA(pkey)
	if err != nil {
		return nil, errors.Wrap(err, "ToECDSA failed")
	}
	signedTX, err := types.SignTx(tx, types.NewEIP155Signer(chainID), ECDSAkey)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call types.SignTx()")
	}

	// TODO: baseFee *big.Int param is added in AsMessage method and maybe useful
	msg, err := signedTX.AsMessage(types.NewEIP155Signer(chainID), nil)
	if err != nil {
		return nil, errors.Wrap(err, "fail to cll signedTX.AsMessage()")
	}
	encodedTx, err := ethtx.EncodeTx(signedTX)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call encodeTx()")
	}

	resTx := &ethtx.RawTx{
		From:  msg.From().Hex(),
		To:    msg.To().Hex(),
		Value: *msg.Value(),
		Nonce: msg.Nonce(),
		TxHex: *encodedTx,
		Hash:  signedTX.Hash().Hex(),
	}

	return resTx, nil
}

// GetPrivKey returns keystore.Key object
func (e *Ethereum) GetPrivKeyV2(pKey string) (*ecdsa.PrivateKey, error) {
	pkey := crypto.ToECDSAUnsafe([]byte(pKey))
	log.Println("pkey", pkey)
	return pkey, nil
}
func (e *Ethereum) CreateRawTransactionLocal(fromAddr, toAddr string, amount *big.Int, nonce uint64, gasPrice *big.Int) (*ethtx.RawTx, error) {
	// validation check
	if e.ValidateAddr(fromAddr) != nil || e.ValidateAddr(toAddr) != nil {
		return nil, errors.New("address validation error")
	}
	log.Println("eth.CreateRawTransaction()",
		zap.String("fromAddr", fromAddr),
		zap.String("toAddr", toAddr),
		zap.String("amount", amount.String()),
	)

	// create transaction
	tmpToAddr := common.HexToAddress(toAddr)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &tmpToAddr,
		Value:    amount,
		Gas:      GasLimit,
		GasPrice: gasPrice,
	})
	txHash := tx.Hash().Hex()
	rawTxHex, err := ethtx.EncodeTx(tx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call encodeTx()")
	}

	// RawTx
	rawtx := &ethtx.RawTx{
		From:  fromAddr,
		To:    toAddr,
		Value: *amount,
		Nonce: nonce,
		TxHex: *rawTxHex,
		Hash:  txHash,
	}
	return rawtx, nil
}
func (e *Ethereum) calculateFee2(fromAddr, toAddr common.Address, balance, gasPrice, value *big.Int, estimateGas *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	// msg := &ethereum.CallMsg{
	// 	From:     fromAddr,
	// 	To:       &toAddr,
	// 	Gas:      0,
	// 	GasPrice: gasPrice,
	// 	Value:    value,
	// 	Data:     nil,
	// }

	// estimatedGas, err := e.EstimateGas(msg)
	// if err != nil {
	// 	return nil, nil, nil, errors.Wrap(err, "fail to call EstimateGas()")
	// }
	// txFee := gasPrice * estimatedGas
	txFee := new(big.Int).Mul(gasPrice, estimateGas)
	// newValue := value - txFee
	newValue := new(big.Int)
	if value.Uint64() == 0 {
		// receiver pays fee (deposit, transfer(pays all) action)
		newValue = newValue.Sub(balance, txFee)
	} else {
		// sender pays fee (payment, transfer(pays partially)
		newValue = new(big.Int).Add(value, txFee)
		// newValue = newValue.Sub(value, txFee)
		// if balance.Cmp(value) == -1 {
		if balance.Cmp(newValue) == -1 {
			//   -1 if x <  y
			//    0 if x == y
			//   +1 if x >  y
			return nil, nil, nil, errors.Errorf("balance`%s` is insufficient to send `%s`", balance.String(), newValue.String())
		}
	}

	return newValue, txFee, estimateGas, nil
}
