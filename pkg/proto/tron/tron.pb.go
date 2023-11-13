// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.11
// source: tron/tron.proto

package tron

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CommonReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OperationID string `protobuf:"bytes,1,opt,name=OperationID,proto3" json:"OperationID,omitempty"`
}

func (x *CommonReq) Reset() {
	*x = CommonReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tron_tron_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommonReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommonReq) ProtoMessage() {}

func (x *CommonReq) ProtoReflect() protoreflect.Message {
	mi := &file_tron_tron_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommonReq.ProtoReflect.Descriptor instead.
func (*CommonReq) Descriptor() ([]byte, []int) {
	return file_tron_tron_proto_rawDescGZIP(), []int{0}
}

func (x *CommonReq) GetOperationID() string {
	if x != nil {
		return x.OperationID
	}
	return ""
}

type CommonResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrCode int32  `protobuf:"varint,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg  string `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
}

func (x *CommonResp) Reset() {
	*x = CommonResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tron_tron_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommonResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommonResp) ProtoMessage() {}

func (x *CommonResp) ProtoReflect() protoreflect.Message {
	mi := &file_tron_tron_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommonResp.ProtoReflect.Descriptor instead.
func (*CommonResp) Descriptor() ([]byte, []int) {
	return file_tron_tron_proto_rawDescGZIP(), []int{1}
}

func (x *CommonResp) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *CommonResp) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

type GetBalanceReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OperationID string `protobuf:"bytes,1,opt,name=OperationID,proto3" json:"OperationID,omitempty"`
	CoinType    uint32 `protobuf:"varint,2,opt,name=CoinType,proto3" json:"CoinType,omitempty"`
	Address     string `protobuf:"bytes,3,opt,name=Address,proto3" json:"Address,omitempty"`
}

func (x *GetBalanceReq) Reset() {
	*x = GetBalanceReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tron_tron_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBalanceReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBalanceReq) ProtoMessage() {}

func (x *GetBalanceReq) ProtoReflect() protoreflect.Message {
	mi := &file_tron_tron_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBalanceReq.ProtoReflect.Descriptor instead.
func (*GetBalanceReq) Descriptor() ([]byte, []int) {
	return file_tron_tron_proto_rawDescGZIP(), []int{2}
}

func (x *GetBalanceReq) GetOperationID() string {
	if x != nil {
		return x.OperationID
	}
	return ""
}

func (x *GetBalanceReq) GetCoinType() uint32 {
	if x != nil {
		return x.CoinType
	}
	return 0
}

func (x *GetBalanceReq) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type GetBalanceResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrCode int32  `protobuf:"varint,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg  string `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	Balance string `protobuf:"bytes,3,opt,name=Balance,proto3" json:"Balance,omitempty"`
}

func (x *GetBalanceResp) Reset() {
	*x = GetBalanceResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tron_tron_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBalanceResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBalanceResp) ProtoMessage() {}

func (x *GetBalanceResp) ProtoReflect() protoreflect.Message {
	mi := &file_tron_tron_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBalanceResp.ProtoReflect.Descriptor instead.
func (*GetBalanceResp) Descriptor() ([]byte, []int) {
	return file_tron_tron_proto_rawDescGZIP(), []int{3}
}

func (x *GetBalanceResp) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *GetBalanceResp) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

func (x *GetBalanceResp) GetBalance() string {
	if x != nil {
		return x.Balance
	}
	return ""
}

type PostTransferReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OperationID     string `protobuf:"bytes,1,opt,name=OperationID,proto3" json:"OperationID,omitempty"`
	CoinType        uint32 `protobuf:"varint,2,opt,name=CoinType,proto3" json:"CoinType,omitempty"`
	FromAccountUID  string `protobuf:"bytes,3,opt,name=FromAccountUID,proto3" json:"FromAccountUID,omitempty"`
	FromMerchantUID string `protobuf:"bytes,4,opt,name=FromMerchantUID,proto3" json:"FromMerchantUID,omitempty"`
	FromAddress     string `protobuf:"bytes,5,opt,name=FromAddress,proto3" json:"FromAddress,omitempty"`
	ToAddress       string `protobuf:"bytes,6,opt,name=ToAddress,proto3" json:"ToAddress,omitempty"`
	Amount          string `protobuf:"bytes,7,opt,name=Amount,proto3" json:"Amount,omitempty"`
	TxID            string `protobuf:"bytes,8,opt,name=TxID,proto3" json:"TxID,omitempty"`
	TxDataStr       string `protobuf:"bytes,9,opt,name=TxDataStr,proto3" json:"TxDataStr,omitempty"`
	EnergyUsed      string `protobuf:"bytes,10,opt,name=EnergyUsed,proto3" json:"EnergyUsed,omitempty"`
	EnergyPenalty   string `protobuf:"bytes,11,opt,name=EnergyPenalty,proto3" json:"EnergyPenalty,omitempty"`
}

