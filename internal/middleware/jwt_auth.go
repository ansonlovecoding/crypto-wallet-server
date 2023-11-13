package middleware

import (
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/http"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/common/token_verify"
	"Share-Wallet/pkg/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, userID, errInfo := token_verify.GetAdminUserIDFromToken(c.Request.Header.Get("token"), "", false)
		log.NewInfo("0", utils.GetSelfFuncName(), "userID: ", userID)
		c.Set("userID", userID)
		if !ok {
			log.NewError("", "GetUserIDFromToken false ", c.Request.Header.Get("token"))
			c.Abort()
			http.RespHttp200(c, constant.ErrParseToken, nil)
			return
		} else {
			log.NewInfo("0", utils.GetSelfFuncName(), "failed: ", errInfo)
		}
	}
}
