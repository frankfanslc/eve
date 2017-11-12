// Code generated by protoc-gen-go. DO NOT EDIT.
// source: devconfig.proto

package zconfig

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MapServer struct {
	NameOrIp   string `protobuf:"bytes,1,opt,name=NameOrIp,json=nameOrIp" json:"NameOrIp,omitempty"`
	Credential string `protobuf:"bytes,2,opt,name=Credential,json=credential" json:"Credential,omitempty"`
}

func (m *MapServer) Reset()                    { *m = MapServer{} }
func (m *MapServer) String() string            { return proto.CompactTextString(m) }
func (*MapServer) ProtoMessage()               {}
func (*MapServer) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *MapServer) GetNameOrIp() string {
	if m != nil {
		return m.NameOrIp
	}
	return ""
}

func (m *MapServer) GetCredential() string {
	if m != nil {
		return m.Credential
	}
	return ""
}

type ZedServer struct {
	HostName string   `protobuf:"bytes,1,opt,name=HostName,json=hostName" json:"HostName,omitempty"`
	EID      []string `protobuf:"bytes,2,rep,name=EID,json=eID" json:"EID,omitempty"`
}

func (m *ZedServer) Reset()                    { *m = ZedServer{} }
func (m *ZedServer) String() string            { return proto.CompactTextString(m) }
func (*ZedServer) ProtoMessage()               {}
func (*ZedServer) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *ZedServer) GetHostName() string {
	if m != nil {
		return m.HostName
	}
	return ""
}

func (m *ZedServer) GetEID() []string {
	if m != nil {
		return m.EID
	}
	return nil
}

type DeviceLispDetails struct {
	LispMapServers         []*MapServer `protobuf:"bytes,1,rep,name=LispMapServers,json=lispMapServers" json:"LispMapServers,omitempty"`
	LispInstance           uint32       `protobuf:"varint,2,opt,name=LispInstance,json=lispInstance" json:"LispInstance,omitempty"`
	EID                    string       `protobuf:"bytes,4,opt,name=EID,json=eID" json:"EID,omitempty"`
	EIDHashLen             int32        `protobuf:"varint,5,opt,name=EIDHashLen,json=eIDHashLen" json:"EIDHashLen,omitempty"`
	ZedServers             []*ZedServer `protobuf:"bytes,6,rep,name=ZedServers,json=zedServers" json:"ZedServers,omitempty"`
	EidAllocationPrefix    []byte       `protobuf:"bytes,8,opt,name=EidAllocationPrefix,json=eidAllocationPrefix,proto3" json:"EidAllocationPrefix,omitempty"`
	EidAllocationPrefixLen int32        `protobuf:"varint,9,opt,name=EidAllocationPrefixLen,json=eidAllocationPrefixLen" json:"EidAllocationPrefixLen,omitempty"`
	ClientAddr             string       `protobuf:"bytes,10,opt,name=ClientAddr,json=clientAddr" json:"ClientAddr,omitempty"`
}

func (m *DeviceLispDetails) Reset()                    { *m = DeviceLispDetails{} }
func (m *DeviceLispDetails) String() string            { return proto.CompactTextString(m) }
func (*DeviceLispDetails) ProtoMessage()               {}
func (*DeviceLispDetails) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *DeviceLispDetails) GetLispMapServers() []*MapServer {
	if m != nil {
		return m.LispMapServers
	}
	return nil
}

func (m *DeviceLispDetails) GetLispInstance() uint32 {
	if m != nil {
		return m.LispInstance
	}
	return 0
}

func (m *DeviceLispDetails) GetEID() string {
	if m != nil {
		return m.EID
	}
	return ""
}

func (m *DeviceLispDetails) GetEIDHashLen() int32 {
	if m != nil {
		return m.EIDHashLen
	}
	return 0
}

func (m *DeviceLispDetails) GetZedServers() []*ZedServer {
	if m != nil {
		return m.ZedServers
	}
	return nil
}

func (m *DeviceLispDetails) GetEidAllocationPrefix() []byte {
	if m != nil {
		return m.EidAllocationPrefix
	}
	return nil
}

func (m *DeviceLispDetails) GetEidAllocationPrefixLen() int32 {
	if m != nil {
		return m.EidAllocationPrefixLen
	}
	return 0
}

func (m *DeviceLispDetails) GetClientAddr() string {
	if m != nil {
		return m.ClientAddr
	}
	return ""
}

