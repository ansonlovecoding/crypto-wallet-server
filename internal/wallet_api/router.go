package wallet_api

import (
	"Share-Wallet/internal/wallet_api/wallet"
	"Share-Wallet/pkg/utils"

	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	baseRouter := gin.Default()
	baseRouter.Use(utils.CorsHandler())
	router := baseRouter.Group("/api/v1")
	router.Use(utils.CorsHandler())
	walletRouterGroup := router.Group("/wallet")
	{
		walletRouterGroup.POST("/get_balance", wallet.GetAccountBalance)
		walletRouterGroup.POST("/account", wallet.CreateAccountInformation)
		walletRouterGroup.POST("/update-login-info", wallet.UpdateAccountInformation)
		walletRouterGroup.GET("/coins", wallet.GetCoinStatuses)
		walletRouterGroup.GET("/coin_ratio", wallet.GetCoinRatio)
		walletRouterGroup.POST("/update_coin_rates", wallet.UpdateCoinRates)
		walletRouterGroup.POST("/get_user_wallet", wallet.GetUserWallet)
		walletRouterGroup.POST("/test", wallet.TestWalletApi)
		walletRouterGroup.POST("/test_eth", wallet.TestEthApi)
		walletRouterGroup.POST("/get_support_token_addresses", wallet.GetSupportTokenAddresses)
		walletRouterGroup.POST("/funds_log", wallet.GetFundsLog)
		walletRouterGroup.POST("/get_transaction", wallet.GetTransaction)
		walletRouterGroup.POST("/get_transaction_list", wallet.GetTransactionList)
	}
	ethGroup := router.Group("/eth")
	{
		ethGroup.POST("/get_balance", wallet.GetETHBalance)
		ethGroup.POST("/get_gasprice", wallet.GetGasPrice)
		ethGroup.POST("/transfer", wallet.Transfer)
		ethGroup.POST("/get_confirmation", wallet.GetConfirmationTransaction)
		ethGroup.POST("/transferv2", wallet.TransferV2)
		ethGroup.POST("/check_balance_nonce", wallet.CheckBalanceAndNonce)
	}
	btcGroup := router.Group("/btc")
	{
		btcGroup.POST("/test_btc", wallet.TestBTCApi)
		btcGroup.POST("/get_blockchain_info", wallet.GetBlockChainInfo)
	}
	tronGroup := router.Group("/tron")
	{
		tronGroup.POST("/create_tx", wallet.CreateTronTransaction)
		tronGroup.POST("/transfer", wallet.TransferTron)
		tronGroup.POST("/get_confirmation", wallet.GetTronConfirmationTransaction)
	}

	r2 := router.Group("")
	r2.Use(utils.JWTAuth())

	return baseRouter
}
