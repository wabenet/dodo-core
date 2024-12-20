// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.21.12
// source: build/v1alpha1/build.proto

package v1alpha1

import (
	context "context"
	v1alpha6 "github.com/wabenet/dodo-core/api/core/v1alpha6"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StreamId string              `protobuf:"bytes,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"`
	Config   *v1alpha6.BuildInfo `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	Height   uint32              `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	Width    uint32              `protobuf:"varint,4,opt,name=width,proto3" json:"width,omitempty"`
}

func (x *CreateImageRequest) Reset() {
	*x = CreateImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_build_v1alpha1_build_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateImageRequest) ProtoMessage() {}

func (x *CreateImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_build_v1alpha1_build_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateImageRequest.ProtoReflect.Descriptor instead.
func (*CreateImageRequest) Descriptor() ([]byte, []int) {
	return file_build_v1alpha1_build_proto_rawDescGZIP(), []int{0}
}

func (x *CreateImageRequest) GetStreamId() string {
	if x != nil {
		return x.StreamId
	}
	return ""
}

func (x *CreateImageRequest) GetConfig() *v1alpha6.BuildInfo {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *CreateImageRequest) GetHeight() uint32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *CreateImageRequest) GetWidth() uint32 {
	if x != nil {
		return x.Width
	}
	return 0
}

type CreateImageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImageId string `protobuf:"bytes,1,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty"`
}

func (x *CreateImageResponse) Reset() {
	*x = CreateImageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_build_v1alpha1_build_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateImageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateImageResponse) ProtoMessage() {}

func (x *CreateImageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_build_v1alpha1_build_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateImageResponse.ProtoReflect.Descriptor instead.
func (*CreateImageResponse) Descriptor() ([]byte, []int) {
	return file_build_v1alpha1_build_proto_rawDescGZIP(), []int{1}
}

func (x *CreateImageResponse) GetImageId() string {
	if x != nil {
		return x.ImageId
	}
	return ""
}

var File_build_v1alpha1_build_proto protoreflect.FileDescriptor

var file_build_v1alpha1_build_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x63, 0x6f,
	0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x36, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x36, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x36, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xa2, 0x01, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x49, 0x64, 0x12, 0x41, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e,
	0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x36, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05,
	0x77, 0x69, 0x64, 0x74, 0x68, 0x22, 0x30, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x64, 0x32, 0xe1, 0x03, 0x0a, 0x06, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x12, 0x53, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2a, 0x2e, 0x63, 0x6f,
	0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x36, 0x2e, 0x50, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x58, 0x0a, 0x0a, 0x49, 0x6e, 0x69, 0x74, 0x50,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x32, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x36, 0x2e, 0x49,
	0x6e, 0x69, 0x74, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3d, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x65, 0x74, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x12, 0x71, 0x0a, 0x0c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x12, 0x33, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64,
	0x6f, 0x64, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x36, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65,
	0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x36, 0x2e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x30, 0x01, 0x12, 0x76, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x12, 0x32, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74,
	0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x36, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x61, 0x62,
	0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x36, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6d,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x31, 0x5a, 0x2f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65,
	0x74, 0x2f, 0x64, 0x6f, 0x64, 0x6f, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_build_v1alpha1_build_proto_rawDescOnce sync.Once
	file_build_v1alpha1_build_proto_rawDescData = file_build_v1alpha1_build_proto_rawDesc
)

func file_build_v1alpha1_build_proto_rawDescGZIP() []byte {
	file_build_v1alpha1_build_proto_rawDescOnce.Do(func() {
		file_build_v1alpha1_build_proto_rawDescData = protoimpl.X.CompressGZIP(file_build_v1alpha1_build_proto_rawDescData)
	})
	return file_build_v1alpha1_build_proto_rawDescData
}

