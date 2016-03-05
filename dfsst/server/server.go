package server

import (
	"crypto/rsa"
	"crypto/x509"
	"dfss/auth"
	"dfss/dfsst/api"
	"dfss/mgdb"
	"dfss/net"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

const (
	// PkeyFileName is the private key file default name
	PkeyFileName = "dfsst_pkey.pem"
	// CertFileName is the certificate file default name
	CertFileName = "dfsst_cert.pem"
	// RootCAFileName is the root certificate file default name
	RootCAFileName = "dfsst_rootca.pem"
)

type ttpServer struct {
	PKey         *rsa.PrivateKey
	Cert         *x509.Certificate
	RootCA       *x509.Certificate
	DB           *mgdb.MongoManager
	CertDuration int
	Verbose      bool
}

// Alert route for the TTP
func (server *ttpServer) Alert(ctx context.Context, in *api.AlertRequest) (*api.TTPResponse, error) {
	return nil, nil
}

// Recover route for the TTP
func (server *ttpServer) Recover(ctx context.Context, in *api.RecoverRequest) (*api.TTPResponse, error) {
	return nil, nil
}

// GetServer returns the gRPC server
func GetServer(keyPath, db string, certValidity int, verbose bool) *grpc.Server {
	server, err := loadCredentials(keyPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occured during the private key and certificates check:", err)
		os.Exit(1)
	}

	dbManager, err := mgdb.NewManager(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occured during the connection to MongoDB:", err)
		os.Exit(1)
	}

	server.CertDuration = certValidity
	server.Verbose = verbose
	server.DB = dbManager

	netServer := net.NewServer(server.Cert, server.PKey, server.RootCA)
	api.RegisterTTPServer(netServer, server)
	return netServer

}

// Assert the private key and certificates are valid
func loadCredentials(path string) (*ttpServer, error) {
	keyPath := filepath.Join(path, PkeyFileName)
	certPath := filepath.Join(path, CertFileName)
	rootCAPath := filepath.Join(path, RootCAFileName)

	// Recovering the private rsa key from file.
	keyBytes, err := ioutil.ReadFile(keyPath)

	if err != nil {
		return nil, err
	}

	key, err := auth.PEMToPrivateKey(keyBytes)

	if err != nil {
		return nil, err
	}

	// Recovering the certificate from file.
	certBytes, err := ioutil.ReadFile(certPath)

	if err != nil {
		return nil, err
	}

	cert, err := auth.PEMToCertificate(certBytes)

	if err != nil {
		return nil, err
	}

	// Recovering the root certificate from file.
	caBytes, err := ioutil.ReadFile(rootCAPath)

	if err != nil {
		return nil, err
	}

	ca, err := auth.PEMToCertificate(caBytes)

	if err != nil {
		return nil, err
	}

	res := &ttpServer{
		PKey:   key,
		Cert:   cert,
		RootCA: ca,
	}

	return res, nil
}
