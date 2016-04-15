// Code generated by protoc-gen-go.
// source: dfss/dfssd/api/demonstrator.proto
// DO NOT EDIT!

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	dfss/dfssd/api/demonstrator.proto

It has these top-level messages:
	Log
	Ack
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

// Log message to display information
type Log struct {
	Timestamp  int64  `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
	Identifier string `protobuf:"bytes,2,opt,name=identifier" json:"identifier,omitempty"`
	Log        string `protobuf:"bytes,3,opt,name=log" json:"log,omitempty"`
}

func (m *Log) Reset()                    { *m = Log{} }
func (m *Log) String() string            { return proto.CompactTextString(m) }
func (*Log) ProtoMessage()               {}
func (*Log) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Empty ack message
type Ack struct {
}

func (m *Ack) Reset()                    { *m = Ack{} }
func (m *Ack) String() string            { return proto.CompactTextString(m) }
func (*Ack) ProtoMessage()               {}
func (*Ack) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*Log)(nil), "api.Log")
	proto.RegisterType((*Ack)(nil), "api.Ack")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion1

// Client API for Demonstrator service

type DemonstratorClient interface {
	// Log message.
	//
	// Send the UnixNano timetamp, sender's identifier and log message
	// Returns nothing ?
	SendLog(ctx context.Context, in *Log, opts ...grpc.CallOption) (*Ack, error)
}

type demonstratorClient struct {
	cc *grpc.ClientConn
}

func NewDemonstratorClient(cc *grpc.ClientConn) DemonstratorClient {
	return &demonstratorClient{cc}
}

func (c *demonstratorClient) SendLog(ctx context.Context, in *Log, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/api.Demonstrator/SendLog", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Demonstrator service

type DemonstratorServer interface {
	// Log message.
	//
	// Send the UnixNano timetamp, sender's identifier and log message
	// Returns nothing ?
	SendLog(context.Context, *Log) (*Ack, error)
}

func RegisterDemonstratorServer(s *grpc.Server, srv DemonstratorServer) {
	s.RegisterService(&_Demonstrator_serviceDesc, srv)
}

func _Demonstrator_SendLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Log)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(DemonstratorServer).SendLog(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Demonstrator_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Demonstrator",
	HandlerType: (*DemonstratorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendLog",
			Handler:    _Demonstrator_SendLog_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 167 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x52, 0x4c, 0x49, 0x2b, 0x2e,
	0xd6, 0x07, 0x11, 0x29, 0xfa, 0x89, 0x05, 0x99, 0xfa, 0x29, 0xa9, 0xb9, 0xf9, 0x79, 0xc5, 0x25,
	0x45, 0x89, 0x25, 0xf9, 0x45, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0xcc, 0x40, 0x71, 0xa5,
	0x50, 0x2e, 0x66, 0x9f, 0xfc, 0x74, 0x21, 0x19, 0x2e, 0xce, 0x92, 0xcc, 0xdc, 0xd4, 0xe2, 0x92,
	0xc4, 0xdc, 0x02, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xe6, 0x20, 0x84, 0x80, 0x90, 0x1c, 0x17, 0x57,
	0x66, 0x4a, 0x6a, 0x5e, 0x49, 0x66, 0x5a, 0x66, 0x6a, 0x91, 0x04, 0x13, 0x50, 0x9a, 0x33, 0x08,
	0x49, 0x44, 0x48, 0x80, 0x8b, 0x39, 0x27, 0x3f, 0x5d, 0x82, 0x19, 0x2c, 0x01, 0x62, 0x2a, 0xb1,
	0x72, 0x31, 0x3b, 0x26, 0x67, 0x1b, 0xe9, 0x73, 0xf1, 0xb8, 0x20, 0x59, 0x2c, 0x24, 0xcf, 0xc5,
	0x1e, 0x9c, 0x9a, 0x97, 0x02, 0xb2, 0x91, 0x43, 0x0f, 0x68, 0xbd, 0x1e, 0x90, 0x25, 0x05, 0x61,
	0x01, 0x95, 0x2b, 0x31, 0x24, 0xb1, 0x81, 0x9d, 0x66, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xc3,
	0xc3, 0x86, 0x7b, 0xbf, 0x00, 0x00, 0x00,
}