func (x *PostTransferReq) Reset() {
	*x = PostTransferReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tron_tron_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostTransferReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostTransferReq) ProtoMessage() {}

func (x *PostTransferReq) ProtoReflect() protoreflect.Message {
	mi := &file_tron_tron_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostTransferReq.ProtoReflect.Descriptor instead.
func (*PostTransferReq) Descriptor() ([]byte, []int) {
	return file_tron_tron_proto_rawDescGZIP(), []int{4}
}

func (x *PostTransferReq) GetOperationID() string {
	if x != nil {
		return x.OperationID
	}
	return ""
}

func (x *PostTransferReq) GetCoinType() uint32 {
	if x != nil {
		return x.CoinType
	}
	return 0
}

func (x *PostTransferReq) GetFromAccountUID() string {
	if x != nil {
		return x.FromAccountUID
	}
	return ""
}

func (x *PostTransferReq) GetFromMerchantUID() string {
	if x != nil {
		return x.FromMerchantUID
	}
	return ""
}

func (x *PostTransferReq) GetFromAddress() string {
	if x != nil {
		return x.FromAddress
	}
	return ""
}

func (x *PostTransferReq) GetToAddress() string {
	if x != nil {
		return x.ToAddress
	}
	return ""
}

func (x *PostTransferReq) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *PostTransferReq) GetTxID() string {
	if x != nil {
		return x.TxID
	}
	return ""
}

func (x *PostTransferReq) GetTxDataStr() string {
	if x != nil {
		return x.TxDataStr
	}
	return ""
}

func (x *PostTransferReq) GetEnergyUsed() string {
	if x != nil {
		return x.EnergyUsed
	}
	return ""
}

func (x *PostTransferReq) GetEnergyPenalty() string {
	if x != nil {
		return x.EnergyPenalty
	}
	return ""
}

type TransferRPCResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrCode int32  `protobuf:"varint,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg  string `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
}

func (x *TransferRPCResponse) Reset() {
	*x = TransferRPCResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tron_tron_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransferRPCResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferRPCResponse) ProtoMessage() {}

func (x *TransferRPCResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tron_tron_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferRPCResponse.ProtoReflect.Descriptor instead.
func (*TransferRPCResponse) Descriptor() ([]byte, []int) {
	return file_tron_tron_proto_rawDescGZIP(), []int{5}
}

func (x *TransferRPCResponse) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *TransferRPCResponse) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

type CreateTransactionReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OperationID string `protobuf:"bytes,1,opt,name=OperationID,proto3" json:"OperationID,omitempty"`
	CoinType    uint32 `protobuf:"varint,2,opt,name=CoinType,proto3" json:"CoinType,omitempty"`
	FromAddress string `protobuf:"bytes,3,opt,name=FromAddress,proto3" json:"FromAddress,omitempty"`
	ToAddress   string `protobuf:"bytes,4,opt,name=ToAddress,proto3" json:"ToAddress,omitempty"`
	Amount      string `protobuf:"bytes,5,opt,name=Amount,proto3" json:"Amount,omitempty"`
}

func (x *CreateTransactionReq) Reset() {
	*x = CreateTransactionReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tron_tron_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTransactionReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTransactionReq) ProtoMessage() {}

