// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: product_service.proto

package inventories

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type ListProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination        *Pagination `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	ProductCategoryId string      `protobuf:"bytes,2,opt,name=product_category_id,json=productCategoryId,proto3" json:"product_category_id,omitempty"`
	BrandId           string      `protobuf:"bytes,3,opt,name=brand_id,json=brandId,proto3" json:"brand_id,omitempty"`
}

func (x *ListProductRequest) Reset() {
	*x = ListProductRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListProductRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProductRequest) ProtoMessage() {}

func (x *ListProductRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListProductRequest.ProtoReflect.Descriptor instead.
func (*ListProductRequest) Descriptor() ([]byte, []int) {
	return file_product_service_proto_rawDescGZIP(), []int{0}
}

func (x *ListProductRequest) GetPagination() *Pagination {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListProductRequest) GetProductCategoryId() string {
	if x != nil {
		return x.ProductCategoryId
	}
	return ""
}

func (x *ListProductRequest) GetBrandId() string {
	if x != nil {
		return x.BrandId
	}
	return ""
}

type ProductPaginationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination        *Pagination `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	ProductCategoryId string      `protobuf:"bytes,2,opt,name=product_category_id,json=productCategoryId,proto3" json:"product_category_id,omitempty"`
	BrandId           string      `protobuf:"bytes,3,opt,name=brand_id,json=brandId,proto3" json:"brand_id,omitempty"`
	Count             uint32      `protobuf:"varint,4,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *ProductPaginationResponse) Reset() {
	*x = ProductPaginationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProductPaginationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductPaginationResponse) ProtoMessage() {}

func (x *ProductPaginationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_product_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductPaginationResponse.ProtoReflect.Descriptor instead.
func (*ProductPaginationResponse) Descriptor() ([]byte, []int) {
	return file_product_service_proto_rawDescGZIP(), []int{1}
}

func (x *ProductPaginationResponse) GetPagination() *Pagination {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ProductPaginationResponse) GetProductCategoryId() string {
	if x != nil {
		return x.ProductCategoryId
	}
	return ""
}

func (x *ProductPaginationResponse) GetBrandId() string {
	if x != nil {
		return x.BrandId
	}
	return ""
}

func (x *ProductPaginationResponse) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type ListProductResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *ProductPaginationResponse `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Product    *Product                   `protobuf:"bytes,2,opt,name=product,proto3" json:"product,omitempty"`
}

func (x *ListProductResponse) Reset() {
	*x = ListProductResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListProductResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProductResponse) ProtoMessage() {}

func (x *ListProductResponse) ProtoReflect() protoreflect.Message {
	mi := &file_product_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListProductResponse.ProtoReflect.Descriptor instead.
func (*ListProductResponse) Descriptor() ([]byte, []int) {
	return file_product_service_proto_rawDescGZIP(), []int{2}
}

func (x *ListProductResponse) GetPagination() *ProductPaginationResponse {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListProductResponse) GetProduct() *Product {
	if x != nil {
		return x.Product
	}
	return nil
}

var File_product_service_proto protoreflect.FileDescriptor

var file_product_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x1a, 0x15, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x5f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa1, 0x01, 0x0a, 0x12,
	0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x40, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x50, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2e, 0x0a, 0x13, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f,
	0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x11, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x49, 0x64, 0x22,
	0xbe, 0x01, 0x0a, 0x19, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x50, 0x61, 0x67, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a,
	0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x20, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76,
	0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x2e, 0x0a, 0x13, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x63, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x64, 0x12,
	0x19, 0x0a, 0x08, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x22, 0x9f, 0x01, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4f, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x77,
	0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x69, 0x65, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x50, 0x61, 0x67, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x0a, 0x70,
	0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x37, 0x0a, 0x07, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x77, 0x69, 0x72,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65,
	0x73, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x32, 0x8d, 0x03, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12,
	0x1d, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e,
	0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x1a, 0x1d,
	0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74,
	0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x22, 0x00, 0x12,
	0x48, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x77, 0x69, 0x72, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73,
	0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x1a, 0x1d, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e,
	0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x04, 0x56, 0x69, 0x65,
	0x77, 0x12, 0x18, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76,
	0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x49, 0x64, 0x1a, 0x1d, 0x2e, 0x77, 0x69,
	0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69,
	0x65, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x06,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x49, 0x64,
	0x1a, 0x1d, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65,
	0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x22,
	0x00, 0x12, 0x5f, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x28, 0x2e, 0x77, 0x69, 0x72, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x69,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x30, 0x01, 0x42, 0x1c, 0x5a, 0x1a, 0x70, 0x62, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x69, 0x65, 0x73, 0x3b, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_product_service_proto_rawDescOnce sync.Once
	file_product_service_proto_rawDescData = file_product_service_proto_rawDesc
)

