package net

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// Create a new grpc server with given tls creds
//
// cert/key/ca are PEM-encoded array of byte
//
// The returned grpcServer must be used in association with server{} to
// register APIs before calling Listen()
func NewServer(cert, key, ca []byte) *Server {
	// configure gRPC
	var opts []grpc.ServerOption

	serverCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatal("Load peer cert/key error: %v", err)
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

// Listen with specified server on addr:port TCP
//
// addr is the addr to bind to
// port is the port to listen on
func Listen(addr, port string, grpcServer *Server) {
	// open tcp socket
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatalf("Failed to open tcp socket: %v", err)
	}

	grpcServer.Serve(lis)
}
