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

// Alert message sent by a signer
type AlertRequest struct {
}

func (m *AlertRequest) Reset()                    { *m = AlertRequest{} }
func (m *AlertRequest) String() string            { return proto.CompactTextString(m) }
func (*AlertRequest) ProtoMessage()               {}
func (*AlertRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// RecoverRequest message to get a signed contract
type RecoverRequest struct {
}

func (m *RecoverRequest) Reset()                    { *m = RecoverRequest{} }
func (m *RecoverRequest) String() string            { return proto.CompactTextString(m) }
func (*RecoverRequest) ProtoMessage()               {}
func (*RecoverRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// TTPResponse message to answer an Alert or Recover message
type TTPResponse struct {
	// Yes for abort token, No for signed contract
	Abort bool `protobuf:"varint,1,opt,name=abort" json:"abort,omitempty"`
	// Nil for abort token, non-empty otherwise
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
const _ = grpc.SupportPackageIsVersion1

// Client API for TTP service

type TTPClient interface {
	Alert(ctx context.Context, in *AlertRequest, opts ...grpc.CallOption) (*TTPResponse, error)
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
	Alert(context.Context, *AlertRequest) (*TTPResponse, error)
	Recover(context.Context, *RecoverRequest) (*TTPResponse, error)
}

func RegisterTTPServer(s *grpc.Server, srv TTPServer) {
	s.RegisterService(&_TTP_serviceDesc, srv)
}

func _TTP_Alert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(AlertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(TTPServer).Alert(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _TTP_Recover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(RecoverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(TTPServer).Recover(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
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
	// 181 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x8f, 0x41, 0xae, 0x82, 0x40,
	0x0c, 0x86, 0x1f, 0x8f, 0xa0, 0xa4, 0x22, 0xc1, 0xba, 0x21, 0x6e, 0x34, 0xac, 0x5c, 0x0d, 0x09,
	0x9e, 0xc0, 0x1b, 0x18, 0xc2, 0x05, 0x00, 0x6b, 0x42, 0x42, 0x28, 0xce, 0x14, 0xcf, 0xef, 0x38,
	0xd1, 0x04, 0x13, 0x37, 0x5d, 0x7c, 0x6d, 0xff, 0xaf, 0x85, 0xfd, 0xf5, 0x66, 0x4c, 0xfe, 0x2a,
	0x92, 0xd7, 0x63, 0x97, 0x6b, 0x32, 0xdc, 0x4f, 0xd2, 0xf1, 0xa0, 0x46, 0xcd, 0xc2, 0xe8, 0x5b,
	0x9a, 0xc5, 0x10, 0x9d, 0x7b, 0xd2, 0x52, 0xd2, 0x7d, 0x22, 0x23, 0x59, 0x02, 0x71, 0x49, 0x2d,
	0x3f, 0x48, 0x7f, 0x88, 0x82, 0x55, 0x55, 0x5d, 0x4a, 0x32, 0x23, 0x0f, 0x86, 0x70, 0x0d, 0x41,
	0xdd, 0xb0, 0x96, 0xd4, 0x3b, 0x78, 0xc7, 0x10, 0x13, 0x08, 0x5b, 0x1e, 0x44, 0xd7, 0xad, 0xa4,
	0xff, 0x96, 0x44, 0x45, 0x07, 0xbe, 0x9d, 0x47, 0x05, 0x81, 0x0b, 0xc6, 0x8d, 0xb2, 0x1e, 0x35,
	0x97, 0xec, 0x12, 0x87, 0x66, 0xa9, 0xd9, 0x1f, 0x16, 0xb0, 0x7c, 0x8b, 0x71, 0xeb, 0xda, 0xdf,
	0x67, 0xfc, 0xda, 0x69, 0x16, 0xee, 0x91, 0xd3, 0x33, 0x00, 0x00, 0xff, 0xff, 0x08, 0xb8, 0xe6,
	0xa5, 0xeb, 0x00, 0x00, 0x00,
}
