// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.21.12
// source: core/v1alpha6/build.proto

package v1alpha6

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

type BuildInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Builder          string           `protobuf:"bytes,1,opt,name=builder,proto3" json:"builder,omitempty"`
	ImageName        string           `protobuf:"bytes,2,opt,name=image_name,json=imageName,proto3" json:"image_name,omitempty"`
	Context          string           `protobuf:"bytes,3,opt,name=context,proto3" json:"context,omitempty"`
	Dockerfile       string           `protobuf:"bytes,4,opt,name=dockerfile,proto3" json:"dockerfile,omitempty"`
	InlineDockerfile []string         `protobuf:"bytes,5,rep,name=inline_dockerfile,json=inlineDockerfile,proto3" json:"inline_dockerfile,omitempty"`
	Arguments        []*BuildArgument `protobuf:"bytes,6,rep,name=arguments,proto3" json:"arguments,omitempty"`
	Secrets          []*BuildSecret   `protobuf:"bytes,7,rep,name=secrets,proto3" json:"secrets,omitempty"`
	SshAgents        []*SshAgent      `protobuf:"bytes,8,rep,name=ssh_agents,json=sshAgents,proto3" json:"ssh_agents,omitempty"`
	NoCache          bool             `protobuf:"varint,9,opt,name=no_cache,json=noCache,proto3" json:"no_cache,omitempty"`
	ForceRebuild     bool             `protobuf:"varint,10,opt,name=force_rebuild,json=forceRebuild,proto3" json:"force_rebuild,omitempty"`
	ForcePull        bool             `protobuf:"varint,11,opt,name=force_pull,json=forcePull,proto3" json:"force_pull,omitempty"`
	Dependencies     []string         `protobuf:"bytes,12,rep,name=dependencies,proto3" json:"dependencies,omitempty"`
}

func (x *BuildInfo) Reset() {
	*x = BuildInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_v1alpha6_build_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildInfo) ProtoMessage() {}

func (x *BuildInfo) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1alpha6_build_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildInfo.ProtoReflect.Descriptor instead.
func (*BuildInfo) Descriptor() ([]byte, []int) {
	return file_core_v1alpha6_build_proto_rawDescGZIP(), []int{0}
}

func (x *BuildInfo) GetBuilder() string {
	if x != nil {
		return x.Builder
	}
	return ""
}

func (x *BuildInfo) GetImageName() string {
	if x != nil {
		return x.ImageName
	}
	return ""
}

func (x *BuildInfo) GetContext() string {
	if x != nil {
		return x.Context
	}
	return ""
}

func (x *BuildInfo) GetDockerfile() string {
	if x != nil {
		return x.Dockerfile
	}
	return ""
}

func (x *BuildInfo) GetInlineDockerfile() []string {
	if x != nil {
		return x.InlineDockerfile
	}
	return nil
}

func (x *BuildInfo) GetArguments() []*BuildArgument {
	if x != nil {
		return x.Arguments
	}
	return nil
}

func (x *BuildInfo) GetSecrets() []*BuildSecret {
	if x != nil {
		return x.Secrets
	}
	return nil
}

func (x *BuildInfo) GetSshAgents() []*SshAgent {
	if x != nil {
		return x.SshAgents
	}
	return nil
}

func (x *BuildInfo) GetNoCache() bool {
	if x != nil {
		return x.NoCache
	}
	return false
}

func (x *BuildInfo) GetForceRebuild() bool {
	if x != nil {
		return x.ForceRebuild
	}
	return false
}

func (x *BuildInfo) GetForcePull() bool {
	if x != nil {
		return x.ForcePull
	}
	return false
}

func (x *BuildInfo) GetDependencies() []string {
	if x != nil {
		return x.Dependencies
	}
	return nil
}

type BuildArgument struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *BuildArgument) Reset() {
	*x = BuildArgument{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_v1alpha6_build_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildArgument) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildArgument) ProtoMessage() {}

func (x *BuildArgument) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1alpha6_build_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildArgument.ProtoReflect.Descriptor instead.
func (*BuildArgument) Descriptor() ([]byte, []int) {
	return file_core_v1alpha6_build_proto_rawDescGZIP(), []int{1}
}

