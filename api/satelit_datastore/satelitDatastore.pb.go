// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.11.4
// source: satelitDatastore.proto

package satelit_datastore

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

type DHCPLease struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MacAddress     string `protobuf:"bytes,1,opt,name=mac_address,json=macAddress,proto3" json:"mac_address,omitempty"`
	Ip             string `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	Network        string `protobuf:"bytes,3,opt,name=network,proto3" json:"network,omitempty"`
	Gateway        string `protobuf:"bytes,4,opt,name=gateway,proto3" json:"gateway,omitempty"`
	DnsServer      string `protobuf:"bytes,5,opt,name=dns_server,json=dnsServer,proto3" json:"dns_server,omitempty"`
	MetadataServer string `protobuf:"bytes,6,opt,name=metadata_server,json=metadataServer,proto3" json:"metadata_server,omitempty"`
}

func (x *DHCPLease) Reset() {
	*x = DHCPLease{}
	if protoimpl.UnsafeEnabled {
		mi := &file_satelitDatastore_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DHCPLease) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DHCPLease) ProtoMessage() {}

func (x *DHCPLease) ProtoReflect() protoreflect.Message {
	mi := &file_satelitDatastore_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DHCPLease.ProtoReflect.Descriptor instead.
func (*DHCPLease) Descriptor() ([]byte, []int) {
	return file_satelitDatastore_proto_rawDescGZIP(), []int{0}
}

func (x *DHCPLease) GetMacAddress() string {
	if x != nil {
		return x.MacAddress
	}
	return ""
}

func (x *DHCPLease) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *DHCPLease) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *DHCPLease) GetGateway() string {
	if x != nil {
		return x.Gateway
	}
	return ""
}

func (x *DHCPLease) GetDnsServer() string {
	if x != nil {
		return x.DnsServer
	}
	return ""
}

func (x *DHCPLease) GetMetadataServer() string {
	if x != nil {
		return x.MetadataServer
	}
	return ""
}

type GetDHCPLeaseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MacAddress string `protobuf:"bytes,1,opt,name=mac_address,json=macAddress,proto3" json:"mac_address,omitempty"`
}

func (x *GetDHCPLeaseRequest) Reset() {
	*x = GetDHCPLeaseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_satelitDatastore_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDHCPLeaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDHCPLeaseRequest) ProtoMessage() {}

func (x *GetDHCPLeaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_satelitDatastore_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDHCPLeaseRequest.ProtoReflect.Descriptor instead.
func (*GetDHCPLeaseRequest) Descriptor() ([]byte, []int) {
	return file_satelitDatastore_proto_rawDescGZIP(), []int{1}
}

func (x *GetDHCPLeaseRequest) GetMacAddress() string {
	if x != nil {
		return x.MacAddress
	}
	return ""
}

type GetDHCPLeaseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lease *DHCPLease `protobuf:"bytes,1,opt,name=lease,proto3" json:"lease,omitempty"`
}

func (x *GetDHCPLeaseResponse) Reset() {
	*x = GetDHCPLeaseResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_satelitDatastore_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDHCPLeaseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDHCPLeaseResponse) ProtoMessage() {}

func (x *GetDHCPLeaseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_satelitDatastore_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDHCPLeaseResponse.ProtoReflect.Descriptor instead.
func (*GetDHCPLeaseResponse) Descriptor() ([]byte, []int) {
	return file_satelitDatastore_proto_rawDescGZIP(), []int{2}
}

func (x *GetDHCPLeaseResponse) GetLease() *DHCPLease {
	if x != nil {
		return x.Lease
	}
	return nil
}

var File_satelitDatastore_proto protoreflect.FileDescriptor

var file_satelitDatastore_proto_rawDesc = []byte{
	0x0a, 0x16, 0x73, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x74, 0x44, 0x61, 0x74, 0x61, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x61, 0x74, 0x65, 0x6c, 0x69,
	0x74, 0x22, 0xb8, 0x01, 0x0a, 0x09, 0x44, 0x48, 0x43, 0x50, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x6d, 0x61, 0x63, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x61, 0x63, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70,
	0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x61, 0x74,
	0x65, 0x77, 0x61, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x6e, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x6e, 0x73, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x12, 0x27, 0x0a, 0x0f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x22, 0x36, 0x0a, 0x13,
	0x47, 0x65, 0x74, 0x44, 0x48, 0x43, 0x50, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x61, 0x63, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x61, 0x63, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x22, 0x40, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x44, 0x48, 0x43, 0x50, 0x4c,
	0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x05,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x61,
	0x74, 0x65, 0x6c, 0x69, 0x74, 0x2e, 0x44, 0x48, 0x43, 0x50, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x52,
	0x05, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x32, 0x61, 0x0a, 0x10, 0x53, 0x61, 0x74, 0x65, 0x6c, 0x69,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x4d, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x44, 0x48, 0x43, 0x50, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x1c, 0x2e, 0x73, 0x61, 0x74,
	0x65, 0x6c, 0x69, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x48, 0x43, 0x50, 0x4c, 0x65, 0x61, 0x73,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x73, 0x61, 0x74, 0x65, 0x6c,
	0x69, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x48, 0x43, 0x50, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x15, 0x5a, 0x13, 0x2e, 0x3b, 0x73,
	0x61, 0x74, 0x65, 0x6c, 0x69, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_satelitDatastore_proto_rawDescOnce sync.Once
	file_satelitDatastore_proto_rawDescData = file_satelitDatastore_proto_rawDesc
)

func file_satelitDatastore_proto_rawDescGZIP() []byte {
	file_satelitDatastore_proto_rawDescOnce.Do(func() {
		file_satelitDatastore_proto_rawDescData = protoimpl.X.CompressGZIP(file_satelitDatastore_proto_rawDescData)
	})
	return file_satelitDatastore_proto_rawDescData
}

var file_satelitDatastore_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_satelitDatastore_proto_goTypes = []interface{}{
	(*DHCPLease)(nil),            // 0: satelit.DHCPLease
	(*GetDHCPLeaseRequest)(nil),  // 1: satelit.GetDHCPLeaseRequest
	(*GetDHCPLeaseResponse)(nil), // 2: satelit.GetDHCPLeaseResponse
}
var file_satelitDatastore_proto_depIdxs = []int32{
	0, // 0: satelit.GetDHCPLeaseResponse.lease:type_name -> satelit.DHCPLease
	1, // 1: satelit.SatelitDatastore.GetDHCPLease:input_type -> satelit.GetDHCPLeaseRequest
	2, // 2: satelit.SatelitDatastore.GetDHCPLease:output_type -> satelit.GetDHCPLeaseResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_satelitDatastore_proto_init() }
func file_satelitDatastore_proto_init() {
	if File_satelitDatastore_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_satelitDatastore_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DHCPLease); i {
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
		file_satelitDatastore_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDHCPLeaseRequest); i {
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
		file_satelitDatastore_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDHCPLeaseResponse); i {
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
			RawDescriptor: file_satelitDatastore_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_satelitDatastore_proto_goTypes,
		DependencyIndexes: file_satelitDatastore_proto_depIdxs,
		MessageInfos:      file_satelitDatastore_proto_msgTypes,
	}.Build()
	File_satelitDatastore_proto = out.File
	file_satelitDatastore_proto_rawDesc = nil
	file_satelitDatastore_proto_goTypes = nil
	file_satelitDatastore_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SatelitDatastoreClient is the client API for SatelitDatastore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SatelitDatastoreClient interface {
	GetDHCPLease(ctx context.Context, in *GetDHCPLeaseRequest, opts ...grpc.CallOption) (*GetDHCPLeaseResponse, error)
}

type satelitDatastoreClient struct {
	cc grpc.ClientConnInterface
}

func NewSatelitDatastoreClient(cc grpc.ClientConnInterface) SatelitDatastoreClient {
	return &satelitDatastoreClient{cc}
}

func (c *satelitDatastoreClient) GetDHCPLease(ctx context.Context, in *GetDHCPLeaseRequest, opts ...grpc.CallOption) (*GetDHCPLeaseResponse, error) {
	out := new(GetDHCPLeaseResponse)
	err := c.cc.Invoke(ctx, "/satelit.SatelitDatastore/GetDHCPLease", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SatelitDatastoreServer is the server API for SatelitDatastore service.
type SatelitDatastoreServer interface {
	GetDHCPLease(context.Context, *GetDHCPLeaseRequest) (*GetDHCPLeaseResponse, error)
}

// UnimplementedSatelitDatastoreServer can be embedded to have forward compatible implementations.
type UnimplementedSatelitDatastoreServer struct {
}

func (*UnimplementedSatelitDatastoreServer) GetDHCPLease(context.Context, *GetDHCPLeaseRequest) (*GetDHCPLeaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDHCPLease not implemented")
}

func RegisterSatelitDatastoreServer(s *grpc.Server, srv SatelitDatastoreServer) {
	s.RegisterService(&_SatelitDatastore_serviceDesc, srv)
}

func _SatelitDatastore_GetDHCPLease_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDHCPLeaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SatelitDatastoreServer).GetDHCPLease(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/satelit.SatelitDatastore/GetDHCPLease",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SatelitDatastoreServer).GetDHCPLease(ctx, req.(*GetDHCPLeaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SatelitDatastore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "satelit.SatelitDatastore",
	HandlerType: (*SatelitDatastoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDHCPLease",
			Handler:    _SatelitDatastore_GetDHCPLease_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "satelitDatastore.proto",
}
