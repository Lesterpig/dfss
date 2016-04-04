package sign

import (
	"dfss"
	cAPI "dfss/dfssc/api"
	pAPI "dfss/dfssp/api"
	"dfss/net"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type clientServer struct{}

// TreatPromise handler
//
// Handle incoming TreatPromise messages
func (s *clientServer) TreatPromise(ctx context.Context, in *cAPI.Promise) (*pAPI.ErrorCode, error) {
	// Pass the message to Sign()
	if incomingPromises != nil {
		incomingPromises <- in
		// Maybe we can add another channel here for better error management
		return &pAPI.ErrorCode{Code: pAPI.ErrorCode_SUCCESS}, nil
	}

	return &pAPI.ErrorCode{Code: pAPI.ErrorCode_INVARG}, fmt.Errorf("Cannot pass incoming promise")
}

// TreatSignature handler
//
// Handle incoming TreatSignature messages
func (s *clientServer) TreatSignature(ctx context.Context, in *cAPI.Signature) (*pAPI.ErrorCode, error) {
	if incomingSignatures != nil {
		incomingSignatures <- in
		// Maybe we can add another channel here for better error management
		return &pAPI.ErrorCode{Code: pAPI.ErrorCode_SUCCESS}, nil
	}

	return &pAPI.ErrorCode{Code: pAPI.ErrorCode_INVARG}, fmt.Errorf("Cannot pass incoming signature")

	return nil, nil
}

// Discover handler
//
// Handle incoming Discover messages
func (s *clientServer) Discover(ctx context.Context, in *cAPI.Hello) (*cAPI.Hello, error) {
	return &cAPI.Hello{Version: dfss.Version}, nil
}

// GetServer create and registers a ClientServer, returning the associated GRPC server
func (m *SignatureManager) GetServer() *grpc.Server {
	server := net.NewServer(m.auth.Cert, m.auth.Key, m.auth.CA)
	cAPI.RegisterClientServer(server, &clientServer{})
	return server
}
