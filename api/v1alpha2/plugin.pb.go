// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.2
// source: api/v1alpha2/plugin.proto

package v1alpha2

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

type PluginName struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *PluginName) Reset() {
	*x = PluginName{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1alpha2_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginName) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginName) ProtoMessage() {}

func (x *PluginName) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1alpha2_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginName.ProtoReflect.Descriptor instead.
func (*PluginName) Descriptor() ([]byte, []int) {
	return file_api_v1alpha2_plugin_proto_rawDescGZIP(), []int{0}
}

func (x *PluginName) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *PluginName) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type PluginInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         *PluginName       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Dependencies []*PluginName     `protobuf:"bytes,2,rep,name=dependencies,proto3" json:"dependencies,omitempty"`
	Fields       map[string]string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *PluginInfo) Reset() {
	*x = PluginInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1alpha2_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginInfo) ProtoMessage() {}

func (x *PluginInfo) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1alpha2_plugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginInfo.ProtoReflect.Descriptor instead.
func (*PluginInfo) Descriptor() ([]byte, []int) {
	return file_api_v1alpha2_plugin_proto_rawDescGZIP(), []int{1}
}

func (x *PluginInfo) GetName() *PluginName {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *PluginInfo) GetDependencies() []*PluginName {
	if x != nil {
		return x.Dependencies
	}
	return nil
}

func (x *PluginInfo) GetFields() map[string]string {
	if x != nil {
		return x.Fields
	}
	return nil
}

type InitPluginResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config map[string]string `protobuf:"bytes,1,rep,name=config,proto3" json:"config,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *InitPluginResponse) Reset() {
	*x = InitPluginResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1alpha2_plugin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitPluginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitPluginResponse) ProtoMessage() {}

func (x *InitPluginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1alpha2_plugin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitPluginResponse.ProtoReflect.Descriptor instead.
func (*InitPluginResponse) Descriptor() ([]byte, []int) {
	return file_api_v1alpha2_plugin_proto_rawDescGZIP(), []int{2}
}

func (x *InitPluginResponse) GetConfig() map[string]string {
	if x != nil {
		return x.Config
	}
	return nil
}

type GetStreamingConnectionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetStreamingConnectionRequest) Reset() {
	*x = GetStreamingConnectionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1alpha2_plugin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStreamingConnectionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStreamingConnectionRequest) ProtoMessage() {}

func (x *GetStreamingConnectionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1alpha2_plugin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStreamingConnectionRequest.ProtoReflect.Descriptor instead.
func (*GetStreamingConnectionRequest) Descriptor() ([]byte, []int) {
	return file_api_v1alpha2_plugin_proto_rawDescGZIP(), []int{3}
}

type GetStreamingConnectionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *GetStreamingConnectionResponse) Reset() {
	*x = GetStreamingConnectionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1alpha2_plugin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStreamingConnectionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStreamingConnectionResponse) ProtoMessage() {}

func (x *GetStreamingConnectionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1alpha2_plugin_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStreamingConnectionResponse.ProtoReflect.Descriptor instead.
func (*GetStreamingConnectionResponse) Descriptor() ([]byte, []int) {
	return file_api_v1alpha2_plugin_proto_rawDescGZIP(), []int{4}
}

func (x *GetStreamingConnectionResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_api_v1alpha2_plugin_proto protoreflect.FileDescriptor

var file_api_v1alpha2_plugin_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x2f, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x2a, 0x63, 0x6f, 0x6d,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x5f, 0x63, 0x6c, 0x69,
	0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x5f, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x22, 0x34, 0x0a, 0x0a, 0x50, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0xcb, 0x02,
	0x0a, 0x0a, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x4a, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x36, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x5f, 0x63, 0x6c, 0x69,
	0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x5f, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x2e, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x4e, 0x61,
	0x6d, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x5a, 0x0a, 0x0c, 0x64, 0x65, 0x70, 0x65,
	0x6e, 0x64, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x36,
	0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x64, 0x6f, 0x64, 0x6f,
	0x5f, 0x63, 0x6c, 0x69, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x5f, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x2e, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x0c, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e,
	0x63, 0x69, 0x65, 0x73, 0x12, 0x5a, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x42, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x5f, 0x63, 0x6c, 0x69, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x5f,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x32, 0x2e, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73,
	0x1a, 0x39, 0x0a, 0x0b, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xb3, 0x01, 0x0a, 0x12,
	0x49, 0x6e, 0x69, 0x74, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x62, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x4a, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x64, 0x6f, 0x64, 0x6f, 0x5f, 0x63, 0x6c, 0x69, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x5f, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x2e,
	0x49, 0x6e, 0x69, 0x74, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x39, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x1f, 0x0a, 0x1d, 0x47, 0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e,
	0x67, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x32, 0x0a, 0x1e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69,
	0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x64, 0x6f, 0x2d, 0x63, 0x6c, 0x69, 0x2f, 0x64, 0x6f,
	0x64, 0x6f, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1alpha2_plugin_proto_rawDescOnce sync.Once
	file_api_v1alpha2_plugin_proto_rawDescData = file_api_v1alpha2_plugin_proto_rawDesc
)

func file_api_v1alpha2_plugin_proto_rawDescGZIP() []byte {
	file_api_v1alpha2_plugin_proto_rawDescOnce.Do(func() {
		file_api_v1alpha2_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1alpha2_plugin_proto_rawDescData)
	})
	return file_api_v1alpha2_plugin_proto_rawDescData
}

var file_api_v1alpha2_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_v1alpha2_plugin_proto_goTypes = []interface{}{
	(*PluginName)(nil),                     // 0: com.github.dodo_cli.dodo_core.api.v1alpha2.PluginName
	(*PluginInfo)(nil),                     // 1: com.github.dodo_cli.dodo_core.api.v1alpha2.PluginInfo
	(*InitPluginResponse)(nil),             // 2: com.github.dodo_cli.dodo_core.api.v1alpha2.InitPluginResponse
	(*GetStreamingConnectionRequest)(nil),  // 3: com.github.dodo_cli.dodo_core.api.v1alpha2.GetStreamingConnectionRequest
	(*GetStreamingConnectionResponse)(nil), // 4: com.github.dodo_cli.dodo_core.api.v1alpha2.GetStreamingConnectionResponse
	nil,                                    // 5: com.github.dodo_cli.dodo_core.api.v1alpha2.PluginInfo.FieldsEntry
	nil,                                    // 6: com.github.dodo_cli.dodo_core.api.v1alpha2.InitPluginResponse.ConfigEntry
}
var file_api_v1alpha2_plugin_proto_depIdxs = []int32{
	0, // 0: com.github.dodo_cli.dodo_core.api.v1alpha2.PluginInfo.name:type_name -> com.github.dodo_cli.dodo_core.api.v1alpha2.PluginName
	0, // 1: com.github.dodo_cli.dodo_core.api.v1alpha2.PluginInfo.dependencies:type_name -> com.github.dodo_cli.dodo_core.api.v1alpha2.PluginName
	5, // 2: com.github.dodo_cli.dodo_core.api.v1alpha2.PluginInfo.fields:type_name -> com.github.dodo_cli.dodo_core.api.v1alpha2.PluginInfo.FieldsEntry
	6, // 3: com.github.dodo_cli.dodo_core.api.v1alpha2.InitPluginResponse.config:type_name -> com.github.dodo_cli.dodo_core.api.v1alpha2.InitPluginResponse.ConfigEntry
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_v1alpha2_plugin_proto_init() }
func file_api_v1alpha2_plugin_proto_init() {
	if File_api_v1alpha2_plugin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1alpha2_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PluginName); i {
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
		file_api_v1alpha2_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PluginInfo); i {
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
		file_api_v1alpha2_plugin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitPluginResponse); i {
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
		file_api_v1alpha2_plugin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStreamingConnectionRequest); i {
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
		file_api_v1alpha2_plugin_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStreamingConnectionResponse); i {
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
			RawDescriptor: file_api_v1alpha2_plugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1alpha2_plugin_proto_goTypes,
		DependencyIndexes: file_api_v1alpha2_plugin_proto_depIdxs,
		MessageInfos:      file_api_v1alpha2_plugin_proto_msgTypes,
	}.Build()
	File_api_v1alpha2_plugin_proto = out.File
	file_api_v1alpha2_plugin_proto_rawDesc = nil
	file_api_v1alpha2_plugin_proto_goTypes = nil
	file_api_v1alpha2_plugin_proto_depIdxs = nil
}