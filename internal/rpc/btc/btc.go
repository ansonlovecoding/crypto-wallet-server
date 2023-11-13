package btc

import (
	"Share-Wallet/pkg/coingrpc"
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/btc"
	"Share-Wallet/pkg/utils"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type btcRPCServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewBtcRPCServer(port int) *btcRPCServer {
	return &btcRPCServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.BitcoinRPC,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *btcRPCServer) Run() {
	log.NewInfo("0", "btcRPCServer start ")
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)

	// listener network
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + s.rpcRegisterName)
	}
	log.NewInfo("0", "listen network success, ", address, listener)
	defer listener.Close()
	// grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()

	// Service registers with etcd
	btc.RegisterBtcServer(srv, s)
	rpcRegisterIP := ""
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "message cms rpc success")
}

func (s *btcRPCServer) TestBtcRPC(_ context.Context, req *btc.CommonReq) (*btc.CommonResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TestBtcRPC!!", req.String())
	resp := &btc.CommonResp{
		ErrCode: 0,
		ErrMsg:  "Test Success!",
	}
	return resp, nil
}

func (s *btcRPCServer) GetBlockChainInfoRPC(_ context.Context, req *btc.CommonReq) (*btc.GetBlockChainInfoResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetBlockChainInfoRPC!!", req.String())

	btcClient, err := coingrpc.GetBTCInstance()
	if err != nil {
		resp := &btc.GetBlockChainInfoResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  err.Error(),
		}
		return resp, nil
	}

	blockInfo, err := btcClient.GetBlockchainInfo()
	if err != nil {
		resp := &btc.GetBlockChainInfoResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  err.Error(),
		}
		return resp, nil
	}

	blockInfoByte, err := json.Marshal(blockInfo)
	if err != nil {
		resp := &btc.GetBlockChainInfoResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  err.Error(),
		}
		return resp, nil
	}

	resp := &btc.GetBlockChainInfoResp{
		ErrCode: 0,
		ErrMsg:  "get balance success!",
		Data: string(blockInfoByte),
	}
	return resp, nil
}
