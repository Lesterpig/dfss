package sign

import (
	"dfss"
	cAPI "dfss/dfssc/api"
	pAPI "dfss/dfssp/api"
	"dfss/dfsst/entities"
	"dfss/net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type clientServer struct {
	incomingPromises   chan interface{}
	incomingSignatures chan interface{}
}

func getServerErrorCode(c chan interface{}, in interface{}) *pAPI.ErrorCode {
	if c != nil {
		c <- in
		return &pAPI.ErrorCode{Code: pAPI.ErrorCode_SUCCESS}
	}
	return &pAPI.ErrorCode{Code: pAPI.ErrorCode_INTERR} // server not ready
}

// TreatPromise handler
//
// Handle incoming TreatPromise messages
func (s *clientServer) TreatPromise(ctx context.Context, in *cAPI.Promise) (*pAPI.ErrorCode, error) {
	// we check that the incoming promise is valid (ie. no data inconsistency)
	// we do not check that we expected that promise
	valid, _, _, _ := entities.IsRequestValid(ctx, []*cAPI.Promise{in})
	if !valid {
		return &pAPI.ErrorCode{Code: pAPI.ErrorCode_SUCCESS}, nil
	}
	return getServerErrorCode(s.incomingPromises, in), nil
}

// TreatSignature handler
//
// Handle incoming TreatSignature messages
func (s *clientServer) TreatSignature(ctx context.Context, in *cAPI.Signature) (*pAPI.ErrorCode, error) {
	return getServerErrorCode(s.incomingSignatures, in), nil
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
	m.cServerIface = clientServer{}
	cAPI.RegisterClientServer(server, &m.cServerIface)
	return server
}