func (x *CreateTransactionReq) ProtoReflect() protoreflect.Message {
	mi := &file_tron_tron_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTransactionReq.ProtoReflect.Descriptor instead.
func (*CreateTransactionReq) Descriptor() ([]byte, []int) {
	return file_tron_tron_proto_rawDescGZIP(), []int{6}
}

func (x *CreateTransactionReq) GetOperationID() string {
	if x != nil {
		return x.OperationID
	}
	return ""
}

func (x *CreateTransactionReq) GetCoinType() uint32 {
	if x != nil {
		return x.CoinType
	}
	return 0
}

func (x *CreateTransactionReq) GetFromAddress() string {
	if x != nil {
		return x.FromAddress
	}
	return ""
}

func (x *CreateTransactionReq) GetToAddress() string {
	if x != nil {
		return x.ToAddress
	}
	return ""
}

func (x *CreateTransactionReq) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

type CreateTransactionResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrCode   int32  `protobuf:"varint,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg    string `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	TxID      string `protobuf:"bytes,3,opt,name=TxID,proto3" json:"TxID,omitempty"`
	RawTXData string `protobuf:"bytes,4,opt,name=RawTXData,proto3" json:"RawTXData,omitempty"`
}

func (x *CreateTransactionResp) Reset() {
	*x = CreateTransactionResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tron_tron_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTransactionResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTransactionResp) ProtoMessage() {}

func (x *CreateTransactionResp) ProtoReflect() protoreflect.Message {
	mi := &file_tron_tron_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTransactionResp.ProtoReflect.Descriptor instead.
func (*CreateTransactionResp) Descriptor() ([]byte, []int) {
	return file_tron_tron_proto_rawDescGZIP(), []int{7}
}

func (x *CreateTransactionResp) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *CreateTransactionResp) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

func (x *CreateTransactionResp) GetTxID() string {
	if x != nil {
		return x.TxID
	}
	return ""
}

func (x *CreateTransactionResp) GetRawTXData() string {
	if x != nil {
		return x.RawTXData
	}
	return ""
}

type GetTronConfirmationReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OperationID     string `protobuf:"bytes,1,opt,name=OperationID,proto3" json:"OperationID,omitempty"`
	CoinType        uint32 `protobuf:"varint,2,opt,name=CoinType,proto3" json:"CoinType,omitempty"`
	TransactionHash string `protobuf:"bytes,3,opt,name=TransactionHash,proto3" json:"TransactionHash,omitempty"`
	MessageID       string `protobuf:"bytes,4,opt,name=MessageID,proto3" json:"MessageID,omitempty"`
}

func (x *GetTronConfirmationReq) Reset() {
	*x = GetTronConfirmationReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tron_tron_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTronConfirmationReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTronConfirmationReq) ProtoMessage() {}

func (x *GetTronConfirmationReq) ProtoReflect() protoreflect.Message {
	mi := &file_tron_tron_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTronConfirmationReq.ProtoReflect.Descriptor instead.
func (*GetTronConfirmationReq) Descriptor() ([]byte, []int) {
	return file_tron_tron_proto_rawDescGZIP(), []int{8}
}

func (x *GetTronConfirmationReq) GetOperationID() string {
	if x != nil {
		return x.OperationID
	}
	return ""
}

func (x *GetTronConfirmationReq) GetCoinType() uint32 {
	if x != nil {
		return x.CoinType
	}
	return 0
}

func (x *GetTronConfirmationReq) GetTransactionHash() string {
	if x != nil {
		return x.TransactionHash
	}
	return ""
}

func (x *GetTronConfirmationReq) GetMessageID() string {
	if x != nil {
		return x.MessageID
	}
	return ""
}

type GetTronConfirmationRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockNum    string `protobuf:"bytes,1,opt,name=BlockNum,proto3" json:"BlockNum,omitempty"`
	Status      int32  `protobuf:"varint,2,opt,name=Status,proto3" json:"Status,omitempty"`
	ConfirmTime uint64 `protobuf:"varint,3,opt,name=ConfirmTime,proto3" json:"ConfirmTime,omitempty"`
	NetFee      string `protobuf:"bytes,4,opt,name=NetFee,proto3" json:"NetFee,omitempty"`
	EnergyUsage string `protobuf:"bytes,5,opt,name=EnergyUsage,proto3" json:"EnergyUsage,omitempty"`
}

func (x *GetTronConfirmationRes) Reset() {
	*x = GetTronConfirmationRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tron_tron_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTronConfirmationRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTronConfirmationRes) ProtoMessage() {}