var file_build_v1alpha1_build_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_build_v1alpha1_build_proto_goTypes = []interface{}{
	(*CreateImageRequest)(nil),           // 0: com.wabenet.dodo.core.v1alpha6.CreateImageRequest
	(*CreateImageResponse)(nil),          // 1: com.wabenet.dodo.core.v1alpha6.CreateImageResponse
	(*v1alpha6.BuildInfo)(nil),           // 2: com.wabenet.dodo.core.v1alpha6.BuildInfo
	(*emptypb.Empty)(nil),                // 3: google.protobuf.Empty
	(*v1alpha6.StreamOutputRequest)(nil), // 4: com.wabenet.dodo.core.v1alpha6.StreamOutputRequest
	(*v1alpha6.PluginInfo)(nil),          // 5: com.wabenet.dodo.core.v1alpha6.PluginInfo
	(*v1alpha6.InitPluginResponse)(nil),  // 6: com.wabenet.dodo.core.v1alpha6.InitPluginResponse
	(*v1alpha6.OutputData)(nil),          // 7: com.wabenet.dodo.core.v1alpha6.OutputData
}
var file_build_v1alpha1_build_proto_depIdxs = []int32{
	2, // 0: com.wabenet.dodo.core.v1alpha6.CreateImageRequest.config:type_name -> com.wabenet.dodo.core.v1alpha6.BuildInfo
	3, // 1: com.wabenet.dodo.core.v1alpha6.Plugin.GetPluginInfo:input_type -> google.protobuf.Empty
	3, // 2: com.wabenet.dodo.core.v1alpha6.Plugin.InitPlugin:input_type -> google.protobuf.Empty
	3, // 3: com.wabenet.dodo.core.v1alpha6.Plugin.ResetPlugin:input_type -> google.protobuf.Empty
	4, // 4: com.wabenet.dodo.core.v1alpha6.Plugin.StreamOutput:input_type -> com.wabenet.dodo.core.v1alpha6.StreamOutputRequest
	0, // 5: com.wabenet.dodo.core.v1alpha6.Plugin.CreateImage:input_type -> com.wabenet.dodo.core.v1alpha6.CreateImageRequest
	5, // 6: com.wabenet.dodo.core.v1alpha6.Plugin.GetPluginInfo:output_type -> com.wabenet.dodo.core.v1alpha6.PluginInfo
	6, // 7: com.wabenet.dodo.core.v1alpha6.Plugin.InitPlugin:output_type -> com.wabenet.dodo.core.v1alpha6.InitPluginResponse
	3, // 8: com.wabenet.dodo.core.v1alpha6.Plugin.ResetPlugin:output_type -> google.protobuf.Empty
	7, // 9: com.wabenet.dodo.core.v1alpha6.Plugin.StreamOutput:output_type -> com.wabenet.dodo.core.v1alpha6.OutputData
	1, // 10: com.wabenet.dodo.core.v1alpha6.Plugin.CreateImage:output_type -> com.wabenet.dodo.core.v1alpha6.CreateImageResponse
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_build_v1alpha1_build_proto_init() }
func file_build_v1alpha1_build_proto_init() {
	if File_build_v1alpha1_build_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_build_v1alpha1_build_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateImageRequest); i {
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
		file_build_v1alpha1_build_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateImageResponse); i {
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
			RawDescriptor: file_build_v1alpha1_build_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_build_v1alpha1_build_proto_goTypes,
		DependencyIndexes: file_build_v1alpha1_build_proto_depIdxs,
		MessageInfos:      file_build_v1alpha1_build_proto_msgTypes,
	}.Build()
	File_build_v1alpha1_build_proto = out.File
	file_build_v1alpha1_build_proto_rawDesc = nil
	file_build_v1alpha1_build_proto_goTypes = nil
	file_build_v1alpha1_build_proto_depIdxs = nil
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
	GetPluginInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1alpha6.PluginInfo, error)
	InitPlugin(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1alpha6.InitPluginResponse, error)
	ResetPlugin(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	StreamOutput(ctx context.Context, in *v1alpha6.StreamOutputRequest, opts ...grpc.CallOption) (Plugin_StreamOutputClient, error)
	CreateImage(ctx context.Context, in *CreateImageRequest, opts ...grpc.CallOption) (*CreateImageResponse, error)
}

type pluginClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginClient(cc grpc.ClientConnInterface) PluginClient {
	return &pluginClient{cc}
}

