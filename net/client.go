package net

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Connect to a peer.
//
// Given parameters cert/key/ca are PEM-encoded array of bytes.
// Closing must be defered after call.
func Connect(addrPort string, cert *x509.Certificate, key *rsa.PrivateKey, ca *x509.Certificate) (*grpc.ClientConn, error) {

	var certificates = make([]tls.Certificate, 1)

	if key != nil && cert != nil {
		peerCert := tls.Certificate{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  key,
		}
		certificates = append(certificates, peerCert)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(ca)

	// configure transport authentificator
	ta := credentials.NewTLS(&tls.Config{
		Certificates: certificates,
		RootCAs:      caCertPool,
	})

	// let's do the dialing !
	return grpc.Dial(addrPort, grpc.WithTransportCredentials(ta))
}
