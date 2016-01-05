package net

import (
	"crypto/tls"
	"crypto/x509"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// Connect to a peer.
//
// Given parameters cert/key/ca are PEM-encoded array of byte
// Closing must be defered after call
func Connect(addrPort string, cert, key, ca []byte) *grpc.ClientConn {
	// load peer cert/key, ca as PEM buffers
	peerCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatalf("Load peer cert/key error: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(ca)

	// configure transport authentificator
	ta := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{peerCert},
		RootCAs:      caCertPool,
	})

	// let's do the dialing !
	con, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(ta))
	if err != nil {
		grpclog.Fatalf("Fail to dial: %v", err)
	}

	return con
}