func (x *GetTronConfirmationRes) ProtoReflect() protoreflect.Message {
	mi := &file_tron_tron_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTronConfirmationRes.ProtoReflect.Descriptor instead.
func (*GetTronConfirmationRes) Descriptor() ([]byte, []int) {
	return file_tron_tron_proto_rawDescGZIP(), []int{9}
}

func (x *GetTronConfirmationRes) GetBlockNum() string {
	if x != nil {
		return x.BlockNum
	}
	return ""
}

func (x *GetTronConfirmationRes) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *GetTronConfirmationRes) GetConfirmTime() uint64 {
	if x != nil {
		return x.ConfirmTime
	}
	return 0
}

func (x *GetTronConfirmationRes) GetNetFee() string {
	if x != nil {
		return x.NetFee
	}
	return ""
}

func (x *GetTronConfirmationRes) GetEnergyUsage() string {
	if x != nil {
		return x.EnergyUsage
	}
	return ""
}

var File_tron_tron_proto protoreflect.FileDescriptor

var file_tron_tron_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x74, 0x72, 0x6f, 0x6e, 0x2f, 0x74, 0x72, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x74, 0x72, 0x6f, 0x6e, 0x22, 0x2d, 0x0a, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x12, 0x20, 0x0a, 0x0b, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x22, 0x3e, 0x0a, 0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x22, 0x67, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x42, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x12, 0x20, 0x0a, 0x0b, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6f, 0x69,
	0x6e, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x43, 0x6f, 0x69,
	0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22,
	0x5c, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x45,
	0x72, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x45, 0x72, 0x72,
	0x4d, 0x73, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x22, 0xf1, 0x02,
	0x0a, 0x0f, 0x50, 0x6f, 0x73, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x12, 0x20, 0x0a, 0x0b, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6f, 0x69, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x43, 0x6f, 0x69, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x26, 0x0a, 0x0e, 0x46, 0x72, 0x6f, 0x6d, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x55, 0x49,
	0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x46, 0x72, 0x6f, 0x6d, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x55, 0x49, 0x44, 0x12, 0x28, 0x0a, 0x0f, 0x46, 0x72, 0x6f, 0x6d, 0x4d,
	0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x55, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0f, 0x46, 0x72, 0x6f, 0x6d, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x55, 0x49,
	0x44, 0x12, 0x20, 0x0a, 0x0b, 0x46, 0x72, 0x6f, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x46, 0x72, 0x6f, 0x6d, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x78, 0x49,
	0x44, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x78, 0x49, 0x44, 0x12, 0x1c, 0x0a,
	0x09, 0x54, 0x78, 0x44, 0x61, 0x74, 0x61, 0x53, 0x74, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x54, 0x78, 0x44, 0x61, 0x74, 0x61, 0x53, 0x74, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x45,
	0x6e, 0x65, 0x72, 0x67, 0x79, 0x55, 0x73, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x45, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x55, 0x73, 0x65, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x45,
	0x6e, 0x65, 0x72, 0x67, 0x79, 0x50, 0x65, 0x6e, 0x61, 0x6c, 0x74, 0x79, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x45, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x50, 0x65, 0x6e, 0x61, 0x6c, 0x74,
	0x79, 0x22, 0x47, 0x0a, 0x13, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52, 0x50, 0x43,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f,
	0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x22, 0xac, 0x01, 0x0a, 0x14, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x12, 0x20, 0x0a, 0x0b, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6f, 0x69, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x43, 0x6f, 0x69, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x46, 0x72, 0x6f, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x46, 0x72, 0x6f, 0x6d, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x7b, 0x0a, 0x15, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x45, 0x72,
	0x72, 0x4d, 0x73, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x78, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x54, 0x78, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x61, 0x77, 0x54,
	0x58, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x52, 0x61, 0x77,
	0x54, 0x58, 0x44, 0x61, 0x74, 0x61, 0x22, 0x9e, 0x01, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x54, 0x72,
	0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x12, 0x20, 0x0a, 0x0b, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6f, 0x69, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x43, 0x6f, 0x69, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x28, 0x0a, 0x0f, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x61,
	0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x22, 0xa8, 0x01, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x54,
	0x72, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x12, 0x16,
	0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72,
	0x6d, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x72, 0x6d, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x65, 0x74, 0x46,
	0x65, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x65, 0x74, 0x46, 0x65, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x45, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x45, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x55, 0x73, 0x61,
	0x67, 0x65, 0x32, 0xaa, 0x02, 0x0a, 0x04, 0x74, 0x72, 0x6f, 0x6e, 0x12, 0x3e, 0x0a, 0x11, 0x47,
	0x65, 0x74, 0x54, 0x72, 0x6f, 0x6e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x50, 0x43,
	0x12, 0x13, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74,
	0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x4f, 0x0a, 0x14, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x50, 0x43, 0x12, 0x1a, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x1a,
	0x1b, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x3f, 0x0a, 0x0b,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52, 0x50, 0x43, 0x12, 0x15, 0x2e, 0x74, 0x72,
	0x6f, 0x6e, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x1a, 0x19, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66,
	0x65, 0x72, 0x52, 0x50, 0x43, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a,
	0x12, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x50, 0x43, 0x12, 0x1c, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x72,
	0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x1a, 0x1c, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x72, 0x6f, 0x6e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x42,
	0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x74, 0x72, 0x6f, 0x6e, 0x3b, 0x74, 0x72, 0x6f, 0x6e, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tron_tron_proto_rawDescOnce sync.Once
	file_tron_tron_proto_rawDescData = file_tron_tron_proto_rawDesc
)