type EdgeDevConfig struct {
	Id                 *UUIDandVersion      `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	DevConfigSha256    []byte               `protobuf:"bytes,2,opt,name=devConfigSha256,proto3" json:"devConfigSha256,omitempty"`
	DevConfigSignature []byte               `protobuf:"bytes,3,opt,name=devConfigSignature,proto3" json:"devConfigSignature,omitempty"`
	Apps               []*AppInstanceConfig `protobuf:"bytes,4,rep,name=apps" json:"apps,omitempty"`
	Networks           []*NetworkConfig     `protobuf:"bytes,5,rep,name=networks" json:"networks,omitempty"`
	Datastores         []*DatastoreConfig   `protobuf:"bytes,6,rep,name=datastores" json:"datastores,omitempty"`
	LispInfo           *DeviceLispDetails   `protobuf:"bytes,7,opt,name=lispInfo" json:"lispInfo,omitempty"`
}

func (m *EdgeDevConfig) Reset()                    { *m = EdgeDevConfig{} }
func (m *EdgeDevConfig) String() string            { return proto.CompactTextString(m) }
func (*EdgeDevConfig) ProtoMessage()               {}
func (*EdgeDevConfig) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

func (m *EdgeDevConfig) GetId() *UUIDandVersion {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *EdgeDevConfig) GetDevConfigSha256() []byte {
	if m != nil {
		return m.DevConfigSha256
	}
	return nil
}

func (m *EdgeDevConfig) GetDevConfigSignature() []byte {
	if m != nil {
		return m.DevConfigSignature
	}
	return nil
}

func (m *EdgeDevConfig) GetApps() []*AppInstanceConfig {
	if m != nil {
		return m.Apps
	}
	return nil
}

func (m *EdgeDevConfig) GetNetworks() []*NetworkConfig {
	if m != nil {
		return m.Networks
	}
	return nil
}

func (m *EdgeDevConfig) GetDatastores() []*DatastoreConfig {
	if m != nil {
		return m.Datastores
	}
	return nil
}

func (m *EdgeDevConfig) GetLispInfo() *DeviceLispDetails {
	if m != nil {
		return m.LispInfo
	}
	return nil
}

func init() {
	proto.RegisterType((*MapServer)(nil), "MapServer")
	proto.RegisterType((*ZedServer)(nil), "ZedServer")
	proto.RegisterType((*DeviceLispDetails)(nil), "DeviceLispDetails")
	proto.RegisterType((*EdgeDevConfig)(nil), "EdgeDevConfig")
}

func init() { proto.RegisterFile("devconfig.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 523 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x93, 0xd1, 0x4e, 0xdb, 0x30,
	0x14, 0x86, 0xd5, 0x14, 0x58, 0x7b, 0x28, 0x85, 0x19, 0x09, 0x59, 0x48, 0x1b, 0x55, 0x2f, 0xa6,
	0x88, 0x0b, 0x17, 0x75, 0x1a, 0xd2, 0xee, 0xd6, 0x91, 0x6a, 0x44, 0x62, 0x6c, 0x0a, 0x62, 0x17,
	0xdc, 0x99, 0xf8, 0xb4, 0xb5, 0x96, 0xda, 0x91, 0xed, 0x76, 0x53, 0x9f, 0x6e, 0xef, 0xb0, 0x17,
	0x9a, 0xe2, 0xa4, 0x81, 0x8e, 0x5e, 0x9e, 0xef, 0xff, 0x6d, 0xff, 0xf6, 0x39, 0x86, 0x43, 0x81,
	0xcb, 0x54, 0xab, 0x89, 0x9c, 0xb2, 0xdc, 0x68, 0xa7, 0x4f, 0x4b, 0x30, 0x9f, 0x6b, 0xb5, 0x06,
	0x3c, 0xcf, 0x37, 0x1d, 0x0a, 0xdd, 0x06, 0x38, 0xb0, 0x4e, 0x1b, 0x3e, 0xc5, 0xb2, 0xec, 0x7f,
	0x81, 0xf6, 0x57, 0x9e, 0xdf, 0xa1, 0x59, 0xa2, 0x21, 0xa7, 0xd0, 0xba, 0xe5, 0x73, 0xfc, 0x66,
	0xe2, 0x9c, 0x36, 0x7a, 0x8d, 0xb0, 0x9d, 0xb4, 0x54, 0x55, 0x93, 0xb7, 0x00, 0x57, 0x06, 0x05,
	0x2a, 0x27, 0x79, 0x46, 0x03, 0xaf, 0x42, 0x5a, 0x93, 0xfe, 0x47, 0x68, 0x3f, 0xa0, 0x78, 0xda,
	0xe8, 0x5a, 0x5b, 0x57, 0x6c, 0xb6, 0xde, 0x68, 0x56, 0xd5, 0xe4, 0x08, 0x9a, 0xe3, 0x38, 0xa2,
	0x41, 0xaf, 0x19, 0xb6, 0x93, 0x26, 0xc6, 0x51, 0xff, 0x6f, 0x00, 0xaf, 0x23, 0x5c, 0xca, 0x14,
	0x6f, 0xa4, 0xcd, 0x23, 0x74, 0x5c, 0x66, 0x96, 0x0c, 0xa1, 0x5b, 0x94, 0x75, 0x3a, 0x4b, 0x1b,
	0xbd, 0x66, 0xb8, 0x3f, 0x04, 0x56, 0xa3, 0xa4, 0x9b, 0x6d, 0x38, 0x48, 0x1f, 0x3a, 0xc5, 0x9a,
	0x58, 0x59, 0xc7, 0x55, 0x8a, 0x3e, 0xe6, 0x41, 0xd2, 0xc9, 0x9e, 0xb1, 0xf5, 0xf9, 0x3b, 0x3e,
	0x56, 0x71, 0x7e, 0x71, 0xb5, 0x71, 0x1c, 0x5d, 0x73, 0x3b, 0xbb, 0x41, 0x45, 0x77, 0x7b, 0x8d,
	0x70, 0x37, 0x01, 0xac, 0x09, 0x39, 0x07, 0xa8, 0xaf, 0x66, 0xe9, 0x5e, 0x95, 0xa2, 0x46, 0x09,
	0xac, 0x6a, 0x95, 0x5c, 0xc0, 0xf1, 0x58, 0x8a, 0x51, 0x96, 0xe9, 0x94, 0x3b, 0xa9, 0xd5, 0x77,
	0x83, 0x13, 0xf9, 0x9b, 0xb6, 0x7a, 0x8d, 0xb0, 0x93, 0x1c, 0xe3, 0x4b, 0x89, 0x5c, 0xc2, 0xc9,
	0x96, 0x15, 0x45, 0x92, 0xb6, 0x4f, 0x72, 0x82, 0x5b, 0x55, 0xdf, 0x90, 0x4c, 0xa2, 0x72, 0x23,
	0x21, 0x0c, 0x85, 0xaa, 0x21, 0x35, 0xe9, 0xff, 0x09, 0xe0, 0x60, 0x2c, 0xa6, 0x18, 0xe1, 0xf2,
	0xca, 0x0f, 0x00, 0x39, 0x83, 0x40, 0x0a, 0xdf, 0x8f, 0xfd, 0xe1, 0x21, 0xbb, 0xbf, 0x8f, 0x23,
	0xae, 0xc4, 0x0f, 0x34, 0x56, 0x6a, 0x95, 0x04, 0x52, 0x90, 0xd0, 0x4f, 0x58, 0xe9, 0xbe, 0x9b,
	0xf1, 0xe1, 0x87, 0x4b, 0xff, 0x82, 0x9d, 0xe4, 0x7f, 0x4c, 0x18, 0x90, 0x27, 0x24, 0xa7, 0x8a,
	0xbb, 0x85, 0x41, 0xda, 0xf4, 0xe6, 0x2d, 0x0a, 0x79, 0x07, 0x3b, 0x3c, 0xcf, 0x2d, 0xdd, 0xf1,
	0x8f, 0x47, 0xd8, 0x28, 0xaf, 0x1b, 0x52, 0x5a, 0x13, 0xaf, 0x93, 0x73, 0x68, 0x29, 0x74, 0xbf,
	0xb4, 0xf9, 0x69, 0xe9, 0xae, 0xf7, 0x76, 0xd9, 0x6d, 0x09, 0x2a, 0x5f, 0xad, 0x93, 0x0b, 0x00,
	0xc1, 0x1d, 0x2f, 0xe6, 0x19, 0xd7, 0x6d, 0x39, 0x62, 0xd1, 0x1a, 0x55, 0xfe, 0x67, 0x1e, 0xc2,
	0xa0, 0x55, 0x8e, 0xc2, 0x44, 0xd3, 0x57, 0xfe, 0x19, 0x08, 0x7b, 0x31, 0x78, 0x49, 0xed, 0xf9,
	0xfc, 0x09, 0xce, 0x52, 0x3d, 0x67, 0x2b, 0x14, 0x28, 0x38, 0x4b, 0x33, 0xbd, 0x10, 0x6c, 0x61,
	0xd1, 0x14, 0x2b, 0xca, 0xff, 0xf3, 0xf0, 0x66, 0x2a, 0xdd, 0x6c, 0xf1, 0xc8, 0x52, 0x3d, 0x1f,
	0x94, 0xbe, 0x01, 0xcf, 0xe5, 0x60, 0x55, 0xfe, 0xb9, 0xc7, 0x3d, 0xef, 0x7a, 0xff, 0x2f, 0x00,
	0x00, 0xff, 0xff, 0x89, 0x9b, 0x6f, 0x2b, 0xba, 0x03, 0x00, 0x00,
}
