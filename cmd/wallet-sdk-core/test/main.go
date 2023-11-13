package main

import (
	cold_wallet "Share-Wallet/cmd/wallet-sdk-core/cold-wallet"
	hot_wallet "Share-Wallet/cmd/wallet-sdk-core/hot-wallet"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/utils"
	"fmt"
)

func main() {
	testTransaction1()
}

func testTransaction1() {

	userID := "08424"
	password := "123456"

	fmt.Println("\nSDK Version: ", cold_wallet.SDKVersion())

	config := "{\"platform\":1,\"api_addr\":\"http://api.devwallet.com\",\"data_dir\":\"db/sdk\",\"log_level\":1}"
	//config := "{\"platform\":1,\"api_addr\":\"http://api.wallet.com\",\"data_dir\":\"db/sdk\",\"log_level\":1}"

	flag := cold_wallet.InitWalletSDK(userID, config)
	fmt.Println("\nSDK Initialisation Status:", flag)

	isExist := cold_wallet.CheckUserExist(userID)
	fmt.Println("\nCheckUserExist:", isExist)

	//generate the seed phrase and accounts
	if !isExist {
		//testSeedPhrase := "boring invest empty dress juice poem renew pledge ribbon damp lava catalog"
		//isSuccess, err := cold_wallet.RecoverUserAccount(password, testSeedPhrase)
		//if isSuccess {
		//	fmt.Println("\nRecoverUserAccount success!")
		//} else {
		//	fmt.Println("\nRecoverUserAccount failed!", err)
		//}

		seedphrase, err := cold_wallet.GenerateUserAccount("en", password)
		if err != nil {
			fmt.Println("\nGenerateUserAccount failed:", err)
		} else {
			fmt.Println("\nGenerateUserAccount success:", seedphrase)
		}

	}

	//checking seed phrase
	mySeedPhrase := cold_wallet.GetUserSeedPhrase(password)
	fmt.Println("\nMySeedPhrase:", mySeedPhrase)

	//checking address
	btcAddress := hot_wallet.GetPublicAddress(constant.BTCCoin)
	fmt.Println("\nBTC Address:", btcAddress)

	ethAddress := hot_wallet.GetPublicAddress(constant.ETHCoin)
	fmt.Println("\nETH Address:", ethAddress)

	erc20Address := hot_wallet.GetPublicAddress(constant.USDTERC20)
	fmt.Println("\nUSDT-ERC20 Address:", erc20Address)

	trxAddress := hot_wallet.GetPublicAddress(constant.TRX)
	fmt.Println("\nTRX Address:", trxAddress)

	trc20Address := hot_wallet.GetPublicAddress(constant.USDTERC20)
	fmt.Println("\nUSDT-TRC20 Address:", trc20Address)

	//checking private key
	ethPrivateKey := cold_wallet.GetUserPrivateKey(password, constant.ETHCoin, constant.PrivateKeyHex)
	fmt.Println("\nETH Private key: ", ethPrivateKey)

	trxPrivateKey := cold_wallet.GetUserPrivateKey(password, constant.TRX, constant.PrivateKeyHex)
	fmt.Println("\nTRX Private key: ", trxPrivateKey)

	//checking balance
	ethBalance := hot_wallet.GetBalance(constant.ETHCoin, ethAddress)
	fmt.Println("\nETH Balance: ", ethBalance)

	usdtBalance := hot_wallet.GetBalance(constant.USDTERC20, ethAddress)
	fmt.Println("\nUSDT Balance: ", usdtBalance)

	trxBalance := hot_wallet.GetBalance(constant.TRX, "TTy7o4hXwuiztVe24EesKAB8haMcE5Keyo")
	fmt.Println("\nTRX Balance: ", trxBalance)

	trcBalance := hot_wallet.GetBalance(constant.USDTTRC20, "TTy7o4hXwuiztVe24EesKAB8haMcE5Keyo")
	fmt.Println("\nUSDT-TRC20 Balance: ", trcBalance)

	result, err := cold_wallet.AddAddressBook(constant.TRX, "TTy7o4hXwuiztVe24EesKAB8haMcE5Keyo", "Test1")
	fmt.Println("\nAddAddressBook1: result", result, "error:", err)

	result, err = cold_wallet.AddAddressBook(constant.USDTERC20, "0xBE41F2dd194e8C8C6801F8830B50aF4F066cf3b6", "Test1")
	fmt.Println("\nAddAddressBook USDT-ERC20 From Address: result", result, "error:", err)

	result, err = cold_wallet.AddAddressBook(constant.USDTERC20, "0x1232cd78885619C93bCce761f7c4B53c8B3b6CC3", "Test2")
	fmt.Println("\nAddAddressBook USDT-ERC20 To Address: result", result, "error:", err)

	//checking friend address
	myfriendTrxAddress := hot_wallet.FetchFriendAddress("08421", constant.TRX)
	fmt.Println("\nMy friend trx address: ", myfriendTrxAddress)

	myfriendETHAddress := hot_wallet.FetchFriendAddress("08421", constant.ETHCoin)
	fmt.Println("\nMy friend eth address: ", myfriendETHAddress)

	fromTrxAddress := "TTy7o4hXwuiztVe24EesKAB8haMcE5Keyo"
	isValid := utils.IsValidTRXAddress(fromTrxAddress)
	if isValid {
		fmt.Println("\nFrom address is valid address")
	} else {
		fmt.Println("\nFrom address is invalid address")
	}
	toTrxAddress := "TMbsRYyymxA54JjU1H9kxvtcbK5ZiUabch"
	isValid = utils.IsValidTRXAddress(toTrxAddress)
	if isValid {
		fmt.Println("\nTo address is valid address")
	} else {
		fmt.Println("\nTo address is invalid address")
	}

	//fromAddress := "0xBE41F2dd194e8C8C6801F8830B50aF4F066cf3b6"
	//fromAddress := "0x4de065E8Af7EA91a2209Fe2bB9cee186ccE38C70"

	//testing transaction
	/*
		trxAmount := 1.291234
		trxHashTX, err := hot_wallet.Transfer(constant.TRX, fromTrxAddress, toTrxAddress, password, float64(trxAmount), 0)
		fmt.Println("\nTransfer TRX hashTX: ", trxHashTX, "error:", err)

		gasPrice := hot_wallet.GetGasPrice(constant.ETHCoin)
		fmt.Println("\nGas price: ", gasPrice)

		gasPriceFloat64, _ := strconv.ParseFloat(gasPrice, 64)
		ethNetworkFee := hot_wallet.GetTransactionFee(constant.ETHCoin, gasPriceFloat64)
		fmt.Println("\nEth NetworkFee: ", ethNetworkFee)

		usdtNetworkFee := hot_wallet.GetTransactionFee(constant.USDTERC20, gasPriceFloat64)
		fmt.Println("\nUSDT NetworkFee: ", usdtNetworkFee)

		ethAmount := 0.00123
		ethAmountForWei := transfer.FromFloatEther(ethAmount)
		fmt.Println("\nEthAmountForWei: ", float64(ethAmount), ethAmountForWei)

		ethGasPrice := hot_wallet.GetGasPrice(constant.ETHCoin)
		ethGasPriceFloat, _ := strconv.ParseFloat(ethGasPrice, 64)
		fmt.Println("\nETH Gas price: ", ethGasPrice, ethGasPriceFloat)

		usdtGasPrice := hot_wallet.GetGasPrice(constant.USDTERC20)
		usdtGasPriceFloat, _ := strconv.ParseFloat(usdtGasPrice, 64)
		fmt.Println("\nUSDT-ERC20 Gas price: ", usdtGasPrice, usdtGasPriceFloat)


	*/
	//fromEthAddress := "0xBE41F2dd194e8C8C6801F8830B50aF4F066cf3b6"
	//toEthAdDress := "0x1232cd78885619C93bCce761f7c4B53c8B3b6CC3"
	//ethHashTX, err := hot_wallet.Transfer(constant.ETHCoin, fromEthAddress, toEthAdDress, password, float64(ethAmount), ethGasPriceFloat)
	//fmt.Println("\nTransfer ETH hashTX: ", ethHashTX, "error:", err)

	//usdtHashTX, err := hot_wallet.Transfer(constant.USDTERC20, fromAddress, toAddress, password, ethAmount, usdtGasPriceFloat)
	//fmt.Println("\nTransfer USDT-ERC20 hashTX: ", usdtHashTX, "error:", err)
	//ethHashTX := usdtHashTX

	//
	//trcAmount := 10.291234
	//trcHashTX, err := hot_wallet.Transfer(constant.USDTTRC20, fromAddress, toAddress, password, float64(trcAmount), 0)
	//fmt.Println("\nTransfer USDT-TRC20 hashTX: ", trcHashTX, "error:", err)
	//
	//trxHashTX := "3c572cb9a14efcde0cc621caa41f8bb89f123e26235fd688bcabde4baccfe8b9"
	//confirmationBlockNumber := hot_wallet.GetConfirmation(constant.TRX, trxHashTX)
	//fmt.Println("\nGetConfirmation: ", confirmationBlockNumber)
	//
	//transactionList := hot_wallet.GetTransactionList(constant.TRX, fromAddress, constant.TransactionTypeSend, 1, 100, "")
	//fmt.Println("\nTransactionList: ", transactionList)

	//transactionList := hot_wallet.GetTransactionList(constant.ETHCoin, fromAddress, constant.TransactionTypeAll, 0, 3, "create_time:desc")
	//fmt.Println("\nTransactionList: ", transactionList)

	//tran := hot_wallet.GetTransaction(constant.USDTERC20, fromEthAddress, "0x087b8a2718bc0f27bd3a6da531b7c746d9495dcd6b33372e7f449a63110ac7d7")
	//fmt.Println("\nGetTransaction: ", tran)

	//recentRecords := hot_wallet.GetRecentTransactions(1, 10)
	//fmt.Println("\nRecent Records: ", recentRecords)

	//0x84615940e25da849e0128a0fa34ab835dd7aa69f3324719763b5e554f4187bdf

	/*
		tran := hot_wallet.GetTransaction(constant.TRX, fromAddress, "3c572cb9a14efcde0cc621caa41f8bb89f123e26235fd688bcabde4baccfe8b9")
		fmt.Println("\nGetTransaction: ", tran)

		recentRecords := hot_wallet.GetRecentTransactions(1, 10)
		fmt.Println("\nRecent Records: ", recentRecords)
		//transactionList := hot_wallet.GetTransactionList(constant.ETHCoin, "0x86bdb0AE8aC56b6C70E4703336916364ab73847e", constant.TransactionTypeReceive, 1, 100, "")
		//fmt.Println("\nTransactionList: ", transactionList)
		//
		//ethBalance := hot_wallet.GetBalance(constant.ETHCoin, "0xAb38513d16C153c98af15CB1C50aBEcDFCf823b7")
		//fmt.Println("\nETH Balance: ", ethBalance)
		//
		//usdtBaLance := utils.ConvertUSDTBalance("3021292.999999999900")
		//fmt.Println("\nUSDT Balance: ", usdtBaLance)





		isVerify := cold_wallet.VerifyPasscode(password)
		fmt.Println("\nVerifyPasscode: ", isVerify)

		ethPublicKeyAddress := hot_wallet.GetPublicAddress(constant.ETHCoin)
		fmt.Println("\nETH Address: ", ethPublicKeyAddress)

		usdtERC20Address := hot_wallet.GetPublicAddress(constant.USDTERC20)
		fmt.Println("\nUSDT-ERC20 Address: ", usdtERC20Address)



		allAddress := hot_wallet.GetPublicAddress(0)
		fmt.Println("\nAll Address: ", allAddress)

		//ethHashTX := "0x089f4d33126a4ef728dad509da3e89bf41383d1fcd570667e1db03373ff7b216"

		//confirmationBlockNumber := hot_wallet.GetConfirmation(constant.ETHCoin, ethHashTX)
		//fmt.Println("\nGetConfirmation: ", confirmationBlockNumber)
		//
		//tran := hot_wallet.GetTransaction(constant.ETHCoin, fromAddress, ethHashTX)
		//fmt.Println("\nGetTransaction: ", tran)
		//
		//ethBalance2 := hot_wallet.GetBalance(constant.ETHCoin, fromAddress)
		//fmt.Println("\nNew ETH Balance: ", ethBalance2)
		//
		//usdtBalance2 := hot_wallet.GetBalance(constant.USDTERC20, fromAddress)
		//fmt.Println("\nNew USDT Balance: ", usdtBalance2)

	*/

}

