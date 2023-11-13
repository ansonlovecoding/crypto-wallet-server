package hot_wallet

import (
	syncr "Share-Wallet/internal/wallet-sdk-core/synchronization"
	"Share-Wallet/internal/wallet-sdk-core/transfer"
)

func TestHot() string {
	return "I'm hot interface!"
}

/**
coinType: 1 BTC, 2 ETH, 3 USDT-ERC20, 4 TRX, 5 USDT-TRC20
*/
func GetBalance(coinType int, publicAddress string) string {
	return transfer.TransferMgr.GetBalance(coinType, publicAddress)
}

//only for ETH, USDT-ERC20
func GetGasPrice(coinType int) string {
	return transfer.TransferMgr.GetGasPrice(coinType)
}

func Transfer(coinType int, fromAddress, toAddress, secret string, amount, gasPrice float64) (txHash string, error error) {
	return transfer.TransferMgr.Transferfn(coinType, fromAddress, toAddress, secret, amount, gasPrice)
}

func GetTransaction(coinType int, publicAddress, transactionHash string) string {
	return transfer.TransferMgr.GetTransactionDetails(coinType, publicAddress, transactionHash)
}

func GetConfirmation(coinType int, transactionHash string) string {
	return transfer.TransferMgr.GetConfirmation(coinType, transactionHash)
}

func GetPublicAddress(coinType int) string {
	return transfer.TransferMgr.GetPublicAddress(coinType)
}

func GetTransactionList(coinType int, publicAddress string, transactionType, page, pageSize int, orderBy string) string {
	return transfer.TransferMgr.GetTransactionList(coinType, publicAddress, transactionType, page, pageSize, orderBy, "")
}

func GetRecentTransactions(page, pageSize int) string {
	return transfer.TransferMgr.GetRecentTransactions(page, pageSize)
}

func GetTransactionFee(coinType int64, gasPrice float64) string {
	return transfer.TransferMgr.GetTransactionFee(coinType, gasPrice)
}

// To do:
//func TransferSignOffline(coinType int, toAddress, secret string, amount, gasPrice uint64) string {
//	return transfer.TransferMgr.TransferfnOffline(coinType, toAddress, secret, amount, gasPrice)
//}

func GetCoinStatuses() bool {
	return syncr.Sg.GetCoinStatuses()
}

func GetCoinRatio() string {
	return syncr.Sg.GetCoinRatio()
}

func CheckUserWallet(userID string) bool {
	hasWallet, _ := syncr.Sg.GetWallet(userID, 0)
	return hasWallet
}

//fetch the specific coin address of the friend
func FetchFriendAddress(userID string, coinType int64) string {
	_, address := syncr.Sg.GetWallet(userID, uint32(coinType))
	return address
}
