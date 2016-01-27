package net

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Connect to a peer.
//
// Given parameters cert/key/ca are PEM-encoded array of bytes.
// Closing must be defered after call.
func Connect(addrPort string, cert, key, ca []byte) (*grpc.ClientConn, error) {

	var certificates = make([]tls.Certificate, 1)

	if len(key) > 0 && len(cert) > 0 {
		// load peer cert/key, ca as PEM buffers
		peerCert, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, peerCert)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(ca)
	if !ok {
		return nil, errors.New("Bad format for CA")
	}

	// configure transport authentificator
	ta := credentials.NewTLS(&tls.Config{
		Certificates: certificates,
		RootCAs:      caCertPool,
	})

	// let's do the dialing !
	return grpc.Dial(addrPort, grpc.WithTransportCredentials(ta))
}
