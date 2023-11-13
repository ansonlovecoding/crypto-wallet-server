package eth

import (
	"Share-Wallet/pkg/common/http"
	db "Share-Wallet/pkg/db/local_db"
)

type EthMgr struct {
	Db            *db.DataBase
	LoginUserID   string
	LoginTime     int64
	PostWalletAPI *http.PostAPI
}

var EthManager *EthMgr

func NewEthMgr(dataBase *db.DataBase, loginUserID string, loginTime int32, p *http.PostAPI) (w *EthMgr) {
	return &EthMgr{Db: dataBase, LoginUserID: loginUserID, LoginTime: int64(loginTime), PostWalletAPI: p}
}
