package net

import (
	"crypto/tls"
	"crypto/x509"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"log"
	"net"
)

// NewServer creates a new grpc server with given tls credentials.
//
// cert/key/ca are PEM-encoded array of bytes.
//
// The returned grpcServer must be used in association with server{} to
// register APIs before calling Listen().
func NewServer(cert, key, ca []byte) *grpc.Server {
	// configure gRPC
	var opts []grpc.ServerOption

	serverCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatalf("Load peer cert/key error: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(ca)

	// configure transport authentificator
	ta := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		RootCAs:      caCertPool,
		ClientCAs:    caCertPool,
		ClientAuth:   tls.VerifyClientCertIfGiven,
	})

	opts = []grpc.ServerOption{grpc.Creds(ta)}
	return grpc.NewServer(opts...)
}

// Listen with specified server on addr:port.
//
// addrPort is formated as 127.0.0.1:8001.
func Listen(addrPort string, grpcServer *grpc.Server) error {
	// open tcp socket
	lis, err := net.Listen("tcp", addrPort)
	if err != nil {
		return err
	}
	return grpcServer.Serve(lis)
}

// GetTLSState returns the current tls connection state from a grpc context.
// If you just need to check that the connected peer provides its certificate, use `GetCN`.
func GetTLSState(ctx *context.Context) (tls.ConnectionState, bool) {
	p, ok := peer.FromContext(*ctx)
	if !ok {
		return tls.ConnectionState{}, false
	}
	return p.AuthInfo.(credentials.TLSInfo).State, true
}

// GetCN returns the current common name of connected peer from grpc context.
// The returned string is empty if encountering a non-auth peer.
func GetCN(ctx *context.Context) string {
	state, ok := GetTLSState(ctx)
	if !ok || len(state.VerifiedChains) == 0 {
		return ""
	}
	return state.VerifiedChains[0][0].Subject.CommonName
}
