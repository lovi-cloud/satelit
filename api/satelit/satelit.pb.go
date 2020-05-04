// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        v3.9.2
// source: satelit.proto

package satelit

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

// type
type Volume struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Volume) Reset() {
	*x = Volume{}
	if protoimpl.UnsafeEnabled {
		mi := &file_satelit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Volume) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Volume) ProtoMessage() {}

func (x *Volume) ProtoReflect() protoreflect.Message {
	mi := &file_satelit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Volume.ProtoReflect.Descriptor instead.
func (*Volume) Descriptor() ([]byte, []int) {
	return file_satelit_proto_rawDescGZIP(), []int{0}
}

func (x *Volume) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Request / Response
type GetVolumesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetVolumesRequest) Reset() {
	*x = GetVolumesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_satelit_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVolumesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVolumesRequest) ProtoMessage() {}

func (x *GetVolumesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_satelit_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVolumesRequest.ProtoReflect.Descriptor instead.
func (*GetVolumesRequest) Descriptor() ([]byte, []int) {
	return file_satelit_proto_rawDescGZIP(), []int{1}
}

type GetVolumesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Volumes []*Volume `protobuf:"bytes,1,rep,name=volumes,proto3" json:"volumes,omitempty"`
}

func (x *GetVolumesResponse) Reset() {
	*x = GetVolumesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_satelit_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVolumesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVolumesResponse) ProtoMessage() {}

func (x *GetVolumesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_satelit_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVolumesResponse.ProtoReflect.Descriptor instead.
func (*GetVolumesResponse) Descriptor() ([]byte, []int) {
	return file_satelit_proto_rawDescGZIP(), []int{2}
}

func (x *GetVolumesResponse) GetVolumes() []*Volume {
	if x != nil {
		return x.Volumes
	}
	return nil
}

type AddVolumeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddVolumeRequest) Reset() {
	*x = AddVolumeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_satelit_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddVolumeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddVolumeRequest) ProtoMessage() {}

func (x *AddVolumeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_satelit_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddVolumeRequest.ProtoReflect.Descriptor instead.
func (*AddVolumeRequest) Descriptor() ([]byte, []int) {
	return file_satelit_proto_rawDescGZIP(), []int{3}
}

type AddVolumeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddVolumeResponse) Reset() {
	*x = AddVolumeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_satelit_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddVolumeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddVolumeResponse) ProtoMessage() {}

func (x *AddVolumeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_satelit_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddVolumeResponse.ProtoReflect.Descriptor instead.
func (*AddVolumeResponse) Descriptor() ([]byte, []int) {
	return file_satelit_proto_rawDescGZIP(), []int{4}
}

type DeleteVolumeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteVolumeRequest) Reset() {
	*x = DeleteVolumeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_satelit_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteVolumeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteVolumeRequest) ProtoMessage() {}

func (x *DeleteVolumeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_satelit_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteVolumeRequest.ProtoReflect.Descriptor instead.
func (*DeleteVolumeRequest) Descriptor() ([]byte, []int) {
	return file_satelit_proto_rawDescGZIP(), []int{5}
}

type DeleteVolumeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteVolumeResponse) Reset() {
	*x = DeleteVolumeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_satelit_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteVolumeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteVolumeResponse) ProtoMessage() {}

func (x *DeleteVolumeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_satelit_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteVolumeResponse.ProtoReflect.Descriptor instead.
func (*DeleteVolumeResponse) Descriptor() ([]byte, []int) {
	return file_satelit_proto_rawDescGZIP(), []int{6}
}

var File_satelit_proto protoreflect.FileDescriptor

var file_satelit_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x73, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x74, 0x22, 0x18, 0x0a, 0x06, 0x56, 0x6f, 0x6c, 0x75,
	0x6d, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x13, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3f, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x56, 0x6f,
	0x6c, 0x75, 0x6d, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a,
	0x07, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x73, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x74, 0x2e, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52,
	0x07, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x73, 0x22, 0x12, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x56,
	0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x13, 0x0a, 0x11,
	0x41, 0x64, 0x64, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x15, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x56, 0x6f, 0x6c, 0x75, 0x6d,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x16, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x32, 0x52, 0x0a, 0x07, 0x53, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x74, 0x12, 0x47, 0x0a, 0x0a, 0x47,
	0x65, 0x74, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x73, 0x12, 0x1a, 0x2e, 0x73, 0x61, 0x74, 0x65,
	0x6c, 0x69, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x73, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x74, 0x2e,
	0x47, 0x65, 0x74, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x3b, 0x73, 0x61, 0x74, 0x65, 0x6c, 0x69,
	0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_satelit_proto_rawDescOnce sync.Once
	file_satelit_proto_rawDescData = file_satelit_proto_rawDesc
)

