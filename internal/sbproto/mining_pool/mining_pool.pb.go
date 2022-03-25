// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: protolib/mining_pool.proto

package mining_pool

import (
	snowblossom "go.snowblossom/internal/sbproto/snowblossom"
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

type GetWorkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId     string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	PayToAddress string `protobuf:"bytes,2,opt,name=pay_to_address,json=payToAddress,proto3" json:"pay_to_address,omitempty"`
}

func (x *GetWorkRequest) Reset() {
	*x = GetWorkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protolib_mining_pool_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWorkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWorkRequest) ProtoMessage() {}

func (x *GetWorkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protolib_mining_pool_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWorkRequest.ProtoReflect.Descriptor instead.
func (*GetWorkRequest) Descriptor() ([]byte, []int) {
	return file_protolib_mining_pool_proto_rawDescGZIP(), []int{0}
}

func (x *GetWorkRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *GetWorkRequest) GetPayToAddress() string {
	if x != nil {
		return x.PayToAddress
	}
	return ""
}

type WorkUnit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header       *snowblossom.BlockHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	WorkId       int32                    `protobuf:"varint,2,opt,name=work_id,json=workId,proto3" json:"work_id,omitempty"`
	ReportTarget []byte                   `protobuf:"bytes,3,opt,name=report_target,json=reportTarget,proto3" json:"report_target,omitempty"`
}

func (x *WorkUnit) Reset() {
	*x = WorkUnit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protolib_mining_pool_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkUnit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkUnit) ProtoMessage() {}

func (x *WorkUnit) ProtoReflect() protoreflect.Message {
	mi := &file_protolib_mining_pool_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkUnit.ProtoReflect.Descriptor instead.
func (*WorkUnit) Descriptor() ([]byte, []int) {
	return file_protolib_mining_pool_proto_rawDescGZIP(), []int{1}
}

func (x *WorkUnit) GetHeader() *snowblossom.BlockHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *WorkUnit) GetWorkId() int32 {
	if x != nil {
		return x.WorkId
	}
	return 0
}

func (x *WorkUnit) GetReportTarget() []byte {
	if x != nil {
		return x.ReportTarget
	}
	return nil
}

type WorkSubmitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header *snowblossom.BlockHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	WorkId int32                    `protobuf:"varint,2,opt,name=work_id,json=workId,proto3" json:"work_id,omitempty"`
}

func (x *WorkSubmitRequest) Reset() {
	*x = WorkSubmitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protolib_mining_pool_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkSubmitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkSubmitRequest) ProtoMessage() {}

func (x *WorkSubmitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protolib_mining_pool_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkSubmitRequest.ProtoReflect.Descriptor instead.
func (*WorkSubmitRequest) Descriptor() ([]byte, []int) {
	return file_protolib_mining_pool_proto_rawDescGZIP(), []int{2}
}

func (x *WorkSubmitRequest) GetHeader() *snowblossom.BlockHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *WorkSubmitRequest) GetWorkId() int32 {
	if x != nil {
		return x.WorkId
	}
	return 0
}

type PPLNSState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShareEntries []*ShareEntry `protobuf:"bytes,1,rep,name=share_entries,json=shareEntries,proto3" json:"share_entries,omitempty"`
}

func (x *PPLNSState) Reset() {
	*x = PPLNSState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protolib_mining_pool_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PPLNSState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PPLNSState) ProtoMessage() {}

func (x *PPLNSState) ProtoReflect() protoreflect.Message {
	mi := &file_protolib_mining_pool_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PPLNSState.ProtoReflect.Descriptor instead.
func (*PPLNSState) Descriptor() ([]byte, []int) {
	return file_protolib_mining_pool_proto_rawDescGZIP(), []int{3}
}

func (x *PPLNSState) GetShareEntries() []*ShareEntry {
	if x != nil {
		return x.ShareEntries
	}
	return nil
}

type ShareEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address    string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	ShareCount int64  `protobuf:"varint,2,opt,name=share_count,json=shareCount,proto3" json:"share_count,omitempty"`
}

