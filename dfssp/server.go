package main

import (
	api "dfss/dfssp/api"
	"golang.org/x/net/context"
)

type server struct{}

// Register handler
//
// Handle incoming RegisterRequest messages
func (s *server) Register(ctx context.Context, in *api.RegisterRequest) (*api.ErrorCode, error) {
	// TODO
	_ = new(server)
	return nil, nil
}

// Auth handler
//
// Handle incoming AuthRequest messages
func (s *server) Auth(ctx context.Context, in *api.AuthRequest) (*api.RegisteredUser, error) {
	// TODO
	return nil, nil
}

// Unregister handler
//
// Handle incoming UnregisterRequest messages
func (s *server) Unegister(ctx context.Context, in *api.Empty) (*api.ErrorCode, error) {
	// TODO
	return nil, nil
}

// PostContract handler
//
// Handle incoming PostContractRequest messages
func (s *server) PostContract(ctx context.Context, in *api.PostContractRequest) (*api.ErrorCode, error) {
	// TODO
	return nil, nil
}

// JoinSignature handler
//
// Handle incoming JoinSignatureRequest messages
func (s *server) JoinSignature(ctx context.Context, in *api.JoinSignatureRequest) (*api.ErrorCode, error) {
	// TODO
	return nil, nil
}

// ReadySign handler
//
// Handle incoming ReadySignRequest messages
func (s *server) ReadySign(ctx context.Context, in *api.ReadySignRequest) (*api.ErrorCode, error) {
	// TODO
	return nil, nil
}
