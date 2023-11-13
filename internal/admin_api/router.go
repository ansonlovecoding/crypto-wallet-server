package cms_api

import (
	"Share-Wallet/internal/admin_api/admin"
	"Share-Wallet/internal/middleware"
	"Share-Wallet/pkg/utils"

	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	baseRouter := gin.Default()
	baseRouter.Use(utils.CorsHandler())
	router := baseRouter.Group("/cms/v1")
	router.Use(utils.CorsHandler())
	adminRouterGroup := router.Group("/admin")
	adminRouterGroup.POST("/test", admin.TestAdminApi)
	adminRouterGroup.POST("/login", admin.AdminLogin)
	adminRouterGroup.POST("/admin-verify-totp", admin.VerifyTOTPAdminUser)

	adminRouterGroup.Use(middleware.JWTAuth())
	{
		// POST
		adminRouterGroup.POST("/reset-password", admin.ChangePassword)
		adminRouterGroup.POST("/users", admin.AddAdminUser)
		adminRouterGroup.POST("/roles", admin.AddAdminUserRole)

		adminRouterGroup.POST("/user-delete", admin.DeleteUserAPI)
		adminRouterGroup.POST("/role-delete", admin.DeleteRole)

		adminRouterGroup.POST("/users-update", admin.UpdateAdmin)
		adminRouterGroup.POST("/roles-update", admin.UpdateAdminRole)
		adminRouterGroup.POST("/reset-google-key", admin.ResetGoogleKey)
		adminRouterGroup.POST("/role-actions", admin.GetRoleActions)
		adminRouterGroup.POST("/update-currency", admin.UpdateCurrency)
		adminRouterGroup.POST("/confirm_tx", admin.ConfirmTransaction)
		adminRouterGroup.POST("/update_account_balance", admin.UpdateAccountBalance)

		//GET
		adminRouterGroup.GET("/users", admin.AdminUserList)
		adminRouterGroup.GET("/roles", admin.AdminUserRole)
		adminRouterGroup.GET("/actions", admin.AdminUserActions)
		adminRouterGroup.GET("/user-info", admin.GetAdminUser)
		adminRouterGroup.GET("/account-management", admin.GetAccountInformation)
		adminRouterGroup.GET("/funds-log", admin.GetFundsLog)
		adminRouterGroup.GET("/receive-details", admin.GetReceiveDetails)
		adminRouterGroup.GET("/transfer-details", admin.GetTransferDetails)
		adminRouterGroup.GET("/currencies", admin.GetCurrencies)
		adminRouterGroup.GET("/operational-report", admin.GetOperationalReport)
	}

	r2 := router.Group("")
	r2.Use(utils.JWTAuth())

	return baseRouter
}
