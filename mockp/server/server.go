// Package server is the mock server, all functions must be redefined
package server

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"dfss/dfssp/api"
	"dfss/mockp/fixtures"
	"dfss/net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// mockServer empty struct to use grpc server
type mockServer struct{}

// Register handler
//
// Handle incoming RegisterRequest messages
func (s *mockServer) Register(ctx context.Context, in *api.RegisterRequest) (*api.ErrorCode, error) {
	if response, ok := fixtures.RegisterFixture[in.Email]; ok {
		return response, nil
	}
	return fixtures.RegisterFixture["default"], nil
}

// Auth handler
//
// Handle incoming AuthRequest messages
func (s *mockServer) Auth(ctx context.Context, in *api.AuthRequest) (*api.RegisteredUser, error) {
	return fixtures.AuthFixture["default"], nil
}

// Unregister handler
//
// Handle incoming UnregisterRequest messages
func (s *mockServer) Unregister(ctx context.Context, in *api.Empty) (*api.ErrorCode, error) {
	// TODO
	return nil, nil
}

// PostContract handler
//
// Handle incoming PostContractRequest messages
func (s *mockServer) PostContract(ctx context.Context, in *api.PostContractRequest) (*api.ErrorCode, error) {
	return fixtures.CreateFixture[in.Comment], nil
}

// GetContract handler
//
// Handle incoming GetContractRequest messages
func (s *mockServer) GetContract(ctx context.Context, in *api.GetContractRequest) (*api.Contract, error) {
	return fixtures.FetchFixture[in.Uuid], nil
}

// JoinSignature handler
//
// Handle incoming JoinSignatureRequest messages
func (s *mockServer) JoinSignature(in *api.JoinSignatureRequest, stream api.Platform_JoinSignatureServer) error {
	// TODO
	return nil
}

// ReadySign handler
//
// Handle incoming ReadySignRequest messages
func (s *mockServer) ReadySign(ctx context.Context, in *api.ReadySignRequest) (*api.LaunchSignature, error) {
	// TODO
	return nil, nil
}

// GetServer returns the GRPC server associated with the platform
func GetServer(ca *x509.Certificate, pkey *rsa.PrivateKey) *grpc.Server {
	server := net.NewServer(ca, pkey, ca)
	api.RegisterPlatformServer(server, &mockServer{})
	return server
}

// Run the mock server on provided address and with provided ca
func Run(ca *x509.Certificate, pkey *rsa.PrivateKey, addrPort string) {
	srv := GetServer(ca, pkey)
	err := net.Listen(addrPort, srv)
	if err != nil {
		fmt.Println(err)
	}
}
