package net

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// Connect to a peer
//
// Closing must be defered after call
func Connect(addr_port string, cert, key, ca []byte) *ClientConn {
	// load peer cert/key, ca as PEM buffers
	peerCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatal("Load peer cert/key error: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(ca)

	// configure transport authentificator
	ta := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{peerCert},
		RootCAs:      caCertPool,
	})

	// let's do the dialing !
	con, err := grpc.Dial(addr_port, grpc.WithTransportCredentials(ta))
	if err != nil {
		grpclog.Fatalf("Fail to dial: %v", err)
	}

	return con
}
