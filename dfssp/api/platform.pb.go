// Code generated by protoc-gen-go.
// source: dfss/dfssp/api/platform.proto
// DO NOT EDIT!

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	dfss/dfssp/api/platform.proto

It has these top-level messages:
	RegisterRequest
	ErrorCode
	AuthRequest
	RegisteredUser
	Empty
	PostContractRequest
	GetContractRequest
	Contract
	JoinSignatureRequest
	UserConnected
	User
	ReadySignRequest
	LaunchSignature
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

type ErrorCode_Code int32

const (
	// SUCCESS is the error code for a successful request
	ErrorCode_SUCCESS ErrorCode_Code = 0
	// INVARG is the error code for an invalid argument
	ErrorCode_INVARG ErrorCode_Code = 1
	// BADAUTH is the error code for a bad authentication
	ErrorCode_BADAUTH ErrorCode_Code = 2
	// WARNING is the error code for a success state containing a specific warning message
	ErrorCode_WARNING ErrorCode_Code = 3
	// INTERR is the error code for an internal server error
	ErrorCode_INTERR ErrorCode_Code = -1
	// TIMEOUT is the error code for a timeout or unreacheable target
	ErrorCode_TIMEOUT ErrorCode_Code = -2
)

var ErrorCode_Code_name = map[int32]string{
	0:  "SUCCESS",
	1:  "INVARG",
	2:  "BADAUTH",
	3:  "WARNING",
	-1: "INTERR",
	-2: "TIMEOUT",
}
var ErrorCode_Code_value = map[string]int32{
	"SUCCESS": 0,
	"INVARG":  1,
	"BADAUTH": 2,
	"WARNING": 3,
	"INTERR":  -1,
	"TIMEOUT": -2,
}

func (x ErrorCode_Code) String() string {
	return proto.EnumName(ErrorCode_Code_name, int32(x))
}
func (ErrorCode_Code) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

// RegisterRequest message contains the client's email adress and his
// request (ie the PEM-encoded certificate request)
type RegisterRequest struct {
	Email   string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	Request string `protobuf:"bytes,2,opt,name=request" json:"request,omitempty"`
}

