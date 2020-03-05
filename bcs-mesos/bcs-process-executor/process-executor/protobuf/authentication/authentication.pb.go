// Code generated by protoc-gen-go.
// source: mesos/authentication/authentication.proto
// DO NOT EDIT!

/*
Package mesos_internal is a generated protocol buffer package.

It is generated from these files:
	mesos/authentication/authentication.proto

It has these top-level messages:
	AuthenticateMessage
	AuthenticationMechanismsMessage
	AuthenticationStartMessage
	AuthenticationStepMessage
	AuthenticationCompletedMessage
	AuthenticationFailedMessage
	AuthenticationErrorMessage
*/
package authentication

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type AuthenticateMessage struct {
	Pid              *string `protobuf:"bytes,1,req,name=pid" json:"pid,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AuthenticateMessage) Reset()                    { *m = AuthenticateMessage{} }
func (m *AuthenticateMessage) String() string            { return proto.CompactTextString(m) }
func (*AuthenticateMessage) ProtoMessage()               {}
func (*AuthenticateMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AuthenticateMessage) GetPid() string {
	if m != nil && m.Pid != nil {
		return *m.Pid
	}
	return ""
}

type AuthenticationMechanismsMessage struct {
	Mechanisms       []string `protobuf:"bytes,1,rep,name=mechanisms" json:"mechanisms,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *AuthenticationMechanismsMessage) Reset()         { *m = AuthenticationMechanismsMessage{} }
func (m *AuthenticationMechanismsMessage) String() string { return proto.CompactTextString(m) }
func (*AuthenticationMechanismsMessage) ProtoMessage()    {}
func (*AuthenticationMechanismsMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{1}
}

func (m *AuthenticationMechanismsMessage) GetMechanisms() []string {
	if m != nil {
		return m.Mechanisms
	}
	return nil
}

type AuthenticationStartMessage struct {
	Mechanism        *string `protobuf:"bytes,1,req,name=mechanism" json:"mechanism,omitempty"`
	Data             []byte  `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AuthenticationStartMessage) Reset()                    { *m = AuthenticationStartMessage{} }
func (m *AuthenticationStartMessage) String() string            { return proto.CompactTextString(m) }
func (*AuthenticationStartMessage) ProtoMessage()               {}
func (*AuthenticationStartMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *AuthenticationStartMessage) GetMechanism() string {
	if m != nil && m.Mechanism != nil {
		return *m.Mechanism
	}
	return ""
}

func (m *AuthenticationStartMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type AuthenticationStepMessage struct {
	Data             []byte `protobuf:"bytes,1,req,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *AuthenticationStepMessage) Reset()                    { *m = AuthenticationStepMessage{} }
func (m *AuthenticationStepMessage) String() string            { return proto.CompactTextString(m) }
func (*AuthenticationStepMessage) ProtoMessage()               {}
func (*AuthenticationStepMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AuthenticationStepMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type AuthenticationCompletedMessage struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *AuthenticationCompletedMessage) Reset()                    { *m = AuthenticationCompletedMessage{} }
func (m *AuthenticationCompletedMessage) String() string            { return proto.CompactTextString(m) }
func (*AuthenticationCompletedMessage) ProtoMessage()               {}
func (*AuthenticationCompletedMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type AuthenticationFailedMessage struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *AuthenticationFailedMessage) Reset()                    { *m = AuthenticationFailedMessage{} }
func (m *AuthenticationFailedMessage) String() string            { return proto.CompactTextString(m) }
func (*AuthenticationFailedMessage) ProtoMessage()               {}
func (*AuthenticationFailedMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type AuthenticationErrorMessage struct {
	Error            *string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AuthenticationErrorMessage) Reset()                    { *m = AuthenticationErrorMessage{} }
func (m *AuthenticationErrorMessage) String() string            { return proto.CompactTextString(m) }
func (*AuthenticationErrorMessage) ProtoMessage()               {}
func (*AuthenticationErrorMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *AuthenticationErrorMessage) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*AuthenticateMessage)(nil), "mesos.internal.AuthenticateMessage")
	proto.RegisterType((*AuthenticationMechanismsMessage)(nil), "mesos.internal.AuthenticationMechanismsMessage")
	proto.RegisterType((*AuthenticationStartMessage)(nil), "mesos.internal.AuthenticationStartMessage")
	proto.RegisterType((*AuthenticationStepMessage)(nil), "mesos.internal.AuthenticationStepMessage")
	proto.RegisterType((*AuthenticationCompletedMessage)(nil), "mesos.internal.AuthenticationCompletedMessage")
	proto.RegisterType((*AuthenticationFailedMessage)(nil), "mesos.internal.AuthenticationFailedMessage")
	proto.RegisterType((*AuthenticationErrorMessage)(nil), "mesos.internal.AuthenticationErrorMessage")
}

