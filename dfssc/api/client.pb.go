// Code generated by protoc-gen-go.
// source: dfss/dfssc/api/client.proto
// DO NOT EDIT!

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	dfss/dfssc/api/client.proto

It has these top-level messages:
	Context
	Promise
	Signature
	Hello
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import api1 "dfss/dfssp/api"

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

type Context struct {
	RecipientKeyHash     []byte   `protobuf:"bytes,1,opt,name=recipientKeyHash,proto3" json:"recipientKeyHash,omitempty"`
	SenderKeyHash        []byte   `protobuf:"bytes,2,opt,name=senderKeyHash,proto3" json:"senderKeyHash,omitempty"`
	Sequence             []uint32 `protobuf:"varint,3,rep,name=sequence" json:"sequence,omitempty"`
	Signers              [][]byte `protobuf:"bytes,4,rep,name=signers,proto3" json:"signers,omitempty"`
	ContractDocumentHash []byte   `protobuf:"bytes,5,opt,name=contractDocumentHash,proto3" json:"contractDocumentHash,omitempty"`
	SignatureUUID        string   `protobuf:"bytes,6,opt,name=signatureUUID" json:"signatureUUID,omitempty"`
	SignedHash           []byte   `protobuf:"bytes,7,opt,name=signedHash,proto3" json:"signedHash,omitempty"`
}

func (m *Context) Reset()                    { *m = Context{} }
func (m *Context) String() string            { return proto.CompactTextString(m) }
func (*Context) ProtoMessage()               {}
func (*Context) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Promise message contains all the required information to verify
// the identity of the sender and receiver, and the actual promise
type Promise struct {
	Context *Context `protobuf:"bytes,1,opt,name=context" json:"context,omitempty"`
	Index   uint32   `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	Payload []byte   `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *Promise) Reset()                    { *m = Promise{} }
func (m *Promise) String() string            { return proto.CompactTextString(m) }
func (*Promise) ProtoMessage()               {}
func (*Promise) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Promise) GetContext() *Context {
	if m != nil {
		return m.Context
	}
	return nil
}

// Signature message contains all the required information to verify
// the identity of the sender and receiver, and the actual signature
type Signature struct {
	Context *Context `protobuf:"bytes,1,opt,name=context" json:"context,omitempty"`
	Payload []byte   `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *Signature) Reset()                    { *m = Signature{} }
func (m *Signature) String() string            { return proto.CompactTextString(m) }
func (*Signature) ProtoMessage()               {}
func (*Signature) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Signature) GetContext() *Context {
	if m != nil {
		return m.Context
	}
	return nil
}

// Hello message is used when discovering peers.
// It contains the current version of the software.
type Hello struct {
	Version string `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
}

func (m *Hello) Reset()                    { *m = Hello{} }
func (m *Hello) String() string            { return proto.CompactTextString(m) }
func (*Hello) ProtoMessage()               {}
func (*Hello) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*Context)(nil), "api.Context")
	proto.RegisterType((*Promise)(nil), "api.Promise")
	proto.RegisterType((*Signature)(nil), "api.Signature")
	proto.RegisterType((*Hello)(nil), "api.Hello")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for Client service

type ClientClient interface {
	TreatPromise(ctx context.Context, in *Promise, opts ...grpc.CallOption) (*api1.ErrorCode, error)
	TreatSignature(ctx context.Context, in *Signature, opts ...grpc.CallOption) (*api1.ErrorCode, error)
	Discover(ctx context.Context, in *Hello, opts ...grpc.CallOption) (*Hello, error)
}

type clientClient struct {
	cc *grpc.ClientConn
}

func NewClientClient(cc *grpc.ClientConn) ClientClient {
	return &clientClient{cc}
}

