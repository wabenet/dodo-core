// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/types/types.proto

package types

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Backdrop struct {
	Name                 string         `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Aliases              []string       `protobuf:"bytes,2,rep,name=aliases,proto3" json:"aliases,omitempty"`
	ContainerName        string         `protobuf:"bytes,3,opt,name=container_name,json=containerName,proto3" json:"container_name,omitempty"`
	ImageId              string         `protobuf:"bytes,4,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty"`
	Entrypoint           *Entrypoint    `protobuf:"bytes,5,opt,name=entrypoint,proto3" json:"entrypoint,omitempty"`
	Environment          []*Environment `protobuf:"bytes,6,rep,name=environment,proto3" json:"environment,omitempty"`
	Volumes              []*Volume      `protobuf:"bytes,7,rep,name=volumes,proto3" json:"volumes,omitempty"`
	Devices              []*Device      `protobuf:"bytes,8,rep,name=devices,proto3" json:"devices,omitempty"`
	Ports                []*Port        `protobuf:"bytes,9,rep,name=ports,proto3" json:"ports,omitempty"`
	User                 string         `protobuf:"bytes,10,opt,name=user,proto3" json:"user,omitempty"`
	WorkingDir           string         `protobuf:"bytes,11,opt,name=working_dir,json=workingDir,proto3" json:"working_dir,omitempty"`
	Attach               bool           `protobuf:"varint,12,opt,name=attach,proto3" json:"attach,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Backdrop) Reset()         { *m = Backdrop{} }
func (m *Backdrop) String() string { return proto.CompactTextString(m) }
func (*Backdrop) ProtoMessage()    {}
func (*Backdrop) Descriptor() ([]byte, []int) {
	return fileDescriptor_6579d14f41ea6320, []int{0}
}

func (m *Backdrop) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Backdrop.Unmarshal(m, b)
}
func (m *Backdrop) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Backdrop.Marshal(b, m, deterministic)
}
func (m *Backdrop) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backdrop.Merge(m, src)
}
func (m *Backdrop) XXX_Size() int {
	return xxx_messageInfo_Backdrop.Size(m)
}
func (m *Backdrop) XXX_DiscardUnknown() {
	xxx_messageInfo_Backdrop.DiscardUnknown(m)
}

var xxx_messageInfo_Backdrop proto.InternalMessageInfo

func (m *Backdrop) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Backdrop) GetAliases() []string {
	if m != nil {
		return m.Aliases
	}
	return nil
}

func (m *Backdrop) GetContainerName() string {
	if m != nil {
		return m.ContainerName
	}
	return ""
}

func (m *Backdrop) GetImageId() string {
	if m != nil {
		return m.ImageId
	}
	return ""
}

func (m *Backdrop) GetEntrypoint() *Entrypoint {
	if m != nil {
		return m.Entrypoint
	}
	return nil
}

func (m *Backdrop) GetEnvironment() []*Environment {
	if m != nil {
		return m.Environment
	}
	return nil
}

func (m *Backdrop) GetVolumes() []*Volume {
	if m != nil {
		return m.Volumes
	}
	return nil
}

func (m *Backdrop) GetDevices() []*Device {
	if m != nil {
		return m.Devices
	}
	return nil
}

func (m *Backdrop) GetPorts() []*Port {
	if m != nil {
		return m.Ports
	}
	return nil
}

func (m *Backdrop) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *Backdrop) GetWorkingDir() string {
	if m != nil {
		return m.WorkingDir
	}
	return ""
}

func (m *Backdrop) GetAttach() bool {
	if m != nil {
		return m.Attach
	}
	return false
}