func init() { proto.RegisterFile("mesos/authentication/authentication.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x91, 0x51, 0x4b, 0xc3, 0x30,
	0x14, 0x85, 0xe9, 0xa6, 0x42, 0x2f, 0x22, 0xa3, 0xfa, 0x50, 0xa7, 0xce, 0x90, 0x17, 0xeb, 0x4b,
	0x07, 0xfe, 0x83, 0x4e, 0xf4, 0x6d, 0x22, 0xf5, 0x17, 0x5c, 0xda, 0xcb, 0x1a, 0x68, 0x93, 0x90,
	0x5c, 0xff, 0xbf, 0xb4, 0x2e, 0xeb, 0xda, 0xb7, 0x7b, 0x4e, 0xce, 0x77, 0x2e, 0x49, 0xe0, 0xb5,
	0x23, 0x6f, 0xfc, 0x16, 0x7f, 0xb9, 0x21, 0xcd, 0xaa, 0x42, 0x56, 0x46, 0xcf, 0x64, 0x6e, 0x9d,
	0x61, 0x93, 0xdc, 0x0c, 0xd1, 0x5c, 0x69, 0x26, 0xa7, 0xb1, 0x95, 0x2f, 0x70, 0x5b, 0x8c, 0x39,
	0xda, 0x93, 0xf7, 0x78, 0xa0, 0x64, 0x05, 0x4b, 0xab, 0xea, 0x34, 0x12, 0x8b, 0x2c, 0x2e, 0xfb,
	0x51, 0x16, 0xf0, 0x5c, 0x4c, 0x0a, 0xf7, 0x54, 0x35, 0xa8, 0x95, 0xef, 0x7c, 0x80, 0x36, 0x00,
	0xdd, 0xc9, 0x4c, 0x23, 0xb1, 0xcc, 0xe2, 0xf2, 0xcc, 0x91, 0x5f, 0xb0, 0x9e, 0x56, 0xfc, 0x30,
	0x3a, 0x0e, 0xf4, 0x23, 0xc4, 0xa7, 0xec, 0x71, 0xf1, 0x68, 0x24, 0x09, 0x5c, 0xd4, 0xc8, 0x98,
	0x2e, 0x44, 0x94, 0x5d, 0x97, 0xc3, 0x2c, 0xb7, 0x70, 0x3f, 0xef, 0x23, 0x1b, 0xea, 0x02, 0xd0,
	0x37, 0x05, 0x40, 0xc0, 0x66, 0x0a, 0xbc, 0x9b, 0xce, 0xb6, 0xc4, 0x54, 0x1f, 0x29, 0xf9, 0x04,
	0x0f, 0xd3, 0xc4, 0x27, 0xaa, 0x76, 0x3c, 0x7e, 0x9b, 0xdf, 0xe0, 0xc3, 0x39, 0xe3, 0xc2, 0xca,
	0x3b, 0xb8, 0xa4, 0x5e, 0xa7, 0x91, 0x88, 0xb2, 0xb8, 0xfc, 0x17, 0xbb, 0x35, 0xac, 0x8c, 0x3b,
	0xe4, 0x68, 0xb1, 0x6a, 0x28, 0x1f, 0x9e, 0x7f, 0x77, 0xf5, 0xdd, 0x7f, 0x86, 0xff, 0x0b, 0x00,
	0x00, 0xff, 0xff, 0x06, 0xda, 0x2b, 0x92, 0xb9, 0x01, 0x00, 0x00,
}
