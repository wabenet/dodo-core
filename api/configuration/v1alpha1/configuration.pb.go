// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.12.4
// source: configuration/v1alpha1/configuration.proto

package v1alpha1

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	v1alpha5 "github.com/wabenet/dodo-core/api/core/v1alpha5"
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

type ListBackdropsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Backdrops []*v1alpha5.Backdrop `protobuf:"bytes,1,rep,name=backdrops,proto3" json:"backdrops,omitempty"`
}

func (x *ListBackdropsResponse) Reset() {
	*x = ListBackdropsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_configuration_v1alpha1_configuration_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListBackdropsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListBackdropsResponse) ProtoMessage() {}

func (x *ListBackdropsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_configuration_v1alpha1_configuration_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListBackdropsResponse.ProtoReflect.Descriptor instead.
func (*ListBackdropsResponse) Descriptor() ([]byte, []int) {
	return file_configuration_v1alpha1_configuration_proto_rawDescGZIP(), []int{0}
}

func (x *ListBackdropsResponse) GetBackdrops() []*v1alpha5.Backdrop {
	if x != nil {
		return x.Backdrops
	}
	return nil
}

type GetBackdropRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Alias string `protobuf:"bytes,1,opt,name=alias,proto3" json:"alias,omitempty"`
}

func (x *GetBackdropRequest) Reset() {
	*x = GetBackdropRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_configuration_v1alpha1_configuration_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBackdropRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBackdropRequest) ProtoMessage() {}

func (x *GetBackdropRequest) ProtoReflect() protoreflect.Message {
	mi := &file_configuration_v1alpha1_configuration_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBackdropRequest.ProtoReflect.Descriptor instead.
func (*GetBackdropRequest) Descriptor() ([]byte, []int) {
	return file_configuration_v1alpha1_configuration_proto_rawDescGZIP(), []int{1}
}

func (x *GetBackdropRequest) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

var File_configuration_v1alpha1_configuration_proto protoreflect.FileDescriptor

var file_configuration_v1alpha1_configuration_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x27, 0x63, 0x6f,
	0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1a, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x35, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x35, 0x2f, 0x62, 0x61,
	0x63, 0x6b, 0x64, 0x72, 0x6f, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5f, 0x0a, 0x15,
	0x4c, 0x69, 0x73, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x64, 0x72, 0x6f, 0x70, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x09, 0x62, 0x61, 0x63, 0x6b, 0x64, 0x72, 0x6f,
	0x70, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x77,
	0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x35, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x64, 0x72,
	0x6f, 0x70, 0x52, 0x09, 0x62, 0x61, 0x63, 0x6b, 0x64, 0x72, 0x6f, 0x70, 0x73, 0x22, 0x2a, 0x0a,
	0x12, 0x47, 0x65, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x64, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x32, 0xd5, 0x03, 0x0a, 0x06, 0x50, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x12, 0x53, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2a, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x35, 0x2e, 0x50,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x58, 0x0a, 0x0a, 0x49, 0x6e, 0x69,
	0x74, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x32, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f,
	0x64, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x35,
	0x2e, 0x49, 0x6e, 0x69, 0x74, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x65, 0x74, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x67, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x64, 0x72,
	0x6f, 0x70, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x3e, 0x2e, 0x63, 0x6f,
	0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x64, 0x72,
	0x6f, 0x70, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x74, 0x0a, 0x0b, 0x47,
	0x65, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x64, 0x72, 0x6f, 0x70, 0x12, 0x3b, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x64, 0x72, 0x6f, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x61,
	0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x35, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x64, 0x72, 0x6f,
	0x70, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x64, 0x6f, 0x64, 0x6f, 0x2d, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_configuration_v1alpha1_configuration_proto_rawDescOnce sync.Once
	file_configuration_v1alpha1_configuration_proto_rawDescData = file_configuration_v1alpha1_configuration_proto_rawDesc
)

func file_configuration_v1alpha1_configuration_proto_rawDescGZIP() []byte {
	file_configuration_v1alpha1_configuration_proto_rawDescOnce.Do(func() {
		file_configuration_v1alpha1_configuration_proto_rawDescData = protoimpl.X.CompressGZIP(file_configuration_v1alpha1_configuration_proto_rawDescData)
	})
	return file_configuration_v1alpha1_configuration_proto_rawDescData
}

