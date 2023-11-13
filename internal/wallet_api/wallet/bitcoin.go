package wallet

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	http2 "Share-Wallet/pkg/common/http"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/btc"
	walletStruct "Share-Wallet/pkg/struct/wallet_api"
	"Share-Wallet/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// TestBTCApi godoc
// @Summary      Testing btc-rpc server
// @Description  Testing btc-rpc server
// @Tags         BTC
// @Accept       json
// @Produce      json
// @Param        req body wallet_api.TestRequest true "operationID is only for tracking"
// @Success      200  {object}  wallet_api.TestResponse
// @Failure      400  {object}  wallet_api.TestResponse
// @Router       /btc/test_btc [post]
func TestBTCApi(c *gin.Context) {
	var (
		req   walletStruct.TestRequest
		resp  walletStruct.TestResponse
		reqPb btc.CommonReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TestBTCApi!!")

	reqPb.OperationID = req.OperationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.BitcoinRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := btc.NewBtcClient(etcdConn)
	_, err := client.TestBtcRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	resp.Name = req.Name
	http2.RespHttp200(c, constant.OK, resp)
}

// GetBlockChainInfo godoc
// @Summary      Get blockchain info
// @Description  Get blockchain info
// @Tags         BTC
// @Accept       json
// @Produce      json
// @Param        req body wallet_api.TestRequest true "operationID is only for tracking"
// @Success      200  {object}  wallet_api.GetEthBalanceRequest
// @Failure      400  {object}  wallet_api.GetEthBalanceResponse
// @Router       /btc/get_balance [post]
func GetBlockChainInfo(c *gin.Context) {
	var (
		req   walletStruct.GetBlockChainInfoRequest
		resp  walletStruct.GetBlockChainInfoResponse
		reqPb btc.CommonReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.BitcoinRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := btc.NewBtcClient(etcdConn)
	respPb, err := client.GetBlockChainInfoRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}

	log.NewInfo(req.OperationID, "blockchain info:", respPb.Data)
	resp.Data = respPb.Data
	http2.RespHttp200(c, constant.OK, resp)
}