//func testTransaction2() {
//
//	fromAddress := "0x36587c80f8652875bcb4bb85de44409ef9a35245"
//	toAddress := "0xb7570D5034E8EEd9A98637a02351d2f78a8A8651" //private key: 0x2d45cef7695783819c48a7f21604c66c04625fe5d0012c58a19b09151b41a257
//
//	keystoreString := `{"address":"36587c80f8652875bcb4bb85de44409ef9a35245","crypto":{"cipher":"aes-128-ctr","ciphertext":"a812d3f4415b13b17986f4a5a629010a9e4b67fff9f08e2160f7b16e8b70aa72","cipherparams":{"iv":"dddf9409492770d84e6b5daccc75915c"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"2c5309821e5d662dca538392c7de3a6ff9e1e158289e9b4b8f5a72e451e403b3"},"mac":"62833d5deebe480e4623c2a05a2e417e130038ca72420819962f29c4e7878aa1"},"id":"9a890b91-6300-4bf5-920c-218ed2b95160","version":3}`
//	password := "123456"
//
//	//create transaction
//	eth := eth.Ethereum{}
//
//	gasPrice := big.NewInt(1000)
//	estimatedGas := big.NewInt(136)
//	txFee := new(big.Int).Mul(gasPrice, estimatedGas)
//	rawTx, _, err := eth.CreateRawTransactionLocal(fromAddress, toAddress, big.NewInt(12400000), gasPrice, estimatedGas, txFee, 1)
//	if err != nil {
//		fmt.Println("CreateRawTransactionLocal err", err)
//		return
//	}
//
//	//sign transaction
//	//chainID need to get it from chain info
//	signTx, err := eth.SignOnRawTransactionLocal(rawTx, keystoreString, password, big.NewInt(10086))
//	if err != nil {
//		fmt.Println("SignOnRawTransactionLocal err", err)
//		return
//	}
//
//	fmt.Println("signTx", signTx.Hash)
//
//	//submit transaction
//	// TODO
//
//}
