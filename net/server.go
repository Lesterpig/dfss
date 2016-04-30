// Package net wraps TLS and GRPC client/server to simplify connections.
package net

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"net"

	"dfss/auth"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// NewServer creates a new grpc server with given tls credentials.
//
// cert/key/ca are PEM-encoded array of bytes.
//
// The returned grpcServer must be used in association with server{} to
// register APIs before calling Listen().
func NewServer(cert *x509.Certificate, key *rsa.PrivateKey, ca *x509.Certificate) *grpc.Server {
	// configure gRPC
	var opts []grpc.ServerOption

	serverCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(ca)

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
func GetTLSState(ctx *context.Context) (tls.ConnectionState, net.Addr, bool) {
	p, ok := peer.FromContext(*ctx)
	if !ok {
		return tls.ConnectionState{}, nil, false
	}
	return p.AuthInfo.(credentials.TLSInfo).State, p.Addr, true
}

// GetCN returns the current common name of connected peer from grpc context.
// The returned string is empty if encountering a non-auth peer.
func GetCN(ctx *context.Context) string {
	state, _, ok := GetTLSState(ctx)
	if !ok || len(state.VerifiedChains) == 0 {
		return ""
	}
	return state.VerifiedChains[0][0].Subject.CommonName
}

// GetClientHash returns the current certificate hash of connected peer from grpc context.
// The returned slice is nil if encoutering a non-auth peer.
func GetClientHash(ctx *context.Context) []byte {
	state, _, ok := GetTLSState(ctx)
	if !ok || len(state.VerifiedChains) == 0 {
		return nil
	}
	return auth.GetCertificateHash(state.VerifiedChains[0][0])
}