func file_tron_tron_proto_rawDescGZIP() []byte {
	file_tron_tron_proto_rawDescOnce.Do(func() {
		file_tron_tron_proto_rawDescData = protoimpl.X.CompressGZIP(file_tron_tron_proto_rawDescData)
	})
	return file_tron_tron_proto_rawDescData
}

var file_tron_tron_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_tron_tron_proto_goTypes = []interface{}{
	(*CommonReq)(nil),              // 0: tron.CommonReq
	(*CommonResp)(nil),             // 1: tron.CommonResp
	(*GetBalanceReq)(nil),          // 2: tron.GetBalanceReq
	(*GetBalanceResp)(nil),         // 3: tron.GetBalanceResp
	(*PostTransferReq)(nil),        // 4: tron.PostTransferReq
	(*TransferRPCResponse)(nil),    // 5: tron.TransferRPCResponse
	(*CreateTransactionReq)(nil),   // 6: tron.CreateTransactionReq
	(*CreateTransactionResp)(nil),  // 7: tron.CreateTransactionResp
	(*GetTronConfirmationReq)(nil), // 8: tron.GetTronConfirmationReq
	(*GetTronConfirmationRes)(nil), // 9: tron.GetTronConfirmationRes
}
var file_tron_tron_proto_depIdxs = []int32{
	2, // 0: tron.tron.GetTronBalanceRPC:input_type -> tron.GetBalanceReq
	6, // 1: tron.tron.CreateTransactionRPC:input_type -> tron.CreateTransactionReq
	4, // 2: tron.tron.TransferRPC:input_type -> tron.PostTransferReq
	8, // 3: tron.tron.GetConfirmationRPC:input_type -> tron.GetTronConfirmationReq
	3, // 4: tron.tron.GetTronBalanceRPC:output_type -> tron.GetBalanceResp
	7, // 5: tron.tron.CreateTransactionRPC:output_type -> tron.CreateTransactionResp
	5, // 6: tron.tron.TransferRPC:output_type -> tron.TransferRPCResponse
	9, // 7: tron.tron.GetConfirmationRPC:output_type -> tron.GetTronConfirmationRes
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_tron_tron_proto_init() }
func file_tron_tron_proto_init() {
	if File_tron_tron_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tron_tron_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommonReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tron_tron_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommonResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tron_tron_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBalanceReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tron_tron_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBalanceResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tron_tron_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostTransferReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tron_tron_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransferRPCResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tron_tron_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTransactionReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tron_tron_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTransactionResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tron_tron_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTronConfirmationReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tron_tron_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTronConfirmationRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tron_tron_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tron_tron_proto_goTypes,
		DependencyIndexes: file_tron_tron_proto_depIdxs,
		MessageInfos:      file_tron_tron_proto_msgTypes,
	}.Build()
	File_tron_tron_proto = out.File
	file_tron_tron_proto_rawDesc = nil
	file_tron_tron_proto_goTypes = nil
	file_tron_tron_proto_depIdxs = nil
}