func (x *BuildArgument) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *BuildArgument) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type BuildSecret struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Path string `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *BuildSecret) Reset() {
	*x = BuildSecret{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_v1alpha6_build_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildSecret) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildSecret) ProtoMessage() {}

func (x *BuildSecret) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1alpha6_build_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildSecret.ProtoReflect.Descriptor instead.
func (*BuildSecret) Descriptor() ([]byte, []int) {
	return file_core_v1alpha6_build_proto_rawDescGZIP(), []int{2}
}

func (x *BuildSecret) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BuildSecret) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type SshAgent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	IdentityFile string `protobuf:"bytes,2,opt,name=identity_file,json=identityFile,proto3" json:"identity_file,omitempty"`
}

func (x *SshAgent) Reset() {
	*x = SshAgent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_v1alpha6_build_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SshAgent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SshAgent) ProtoMessage() {}

func (x *SshAgent) ProtoReflect() protoreflect.Message {
	mi := &file_core_v1alpha6_build_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SshAgent.ProtoReflect.Descriptor instead.
func (*SshAgent) Descriptor() ([]byte, []int) {
	return file_core_v1alpha6_build_proto_rawDescGZIP(), []int{3}
}

func (x *SshAgent) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SshAgent) GetIdentityFile() string {
	if x != nil {
		return x.IdentityFile
	}
	return ""
}

var File_core_v1alpha6_build_proto protoreflect.FileDescriptor

var file_core_v1alpha6_build_proto_rawDesc = []byte{
	0x0a, 0x19, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x36, 0x2f,
	0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x63, 0x6f, 0x6d,
	0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x36, 0x22, 0x8b, 0x04, 0x0a, 0x09,
	0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x75, 0x69,
	0x6c, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x75, 0x69, 0x6c,
	0x64, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x1e, 0x0a, 0x0a,
	0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x2b, 0x0a, 0x11,
	0x69, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x66, 0x69, 0x6c,
	0x65, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x10, 0x69, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x44,
	0x6f, 0x63, 0x6b, 0x65, 0x72, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x4b, 0x0a, 0x09, 0x61, 0x72, 0x67,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x36, 0x2e, 0x42, 0x75,
	0x69, 0x6c, 0x64, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x61, 0x72, 0x67,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x45, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x61,
	0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e, 0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x36, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x53, 0x65,
	0x63, 0x72, 0x65, 0x74, 0x52, 0x07, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x12, 0x47, 0x0a,
	0x0a, 0x73, 0x73, 0x68, 0x5f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x28, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2e,
	0x64, 0x6f, 0x64, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x36, 0x2e, 0x53, 0x73, 0x68, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x73, 0x73, 0x68,
	0x41, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x6e, 0x6f, 0x5f, 0x63, 0x61, 0x63,
	0x68, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x6e, 0x6f, 0x43, 0x61, 0x63, 0x68,
	0x65, 0x12, 0x23, 0x0a, 0x0d, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x5f, 0x72, 0x65, 0x62, 0x75, 0x69,
	0x6c, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x52,
	0x65, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x5f,
	0x70, 0x75, 0x6c, 0x6c, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x66, 0x6f, 0x72, 0x63,
	0x65, 0x50, 0x75, 0x6c, 0x6c, 0x12, 0x22, 0x0a, 0x0c, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65,
	0x6e, 0x63, 0x69, 0x65, 0x73, 0x18, 0x0c, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x65, 0x70,
	0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x22, 0x37, 0x0a, 0x0d, 0x42, 0x75, 0x69,
	0x6c, 0x64, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x31, 0x0a, 0x0b, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x53, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x3f, 0x0a, 0x08, 0x53, 0x73, 0x68, 0x41, 0x67, 0x65, 0x6e,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x23, 0x0a, 0x0d, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x66, 0x69,
	0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x61, 0x62, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x64, 0x6f, 0x64,
	0x6f, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x36, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_core_v1alpha6_build_proto_rawDescOnce sync.Once
	file_core_v1alpha6_build_proto_rawDescData = file_core_v1alpha6_build_proto_rawDesc
)

func file_core_v1alpha6_build_proto_rawDescGZIP() []byte {
	file_core_v1alpha6_build_proto_rawDescOnce.Do(func() {
		file_core_v1alpha6_build_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_v1alpha6_build_proto_rawDescData)
	})
	return file_core_v1alpha6_build_proto_rawDescData
}

var file_core_v1alpha6_build_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_core_v1alpha6_build_proto_goTypes = []interface{}{
	(*BuildInfo)(nil),     // 0: com.wabenet.dodo.core.v1alpha6.BuildInfo
	(*BuildArgument)(nil), // 1: com.wabenet.dodo.core.v1alpha6.BuildArgument
	(*BuildSecret)(nil),   // 2: com.wabenet.dodo.core.v1alpha6.BuildSecret
	(*SshAgent)(nil),      // 3: com.wabenet.dodo.core.v1alpha6.SshAgent
}
var file_core_v1alpha6_build_proto_depIdxs = []int32{
	1, // 0: com.wabenet.dodo.core.v1alpha6.BuildInfo.arguments:type_name -> com.wabenet.dodo.core.v1alpha6.BuildArgument
	2, // 1: com.wabenet.dodo.core.v1alpha6.BuildInfo.secrets:type_name -> com.wabenet.dodo.core.v1alpha6.BuildSecret
	3, // 2: com.wabenet.dodo.core.v1alpha6.BuildInfo.ssh_agents:type_name -> com.wabenet.dodo.core.v1alpha6.SshAgent
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_core_v1alpha6_build_proto_init() }
func file_core_v1alpha6_build_proto_init() {
	if File_core_v1alpha6_build_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_core_v1alpha6_build_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildInfo); i {
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
		file_core_v1alpha6_build_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildArgument); i {
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
		file_core_v1alpha6_build_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildSecret); i {
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
		file_core_v1alpha6_build_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SshAgent); i {
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
			RawDescriptor: file_core_v1alpha6_build_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_v1alpha6_build_proto_goTypes,
		DependencyIndexes: file_core_v1alpha6_build_proto_depIdxs,
		MessageInfos:      file_core_v1alpha6_build_proto_msgTypes,
	}.Build()
	File_core_v1alpha6_build_proto = out.File
	file_core_v1alpha6_build_proto_rawDesc = nil
	file_core_v1alpha6_build_proto_goTypes = nil
	file_core_v1alpha6_build_proto_depIdxs = nil
}
