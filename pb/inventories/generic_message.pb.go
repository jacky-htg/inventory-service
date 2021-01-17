// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: generic_message.proto

package inventories

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Pagination_Sort int32

const (
	Pagination_ASC  Pagination_Sort = 0
	Pagination_DESC Pagination_Sort = 1
)

// Enum value maps for Pagination_Sort.
var (
	Pagination_Sort_name = map[int32]string{
		0: "ASC",
		1: "DESC",
	}
	Pagination_Sort_value = map[string]int32{
		"ASC":  0,
		"DESC": 1,
	}
)

func (x Pagination_Sort) Enum() *Pagination_Sort {
	p := new(Pagination_Sort)
	*p = x
	return p
}

func (x Pagination_Sort) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Pagination_Sort) Descriptor() protoreflect.EnumDescriptor {
	return file_generic_message_proto_enumTypes[0].Descriptor()
}

func (Pagination_Sort) Type() protoreflect.EnumType {
	return &file_generic_message_proto_enumTypes[0]
}

func (x Pagination_Sort) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Pagination_Sort.Descriptor instead.
func (Pagination_Sort) EnumDescriptor() ([]byte, []int) {
	return file_generic_message_proto_rawDescGZIP(), []int{4, 0}
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_generic_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_generic_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_generic_message_proto_rawDescGZIP(), []int{0}
}

type Id struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Id) Reset() {
	*x = Id{}
	if protoimpl.UnsafeEnabled {
		mi := &file_generic_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Id) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Id) ProtoMessage() {}

func (x *Id) ProtoReflect() protoreflect.Message {
	mi := &file_generic_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Id.ProtoReflect.Descriptor instead.
func (*Id) Descriptor() ([]byte, []int) {
	return file_generic_message_proto_rawDescGZIP(), []int{1}
}

func (x *Id) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type String struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	String_ string `protobuf:"bytes,1,opt,name=string,proto3" json:"string,omitempty"`
}

func (x *String) Reset() {
	*x = String{}
	if protoimpl.UnsafeEnabled {
		mi := &file_generic_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *String) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*String) ProtoMessage() {}

func (x *String) ProtoReflect() protoreflect.Message {
	mi := &file_generic_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use String.ProtoReflect.Descriptor instead.
func (*String) Descriptor() ([]byte, []int) {
	return file_generic_message_proto_rawDescGZIP(), []int{2}
}

func (x *String) GetString_() string {
	if x != nil {
		return x.String_
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_generic_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_generic_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_generic_message_proto_rawDescGZIP(), []int{3}
}

func (x *Message) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type Pagination struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page    uint32          `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Limit   uint32          `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset  uint32          `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
	Search  string          `protobuf:"bytes,4,opt,name=search,proto3" json:"search,omitempty"`
	OrderBy string          `protobuf:"bytes,5,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Sort    Pagination_Sort `protobuf:"varint,6,opt,name=sort,proto3,enum=wiradata.inventories.Pagination_Sort" json:"sort,omitempty"`
}

func (x *Pagination) Reset() {
	*x = Pagination{}
	if protoimpl.UnsafeEnabled {
		mi := &file_generic_message_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pagination) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pagination) ProtoMessage() {}

func (x *Pagination) ProtoReflect() protoreflect.Message {
	mi := &file_generic_message_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pagination.ProtoReflect.Descriptor instead.
func (*Pagination) Descriptor() ([]byte, []int) {
	return file_generic_message_proto_rawDescGZIP(), []int{4}
}

func (x *Pagination) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *Pagination) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *Pagination) GetOffset() uint32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *Pagination) GetSearch() string {
	if x != nil {
		return x.Search
	}
	return ""
}

func (x *Pagination) GetOrderBy() string {
	if x != nil {
		return x.OrderBy
	}
	return ""
}

func (x *Pagination) GetSort() Pagination_Sort {
	if x != nil {
		return x.Sort
	}
	return Pagination_ASC
}

type PaginationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *Pagination `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Count      uint32      `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *PaginationResponse) Reset() {
	*x = PaginationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_generic_message_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaginationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaginationResponse) ProtoMessage() {}

