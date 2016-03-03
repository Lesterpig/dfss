package sign

import (
	"dfss"
	cAPI "dfss/dfssc/api"
	pAPI "dfss/dfssp/api"
	"dfss/net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type clientServer struct{}

// TreatPromise handler
//
// Handle incoming TreatPromise messages
func (s *clientServer) TreatPromise(ctx context.Context, in *cAPI.Promise) (*pAPI.ErrorCode, error) {
	// TODO
	return nil, nil
}

// TreatSignature handler
//
// Handle incoming TreatSignature messages
func (s *clientServer) TreatSignature(ctx context.Context, in *cAPI.Signature) (*pAPI.ErrorCode, error) {
	// TODO
	return nil, nil
}

// Discover handler
//
// Handle incoming Discover messages
func (s *clientServer) Discover(ctx context.Context, in *cAPI.Hello) (*cAPI.Hello, error) {
	return &cAPI.Hello{Version: dfss.Version}, nil
}

// GetServer create and registers a ClientServer, returning the associted GRPC server
func (m *SignatureManager) GetServer() *grpc.Server {
	server := net.NewServer(m.cert, m.key, m.ca)
	cAPI.RegisterClientServer(server, &clientServer{})
	return server
}
