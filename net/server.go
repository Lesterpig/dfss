package net

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
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
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	opts = []grpc.ServerOption{grpc.Creds(ta)}
	return grpc.NewServer(opts...)
}

// Listen with specified server on addr:port.
//
// addrPort is formated as 127.0.0.1:8001.
func Listen(addrPort string, grpcServer *grpc.Server) {
	// open tcp socket
	lis, err := net.Listen("tcp", addrPort)
	if err != nil {
		grpclog.Fatalf("Failed to open tcp socket: %v", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		grpclog.Fatalf("Failed to bind gRPC server: %v", err)
	}
}