func (x *ShareEntry) Reset() {
	*x = ShareEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protolib_mining_pool_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShareEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShareEntry) ProtoMessage() {}

func (x *ShareEntry) ProtoReflect() protoreflect.Message {
	mi := &file_protolib_mining_pool_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShareEntry.ProtoReflect.Descriptor instead.
func (*ShareEntry) Descriptor() ([]byte, []int) {
	return file_protolib_mining_pool_proto_rawDescGZIP(), []int{4}
}

func (x *ShareEntry) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *ShareEntry) GetShareCount() int64 {
	if x != nil {
		return x.ShareCount
	}
	return 0
}

type GetWordsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WordIndexes []int64 `protobuf:"varint,1,rep,packed,name=word_indexes,json=wordIndexes,proto3" json:"word_indexes,omitempty"`
	Field       int32   `protobuf:"varint,2,opt,name=field,proto3" json:"field,omitempty"`
}

func (x *GetWordsRequest) Reset() {
	*x = GetWordsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protolib_mining_pool_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWordsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWordsRequest) ProtoMessage() {}

func (x *GetWordsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protolib_mining_pool_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWordsRequest.ProtoReflect.Descriptor instead.
func (*GetWordsRequest) Descriptor() ([]byte, []int) {
	return file_protolib_mining_pool_proto_rawDescGZIP(), []int{5}
}

func (x *GetWordsRequest) GetWordIndexes() []int64 {
	if x != nil {
		return x.WordIndexes
	}
	return nil
}

func (x *GetWordsRequest) GetField() int32 {
	if x != nil {
		return x.Field
	}
	return 0
}

type GetWordsResponce struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Words      [][]byte `protobuf:"bytes,1,rep,name=words,proto3" json:"words,omitempty"`
	WrongField bool     `protobuf:"varint,2,opt,name=wrong_field,json=wrongField,proto3" json:"wrong_field,omitempty"`
}

func (x *GetWordsResponce) Reset() {
	*x = GetWordsResponce{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protolib_mining_pool_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWordsResponce) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWordsResponce) ProtoMessage() {}