func (c *pluginClient) GetPluginInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1alpha6.PluginInfo, error) {
	out := new(v1alpha6.PluginInfo)
	err := c.cc.Invoke(ctx, "/com.wabenet.dodo.core.v1alpha6.Plugin/GetPluginInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) InitPlugin(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1alpha6.InitPluginResponse, error) {
	out := new(v1alpha6.InitPluginResponse)
	err := c.cc.Invoke(ctx, "/com.wabenet.dodo.core.v1alpha6.Plugin/InitPlugin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) ResetPlugin(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/com.wabenet.dodo.core.v1alpha6.Plugin/ResetPlugin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) StreamOutput(ctx context.Context, in *v1alpha6.StreamOutputRequest, opts ...grpc.CallOption) (Plugin_StreamOutputClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Plugin_serviceDesc.Streams[0], "/com.wabenet.dodo.core.v1alpha6.Plugin/StreamOutput", opts...)
	if err != nil {
		return nil, err
	}
	x := &pluginStreamOutputClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Plugin_StreamOutputClient interface {
	Recv() (*v1alpha6.OutputData, error)
	grpc.ClientStream
}

type pluginStreamOutputClient struct {
	grpc.ClientStream
}

func (x *pluginStreamOutputClient) Recv() (*v1alpha6.OutputData, error) {
	m := new(v1alpha6.OutputData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pluginClient) CreateImage(ctx context.Context, in *CreateImageRequest, opts ...grpc.CallOption) (*CreateImageResponse, error) {
	out := new(CreateImageResponse)
	err := c.cc.Invoke(ctx, "/com.wabenet.dodo.core.v1alpha6.Plugin/CreateImage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServer is the server API for Plugin service.
type PluginServer interface {
	GetPluginInfo(context.Context, *emptypb.Empty) (*v1alpha6.PluginInfo, error)
	InitPlugin(context.Context, *emptypb.Empty) (*v1alpha6.InitPluginResponse, error)
	ResetPlugin(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	StreamOutput(*v1alpha6.StreamOutputRequest, Plugin_StreamOutputServer) error
	CreateImage(context.Context, *CreateImageRequest) (*CreateImageResponse, error)
}

// UnimplementedPluginServer can be embedded to have forward compatible implementations.
type UnimplementedPluginServer struct {
}

func (*UnimplementedPluginServer) GetPluginInfo(context.Context, *emptypb.Empty) (*v1alpha6.PluginInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPluginInfo not implemented")
}
func (*UnimplementedPluginServer) InitPlugin(context.Context, *emptypb.Empty) (*v1alpha6.InitPluginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitPlugin not implemented")
}
func (*UnimplementedPluginServer) ResetPlugin(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetPlugin not implemented")
}
func (*UnimplementedPluginServer) StreamOutput(*v1alpha6.StreamOutputRequest, Plugin_StreamOutputServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamOutput not implemented")
}
func (*UnimplementedPluginServer) CreateImage(context.Context, *CreateImageRequest) (*CreateImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateImage not implemented")
}

func RegisterPluginServer(s *grpc.Server, srv PluginServer) {
	s.RegisterService(&_Plugin_serviceDesc, srv)
}

func _Plugin_GetPluginInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetPluginInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.wabenet.dodo.core.v1alpha6.Plugin/GetPluginInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetPluginInfo(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_InitPlugin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).InitPlugin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.wabenet.dodo.core.v1alpha6.Plugin/InitPlugin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).InitPlugin(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_ResetPlugin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).ResetPlugin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.wabenet.dodo.core.v1alpha6.Plugin/ResetPlugin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).ResetPlugin(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_StreamOutput_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(v1alpha6.StreamOutputRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PluginServer).StreamOutput(m, &pluginStreamOutputServer{stream})
}

type Plugin_StreamOutputServer interface {
	Send(*v1alpha6.OutputData) error
	grpc.ServerStream
}

type pluginStreamOutputServer struct {
	grpc.ServerStream
}

func (x *pluginStreamOutputServer) Send(m *v1alpha6.OutputData) error {
	return x.ServerStream.SendMsg(m)
}

func _Plugin_CreateImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).CreateImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.wabenet.dodo.core.v1alpha6.Plugin/CreateImage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).CreateImage(ctx, req.(*CreateImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Plugin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "com.wabenet.dodo.core.v1alpha6.Plugin",
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
			MethodName: "CreateImage",
			Handler:    _Plugin_CreateImage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamOutput",
			Handler:       _Plugin_StreamOutput_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "build/v1alpha1/build.proto",
}