func (x *PaginationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_generic_message_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaginationResponse.ProtoReflect.Descriptor instead.
func (*PaginationResponse) Descriptor() ([]byte, []int) {
	return file_generic_message_proto_rawDescGZIP(), []int{5}
}

func (x *PaginationResponse) GetPagination() *Pagination {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *PaginationResponse) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type Boolean struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Boolean bool `protobuf:"varint,1,opt,name=boolean,proto3" json:"boolean,omitempty"`
}

func (x *Boolean) Reset() {
	*x = Boolean{}
	if protoimpl.UnsafeEnabled {
		mi := &file_generic_message_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Boolean) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Boolean) ProtoMessage() {}

func (x *Boolean) ProtoReflect() protoreflect.Message {
	mi := &file_generic_message_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Boolean.ProtoReflect.Descriptor instead.
func (*Boolean) Descriptor() ([]byte, []int) {
	return file_generic_message_proto_rawDescGZIP(), []int{6}
}

func (x *Boolean) GetBoolean() bool {
	if x != nil {
		return x.Boolean
	}
	return false
}

var File_generic_message_proto protoreflect.FileDescriptor

var file_generic_message_proto_rawDesc = []byte{
	0x0a, 0x15, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x22, 0x07, 0x0a,
	0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x14, 0x0a, 0x02, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x20, 0x0a, 0x06,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x23,
	0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0xd7, 0x01, 0x0a, 0x0a, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x19, 0x0a, 0x08,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12, 0x39, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x50, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x6f, 0x72, 0x74, 0x52, 0x04, 0x73, 0x6f,
	0x72, 0x74, 0x22, 0x19, 0x0a, 0x04, 0x53, 0x6f, 0x72, 0x74, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x53,
	0x43, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x45, 0x53, 0x43, 0x10, 0x01, 0x22, 0x6c, 0x0a,
	0x12, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x50,
	0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x23, 0x0a, 0x07, 0x42,
	0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e,
	0x42, 0x1c, 0x5a, 0x1a, 0x70, 0x62, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69,
	0x65, 0x73, 0x3b, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_generic_message_proto_rawDescOnce sync.Once
	file_generic_message_proto_rawDescData = file_generic_message_proto_rawDesc
)

func file_generic_message_proto_rawDescGZIP() []byte {
	file_generic_message_proto_rawDescOnce.Do(func() {
		file_generic_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_generic_message_proto_rawDescData)
	})
	return file_generic_message_proto_rawDescData
}

var file_generic_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_generic_message_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_generic_message_proto_goTypes = []interface{}{
	(Pagination_Sort)(0),       // 0: wiradata.inventories.Pagination.Sort
	(*Empty)(nil),              // 1: wiradata.inventories.Empty
	(*Id)(nil),                 // 2: wiradata.inventories.Id
	(*String)(nil),             // 3: wiradata.inventories.String
	(*Message)(nil),            // 4: wiradata.inventories.Message
	(*Pagination)(nil),         // 5: wiradata.inventories.Pagination
	(*PaginationResponse)(nil), // 6: wiradata.inventories.PaginationResponse
	(*Boolean)(nil),            // 7: wiradata.inventories.Boolean
}
var file_generic_message_proto_depIdxs = []int32{
	0, // 0: wiradata.inventories.Pagination.sort:type_name -> wiradata.inventories.Pagination.Sort
	5, // 1: wiradata.inventories.PaginationResponse.pagination:type_name -> wiradata.inventories.Pagination
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_generic_message_proto_init() }
func file_generic_message_proto_init() {
	if File_generic_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_generic_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_generic_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Id); i {
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
		file_generic_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*String); i {
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
		file_generic_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_generic_message_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pagination); i {
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
		file_generic_message_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaginationResponse); i {
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
		file_generic_message_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Boolean); i {
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
			RawDescriptor: file_generic_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_generic_message_proto_goTypes,
		DependencyIndexes: file_generic_message_proto_depIdxs,
		EnumInfos:         file_generic_message_proto_enumTypes,
		MessageInfos:      file_generic_message_proto_msgTypes,
	}.Build()
	File_generic_message_proto = out.File
	file_generic_message_proto_rawDesc = nil
	file_generic_message_proto_goTypes = nil
	file_generic_message_proto_depIdxs = nil
}