func (x *GetWordsResponce) ProtoReflect() protoreflect.Message {
	mi := &file_protolib_mining_pool_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWordsResponce.ProtoReflect.Descriptor instead.
func (*GetWordsResponce) Descriptor() ([]byte, []int) {
	return file_protolib_mining_pool_proto_rawDescGZIP(), []int{6}
}

func (x *GetWordsResponce) GetWords() [][]byte {
	if x != nil {
		return x.Words
	}
	return nil
}

func (x *GetWordsResponce) GetWrongField() bool {
	if x != nil {
		return x.WrongField
	}
	return false
}

var File_protolib_mining_pool_proto protoreflect.FileDescriptor

var file_protolib_mining_pool_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6c, 0x69, 0x62, 0x2f, 0x6d, 0x69, 0x6e, 0x69, 0x6e,
	0x67, 0x5f, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x73, 0x6e,
	0x6f, 0x77, 0x62, 0x6c, 0x6f, 0x73, 0x73, 0x6f, 0x6d, 0x1a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x6c, 0x69, 0x62, 0x2f, 0x73, 0x6e, 0x6f, 0x77, 0x62, 0x6c, 0x6f, 0x73, 0x73, 0x6f, 0x6d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x53, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x70, 0x61, 0x79, 0x5f, 0x74, 0x6f, 0x5f, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x61,
	0x79, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x7a, 0x0a, 0x08, 0x57, 0x6f,
	0x72, 0x6b, 0x55, 0x6e, 0x69, 0x74, 0x12, 0x30, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x73, 0x6e, 0x6f, 0x77, 0x62, 0x6c, 0x6f,
	0x73, 0x73, 0x6f, 0x6d, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x17, 0x0a, 0x07, 0x77, 0x6f, 0x72, 0x6b,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x77, 0x6f, 0x72, 0x6b, 0x49,
	0x64, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x22, 0x5e, 0x0a, 0x11, 0x57, 0x6f, 0x72, 0x6b, 0x53, 0x75,
	0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x30, 0x0a, 0x06, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x73, 0x6e,
	0x6f, 0x77, 0x62, 0x6c, 0x6f, 0x73, 0x73, 0x6f, 0x6d, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x17, 0x0a,
	0x07, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x77, 0x6f, 0x72, 0x6b, 0x49, 0x64, 0x22, 0x4a, 0x0a, 0x0a, 0x50, 0x50, 0x4c, 0x4e, 0x53, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x3c, 0x0a, 0x0d, 0x73, 0x68, 0x61, 0x72, 0x65, 0x5f, 0x65, 0x6e,
	0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x6e,
	0x6f, 0x77, 0x62, 0x6c, 0x6f, 0x73, 0x73, 0x6f, 0x6d, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x65, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x69,
	0x65, 0x73, 0x22, 0x47, 0x0a, 0x0a, 0x53, 0x68, 0x61, 0x72, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x73, 0x68, 0x61, 0x72, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x4a, 0x0a, 0x0f, 0x47,
	0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21,
	0x0a, 0x0c, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x03, 0x52, 0x0b, 0x77, 0x6f, 0x72, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x22, 0x49, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x57, 0x6f,
	0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x77,
	0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x05, 0x77, 0x6f, 0x72, 0x64,
	0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x72, 0x6f, 0x6e, 0x67, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x77, 0x72, 0x6f, 0x6e, 0x67, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x32, 0xa0, 0x01, 0x0a, 0x11, 0x4d, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x50, 0x6f, 0x6f,
	0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x41, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x57,
	0x6f, 0x72, 0x6b, 0x12, 0x1b, 0x2e, 0x73, 0x6e, 0x6f, 0x77, 0x62, 0x6c, 0x6f, 0x73, 0x73, 0x6f,
	0x6d, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x15, 0x2e, 0x73, 0x6e, 0x6f, 0x77, 0x62, 0x6c, 0x6f, 0x73, 0x73, 0x6f, 0x6d, 0x2e, 0x57,
	0x6f, 0x72, 0x6b, 0x55, 0x6e, 0x69, 0x74, 0x22, 0x00, 0x30, 0x01, 0x12, 0x48, 0x0a, 0x0a, 0x53,
	0x75, 0x62, 0x6d, 0x69, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x12, 0x1e, 0x2e, 0x73, 0x6e, 0x6f, 0x77,
	0x62, 0x6c, 0x6f, 0x73, 0x73, 0x6f, 0x6d, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x53, 0x75, 0x62, 0x6d,
	0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x6e, 0x6f, 0x77,
	0x62, 0x6c, 0x6f, 0x73, 0x73, 0x6f, 0x6d, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x22, 0x00, 0x32, 0x60, 0x0a, 0x13, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4d,
	0x69, 0x6e, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x49, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x1c, 0x2e, 0x73, 0x6e, 0x6f, 0x77, 0x62,
	0x6c, 0x6f, 0x73, 0x73, 0x6f, 0x6d, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x73, 0x6e, 0x6f, 0x77, 0x62, 0x6c, 0x6f,
	0x73, 0x73, 0x6f, 0x6d, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x63, 0x65, 0x22, 0x00, 0x42, 0x61, 0x0a, 0x18, 0x73, 0x6e, 0x6f, 0x77, 0x62,
	0x6c, 0x6f, 0x73, 0x73, 0x6f, 0x6d, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x42, 0x16, 0x53, 0x6e, 0x6f, 0x77, 0x42, 0x6c, 0x6f, 0x73, 0x73, 0x6f, 0x6d,
	0x4d, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2b, 0x67,
	0x6f, 0x2e, 0x73, 0x6e, 0x6f, 0x77, 0x62, 0x6c, 0x6f, 0x73, 0x73, 0x6f, 0x6d, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x62, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d,
	0x69, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x6f, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_protolib_mining_pool_proto_rawDescOnce sync.Once
	file_protolib_mining_pool_proto_rawDescData = file_protolib_mining_pool_proto_rawDesc
)