func file_satelit_proto_rawDescGZIP() []byte {
	file_satelit_proto_rawDescOnce.Do(func() {
		file_satelit_proto_rawDescData = protoimpl.X.CompressGZIP(file_satelit_proto_rawDescData)
	})
	return file_satelit_proto_rawDescData
}

var file_satelit_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_satelit_proto_goTypes = []interface{}{
	(*Volume)(nil),               // 0: satelit.Volume
	(*GetVolumesRequest)(nil),    // 1: satelit.GetVolumesRequest
	(*GetVolumesResponse)(nil),   // 2: satelit.GetVolumesResponse
	(*AddVolumeRequest)(nil),     // 3: satelit.AddVolumeRequest
	(*AddVolumeResponse)(nil),    // 4: satelit.AddVolumeResponse
	(*DeleteVolumeRequest)(nil),  // 5: satelit.DeleteVolumeRequest
	(*DeleteVolumeResponse)(nil), // 6: satelit.DeleteVolumeResponse
}
var file_satelit_proto_depIdxs = []int32{
	0, // 0: satelit.GetVolumesResponse.volumes:type_name -> satelit.Volume
	1, // 1: satelit.Satelit.GetVolumes:input_type -> satelit.GetVolumesRequest
	2, // 2: satelit.Satelit.GetVolumes:output_type -> satelit.GetVolumesResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_satelit_proto_init() }
func file_satelit_proto_init() {
	if File_satelit_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_satelit_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Volume); i {
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
		file_satelit_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVolumesRequest); i {
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
		file_satelit_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVolumesResponse); i {
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
		file_satelit_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddVolumeRequest); i {
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
		file_satelit_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddVolumeResponse); i {
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
		file_satelit_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteVolumeRequest); i {
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
		file_satelit_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteVolumeResponse); i {
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
			RawDescriptor: file_satelit_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_satelit_proto_goTypes,
		DependencyIndexes: file_satelit_proto_depIdxs,
		MessageInfos:      file_satelit_proto_msgTypes,
	}.Build()
	File_satelit_proto = out.File
	file_satelit_proto_rawDesc = nil
	file_satelit_proto_goTypes = nil
	file_satelit_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SatelitClient is the client API for Satelit service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SatelitClient interface {
	GetVolumes(ctx context.Context, in *GetVolumesRequest, opts ...grpc.CallOption) (*GetVolumesResponse, error)
}

type satelitClient struct {
	cc grpc.ClientConnInterface
}

func NewSatelitClient(cc grpc.ClientConnInterface) SatelitClient {
	return &satelitClient{cc}
}

func (c *satelitClient) GetVolumes(ctx context.Context, in *GetVolumesRequest, opts ...grpc.CallOption) (*GetVolumesResponse, error) {
	out := new(GetVolumesResponse)
	err := c.cc.Invoke(ctx, "/satelit.Satelit/GetVolumes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SatelitServer is the server API for Satelit service.
type SatelitServer interface {
	GetVolumes(context.Context, *GetVolumesRequest) (*GetVolumesResponse, error)
}

// UnimplementedSatelitServer can be embedded to have forward compatible implementations.
type UnimplementedSatelitServer struct {
}

func (*UnimplementedSatelitServer) GetVolumes(context.Context, *GetVolumesRequest) (*GetVolumesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVolumes not implemented")
}

func RegisterSatelitServer(s *grpc.Server, srv SatelitServer) {
	s.RegisterService(&_Satelit_serviceDesc, srv)
}

func _Satelit_GetVolumes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVolumesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SatelitServer).GetVolumes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/satelit.Satelit/GetVolumes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SatelitServer).GetVolumes(ctx, req.(*GetVolumesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Satelit_serviceDesc = grpc.ServiceDesc{
	ServiceName: "satelit.Satelit",
	HandlerType: (*SatelitServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVolumes",
			Handler:    _Satelit_GetVolumes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "satelit.proto",
}