func file_product_service_proto_rawDescGZIP() []byte {
	file_product_service_proto_rawDescOnce.Do(func() {
		file_product_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_product_service_proto_rawDescData)
	})
	return file_product_service_proto_rawDescData
}

var file_product_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_product_service_proto_goTypes = []interface{}{
	(*ListProductRequest)(nil),        // 0: wiradata.inventories.ListProductRequest
	(*ProductPaginationResponse)(nil), // 1: wiradata.inventories.ProductPaginationResponse
	(*ListProductResponse)(nil),       // 2: wiradata.inventories.ListProductResponse
	(*Pagination)(nil),                // 3: wiradata.inventories.Pagination
	(*Product)(nil),                   // 4: wiradata.inventories.Product
	(*Id)(nil),                        // 5: wiradata.inventories.Id
	(*Boolean)(nil),                   // 6: wiradata.inventories.Boolean
}
var file_product_service_proto_depIdxs = []int32{
	3, // 0: wiradata.inventories.ListProductRequest.pagination:type_name -> wiradata.inventories.Pagination
	3, // 1: wiradata.inventories.ProductPaginationResponse.pagination:type_name -> wiradata.inventories.Pagination
	1, // 2: wiradata.inventories.ListProductResponse.pagination:type_name -> wiradata.inventories.ProductPaginationResponse
	4, // 3: wiradata.inventories.ListProductResponse.product:type_name -> wiradata.inventories.Product
	4, // 4: wiradata.inventories.ProductService.Create:input_type -> wiradata.inventories.Product
	4, // 5: wiradata.inventories.ProductService.Update:input_type -> wiradata.inventories.Product
	5, // 6: wiradata.inventories.ProductService.View:input_type -> wiradata.inventories.Id
	5, // 7: wiradata.inventories.ProductService.Delete:input_type -> wiradata.inventories.Id
	0, // 8: wiradata.inventories.ProductService.List:input_type -> wiradata.inventories.ListProductRequest
	4, // 9: wiradata.inventories.ProductService.Create:output_type -> wiradata.inventories.Product
	4, // 10: wiradata.inventories.ProductService.Update:output_type -> wiradata.inventories.Product
	4, // 11: wiradata.inventories.ProductService.View:output_type -> wiradata.inventories.Product
	6, // 12: wiradata.inventories.ProductService.Delete:output_type -> wiradata.inventories.Boolean
	2, // 13: wiradata.inventories.ProductService.List:output_type -> wiradata.inventories.ListProductResponse
	9, // [9:14] is the sub-list for method output_type
	4, // [4:9] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_product_service_proto_init() }
func file_product_service_proto_init() {
	if File_product_service_proto != nil {
		return
	}
	file_product_message_proto_init()
	file_generic_message_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_product_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListProductRequest); i {
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
		file_product_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProductPaginationResponse); i {
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
		file_product_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListProductResponse); i {
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
			RawDescriptor: file_product_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_product_service_proto_goTypes,
		DependencyIndexes: file_product_service_proto_depIdxs,
		MessageInfos:      file_product_service_proto_msgTypes,
	}.Build()
	File_product_service_proto = out.File
	file_product_service_proto_rawDesc = nil
	file_product_service_proto_goTypes = nil
	file_product_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ProductServiceClient is the client API for ProductService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProductServiceClient interface {
	Create(ctx context.Context, in *Product, opts ...grpc.CallOption) (*Product, error)
	Update(ctx context.Context, in *Product, opts ...grpc.CallOption) (*Product, error)
	View(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Product, error)
	Delete(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Boolean, error)
	List(ctx context.Context, in *ListProductRequest, opts ...grpc.CallOption) (ProductService_ListClient, error)
}

type productServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProductServiceClient(cc grpc.ClientConnInterface) ProductServiceClient {
	return &productServiceClient{cc}
}

func (c *productServiceClient) Create(ctx context.Context, in *Product, opts ...grpc.CallOption) (*Product, error) {
	out := new(Product)
	err := c.cc.Invoke(ctx, "/wiradata.inventories.ProductService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) Update(ctx context.Context, in *Product, opts ...grpc.CallOption) (*Product, error) {
	out := new(Product)
	err := c.cc.Invoke(ctx, "/wiradata.inventories.ProductService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) View(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Product, error) {
	out := new(Product)
	err := c.cc.Invoke(ctx, "/wiradata.inventories.ProductService/View", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) Delete(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Boolean, error) {
	out := new(Boolean)
	err := c.cc.Invoke(ctx, "/wiradata.inventories.ProductService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) List(ctx context.Context, in *ListProductRequest, opts ...grpc.CallOption) (ProductService_ListClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ProductService_serviceDesc.Streams[0], "/wiradata.inventories.ProductService/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &productServiceListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ProductService_ListClient interface {
	Recv() (*ListProductResponse, error)
	grpc.ClientStream
}

type productServiceListClient struct {
	grpc.ClientStream
}

func (x *productServiceListClient) Recv() (*ListProductResponse, error) {
	m := new(ListProductResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ProductServiceServer is the server API for ProductService service.
type ProductServiceServer interface {
	Create(context.Context, *Product) (*Product, error)
	Update(context.Context, *Product) (*Product, error)
	View(context.Context, *Id) (*Product, error)
	Delete(context.Context, *Id) (*Boolean, error)
	List(*ListProductRequest, ProductService_ListServer) error
}

// UnimplementedProductServiceServer can be embedded to have forward compatible implementations.
type UnimplementedProductServiceServer struct {
}

func (*UnimplementedProductServiceServer) Create(context.Context, *Product) (*Product, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedProductServiceServer) Update(context.Context, *Product) (*Product, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedProductServiceServer) View(context.Context, *Id) (*Product, error) {
	return nil, status.Errorf(codes.Unimplemented, "method View not implemented")
}
func (*UnimplementedProductServiceServer) Delete(context.Context, *Id) (*Boolean, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (*UnimplementedProductServiceServer) List(*ListProductRequest, ProductService_ListServer) error {
	return status.Errorf(codes.Unimplemented, "method List not implemented")
}

func RegisterProductServiceServer(s *grpc.Server, srv ProductServiceServer) {
	s.RegisterService(&_ProductService_serviceDesc, srv)
}

func _ProductService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Product)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wiradata.inventories.ProductService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).Create(ctx, req.(*Product))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Product)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wiradata.inventories.ProductService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).Update(ctx, req.(*Product))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_View_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).View(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wiradata.inventories.ProductService/View",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).View(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wiradata.inventories.ProductService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).Delete(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListProductRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ProductServiceServer).List(m, &productServiceListServer{stream})
}

type ProductService_ListServer interface {
	Send(*ListProductResponse) error
	grpc.ServerStream
}

type productServiceListServer struct {
	grpc.ServerStream
}

func (x *productServiceListServer) Send(m *ListProductResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _ProductService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "wiradata.inventories.ProductService",
	HandlerType: (*ProductServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _ProductService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ProductService_Update_Handler,
		},
		{
			MethodName: "View",
			Handler:    _ProductService_View_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ProductService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "List",
			Handler:       _ProductService_List_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "product_service.proto",
}