func file_protolib_mining_pool_proto_rawDescGZIP() []byte {
	file_protolib_mining_pool_proto_rawDescOnce.Do(func() {
		file_protolib_mining_pool_proto_rawDescData = protoimpl.X.CompressGZIP(file_protolib_mining_pool_proto_rawDescData)
	})
	return file_protolib_mining_pool_proto_rawDescData
}

var file_protolib_mining_pool_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_protolib_mining_pool_proto_goTypes = []interface{}{
	(*GetWorkRequest)(nil),          // 0: snowblossom.GetWorkRequest
	(*WorkUnit)(nil),                // 1: snowblossom.WorkUnit
	(*WorkSubmitRequest)(nil),       // 2: snowblossom.WorkSubmitRequest
	(*PPLNSState)(nil),              // 3: snowblossom.PPLNSState
	(*ShareEntry)(nil),              // 4: snowblossom.ShareEntry
	(*GetWordsRequest)(nil),         // 5: snowblossom.GetWordsRequest
	(*GetWordsResponce)(nil),        // 6: snowblossom.GetWordsResponce
	(*snowblossom.BlockHeader)(nil), // 7: snowblossom.BlockHeader
	(*snowblossom.SubmitReply)(nil), // 8: snowblossom.SubmitReply
}
var file_protolib_mining_pool_proto_depIdxs = []int32{
	7, // 0: snowblossom.WorkUnit.header:type_name -> snowblossom.BlockHeader
	7, // 1: snowblossom.WorkSubmitRequest.header:type_name -> snowblossom.BlockHeader
	4, // 2: snowblossom.PPLNSState.share_entries:type_name -> snowblossom.ShareEntry
	0, // 3: snowblossom.MiningPoolService.GetWork:input_type -> snowblossom.GetWorkRequest
	2, // 4: snowblossom.MiningPoolService.SubmitWork:input_type -> snowblossom.WorkSubmitRequest
	5, // 5: snowblossom.SharedMiningService.GetWords:input_type -> snowblossom.GetWordsRequest
	1, // 6: snowblossom.MiningPoolService.GetWork:output_type -> snowblossom.WorkUnit
	8, // 7: snowblossom.MiningPoolService.SubmitWork:output_type -> snowblossom.SubmitReply
	6, // 8: snowblossom.SharedMiningService.GetWords:output_type -> snowblossom.GetWordsResponce
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_protolib_mining_pool_proto_init() }
func file_protolib_mining_pool_proto_init() {
	if File_protolib_mining_pool_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protolib_mining_pool_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWorkRequest); i {
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
		file_protolib_mining_pool_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkUnit); i {
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
		file_protolib_mining_pool_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkSubmitRequest); i {
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
		file_protolib_mining_pool_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PPLNSState); i {
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
		file_protolib_mining_pool_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShareEntry); i {
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
		file_protolib_mining_pool_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWordsRequest); i {
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
		file_protolib_mining_pool_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWordsResponce); i {
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
			RawDescriptor: file_protolib_mining_pool_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_protolib_mining_pool_proto_goTypes,
		DependencyIndexes: file_protolib_mining_pool_proto_depIdxs,
		MessageInfos:      file_protolib_mining_pool_proto_msgTypes,
	}.Build()
	File_protolib_mining_pool_proto = out.File
	file_protolib_mining_pool_proto_rawDesc = nil
	file_protolib_mining_pool_proto_goTypes = nil
	file_protolib_mining_pool_proto_depIdxs = nil
}