type Entrypoint struct {
	Interactive          bool     `protobuf:"varint,1,opt,name=interactive,proto3" json:"interactive,omitempty"`
	Script               string   `protobuf:"bytes,2,opt,name=script,proto3" json:"script,omitempty"`
	Interpreter          []string `protobuf:"bytes,3,rep,name=interpreter,proto3" json:"interpreter,omitempty"`
	Arguments            []string `protobuf:"bytes,4,rep,name=arguments,proto3" json:"arguments,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Entrypoint) Reset()         { *m = Entrypoint{} }
func (m *Entrypoint) String() string { return proto.CompactTextString(m) }
func (*Entrypoint) ProtoMessage()    {}
func (*Entrypoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_6579d14f41ea6320, []int{1}
}

func (m *Entrypoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Entrypoint.Unmarshal(m, b)
}
func (m *Entrypoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Entrypoint.Marshal(b, m, deterministic)
}
func (m *Entrypoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entrypoint.Merge(m, src)
}
func (m *Entrypoint) XXX_Size() int {
	return xxx_messageInfo_Entrypoint.Size(m)
}
func (m *Entrypoint) XXX_DiscardUnknown() {
	xxx_messageInfo_Entrypoint.DiscardUnknown(m)
}

var xxx_messageInfo_Entrypoint proto.InternalMessageInfo

func (m *Entrypoint) GetInteractive() bool {
	if m != nil {
		return m.Interactive
	}
	return false
}

func (m *Entrypoint) GetScript() string {
	if m != nil {
		return m.Script
	}
	return ""
}

func (m *Entrypoint) GetInterpreter() []string {
	if m != nil {
		return m.Interpreter
	}
	return nil
}

func (m *Entrypoint) GetArguments() []string {
	if m != nil {
		return m.Arguments
	}
	return nil
}

type Environment struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Environment) Reset()         { *m = Environment{} }
func (m *Environment) String() string { return proto.CompactTextString(m) }
func (*Environment) ProtoMessage()    {}
func (*Environment) Descriptor() ([]byte, []int) {
	return fileDescriptor_6579d14f41ea6320, []int{2}
}

func (m *Environment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Environment.Unmarshal(m, b)
}
func (m *Environment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Environment.Marshal(b, m, deterministic)
}
func (m *Environment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Environment.Merge(m, src)
}
func (m *Environment) XXX_Size() int {
	return xxx_messageInfo_Environment.Size(m)
}
func (m *Environment) XXX_DiscardUnknown() {
	xxx_messageInfo_Environment.DiscardUnknown(m)
}

var xxx_messageInfo_Environment proto.InternalMessageInfo

func (m *Environment) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Environment) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Volume struct {
	Source               string   `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Target               string   `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	Readonly             bool     `protobuf:"varint,3,opt,name=readonly,proto3" json:"readonly,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Volume) Reset()         { *m = Volume{} }
func (m *Volume) String() string { return proto.CompactTextString(m) }
func (*Volume) ProtoMessage()    {}
func (*Volume) Descriptor() ([]byte, []int) {
	return fileDescriptor_6579d14f41ea6320, []int{3}
}

func (m *Volume) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Volume.Unmarshal(m, b)
}
func (m *Volume) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Volume.Marshal(b, m, deterministic)
}
func (m *Volume) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Volume.Merge(m, src)
}
func (m *Volume) XXX_Size() int {
	return xxx_messageInfo_Volume.Size(m)
}
func (m *Volume) XXX_DiscardUnknown() {
	xxx_messageInfo_Volume.DiscardUnknown(m)
}

var xxx_messageInfo_Volume proto.InternalMessageInfo

func (m *Volume) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Volume) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

func (m *Volume) GetReadonly() bool {
	if m != nil {
		return m.Readonly
	}
	return false
}

type Device struct {
	CgroupRule           string   `protobuf:"bytes,1,opt,name=cgroup_rule,json=cgroupRule,proto3" json:"cgroup_rule,omitempty"`
	Source               string   `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Target               string   `protobuf:"bytes,3,opt,name=target,proto3" json:"target,omitempty"`
	Permissions          string   `protobuf:"bytes,4,opt,name=permissions,proto3" json:"permissions,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Device) Reset()         { *m = Device{} }
func (m *Device) String() string { return proto.CompactTextString(m) }
func (*Device) ProtoMessage()    {}
func (*Device) Descriptor() ([]byte, []int) {
	return fileDescriptor_6579d14f41ea6320, []int{4}
}

func (m *Device) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Device.Unmarshal(m, b)
}
func (m *Device) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Device.Marshal(b, m, deterministic)
}
func (m *Device) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Device.Merge(m, src)
}
func (m *Device) XXX_Size() int {
	return xxx_messageInfo_Device.Size(m)
}
func (m *Device) XXX_DiscardUnknown() {
	xxx_messageInfo_Device.DiscardUnknown(m)
}

var xxx_messageInfo_Device proto.InternalMessageInfo

func (m *Device) GetCgroupRule() string {
	if m != nil {
		return m.CgroupRule
	}
	return ""
}

func (m *Device) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Device) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

func (m *Device) GetPermissions() string {
	if m != nil {
		return m.Permissions
	}
	return ""
}

