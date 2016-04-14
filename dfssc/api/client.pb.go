// Code generated by protoc-gen-go.
// source: dfss/dfssc/api/client.proto
// DO NOT EDIT!

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	dfss/dfssc/api/client.proto

It has these top-level messages:
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

// Promise message contains all the required information to verify
// the identity of the sender and receiver, and the actual promise
//
// * sequence is transmitted by platform and identical across clients
// * TODO implement an global signature for content
type Promise struct {
	RecipientKeyHash     []byte `protobuf:"bytes,1,opt,name=recipientKeyHash,proto3" json:"recipientKeyHash,omitempty"`
	SenderKeyHash        []byte `protobuf:"bytes,2,opt,name=senderKeyHash,proto3" json:"senderKeyHash,omitempty"`
	Index                uint32 `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
	ContractDocumentHash string `protobuf:"bytes,4,opt,name=contractDocumentHash" json:"contractDocumentHash,omitempty"`
	SignatureUuid        string `protobuf:"bytes,5,opt,name=signatureUuid" json:"signatureUuid,omitempty"`
	ContractUuid         string `protobuf:"bytes,6,opt,name=contractUuid" json:"contractUuid,omitempty"`
}

func (m *Promise) Reset()                    { *m = Promise{} }
func (m *Promise) String() string            { return proto.CompactTextString(m) }
func (*Promise) ProtoMessage()               {}
func (*Promise) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Signature message contains all the required information to verify
// the identity of the sender and receiver, and the actual signature
type Signature struct {
	RecipientKeyHash []byte `protobuf:"bytes,1,opt,name=recipientKeyHash,proto3" json:"recipientKeyHash,omitempty"`
	SenderKeyHash    []byte `protobuf:"bytes,2,opt,name=senderKeyHash,proto3" json:"senderKeyHash,omitempty"`
	Signature        string `protobuf:"bytes,3,opt,name=signature" json:"signature,omitempty"`
	SignatureUuid    string `protobuf:"bytes,4,opt,name=signatureUuid" json:"signatureUuid,omitempty"`
	ContractUuid     string `protobuf:"bytes,5,opt,name=contractUuid" json:"contractUuid,omitempty"`
}

func (m *Signature) Reset()                    { *m = Signature{} }
func (m *Signature) String() string            { return proto.CompactTextString(m) }
func (*Signature) ProtoMessage()               {}
func (*Signature) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// Hello message is used when discovering peers.
// It contains the current version of the software.
type Hello struct {
	Version string `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
}

func (m *Hello) Reset()                    { *m = Hello{} }
func (m *Hello) String() string            { return proto.CompactTextString(m) }
func (*Hello) ProtoMessage()               {}
func (*Hello) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*Promise)(nil), "api.Promise")
	proto.RegisterType((*Signature)(nil), "api.Signature")
	proto.RegisterType((*Hello)(nil), "api.Hello")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion1

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
	// 303 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x51, 0xdd, 0x4a, 0xf4, 0x30,
	0x10, 0xdd, 0x7e, 0xfb, 0xf7, 0x75, 0x68, 0x57, 0x0d, 0x2b, 0x94, 0xaa, 0x20, 0xc5, 0x0b, 0xaf,
	0x5a, 0x58, 0x1f, 0x61, 0x57, 0x58, 0xf0, 0x46, 0x50, 0x1f, 0x20, 0xa6, 0x59, 0x0d, 0xb4, 0x4d,
	0x99, 0xa4, 0xa2, 0xe0, 0x53, 0x88, 0x0f, 0x6c, 0x3a, 0x6b, 0xab, 0x17, 0xbd, 0xf1, 0x26, 0xe4,
	0x9c, 0x39, 0x33, 0xe7, 0x64, 0x02, 0x27, 0xf9, 0xce, 0x98, 0xac, 0x3d, 0x44, 0xc6, 0x6b, 0x95,
	0x89, 0x42, 0xc9, 0xca, 0xa6, 0x35, 0x6a, 0xab, 0xd9, 0xd8, 0x31, 0xf1, 0x59, 0xaf, 0xa8, 0x49,
	0x51, 0x17, 0xdc, 0xee, 0x34, 0x96, 0x7b, 0x4d, 0xf2, 0xe9, 0xc1, 0xfc, 0x16, 0x75, 0xa9, 0x8c,
	0x64, 0x11, 0x1c, 0xa2, 0x14, 0xaa, 0x6e, 0x47, 0xdc, 0xc8, 0xb7, 0x2d, 0x37, 0xcf, 0x91, 0x77,
	0xee, 0x5d, 0x06, 0xec, 0x18, 0x42, 0x23, 0xab, 0x5c, 0x62, 0x47, 0xff, 0x23, 0x3a, 0x84, 0xa9,
	0x72, 0xec, 0x6b, 0x34, 0x76, 0x30, 0x64, 0xa7, 0xb0, 0x14, 0xba, 0xb2, 0xc8, 0x85, 0xdd, 0x68,
	0xd1, 0x94, 0x6e, 0x0c, 0x89, 0x27, 0xae, 0xea, 0xd3, 0x0c, 0xf5, 0x54, 0x71, 0xdb, 0xa0, 0x7c,
	0x68, 0x54, 0x1e, 0x4d, 0x89, 0x5e, 0x42, 0xd0, 0x35, 0x11, 0x3b, 0x6b, 0xd9, 0xe4, 0x1d, 0xfc,
	0xbb, 0x4e, 0xfc, 0xf7, 0x5c, 0x47, 0xe0, 0xf7, 0x56, 0x94, 0x6d, 0xc0, 0x7d, 0x32, 0xe8, 0x4e,
	0x99, 0x92, 0x08, 0xa6, 0x5b, 0x59, 0x14, 0x9a, 0x1d, 0xc0, 0xfc, 0x45, 0xa2, 0x51, 0xba, 0x22,
	0x43, 0x7f, 0xf5, 0xe1, 0xc1, 0x6c, 0x4d, 0x3b, 0x66, 0x29, 0x04, 0xf7, 0x28, 0xb9, 0xed, 0xb6,
	0x17, 0xa4, 0x6e, 0xbd, 0xe9, 0x37, 0x8a, 0x17, 0x84, 0xae, 0x11, 0x35, 0xae, 0x75, 0x2e, 0x93,
	0x11, 0x5b, 0xc1, 0x82, 0xf4, 0x3f, 0xef, 0xda, 0x6b, 0x7a, 0x3c, 0xd0, 0x73, 0x01, 0xff, 0x37,
	0xca, 0x08, 0xed, 0x42, 0x30, 0xa0, 0x2a, 0xe5, 0x8a, 0x7f, 0xdd, 0x93, 0xd1, 0xe3, 0x8c, 0xbe,
	0xf2, 0xea, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x2d, 0x9f, 0xc1, 0xc3, 0x0d, 0x02, 0x00, 0x00,
}
