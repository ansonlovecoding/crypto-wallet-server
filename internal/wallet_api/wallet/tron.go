package wallet

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	http2 "Share-Wallet/pkg/common/http"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/tron"
	walletStruct "Share-Wallet/pkg/struct/wallet_api"
	"Share-Wallet/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Create Tron Transaction
func CreateTronTransaction(c *gin.Context) {
	var (
		req   walletStruct.CreateTransactionRequest
		resp  walletStruct.CreateTransactionResponse
		reqPb tron.CreateTransactionReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.CoinType = req.CoinType
	reqPb.FromAddress = req.FromAddress
	reqPb.ToAddress = req.ToAddress
	reqPb.Amount = req.Amount

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.TronRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := tron.NewTronClient(etcdConn)
	respPb, err := client.CreateTransactionRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		if respPb != nil && respPb.ErrMsg != "" {
			http2.RespHttp200(c, constant.ErrInfo{ErrCode: respPb.ErrCode, ErrMsg: respPb.ErrMsg}, nil)
		} else {
			http2.RespHttp200(c, err, nil)
		}
		return
	}

	resp.TxID = respPb.TxID
	resp.RawTxData = respPb.RawTXData
	//resp.RawTxData = utils.Base64Decode(respPb.RawTXData)
	//log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TransferRPC response", resp.TxID, resp.RawTxData)
	http2.RespHttp200(c, constant.OK, resp)
}

// Tron transfer
func TransferTron(c *gin.Context) {
	var (
		req   walletStruct.TransferTronRequest
		reqPb tron.PostTransferReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewError("", utils.GetSelfFuncName(), "BindJSON error", err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.CoinType = req.CoinType
	reqPb.FromAccountUID = req.FromAccountUID
	reqPb.FromMerchantUID = req.FromMerchantUID
	reqPb.FromAddress = req.FromAddress
	reqPb.ToAddress = req.ToAddress
	reqPb.Amount = req.Amount
	reqPb.TxID = req.TxID
	reqPb.TxDataStr = req.TxData
	reqPb.EnergyUsed = req.EnergyUsed
	reqPb.EnergyPenalty = req.EnergyPenalty

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.TronRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "getcdv3.GetConn failed", errMsg)
		http2.RespHttp200(c, constant.ErrInfo{ErrCode: constant.ErrServer.ErrCode, ErrMsg: constant.ErrServer.ErrMsg}, nil)
		return
	}
	client := tron.NewTronClient(etcdConn)
	respPb, err := client.TransferRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "client.TransferRPC", err.Error(), "respPb", respPb)
		http2.RespHttp200(c, err, nil)
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TransferRPC response", respPb.ErrCode, respPb.ErrMsg, "error", err)
	http2.RespHttp200(c, constant.OK, nil)
}

func GetTronConfirmationTransaction(c *gin.Context) {
	var (
		req   walletStruct.GetTronConfirmationRequest
		resp  walletStruct.GetTronConfirmationResponse
		reqPb tron.GetTronConfirmationReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.CoinType = uint32(req.CoinType)
	reqPb.TransactionHash = req.TransactionSignedHash
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.TronRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := tron.NewTronClient(etcdConn)
	respPb, err := client.GetConfirmationRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}

	log.NewInfo(req.OperationID, "BlockNumber:", respPb.BlockNum)
	resp.BlockNum = respPb.BlockNum
	resp.ConfirmationTime = respPb.ConfirmTime
	resp.NetFee = respPb.NetFee
	resp.Status = int8(respPb.Status)
	resp.EnergyUsage = respPb.EnergyUsage
	http2.RespHttp200(c, constant.OK, resp)
}
