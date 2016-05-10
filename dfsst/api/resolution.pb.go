// Code generated by protoc-gen-go.
// source: dfss/dfsst/api/resolution.proto
// DO NOT EDIT!

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	dfss/dfsst/api/resolution.proto

It has these top-level messages:
	AlertRequest
	RecoverRequest
	TTPResponse
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import api2 "dfss/dfssc/api"

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

type AlertRequest struct {
	// / Promises obtained at this point of the main protocol
	Promises []*api2.Promise `protobuf:"bytes,1,rep,name=promises" json:"promises,omitempty"`
}

func (m *AlertRequest) Reset()                    { *m = AlertRequest{} }
func (m *AlertRequest) String() string            { return proto.CompactTextString(m) }
func (*AlertRequest) ProtoMessage()               {}
func (*AlertRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AlertRequest) GetPromises() []*api2.Promise {
	if m != nil {
		return m.Promises
	}
	return nil
}

type RecoverRequest struct {
	SignatureUUID string `protobuf:"bytes,1,opt,name=signatureUUID" json:"signatureUUID,omitempty"`
}

func (m *RecoverRequest) Reset()                    { *m = RecoverRequest{} }
func (m *RecoverRequest) String() string            { return proto.CompactTextString(m) }
func (*RecoverRequest) ProtoMessage()               {}
func (*RecoverRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type TTPResponse struct {
	// / True for abort token, False when the TTP was able to generate the fully signed contract
	Abort    bool   `protobuf:"varint,1,opt,name=abort" json:"abort,omitempty"`
	Contract []byte `protobuf:"bytes,2,opt,name=contract,proto3" json:"contract,omitempty"`
}

func (m *TTPResponse) Reset()                    { *m = TTPResponse{} }
func (m *TTPResponse) String() string            { return proto.CompactTextString(m) }
func (*TTPResponse) ProtoMessage()               {}
func (*TTPResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*AlertRequest)(nil), "api.AlertRequest")
	proto.RegisterType((*RecoverRequest)(nil), "api.RecoverRequest")
	proto.RegisterType((*TTPResponse)(nil), "api.TTPResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for TTP service

type TTPClient interface {
	// / Sent by a client when a signature encounters a problem.
	// Triggers the resolve protocol.
	Alert(ctx context.Context, in *AlertRequest, opts ...grpc.CallOption) (*TTPResponse, error)
	// / Sent by a client after a crash or a self-deconnection.
	// Tries to fetch the result of the resolve protocol, if any.
	Recover(ctx context.Context, in *RecoverRequest, opts ...grpc.CallOption) (*TTPResponse, error)
}

type tTPClient struct {
	cc *grpc.ClientConn
}

func NewTTPClient(cc *grpc.ClientConn) TTPClient {
	return &tTPClient{cc}
}

func (c *tTPClient) Alert(ctx context.Context, in *AlertRequest, opts ...grpc.CallOption) (*TTPResponse, error) {
	out := new(TTPResponse)
	err := grpc.Invoke(ctx, "/api.TTP/Alert", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tTPClient) Recover(ctx context.Context, in *RecoverRequest, opts ...grpc.CallOption) (*TTPResponse, error) {
	out := new(TTPResponse)
	err := grpc.Invoke(ctx, "/api.TTP/Recover", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for TTP service

type TTPServer interface {
	// / Sent by a client when a signature encounters a problem.
	// Triggers the resolve protocol.
	Alert(context.Context, *AlertRequest) (*TTPResponse, error)
	// / Sent by a client after a crash or a self-deconnection.
	// Tries to fetch the result of the resolve protocol, if any.
	Recover(context.Context, *RecoverRequest) (*TTPResponse, error)
}

func RegisterTTPServer(s *grpc.Server, srv TTPServer) {
	s.RegisterService(&_TTP_serviceDesc, srv)
}

func _TTP_Alert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TTPServer).Alert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.TTP/Alert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TTPServer).Alert(ctx, req.(*AlertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TTP_Recover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecoverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TTPServer).Recover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.TTP/Recover",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TTPServer).Recover(ctx, req.(*RecoverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TTP_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.TTP",
	HandlerType: (*TTPServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Alert",
			Handler:    _TTP_Alert_Handler,
		},
		{
			MethodName: "Recover",
			Handler:    _TTP_Recover_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 231 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x8f, 0xcd, 0x4a, 0xc3, 0x40,
	0x10, 0xc7, 0xad, 0xa1, 0x1a, 0xa7, 0xa9, 0xd4, 0x11, 0xa1, 0x54, 0x50, 0xc9, 0xc5, 0x9e, 0x36,
	0x10, 0x9f, 0x40, 0xf0, 0xe2, 0xad, 0x84, 0xf6, 0x01, 0xd2, 0x75, 0x94, 0x85, 0xb8, 0x13, 0x77,
	0x26, 0x3e, 0xbf, 0xeb, 0x56, 0xa5, 0x42, 0x2f, 0x7b, 0xf8, 0x7f, 0xed, 0x6f, 0xe0, 0xf6, 0xe5,
	0x55, 0xa4, 0xfa, 0x7e, 0xb4, 0x6a, 0x7b, 0x57, 0x05, 0x12, 0xee, 0x06, 0x75, 0xec, 0x4d, 0x1f,
	0x58, 0x19, 0xb3, 0xa8, 0x2e, 0xae, 0xff, 0x52, 0x36, 0xa5, 0x6c, 0xe7, 0xc8, 0xeb, 0x2e, 0x51,
	0x1a, 0x28, 0x1e, 0x3b, 0x0a, 0xda, 0xd0, 0xc7, 0x40, 0xa2, 0x78, 0x03, 0x79, 0x34, 0xde, 0x9d,
	0x90, 0xcc, 0x47, 0x77, 0xd9, 0x72, 0x52, 0x17, 0x26, 0x96, 0xcc, 0x6a, 0x27, 0x96, 0xf7, 0x70,
	0xde, 0x90, 0xe5, 0x4f, 0x0a, 0xbf, 0x8d, 0x2b, 0x98, 0x8a, 0x7b, 0xf3, 0xad, 0x0e, 0x81, 0x36,
	0x9b, 0xe7, 0xa7, 0x58, 0x1b, 0x2d, 0xcf, 0xe2, 0xf0, 0x64, 0xbd, 0x5e, 0x35, 0x24, 0x3d, 0x7b,
	0x21, 0x9c, 0xc2, 0xb8, 0xdd, 0x72, 0xd0, 0xe4, 0xe6, 0x38, 0x83, 0xdc, 0xb2, 0xd7, 0xd0, 0x5a,
	0x9d, 0x1f, 0x47, 0xa5, 0xa8, 0x1d, 0x64, 0x31, 0x8f, 0x06, 0xc6, 0x89, 0x07, 0x2f, 0xd2, 0xb7,
	0xfb, 0x6c, 0x8b, 0x59, 0x92, 0xf6, 0x56, 0xcb, 0x23, 0xac, 0xe1, 0xf4, 0x87, 0x07, 0x2f, 0x93,
	0xfd, 0x9f, 0xee, 0x50, 0x67, 0x7b, 0x92, 0x4e, 0x7f, 0xf8, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x38,
	0x4a, 0x1c, 0x12, 0x3f, 0x01, 0x00, 0x00,
}
