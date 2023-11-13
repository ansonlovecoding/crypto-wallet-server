package push

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	push2 "Share-Wallet/pkg/common/push"
	push3 "Share-Wallet/pkg/common/push/jpush"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/push"
	push4 "Share-Wallet/pkg/struct/rpc/push"
	"Share-Wallet/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

type pushRPCServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	pusher          push2.OfflinePusher
}

func NewPushRPCServer(port int) *pushRPCServer {
	pushServer := pushRPCServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.PushRPC,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
	if config.Config.Push.Jpns.Enable {
		pushServer.pusher = push3.JPushClient
	}
	return &pushServer
}

func (s *pushRPCServer) Run() {
	log.NewInfo("0", "pushRPCServer start ")
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
	push.RegisterPushMsgServiceServer(srv, s)
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
	log.NewInfo("0", "push rpc success")
}

func (s *pushRPCServer) PushMsg(_ context.Context, pbData *push.PushMsgReq) (*push.PushMsgResp, error) {
	//Call push module to send message to the user
	log.NewError(pbData.OperationID, utils.GetSelfFuncName(), "start pushing message")

	var alert string
	if pbData.Status == constant.TxStatusSuccess {
		if pbData.TransferType == constant.TransactionTypeSend {
			alert = fmt.Sprintf("Sent %s %s successfully", pbData.Amount, utils.GetCoinName(uint8(pbData.CoinType)))
		} else {
			alert = fmt.Sprintf("Received %s %s", pbData.Amount, utils.GetCoinName(uint8(pbData.CoinType)))
		}

	} else {
		if pbData.TransferType == constant.TransactionTypeSend {
			alert = fmt.Sprintf("Failed to send %s %s", pbData.Amount, utils.GetCoinName(uint8(pbData.CoinType)))
		}
	}

	detailContent := push4.PushDetailContent{
		CoinType:        int(pbData.CoinType),
		PublicAddress:   pbData.PublicAddress,
		TransactionHash: pbData.TxHash,
	}
	bDetailContent, _ := json.Marshal(detailContent)
	jsonDetailContent := string(bDetailContent)

	pushResult, err := s.pusher.Push([]string{pbData.PushToUserID}, alert, jsonDetailContent, pbData.OperationID)
	if err != nil {
		return nil, err
	}
	return &push.PushMsgResp{
		PushResult: pushResult,
	}, nil

}
