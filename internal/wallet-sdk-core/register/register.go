package register

import (
	"Share-Wallet/pkg/coingrpc/ethgrpc"
	"Share-Wallet/pkg/common/http"
	db "Share-Wallet/pkg/db/local_db"
	"Share-Wallet/pkg/struct/wallet_api"
)

type Registry struct {
	Db               *db.DataBase
	LoginUserID      string
	LoginTime        int64
	PostWalletAPI    *http.PostAPI
	EthClient        ethgrpc.Ethereumer
	UserRandomSecret []byte
	SupportTokens    []*wallet_api.SupportTokenAddress
}

var RegistryMgr *Registry

func NewRegistry(dataBase *db.DataBase, loginUserID string, loginTime int32, p *http.PostAPI, randomsecret []byte, supportTokens []*wallet_api.SupportTokenAddress) (w *Registry) {

	return &Registry{Db: dataBase,
		LoginUserID:      loginUserID,
		LoginTime:        int64(loginTime),
		PostWalletAPI:    p,
		UserRandomSecret: randomsecret,
		SupportTokens:    supportTokens,
	}
}