var file_configuration_v1alpha1_configuration_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_configuration_v1alpha1_configuration_proto_goTypes = []interface{}{
	(*ListBackdropsResponse)(nil),       // 0: com.wabenet.dodo.configuration.v1alpha1.ListBackdropsResponse
	(*GetBackdropRequest)(nil),          // 1: com.wabenet.dodo.configuration.v1alpha1.GetBackdropRequest
	(*v1alpha5.Backdrop)(nil),           // 2: com.wabenet.dodo.core.v1alpha5.Backdrop
	(*empty.Empty)(nil),                 // 3: google.protobuf.Empty
	(*v1alpha5.PluginInfo)(nil),         // 4: com.wabenet.dodo.core.v1alpha5.PluginInfo
	(*v1alpha5.InitPluginResponse)(nil), // 5: com.wabenet.dodo.core.v1alpha5.InitPluginResponse
}
var file_configuration_v1alpha1_configuration_proto_depIdxs = []int32{
	2, // 0: com.wabenet.dodo.configuration.v1alpha1.ListBackdropsResponse.backdrops:type_name -> com.wabenet.dodo.core.v1alpha5.Backdrop
	3, // 1: com.wabenet.dodo.configuration.v1alpha1.Plugin.GetPluginInfo:input_type -> google.protobuf.Empty
	3, // 2: com.wabenet.dodo.configuration.v1alpha1.Plugin.InitPlugin:input_type -> google.protobuf.Empty
	3, // 3: com.wabenet.dodo.configuration.v1alpha1.Plugin.ResetPlugin:input_type -> google.protobuf.Empty
	3, // 4: com.wabenet.dodo.configuration.v1alpha1.Plugin.ListBackdrops:input_type -> google.protobuf.Empty
	1, // 5: com.wabenet.dodo.configuration.v1alpha1.Plugin.GetBackdrop:input_type -> com.wabenet.dodo.configuration.v1alpha1.GetBackdropRequest
	4, // 6: com.wabenet.dodo.configuration.v1alpha1.Plugin.GetPluginInfo:output_type -> com.wabenet.dodo.core.v1alpha5.PluginInfo
	5, // 7: com.wabenet.dodo.configuration.v1alpha1.Plugin.InitPlugin:output_type -> com.wabenet.dodo.core.v1alpha5.InitPluginResponse
	3, // 8: com.wabenet.dodo.configuration.v1alpha1.Plugin.ResetPlugin:output_type -> google.protobuf.Empty
	0, // 9: com.wabenet.dodo.configuration.v1alpha1.Plugin.ListBackdrops:output_type -> com.wabenet.dodo.configuration.v1alpha1.ListBackdropsResponse
	2, // 10: com.wabenet.dodo.configuration.v1alpha1.Plugin.GetBackdrop:output_type -> com.wabenet.dodo.core.v1alpha5.Backdrop
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_configuration_v1alpha1_configuration_proto_init() }
func file_configuration_v1alpha1_configuration_proto_init() {
	if File_configuration_v1alpha1_configuration_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_configuration_v1alpha1_configuration_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListBackdropsResponse); i {
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
		file_configuration_v1alpha1_configuration_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBackdropRequest); i {
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
			RawDescriptor: file_configuration_v1alpha1_configuration_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_configuration_v1alpha1_configuration_proto_goTypes,
		DependencyIndexes: file_configuration_v1alpha1_configuration_proto_depIdxs,
		MessageInfos:      file_configuration_v1alpha1_configuration_proto_msgTypes,
	}.Build()
	File_configuration_v1alpha1_configuration_proto = out.File
	file_configuration_v1alpha1_configuration_proto_rawDesc = nil
	file_configuration_v1alpha1_configuration_proto_goTypes = nil
	file_configuration_v1alpha1_configuration_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PluginClient is the client API for Plugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PluginClient interface {
	GetPluginInfo(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*v1alpha5.PluginInfo, error)
	InitPlugin(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*v1alpha5.InitPluginResponse, error)
	ResetPlugin(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
	ListBackdrops(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ListBackdropsResponse, error)
	GetBackdrop(ctx context.Context, in *GetBackdropRequest, opts ...grpc.CallOption) (*v1alpha5.Backdrop, error)
}

type pluginClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginClient(cc grpc.ClientConnInterface) PluginClient {
	return &pluginClient{cc}
}

func (c *pluginClient) GetPluginInfo(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*v1alpha5.PluginInfo, error) {
	out := new(v1alpha5.PluginInfo)
	err := c.cc.Invoke(ctx, "/com.wabenet.dodo.configuration.v1alpha1.Plugin/GetPluginInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) InitPlugin(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*v1alpha5.InitPluginResponse, error) {
	out := new(v1alpha5.InitPluginResponse)
	err := c.cc.Invoke(ctx, "/com.wabenet.dodo.configuration.v1alpha1.Plugin/InitPlugin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) ResetPlugin(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/com.wabenet.dodo.configuration.v1alpha1.Plugin/ResetPlugin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) ListBackdrops(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ListBackdropsResponse, error) {
	out := new(ListBackdropsResponse)
	err := c.cc.Invoke(ctx, "/com.wabenet.dodo.configuration.v1alpha1.Plugin/ListBackdrops", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) GetBackdrop(ctx context.Context, in *GetBackdropRequest, opts ...grpc.CallOption) (*v1alpha5.Backdrop, error) {
	out := new(v1alpha5.Backdrop)
	err := c.cc.Invoke(ctx, "/com.wabenet.dodo.configuration.v1alpha1.Plugin/GetBackdrop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServer is the server API for Plugin service.
type PluginServer interface {
	GetPluginInfo(context.Context, *empty.Empty) (*v1alpha5.PluginInfo, error)
	InitPlugin(context.Context, *empty.Empty) (*v1alpha5.InitPluginResponse, error)
	ResetPlugin(context.Context, *empty.Empty) (*empty.Empty, error)
	ListBackdrops(context.Context, *empty.Empty) (*ListBackdropsResponse, error)
	GetBackdrop(context.Context, *GetBackdropRequest) (*v1alpha5.Backdrop, error)
}

// UnimplementedPluginServer can be embedded to have forward compatible implementations.
type UnimplementedPluginServer struct {
}

func (*UnimplementedPluginServer) GetPluginInfo(context.Context, *empty.Empty) (*v1alpha5.PluginInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPluginInfo not implemented")
}
func (*UnimplementedPluginServer) InitPlugin(context.Context, *empty.Empty) (*v1alpha5.InitPluginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitPlugin not implemented")
}
func (*UnimplementedPluginServer) ResetPlugin(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetPlugin not implemented")
}
func (*UnimplementedPluginServer) ListBackdrops(context.Context, *empty.Empty) (*ListBackdropsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBackdrops not implemented")
}
func (*UnimplementedPluginServer) GetBackdrop(context.Context, *GetBackdropRequest) (*v1alpha5.Backdrop, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBackdrop not implemented")
}

func RegisterPluginServer(s *grpc.Server, srv PluginServer) {
	s.RegisterService(&_Plugin_serviceDesc, srv)
}

func _Plugin_GetPluginInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetPluginInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.wabenet.dodo.configuration.v1alpha1.Plugin/GetPluginInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetPluginInfo(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_InitPlugin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).InitPlugin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.wabenet.dodo.configuration.v1alpha1.Plugin/InitPlugin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).InitPlugin(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_ResetPlugin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).ResetPlugin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.wabenet.dodo.configuration.v1alpha1.Plugin/ResetPlugin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).ResetPlugin(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_ListBackdrops_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).ListBackdrops(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.wabenet.dodo.configuration.v1alpha1.Plugin/ListBackdrops",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).ListBackdrops(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_GetBackdrop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBackdropRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetBackdrop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.wabenet.dodo.configuration.v1alpha1.Plugin/GetBackdrop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetBackdrop(ctx, req.(*GetBackdropRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Plugin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "com.wabenet.dodo.configuration.v1alpha1.Plugin",
	HandlerType: (*PluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPluginInfo",
			Handler:    _Plugin_GetPluginInfo_Handler,
		},
		{
			MethodName: "InitPlugin",
			Handler:    _Plugin_InitPlugin_Handler,
		},
		{
			MethodName: "ResetPlugin",
			Handler:    _Plugin_ResetPlugin_Handler,
		},
		{
			MethodName: "ListBackdrops",
			Handler:    _Plugin_ListBackdrops_Handler,
		},
		{
			MethodName: "GetBackdrop",
			Handler:    _Plugin_GetBackdrop_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "configuration/v1alpha1/configuration.proto",
}