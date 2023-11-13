package common

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/push"
	"Share-Wallet/pkg/utils"
	"context"
	"strings"
)

func PushMsg(operationID, senderUid, receiverUid, senderAddress, receiverAddress, txHash, amount string, status, coinType uint32) {
	//push notification to both sides of the transaction
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.PushRPC, operationID)
	if etcdConn == nil {
		log.NewError(operationID, "getcdv3.GetConn == nil")
		return
	}
	client := push.NewPushMsgServiceClient(etcdConn)

	//push to sender
	pushSenderReq := push.PushMsgReq{
		OperationID:   operationID,
		PushToUserID:  senderUid,
		CoinType:      coinType,
		TransferType:  constant.TransactionTypeSend,
		PublicAddress: senderAddress,
		TxHash:        txHash,
		Status:        status,
		Amount:        amount,
	}
	pushReceiverReq := push.PushMsgReq{
		OperationID:   operationID,
		PushToUserID:  receiverUid,
		CoinType:      coinType,
		TransferType:  constant.TransactionTypeReceive,
		PublicAddress: receiverAddress,
		TxHash:        txHash,
		Status:        status,
		Amount:        amount,
	}
	senderResp, err := client.PushMsg(context.Background(), &pushSenderReq)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "client.PushMsg to sender failed", err.Error())
	} else {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "client.PushMsg to sender success", senderResp.PushResult)
	}

	receiverResp, err := client.PushMsg(context.Background(), &pushReceiverReq)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "client.PushMsg to receiver failed", err.Error())
	} else {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "client.PushMsg to receiver success", receiverResp.PushResult)
	}
}
