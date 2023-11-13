// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.11
// source: wallet/wallet.proto

package wallet

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// WalletClient is the client API for Wallet service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WalletClient interface {
	TestWalletRPC(ctx context.Context, in *CommonReq, opts ...grpc.CallOption) (*CommonResp, error)
	GetSupportTokenAddressesRPC(ctx context.Context, in *GetSupportTokenAddressesReq, opts ...grpc.CallOption) (*GetSupportTokenAddressesResp, error)
	CreateAccountInformation(ctx context.Context, in *CreateAccountInfoReq, opts ...grpc.CallOption) (*CreateAccountInfoResp, error)
	UpdateAccountInformation(ctx context.Context, in *UpdateAccountInfoReq, opts ...grpc.CallOption) (*CommonResp, error)
	GetFundsLog(ctx context.Context, in *GetFundsLogReq, opts ...grpc.CallOption) (*GetFundsLogResp, error)
	GetCoinStatuses(ctx context.Context, in *GetCoinStatusesReq, opts ...grpc.CallOption) (*GetCoinStatusesResp, error)
	GetCoinRatio(ctx context.Context, in *GetCoinRatioReq, opts ...grpc.CallOption) (*GetCoinRatioResp, error)
	GetUserWallet(ctx context.Context, in *GetUserWalletReq, opts ...grpc.CallOption) (*GetUserWalletResp, error)
	UpdateCoinRates(ctx context.Context, in *UpdateCoinRatesReq, opts ...grpc.CallOption) (*CommonResp, error)
	GetTransactionDetailRPC(ctx context.Context, in *GetTransactionDetailReq, opts ...grpc.CallOption) (*GetTransactionDetailRes, error)
	GetTransactionListRPC(ctx context.Context, in *GetTransactionListReq, opts ...grpc.CallOption) (*GetTransactionListRes, error)
}

type walletClient struct {
	cc grpc.ClientConnInterface
}

func NewWalletClient(cc grpc.ClientConnInterface) WalletClient {
	return &walletClient{cc}
}