type Port struct {
	Target               string   `protobuf:"bytes,1,opt,name=target,proto3" json:"target,omitempty"`
	Published            string   `protobuf:"bytes,2,opt,name=published,proto3" json:"published,omitempty"`
	Protocol             string   `protobuf:"bytes,3,opt,name=protocol,proto3" json:"protocol,omitempty"`
	HostIp               string   `protobuf:"bytes,4,opt,name=host_ip,json=hostIp,proto3" json:"host_ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Port) Reset()         { *m = Port{} }
func (m *Port) String() string { return proto.CompactTextString(m) }
func (*Port) ProtoMessage()    {}
func (*Port) Descriptor() ([]byte, []int) {
	return fileDescriptor_6579d14f41ea6320, []int{5}
}

func (m *Port) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Port.Unmarshal(m, b)
}
func (m *Port) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Port.Marshal(b, m, deterministic)
}
func (m *Port) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Port.Merge(m, src)
}
func (m *Port) XXX_Size() int {
	return xxx_messageInfo_Port.Size(m)
}
func (m *Port) XXX_DiscardUnknown() {
	xxx_messageInfo_Port.DiscardUnknown(m)
}

var xxx_messageInfo_Port proto.InternalMessageInfo

func (m *Port) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

func (m *Port) GetPublished() string {
	if m != nil {
		return m.Published
	}
	return ""
}

func (m *Port) GetProtocol() string {
	if m != nil {
		return m.Protocol
	}
	return ""
}

func (m *Port) GetHostIp() string {
	if m != nil {
		return m.HostIp
	}
	return ""
}

func init() {
	proto.RegisterType((*Backdrop)(nil), "types.Backdrop")
	proto.RegisterType((*Entrypoint)(nil), "types.Entrypoint")
	proto.RegisterType((*Environment)(nil), "types.Environment")
	proto.RegisterType((*Volume)(nil), "types.Volume")
	proto.RegisterType((*Device)(nil), "types.Device")
	proto.RegisterType((*Port)(nil), "types.Port")
}

func init() { proto.RegisterFile("pkg/types/types.proto", fileDescriptor_6579d14f41ea6320) }

var fileDescriptor_6579d14f41ea6320 = []byte{
	// 514 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x53, 0x4d, 0x8f, 0xd3, 0x30,
	0x10, 0x55, 0x36, 0x6d, 0x9a, 0x4e, 0x58, 0x04, 0x16, 0x1f, 0x06, 0xad, 0x44, 0x88, 0x84, 0xe8,
	0x69, 0x11, 0x0b, 0xfc, 0x01, 0xb4, 0x1c, 0xf6, 0x82, 0x90, 0x85, 0xb8, 0x56, 0xde, 0xc4, 0x4a,
	0xad, 0xa6, 0xb6, 0x19, 0x3b, 0x45, 0x15, 0x77, 0xfe, 0x09, 0xff, 0x13, 0xd9, 0x71, 0xd3, 0xf4,
	0xc0, 0xa5, 0xf2, 0x7b, 0xf3, 0x66, 0x9e, 0xfb, 0x3c, 0x81, 0xa7, 0x66, 0xdb, 0xbe, 0x73, 0x07,
	0x23, 0xec, 0xf0, 0x7b, 0x6d, 0x50, 0x3b, 0x4d, 0xe6, 0x01, 0x54, 0x7f, 0x53, 0xc8, 0x3f, 0xf3,
	0x7a, 0xdb, 0xa0, 0x36, 0x84, 0xc0, 0x4c, 0xf1, 0x9d, 0xa0, 0x49, 0x99, 0xac, 0x96, 0x2c, 0x9c,
	0x09, 0x85, 0x05, 0xef, 0x24, 0xb7, 0xc2, 0xd2, 0x8b, 0x32, 0x5d, 0x2d, 0xd9, 0x11, 0x92, 0x37,
	0xf0, 0xb0, 0xd6, 0xca, 0x71, 0xa9, 0x04, 0xae, 0x43, 0x5f, 0x1a, 0xfa, 0x2e, 0x47, 0xf6, 0xab,
	0x1f, 0xf0, 0x02, 0x72, 0xb9, 0xe3, 0xad, 0x58, 0xcb, 0x86, 0xce, 0x82, 0x60, 0x11, 0xf0, 0x5d,
	0x43, 0xde, 0x03, 0x08, 0xe5, 0xf0, 0x60, 0xb4, 0x54, 0x8e, 0xce, 0xcb, 0x64, 0x55, 0xdc, 0x3c,
	0xbe, 0x1e, 0x6e, 0xf9, 0x65, 0x2c, 0xb0, 0x89, 0x88, 0x7c, 0x84, 0x42, 0xa8, 0xbd, 0x44, 0xad,
	0x76, 0x42, 0x39, 0x9a, 0x95, 0xe9, 0xaa, 0xb8, 0x21, 0x63, 0xcf, 0x58, 0x61, 0x53, 0x19, 0x79,
	0x0b, 0x8b, 0xbd, 0xee, 0xfa, 0x9d, 0xb0, 0x74, 0x11, 0x3a, 0x2e, 0x63, 0xc7, 0x8f, 0xc0, 0xb2,
	0x63, 0xd5, 0x0b, 0x1b, 0xb1, 0x97, 0xb5, 0xb0, 0x34, 0x3f, 0x13, 0xde, 0x06, 0x96, 0x1d, 0xab,
	0xe4, 0x35, 0xcc, 0x8d, 0x46, 0x67, 0xe9, 0x32, 0xc8, 0x8a, 0x28, 0xfb, 0xa6, 0xd1, 0xb1, 0xa1,
	0xe2, 0xd3, 0xec, 0xad, 0x40, 0x0a, 0x43, 0x9a, 0xfe, 0x4c, 0x5e, 0x41, 0xf1, 0x4b, 0xe3, 0x56,
	0xaa, 0x76, 0xdd, 0x48, 0xa4, 0x45, 0x28, 0x41, 0xa4, 0x6e, 0x25, 0x92, 0x67, 0x90, 0x71, 0xe7,
	0x78, 0xbd, 0xa1, 0x0f, 0xca, 0x64, 0x95, 0xb3, 0x88, 0xaa, 0x3f, 0x09, 0xc0, 0x29, 0x12, 0x52,
	0x42, 0x21, 0x95, 0x13, 0xc8, 0x6b, 0x27, 0xf7, 0xc3, 0x83, 0xe5, 0x6c, 0x4a, 0xf9, 0x41, 0xb6,
	0x46, 0x69, 0x1c, 0xbd, 0x08, 0x26, 0x11, 0x8d, 0x9d, 0x06, 0x85, 0x13, 0x48, 0xd3, 0xf0, 0xa6,
	0x53, 0x8a, 0x5c, 0xc1, 0x92, 0x63, 0xdb, 0xfb, 0xe0, 0x2c, 0x9d, 0x85, 0xfa, 0x89, 0xa8, 0x3e,
	0x41, 0x31, 0x89, 0x99, 0x3c, 0x82, 0x74, 0x2b, 0x0e, 0x71, 0x63, 0xfc, 0x91, 0x3c, 0x81, 0xf9,
	0x9e, 0x77, 0xbd, 0x88, 0xbe, 0x03, 0xa8, 0xbe, 0x43, 0x36, 0x64, 0x1d, 0x2e, 0xa6, 0x7b, 0xac,
	0x8f, 0x6b, 0x16, 0x91, 0xe7, 0x1d, 0xc7, 0x56, 0x8c, 0x17, 0x1e, 0x10, 0x79, 0x09, 0x39, 0x0a,
	0xde, 0x68, 0xd5, 0x1d, 0xc2, 0x82, 0xe5, 0x6c, 0xc4, 0xd5, 0x6f, 0xc8, 0x86, 0x87, 0xf1, 0xc1,
	0xd6, 0x2d, 0xea, 0xde, 0xac, 0xb1, 0xef, 0x8e, 0xa3, 0x61, 0xa0, 0x58, 0xdf, 0x4d, 0x6d, 0x2f,
	0xfe, 0x63, 0x9b, 0x9e, 0xd9, 0x96, 0x50, 0x18, 0x81, 0x3b, 0x69, 0xad, 0xd4, 0xca, 0xc6, 0xcd,
	0x9d, 0x52, 0xd5, 0x4f, 0x98, 0xf9, 0xe7, 0x9e, 0x4c, 0x48, 0xce, 0x26, 0x5c, 0xc1, 0xd2, 0xf4,
	0xf7, 0x9d, 0xb4, 0x1b, 0xd1, 0x44, 0xd3, 0x13, 0xe1, 0xff, 0x56, 0xf8, 0x10, 0x6b, 0xdd, 0x45,
	0xe7, 0x11, 0x93, 0xe7, 0xb0, 0xd8, 0x68, 0xeb, 0xd6, 0xd2, 0x44, 0xdf, 0xcc, 0xc3, 0x3b, 0x73,
	0x9f, 0x05, 0xc9, 0x87, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x01, 0x0d, 0x4b, 0x63, 0xd4, 0x03,
	0x00, 0x00,
}
