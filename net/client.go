package net

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"time"

	"dfss/auth"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// DefaultTimeout should be used when a non-critical timeout is used in the application.
var DefaultTimeout = 30 * time.Second

// Connect to a peer.
//
// Given parameters cert/key/ca are PEM-encoded array of bytes.
// Closing must be defered after call.
//
// The cert and key parameters can be set as nil for an unauthentified connection.
// If they are not, they will be provided to the remote server for authentification.
//
// serverCertHash will be matched against the remote server certificate.
// If nil, Connect will consider that the remote server is the root ca.
func Connect(addrPort string, cert *x509.Certificate, key *rsa.PrivateKey, ca *x509.Certificate, serverCertHash []byte) (*grpc.ClientConn, error) {

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
	conf := tls.Config{
		Certificates:       certificates,
		RootCAs:            caCertPool,
		InsecureSkipVerify: true, // Don't panic, it's normal and safe. See tlsCreds structure.
	}

	if serverCertHash == nil {
		serverCertHash = auth.GetCertificateHash(ca)
	}

	// let's do the dialing !
	return grpc.Dial(
		addrPort,
		grpc.WithTransportCredentials(&tlsCreds{config: conf, serverCertHash: serverCertHash}),
		grpc.WithTimeout(DefaultTimeout),
	)
}

// tlsCreds reimplements the default grpc TLS authenticator with no hostname verification.
// It is required because we need to connect to clients with their IP, and there is no IP SANs in our certificates.
//
// We need to enable the "InsecureSkipVerify" to perform this, that's why it's important to check the server certificate
// during the authentication process.
//
// See crypto/tls/handshake_client.go and google.golang.org/grpc/credentials/credentials.go
type tlsCreds struct {
	config         tls.Config
	serverCertHash []byte
}

func (c *tlsCreds) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2",
	}
}

func (c *tlsCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return nil, nil
}

func (c *tlsCreds) RequireTransportSecurity() bool {
	return true
}

func (c *tlsCreds) ClientHandshake(addr string, rawConn net.Conn, timeout time.Duration) (_ net.Conn, _ credentials.AuthInfo, err error) {
	var errChannel chan error
	if timeout != 0 {
		errChannel = make(chan error, 2)
		time.AfterFunc(timeout, func() {
			errChannel <- errors.New("credentials: Dial timed out")
		})
	}

	// Establish a secure connection WITHOUT certificate verification
	conn := tls.Client(rawConn, &c.config)
	if timeout == 0 {
		err = conn.Handshake()
	} else {
		go func() { errChannel <- conn.Handshake() }()
		err = <-errChannel
	}

	if err != nil { // Error during handshake
		_ = rawConn.Close()
		return nil, nil, err
	}

	// Successful handshake, BUT we have to authentify the server NOW
	opts := x509.VerifyOptions{
		Roots:       c.config.RootCAs,
		CurrentTime: time.Now(),
	}

	var chains [][]*x509.Certificate

	state := conn.ConnectionState()
	serverCert := state.PeerCertificates[0]
	chains, err = serverCert.Verify(opts)
	state.VerifiedChains = chains

	if err != nil {
		_ = rawConn.Close()
		return nil, nil, err
	}

	if c.serverCertHash != nil {
		// Additional check for the server cert hash
		if !bytes.Equal(auth.GetCertificateHash(serverCert), c.serverCertHash) {
			_ = rawConn.Close()
			return nil, nil, errors.New("credentials: Bad remote certificate hash")
		}
	}

	return conn, nil, nil
}

func (c *tlsCreds) ServerHandshake(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return nil, nil, errors.New("Server side handshake not implemented")
}