func (c *clientClient) TreatPromise(ctx context.Context, in *Promise, opts ...grpc.CallOption) (*api1.ErrorCode, error) {
	out := new(api1.ErrorCode)
	err := grpc.Invoke(ctx, "/api.Client/TreatPromise", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) TreatSignature(ctx context.Context, in *Signature, opts ...grpc.CallOption) (*api1.ErrorCode, error) {
	out := new(api1.ErrorCode)
	err := grpc.Invoke(ctx, "/api.Client/TreatSignature", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) Discover(ctx context.Context, in *Hello, opts ...grpc.CallOption) (*Hello, error) {
	out := new(Hello)
	err := grpc.Invoke(ctx, "/api.Client/Discover", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Client service

type ClientServer interface {
	TreatPromise(context.Context, *Promise) (*api1.ErrorCode, error)
	TreatSignature(context.Context, *Signature) (*api1.ErrorCode, error)
	Discover(context.Context, *Hello) (*Hello, error)
}

func RegisterClientServer(s *grpc.Server, srv ClientServer) {
	s.RegisterService(&_Client_serviceDesc, srv)
}

func _Client_TreatPromise_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Promise)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ClientServer).TreatPromise(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Client_TreatSignature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Signature)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ClientServer).TreatSignature(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Client_Discover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Hello)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(ClientServer).Discover(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Client_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Client",
	HandlerType: (*ClientServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TreatPromise",
			Handler:    _Client_TreatPromise_Handler,
		},
		{
			MethodName: "TreatSignature",
			Handler:    _Client_TreatSignature_Handler,
		},
		{
			MethodName: "Discover",
			Handler:    _Client_Discover_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 380 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x52, 0x5d, 0x4f, 0xc2, 0x30,
	0x14, 0xe5, 0x43, 0x18, 0x5c, 0x07, 0x31, 0x0d, 0x0f, 0xcb, 0x8c, 0x06, 0x17, 0x62, 0x88, 0x0f,
	0x23, 0xc1, 0x9f, 0x00, 0x26, 0x18, 0x63, 0x62, 0xa6, 0xfc, 0x80, 0xda, 0x15, 0x6d, 0x32, 0xd6,
	0xd9, 0x16, 0x03, 0xbf, 0xc1, 0x37, 0x7f, 0xb1, 0xdd, 0x1d, 0x43, 0x88, 0x3c, 0xf8, 0xb2, 0xec,
	0xdc, 0x7b, 0x7a, 0xce, 0xbd, 0xa7, 0x85, 0xf3, 0x78, 0xa1, 0xf5, 0x28, 0xff, 0xb0, 0x11, 0xcd,
	0xc4, 0x88, 0x25, 0x82, 0xa7, 0x26, 0xcc, 0x94, 0x34, 0x92, 0xd4, 0x6d, 0xc5, 0xbf, 0xd8, 0x31,
	0x32, 0x64, 0x64, 0x09, 0x35, 0x0b, 0xa9, 0x96, 0x05, 0x27, 0xf8, 0xaa, 0x81, 0x33, 0x91, 0xa9,
	0xe1, 0x6b, 0x43, 0x6e, 0xe0, 0x4c, 0x71, 0x26, 0xb2, 0x5c, 0xe2, 0x81, 0x6f, 0x66, 0x54, 0xbf,
	0x7b, 0xd5, 0x7e, 0x75, 0xe8, 0x46, 0x7f, 0xea, 0x64, 0x00, 0x1d, 0xcd, 0xd3, 0x98, 0xab, 0x92,
	0x58, 0x43, 0xe2, 0x61, 0x91, 0xf8, 0xd0, 0xd2, 0xfc, 0x63, 0xc5, 0x53, 0xc6, 0xbd, 0x7a, 0xbf,
	0x3e, 0xec, 0x44, 0x3b, 0x4c, 0x3c, 0x70, 0xb4, 0x78, 0x4b, 0xb9, 0xd2, 0xde, 0x89, 0x6d, 0xb9,
	0x51, 0x09, 0xc9, 0x18, 0x7a, 0xcc, 0x8e, 0xa4, 0x28, 0x33, 0x53, 0xc9, 0x56, 0x4b, 0x6b, 0x8b,
	0x16, 0x0d, 0xb4, 0x38, 0xda, 0xc3, 0x79, 0xec, 0x71, 0x6a, 0x56, 0x8a, 0xcf, 0xe7, 0xf7, 0x53,
	0xaf, 0x69, 0xc9, 0xed, 0xe8, 0xb0, 0x48, 0x2e, 0x01, 0xd0, 0x24, 0x46, 0x3d, 0x07, 0xf5, 0xf6,
	0x2a, 0x01, 0x05, 0xe7, 0x49, 0xc9, 0xa5, 0xd0, 0x9c, 0x5c, 0x83, 0xc3, 0x8a, 0x5c, 0x30, 0x83,
	0xd3, 0xb1, 0x1b, 0xda, 0xf8, 0xc2, 0x6d, 0x56, 0x51, 0xd9, 0x24, 0x3d, 0x68, 0x08, 0xbb, 0xf2,
	0x1a, 0x03, 0xe8, 0x44, 0x05, 0xc8, 0x97, 0xcb, 0xe8, 0x26, 0x91, 0x34, 0xb6, 0x7b, 0xe7, 0x2e,
	0x25, 0x0c, 0x1e, 0xa1, 0xfd, 0x5c, 0xce, 0xf4, 0x6f, 0x93, 0x3d, 0xb9, 0xda, 0xa1, 0xdc, 0x15,
	0x34, 0x66, 0x3c, 0x49, 0x64, 0x4e, 0xf9, 0xb4, 0xe1, 0x09, 0x99, 0xa2, 0x54, 0x3b, 0x2a, 0xe1,
	0xf8, 0xbb, 0x0a, 0xcd, 0x09, 0xbe, 0x0b, 0x12, 0x82, 0xfb, 0xa2, 0x38, 0x35, 0xe5, 0x92, 0x85,
	0xdd, 0x16, 0xf9, 0x5d, 0x44, 0x77, 0x4a, 0x49, 0x35, 0x91, 0x31, 0x0f, 0x2a, 0xf6, 0x26, 0xba,
	0xc8, 0xff, 0x9d, 0xb8, 0xe0, 0xec, 0xf0, 0x91, 0x33, 0x03, 0x68, 0x4d, 0x85, 0x66, 0xd2, 0xda,
	0x13, 0xc0, 0x2e, 0x0e, 0xe8, 0xef, 0xfd, 0x07, 0x95, 0xd7, 0x26, 0x3e, 0xbf, 0xdb, 0x9f, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xb4, 0xf1, 0xe7, 0x80, 0xc1, 0x02, 0x00, 0x00,
}