func (m *RegisterRequest) Reset()                    { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string            { return proto.CompactTextString(m) }
func (*RegisterRequest) ProtoMessage()               {}
func (*RegisterRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// ErrorCode message contains an error code and a message
// Above or zero : target-side error
// Less than 0   : local error
type ErrorCode struct {
	Code    ErrorCode_Code `protobuf:"varint,1,opt,name=code,enum=api.ErrorCode_Code" json:"code,omitempty"`
	Message string         `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *ErrorCode) Reset()                    { *m = ErrorCode{} }
func (m *ErrorCode) String() string            { return proto.CompactTextString(m) }
func (*ErrorCode) ProtoMessage()               {}
func (*ErrorCode) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// AuthRequest message contains the client's email adress and the token used
// for authentication
type AuthRequest struct {
	Email string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	Token string `protobuf:"bytes,2,opt,name=token" json:"token,omitempty"`
}

func (m *AuthRequest) Reset()                    { *m = AuthRequest{} }
func (m *AuthRequest) String() string            { return proto.CompactTextString(m) }
func (*AuthRequest) ProtoMessage()               {}
func (*AuthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// RegisteredUser message contains the generated client certificate
// (PEM-encoded)
type RegisteredUser struct {
	ClientCert string `protobuf:"bytes,1,opt,name=clientCert" json:"clientCert,omitempty"`
}

func (m *RegisteredUser) Reset()                    { *m = RegisteredUser{} }
func (m *RegisteredUser) String() string            { return proto.CompactTextString(m) }
func (*RegisteredUser) ProtoMessage()               {}
func (*RegisteredUser) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// Empty message is an empty message
type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

// PostContractRequest message contains the contract as SHA-512 hash, its filename,
// the list of signers as an array of strings, and a comment
type PostContractRequest struct {
	Hash     []byte   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Filename string   `protobuf:"bytes,2,opt,name=filename" json:"filename,omitempty"`
	Signer   []string `protobuf:"bytes,3,rep,name=signer" json:"signer,omitempty"`
	Comment  string   `protobuf:"bytes,4,opt,name=comment" json:"comment,omitempty"`
}

func (m *PostContractRequest) Reset()                    { *m = PostContractRequest{} }
func (m *PostContractRequest) String() string            { return proto.CompactTextString(m) }
func (*PostContractRequest) ProtoMessage()               {}
func (*PostContractRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

// GetContractRequest message contains the uuid of the asked contract
type GetContractRequest struct {
	Uuid string `protobuf:"bytes,1,opt,name=uuid" json:"uuid,omitempty"`
}

func (m *GetContractRequest) Reset()                    { *m = GetContractRequest{} }
func (m *GetContractRequest) String() string            { return proto.CompactTextString(m) }
func (*GetContractRequest) ProtoMessage()               {}
func (*GetContractRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

// Contract is the return value when a contract is fetched from the platform.
// The contract is in json format to avoid duplicating structures.
type Contract struct {
	ErrorCode *ErrorCode `protobuf:"bytes,1,opt,name=errorCode" json:"errorCode,omitempty"`
	Json      []byte     `protobuf:"bytes,2,opt,name=json,proto3" json:"json,omitempty"`
}

func (m *Contract) Reset()                    { *m = Contract{} }
func (m *Contract) String() string            { return proto.CompactTextString(m) }
func (*Contract) ProtoMessage()               {}
func (*Contract) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Contract) GetErrorCode() *ErrorCode {
	if m != nil {
		return m.ErrorCode
	}
	return nil
}

// JoinSignatureRequest message contains the contract to join unique identifier
// and the port the client will be listening at
type JoinSignatureRequest struct {
	ContractUuid string `protobuf:"bytes,1,opt,name=contractUuid" json:"contractUuid,omitempty"`
	Port         uint32 `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
}

func (m *JoinSignatureRequest) Reset()                    { *m = JoinSignatureRequest{} }
func (m *JoinSignatureRequest) String() string            { return proto.CompactTextString(m) }
func (*JoinSignatureRequest) ProtoMessage()               {}
func (*JoinSignatureRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

// UserConnected is emitted by the platform to the client to announce a new client connection
type UserConnected struct {
	ErrorCode    *ErrorCode `protobuf:"bytes,1,opt,name=errorCode" json:"errorCode,omitempty"`
	ContractUuid string     `protobuf:"bytes,2,opt,name=contractUuid" json:"contractUuid,omitempty"`
	User         *User      `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
}

func (m *UserConnected) Reset()                    { *m = UserConnected{} }
func (m *UserConnected) String() string            { return proto.CompactTextString(m) }
func (*UserConnected) ProtoMessage()               {}
func (*UserConnected) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *UserConnected) GetErrorCode() *ErrorCode {
	if m != nil {
		return m.ErrorCode
	}
	return nil
}

func (m *UserConnected) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type User struct {
	KeyHash []byte `protobuf:"bytes,1,opt,name=keyHash,proto3" json:"keyHash,omitempty"`
	Email   string `protobuf:"bytes,2,opt,name=email" json:"email,omitempty"`
	Ip      string `protobuf:"bytes,3,opt,name=ip" json:"ip,omitempty"`
	Port    uint32 `protobuf:"varint,4,opt,name=port" json:"port,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

// ReadySignRequest contains the contract unique identitier that is ready to be signed
type ReadySignRequest struct {
	ContractUuid string `protobuf:"bytes,1,opt,name=contractUuid" json:"contractUuid,omitempty"`
}

func (m *ReadySignRequest) Reset()                    { *m = ReadySignRequest{} }
func (m *ReadySignRequest) String() string            { return proto.CompactTextString(m) }
func (*ReadySignRequest) ProtoMessage()               {}
func (*ReadySignRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

// LaunchSignature is emitted by the platform when every signers are ready
type LaunchSignature struct {
	ErrorCode     *ErrorCode `protobuf:"bytes,1,opt,name=errorCode" json:"errorCode,omitempty"`
	SignatureUuid string     `protobuf:"bytes,2,opt,name=signatureUuid" json:"signatureUuid,omitempty"`
	KeyHash       [][]byte   `protobuf:"bytes,3,rep,name=keyHash,proto3" json:"keyHash,omitempty"`
	Sequence      []uint32   `protobuf:"varint,4,rep,name=sequence" json:"sequence,omitempty"`
	Hash          []byte     `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (m *LaunchSignature) Reset()                    { *m = LaunchSignature{} }
func (m *LaunchSignature) String() string            { return proto.CompactTextString(m) }
func (*LaunchSignature) ProtoMessage()               {}
func (*LaunchSignature) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *LaunchSignature) GetErrorCode() *ErrorCode {
	if m != nil {
		return m.ErrorCode
	}
	return nil
}

func init() {
	proto.RegisterType((*RegisterRequest)(nil), "api.RegisterRequest")
	proto.RegisterType((*ErrorCode)(nil), "api.ErrorCode")
	proto.RegisterType((*AuthRequest)(nil), "api.AuthRequest")
	proto.RegisterType((*RegisteredUser)(nil), "api.RegisteredUser")
	proto.RegisterType((*Empty)(nil), "api.Empty")
	proto.RegisterType((*PostContractRequest)(nil), "api.PostContractRequest")
	proto.RegisterType((*GetContractRequest)(nil), "api.GetContractRequest")
	proto.RegisterType((*Contract)(nil), "api.Contract")
	proto.RegisterType((*JoinSignatureRequest)(nil), "api.JoinSignatureRequest")
	proto.RegisterType((*UserConnected)(nil), "api.UserConnected")
	proto.RegisterType((*User)(nil), "api.User")
	proto.RegisterType((*ReadySignRequest)(nil), "api.ReadySignRequest")
	proto.RegisterType((*LaunchSignature)(nil), "api.LaunchSignature")
	proto.RegisterEnum("api.ErrorCode_Code", ErrorCode_Code_name, ErrorCode_Code_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion1

// Client API for Platform service

type PlatformClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*ErrorCode, error)
	Auth(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*RegisteredUser, error)
	Unregister(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ErrorCode, error)
	PostContract(ctx context.Context, in *PostContractRequest, opts ...grpc.CallOption) (*ErrorCode, error)
	GetContract(ctx context.Context, in *GetContractRequest, opts ...grpc.CallOption) (*Contract, error)
	JoinSignature(ctx context.Context, in *JoinSignatureRequest, opts ...grpc.CallOption) (Platform_JoinSignatureClient, error)
	ReadySign(ctx context.Context, in *ReadySignRequest, opts ...grpc.CallOption) (*LaunchSignature, error)
}

type platformClient struct {
	cc *grpc.ClientConn
}

func NewPlatformClient(cc *grpc.ClientConn) PlatformClient {
	return &platformClient{cc}
}

func (c *platformClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*ErrorCode, error) {
	out := new(ErrorCode)
	err := grpc.Invoke(ctx, "/api.Platform/Register", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *platformClient) Auth(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*RegisteredUser, error) {
	out := new(RegisteredUser)
	err := grpc.Invoke(ctx, "/api.Platform/Auth", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *platformClient) Unregister(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ErrorCode, error) {
	out := new(ErrorCode)
	err := grpc.Invoke(ctx, "/api.Platform/Unregister", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *platformClient) PostContract(ctx context.Context, in *PostContractRequest, opts ...grpc.CallOption) (*ErrorCode, error) {
	out := new(ErrorCode)
	err := grpc.Invoke(ctx, "/api.Platform/PostContract", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *platformClient) GetContract(ctx context.Context, in *GetContractRequest, opts ...grpc.CallOption) (*Contract, error) {
	out := new(Contract)
	err := grpc.Invoke(ctx, "/api.Platform/GetContract", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *platformClient) JoinSignature(ctx context.Context, in *JoinSignatureRequest, opts ...grpc.CallOption) (Platform_JoinSignatureClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Platform_serviceDesc.Streams[0], c.cc, "/api.Platform/JoinSignature", opts...)
	if err != nil {
		return nil, err
	}
	x := &platformJoinSignatureClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Platform_JoinSignatureClient interface {
	Recv() (*UserConnected, error)
	grpc.ClientStream
}

type platformJoinSignatureClient struct {
	grpc.ClientStream
}

func (x *platformJoinSignatureClient) Recv() (*UserConnected, error) {
	m := new(UserConnected)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *platformClient) ReadySign(ctx context.Context, in *ReadySignRequest, opts ...grpc.CallOption) (*LaunchSignature, error) {
	out := new(LaunchSignature)
	err := grpc.Invoke(ctx, "/api.Platform/ReadySign", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Platform service

type PlatformServer interface {
	Register(context.Context, *RegisterRequest) (*ErrorCode, error)
	Auth(context.Context, *AuthRequest) (*RegisteredUser, error)
	Unregister(context.Context, *Empty) (*ErrorCode, error)
	PostContract(context.Context, *PostContractRequest) (*ErrorCode, error)
	GetContract(context.Context, *GetContractRequest) (*Contract, error)
	JoinSignature(*JoinSignatureRequest, Platform_JoinSignatureServer) error
	ReadySign(context.Context, *ReadySignRequest) (*LaunchSignature, error)
}

func RegisterPlatformServer(s *grpc.Server, srv PlatformServer) {
	s.RegisterService(&_Platform_serviceDesc, srv)
}

func _Platform_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(PlatformServer).Register(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Platform_Auth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(AuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(PlatformServer).Auth(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Platform_Unregister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(PlatformServer).Unregister(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Platform_PostContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(PostContractRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(PlatformServer).PostContract(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Platform_GetContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(GetContractRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(PlatformServer).GetContract(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Platform_JoinSignature_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(JoinSignatureRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PlatformServer).JoinSignature(m, &platformJoinSignatureServer{stream})
}

type Platform_JoinSignatureServer interface {
	Send(*UserConnected) error
	grpc.ServerStream
}

type platformJoinSignatureServer struct {
	grpc.ServerStream
}

func (x *platformJoinSignatureServer) Send(m *UserConnected) error {
	return x.ServerStream.SendMsg(m)
}

func _Platform_ReadySign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(ReadySignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(PlatformServer).ReadySign(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Platform_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Platform",
	HandlerType: (*PlatformServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Platform_Register_Handler,
		},
		{
			MethodName: "Auth",
			Handler:    _Platform_Auth_Handler,
		},
		{
			MethodName: "Unregister",
			Handler:    _Platform_Unregister_Handler,
		},
		{
			MethodName: "PostContract",
			Handler:    _Platform_PostContract_Handler,
		},
		{
			MethodName: "GetContract",
			Handler:    _Platform_GetContract_Handler,
		},
		{
			MethodName: "ReadySign",
			Handler:    _Platform_ReadySign_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "JoinSignature",
			Handler:       _Platform_JoinSignature_Handler,
			ServerStreams: true,
		},
	},
}

var fileDescriptor0 = []byte{
	// 706 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x55, 0x5f, 0x6f, 0x12, 0x4b,
	0x14, 0x67, 0x81, 0x16, 0x38, 0x14, 0xba, 0x19, 0x7a, 0xef, 0xe5, 0x92, 0xd4, 0x34, 0x13, 0x13,
	0x1b, 0x63, 0xa0, 0xc1, 0x44, 0xa3, 0x6f, 0x14, 0x49, 0x5b, 0x53, 0xb1, 0x19, 0x40, 0x13, 0xdf,
	0xd6, 0xdd, 0x69, 0x59, 0xcb, 0xfe, 0x71, 0x66, 0x88, 0xe9, 0x9b, 0x1f, 0xc1, 0x6f, 0xe2, 0x8b,
	0x9f, 0x4f, 0x9d, 0x99, 0xdd, 0x59, 0x16, 0x24, 0x26, 0xe5, 0x61, 0x99, 0x73, 0xe6, 0xcc, 0xef,
	0x9c, 0xf3, 0x9b, 0xf3, 0xdb, 0x85, 0x43, 0xef, 0x9a, 0xf3, 0x9e, 0x7a, 0xc4, 0x3d, 0x27, 0xf6,
	0x7b, 0xf1, 0xc2, 0x11, 0xd7, 0x11, 0x0b, 0xba, 0x31, 0x8b, 0x44, 0x84, 0x4a, 0xd2, 0x87, 0x07,
	0xb0, 0x4f, 0xe8, 0x8d, 0xcf, 0x05, 0x65, 0x84, 0x7e, 0x5e, 0x52, 0x2e, 0xd0, 0x01, 0xec, 0xd0,
	0xc0, 0xf1, 0x17, 0x6d, 0xeb, 0xc8, 0x3a, 0xae, 0x91, 0xc4, 0x40, 0x6d, 0xa8, 0xb0, 0x24, 0xa0,
	0x5d, 0xd4, 0x7e, 0x63, 0xe2, 0x1f, 0x16, 0xd4, 0x46, 0x8c, 0x45, 0x6c, 0x18, 0x79, 0x14, 0x3d,
	0x82, 0xb2, 0x2b, 0xff, 0xf5, 0xe1, 0x66, 0xbf, 0xd5, 0x95, 0x49, 0xba, 0xd9, 0x6e, 0x57, 0x3d,
	0x88, 0x0e, 0x50, 0x80, 0x01, 0xe5, 0xdc, 0xb9, 0xa1, 0x06, 0x30, 0x35, 0xb1, 0x07, 0x65, 0x0d,
	0x55, 0x87, 0xca, 0x64, 0x36, 0x1c, 0x8e, 0x26, 0x13, 0xbb, 0x80, 0x00, 0x76, 0x2f, 0xc6, 0xef,
	0x06, 0xe4, 0xcc, 0xb6, 0xd4, 0xc6, 0xe9, 0xe0, 0xd5, 0x60, 0x36, 0x3d, 0xb7, 0x8b, 0xca, 0x78,
	0x3f, 0x20, 0xe3, 0x8b, 0xf1, 0x99, 0x5d, 0x42, 0x2d, 0x15, 0x35, 0x1d, 0x11, 0x62, 0xff, 0x32,
	0x3f, 0x4b, 0x36, 0x54, 0x99, 0x5e, 0xbc, 0x19, 0xbd, 0x9d, 0x4d, 0xed, 0x9f, 0x99, 0x17, 0xbf,
	0x80, 0xfa, 0x60, 0x29, 0xe6, 0x7f, 0xef, 0x5a, 0x7a, 0x45, 0x74, 0x4b, 0xc3, 0xb4, 0xc4, 0xc4,
	0xc0, 0x27, 0xd0, 0x34, 0xa4, 0x51, 0x6f, 0xc6, 0x29, 0x43, 0x0f, 0x00, 0xdc, 0x85, 0x4f, 0x43,
	0x31, 0xa4, 0x4c, 0xa4, 0x10, 0x39, 0x0f, 0xae, 0xc0, 0xce, 0x28, 0x88, 0xc5, 0x1d, 0xfe, 0x02,
	0xad, 0xab, 0x88, 0x8b, 0x61, 0x14, 0x0a, 0xe6, 0xb8, 0xc2, 0x64, 0x47, 0x50, 0x9e, 0x3b, 0x7c,
	0xae, 0x4f, 0xee, 0x11, 0xbd, 0x46, 0x1d, 0xa8, 0x5e, 0xfb, 0x0b, 0x1a, 0x3a, 0x81, 0x61, 0x28,
	0xb3, 0xd1, 0xbf, 0xb0, 0xcb, 0xfd, 0x9b, 0x90, 0xb2, 0x76, 0xe9, 0xa8, 0x24, 0x77, 0x52, 0x4b,
	0x91, 0xea, 0x46, 0x41, 0x20, 0xd3, 0xb6, 0xcb, 0x09, 0xa9, 0xa9, 0x89, 0x8f, 0x01, 0x9d, 0xd1,
	0x6d, 0x79, 0x97, 0x4b, 0xdf, 0x4b, 0x2b, 0xd6, 0x6b, 0x7c, 0x09, 0x55, 0x13, 0x86, 0x9e, 0x40,
	0x8d, 0x9a, 0xcb, 0xd3, 0x41, 0xf5, 0x7e, 0x73, 0xfd, 0x4a, 0xc9, 0x2a, 0x40, 0xa1, 0x7d, 0xe2,
	0x51, 0x42, 0x96, 0xec, 0x42, 0xad, 0xf1, 0x18, 0x0e, 0x5e, 0x47, 0x7e, 0x38, 0x91, 0xf5, 0x39,
	0x62, 0xc9, 0xa8, 0xc9, 0x8c, 0x61, 0xcf, 0x4d, 0xb3, 0xcc, 0x56, 0x15, 0xac, 0xf9, 0x14, 0x5e,
	0x1c, 0xb1, 0x64, 0xe0, 0x1a, 0x44, 0xaf, 0xf1, 0x57, 0x0b, 0x1a, 0x8a, 0x72, 0x59, 0x62, 0x48,
	0x5d, 0x41, 0xbd, 0x7b, 0xd6, 0xb8, 0x99, 0xb7, 0xb8, 0x25, 0xef, 0xa1, 0x64, 0x85, 0x6b, 0x6e,
	0x15, 0x58, 0x4d, 0x83, 0xa9, 0x9c, 0x44, 0xbb, 0xf1, 0x07, 0x28, 0xeb, 0x4b, 0x97, 0x64, 0xdf,
	0xd2, 0xbb, 0xf3, 0xd5, 0xbd, 0x19, 0x73, 0x35, 0x4c, 0xc5, 0xfc, 0x30, 0x35, 0xa1, 0xe8, 0xc7,
	0x1a, 0xb4, 0x46, 0xe4, 0x2a, 0x6b, 0xaf, 0x9c, 0x6b, 0xef, 0x19, 0xd8, 0x84, 0x3a, 0xde, 0x9d,
	0xe2, 0xeb, 0x1e, 0x54, 0xe1, 0xef, 0x16, 0xec, 0x5f, 0x3a, 0xcb, 0xd0, 0x9d, 0x67, 0x4c, 0xdf,
	0x93, 0x98, 0x87, 0xd0, 0xe0, 0xe6, 0x68, 0x8e, 0x99, 0x75, 0x67, 0xbe, 0x67, 0x35, 0x79, 0xb9,
	0x9e, 0xe5, 0xb8, 0x72, 0x55, 0x70, 0xe8, 0x52, 0xd9, 0x51, 0x49, 0x76, 0x94, 0xd9, 0xd9, 0x78,
	0xef, 0xac, 0xc6, 0xbb, 0xff, 0xad, 0x04, 0xd5, 0xab, 0xf4, 0x8d, 0x84, 0xfa, 0x50, 0x35, 0x8a,
	0x42, 0x07, 0xba, 0xc6, 0x8d, 0xb7, 0x52, 0x67, 0xa3, 0x72, 0x5c, 0x40, 0x3d, 0x28, 0x2b, 0x01,
	0x23, 0x5b, 0xef, 0xe4, 0xb4, 0xdc, 0x69, 0xad, 0x21, 0x24, 0x12, 0x95, 0x07, 0x1e, 0x03, 0xcc,
	0x42, 0x66, 0xd2, 0x40, 0x02, 0xa8, 0x54, 0xb9, 0x05, 0xfc, 0x25, 0xec, 0xe5, 0x75, 0x8a, 0xda,
	0x3a, 0x62, 0x8b, 0x74, 0xb7, 0x9c, 0x7d, 0x0e, 0xf5, 0x9c, 0xd4, 0xd0, 0x7f, 0x3a, 0xe0, 0x4f,
	0xf1, 0x75, 0x1a, 0x7a, 0xc3, 0x78, 0xe5, 0xc1, 0x53, 0x68, 0xac, 0x69, 0x05, 0xfd, 0xaf, 0x23,
	0xb6, 0xe9, 0xa7, 0x83, 0xb2, 0xa9, 0xcc, 0x94, 0x80, 0x0b, 0x27, 0x96, 0x2c, 0xbc, 0x96, 0x0d,
	0x10, 0xfa, 0x27, 0x25, 0x62, 0x7d, 0xa0, 0x3a, 0x09, 0xc3, 0x1b, 0xe3, 0x82, 0x0b, 0x1f, 0x77,
	0xf5, 0x87, 0xe1, 0xe9, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc5, 0x4c, 0xa0, 0x21, 0x39, 0x06,
	0x00, 0x00,
}
