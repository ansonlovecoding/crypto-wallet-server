package eth

import (
	"Share-Wallet/pkg/coingrpc"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/proto/eth"
	"Share-Wallet/pkg/utils"
	"context"
)

func (s *ethRPCServer) GetUSDTERC20BalanceRPC(_ context.Context, req *eth.GetBalanceReq) (*eth.GetBalanceResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetUSDTERC20BalanceRPC!!", req.String())

	if req.Address == "" {
		resp := &eth.GetBalanceResp{
			ErrCode: constant.ErrArgs.ErrCode,
			ErrMsg:  "address is nil!",
		}
		return resp, nil
	}

	erc20Client, err := coingrpc.GetETHInstance()
	if err != nil {
		resp := &eth.GetBalanceResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  err.Error(),
		}
		return resp, nil
	}

	balance, err := erc20Client.GetTokenBalance(req.Address)
	if err != nil {
		resp := &eth.GetBalanceResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  err.Error(),
		}
		return resp, nil
	}
	if balance == nil {
		resp := &eth.GetBalanceResp{
			ErrCode: constant.ErrRPC.ErrCode,
			ErrMsg:  "balance is nil",
		}
		return resp, nil
	}

	resp := &eth.GetBalanceResp{
		ErrCode: 0,
		ErrMsg:  "get balance success!",
		Balance: balance.String(),
	}
	return resp, nil
}
