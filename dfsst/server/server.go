package server

import (
	"fmt"
	"os"

	"dfss/dfssc/security"
	"dfss/dfsst/api"
	"dfss/mgdb"
	"dfss/net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ttpServer struct {
	DB      *mgdb.MongoManager
	Verbose bool
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
func GetServer(fca, fcert, fkey, password, db string, verbose bool) *grpc.Server {
	auth := security.NewAuthContainer(fca, fcert, fkey, "", password)
	ca, cert, key, err := auth.LoadFiles()
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occured during the private key and certificates retrieval:", err)
		os.Exit(1)
	}

	dbManager, err := mgdb.NewManager(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occured during the connection to MongoDB:", err)
		os.Exit(2)
	}

	server := &ttpServer{
		Verbose: verbose,
		DB:      dbManager,
	}
	netServer := net.NewServer(cert, key, ca)
	api.RegisterTTPServer(netServer, server)
	return netServer
}