func (c *walletClient) TestWalletRPC(ctx context.Context, in *CommonReq, opts ...grpc.CallOption) (*CommonResp, error) {
	out := new(CommonResp)
	err := c.cc.Invoke(ctx, "/wallet.wallet/TestWalletRPC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) GetSupportTokenAddressesRPC(ctx context.Context, in *GetSupportTokenAddressesReq, opts ...grpc.CallOption) (*GetSupportTokenAddressesResp, error) {
	out := new(GetSupportTokenAddressesResp)
	err := c.cc.Invoke(ctx, "/wallet.wallet/GetSupportTokenAddressesRPC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) CreateAccountInformation(ctx context.Context, in *CreateAccountInfoReq, opts ...grpc.CallOption) (*CreateAccountInfoResp, error) {
	out := new(CreateAccountInfoResp)
	err := c.cc.Invoke(ctx, "/wallet.wallet/CreateAccountInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) UpdateAccountInformation(ctx context.Context, in *UpdateAccountInfoReq, opts ...grpc.CallOption) (*CommonResp, error) {
	out := new(CommonResp)
	err := c.cc.Invoke(ctx, "/wallet.wallet/UpdateAccountInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) GetFundsLog(ctx context.Context, in *GetFundsLogReq, opts ...grpc.CallOption) (*GetFundsLogResp, error) {
	out := new(GetFundsLogResp)
	err := c.cc.Invoke(ctx, "/wallet.wallet/GetFundsLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) GetCoinStatuses(ctx context.Context, in *GetCoinStatusesReq, opts ...grpc.CallOption) (*GetCoinStatusesResp, error) {
	out := new(GetCoinStatusesResp)
	err := c.cc.Invoke(ctx, "/wallet.wallet/GetCoinStatuses", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) GetCoinRatio(ctx context.Context, in *GetCoinRatioReq, opts ...grpc.CallOption) (*GetCoinRatioResp, error) {
	out := new(GetCoinRatioResp)
	err := c.cc.Invoke(ctx, "/wallet.wallet/GetCoinRatio", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) GetUserWallet(ctx context.Context, in *GetUserWalletReq, opts ...grpc.CallOption) (*GetUserWalletResp, error) {
	out := new(GetUserWalletResp)
	err := c.cc.Invoke(ctx, "/wallet.wallet/GetUserWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) UpdateCoinRates(ctx context.Context, in *UpdateCoinRatesReq, opts ...grpc.CallOption) (*CommonResp, error) {
	out := new(CommonResp)
	err := c.cc.Invoke(ctx, "/wallet.wallet/UpdateCoinRates", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) GetTransactionDetailRPC(ctx context.Context, in *GetTransactionDetailReq, opts ...grpc.CallOption) (*GetTransactionDetailRes, error) {
	out := new(GetTransactionDetailRes)
	err := c.cc.Invoke(ctx, "/wallet.wallet/GetTransactionDetailRPC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) GetTransactionListRPC(ctx context.Context, in *GetTransactionListReq, opts ...grpc.CallOption) (*GetTransactionListRes, error) {
	out := new(GetTransactionListRes)
	err := c.cc.Invoke(ctx, "/wallet.wallet/GetTransactionListRPC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WalletServer is the server API for Wallet service.
// All implementations should embed UnimplementedWalletServer
// for forward compatibility
type WalletServer interface {
	TestWalletRPC(context.Context, *CommonReq) (*CommonResp, error)
	GetSupportTokenAddressesRPC(context.Context, *GetSupportTokenAddressesReq) (*GetSupportTokenAddressesResp, error)
	CreateAccountInformation(context.Context, *CreateAccountInfoReq) (*CreateAccountInfoResp, error)
	UpdateAccountInformation(context.Context, *UpdateAccountInfoReq) (*CommonResp, error)
	GetFundsLog(context.Context, *GetFundsLogReq) (*GetFundsLogResp, error)
	GetCoinStatuses(context.Context, *GetCoinStatusesReq) (*GetCoinStatusesResp, error)
	GetCoinRatio(context.Context, *GetCoinRatioReq) (*GetCoinRatioResp, error)
	GetUserWallet(context.Context, *GetUserWalletReq) (*GetUserWalletResp, error)
	UpdateCoinRates(context.Context, *UpdateCoinRatesReq) (*CommonResp, error)
	GetTransactionDetailRPC(context.Context, *GetTransactionDetailReq) (*GetTransactionDetailRes, error)
	GetTransactionListRPC(context.Context, *GetTransactionListReq) (*GetTransactionListRes, error)
}

// UnimplementedWalletServer should be embedded to have forward compatible implementations.
type UnimplementedWalletServer struct {
}

func (UnimplementedWalletServer) TestWalletRPC(context.Context, *CommonReq) (*CommonResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestWalletRPC not implemented")
}
func (UnimplementedWalletServer) GetSupportTokenAddressesRPC(context.Context, *GetSupportTokenAddressesReq) (*GetSupportTokenAddressesResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSupportTokenAddressesRPC not implemented")
}
func (UnimplementedWalletServer) CreateAccountInformation(context.Context, *CreateAccountInfoReq) (*CreateAccountInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccountInformation not implemented")
}
func (UnimplementedWalletServer) UpdateAccountInformation(context.Context, *UpdateAccountInfoReq) (*CommonResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAccountInformation not implemented")
}
func (UnimplementedWalletServer) GetFundsLog(context.Context, *GetFundsLogReq) (*GetFundsLogResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFundsLog not implemented")
}
func (UnimplementedWalletServer) GetCoinStatuses(context.Context, *GetCoinStatusesReq) (*GetCoinStatusesResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCoinStatuses not implemented")
}
func (UnimplementedWalletServer) GetCoinRatio(context.Context, *GetCoinRatioReq) (*GetCoinRatioResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCoinRatio not implemented")
}
func (UnimplementedWalletServer) GetUserWallet(context.Context, *GetUserWalletReq) (*GetUserWalletResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserWallet not implemented")
}
func (UnimplementedWalletServer) UpdateCoinRates(context.Context, *UpdateCoinRatesReq) (*CommonResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCoinRates not implemented")
}
func (UnimplementedWalletServer) GetTransactionDetailRPC(context.Context, *GetTransactionDetailReq) (*GetTransactionDetailRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransactionDetailRPC not implemented")
}
func (UnimplementedWalletServer) GetTransactionListRPC(context.Context, *GetTransactionListReq) (*GetTransactionListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransactionListRPC not implemented")
}

// UnsafeWalletServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WalletServer will
// result in compilation errors.
type UnsafeWalletServer interface {
	mustEmbedUnimplementedWalletServer()
}

func RegisterWalletServer(s grpc.ServiceRegistrar, srv WalletServer) {
	s.RegisterService(&Wallet_ServiceDesc, srv)
}

func _Wallet_TestWalletRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).TestWalletRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.wallet/TestWalletRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).TestWalletRPC(ctx, req.(*CommonReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_GetSupportTokenAddressesRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSupportTokenAddressesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).GetSupportTokenAddressesRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.wallet/GetSupportTokenAddressesRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).GetSupportTokenAddressesRPC(ctx, req.(*GetSupportTokenAddressesReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_CreateAccountInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).CreateAccountInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.wallet/CreateAccountInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).CreateAccountInformation(ctx, req.(*CreateAccountInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_UpdateAccountInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAccountInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).UpdateAccountInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.wallet/UpdateAccountInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).UpdateAccountInformation(ctx, req.(*UpdateAccountInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_GetFundsLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFundsLogReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).GetFundsLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.wallet/GetFundsLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).GetFundsLog(ctx, req.(*GetFundsLogReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_GetCoinStatuses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCoinStatusesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).GetCoinStatuses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.wallet/GetCoinStatuses",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).GetCoinStatuses(ctx, req.(*GetCoinStatusesReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_GetCoinRatio_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCoinRatioReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).GetCoinRatio(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.wallet/GetCoinRatio",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).GetCoinRatio(ctx, req.(*GetCoinRatioReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_GetUserWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserWalletReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).GetUserWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.wallet/GetUserWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).GetUserWallet(ctx, req.(*GetUserWalletReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_UpdateCoinRates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCoinRatesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).UpdateCoinRates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.wallet/UpdateCoinRates",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).UpdateCoinRates(ctx, req.(*UpdateCoinRatesReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_GetTransactionDetailRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionDetailReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).GetTransactionDetailRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.wallet/GetTransactionDetailRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).GetTransactionDetailRPC(ctx, req.(*GetTransactionDetailReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_GetTransactionListRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).GetTransactionListRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.wallet/GetTransactionListRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).GetTransactionListRPC(ctx, req.(*GetTransactionListReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Wallet_ServiceDesc is the grpc.ServiceDesc for Wallet service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Wallet_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "wallet.wallet",
	HandlerType: (*WalletServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TestWalletRPC",
			Handler:    _Wallet_TestWalletRPC_Handler,
		},
		{
			MethodName: "GetSupportTokenAddressesRPC",
			Handler:    _Wallet_GetSupportTokenAddressesRPC_Handler,
		},
		{
			MethodName: "CreateAccountInformation",
			Handler:    _Wallet_CreateAccountInformation_Handler,
		},
		{
			MethodName: "UpdateAccountInformation",
			Handler:    _Wallet_UpdateAccountInformation_Handler,
		},
		{
			MethodName: "GetFundsLog",
			Handler:    _Wallet_GetFundsLog_Handler,
		},
		{
			MethodName: "GetCoinStatuses",
			Handler:    _Wallet_GetCoinStatuses_Handler,
		},
		{
			MethodName: "GetCoinRatio",
			Handler:    _Wallet_GetCoinRatio_Handler,
		},
		{
			MethodName: "GetUserWallet",
			Handler:    _Wallet_GetUserWallet_Handler,
		},
		{
			MethodName: "UpdateCoinRates",
			Handler:    _Wallet_UpdateCoinRates_Handler,
		},
		{
			MethodName: "GetTransactionDetailRPC",
			Handler:    _Wallet_GetTransactionDetailRPC_Handler,
		},
		{
			MethodName: "GetTransactionListRPC",
			Handler:    _Wallet_GetTransactionListRPC_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wallet/wallet.proto",
}
