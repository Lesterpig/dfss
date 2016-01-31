package server

import (
	"fmt"
	"os"

	"dfss/dfssp/api"
	"dfss/dfssp/authority"
	"dfss/dfssp/contract"
	"dfss/mgdb"
	"dfss/net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type platformServer struct {
	Pid     *authority.PlatformID
	DB      *mgdb.MongoManager
	Verbose bool
}

// Register handler
//
// Handle incoming RegisterRequest messages
func (s *platformServer) Register(ctx context.Context, in *api.RegisterRequest) (*api.ErrorCode, error) {
	// TODO
	_ = new(platformServer)
	return nil, nil
}

// Auth handler
//
// Handle incoming AuthRequest messages
func (s *platformServer) Auth(ctx context.Context, in *api.AuthRequest) (*api.RegisteredUser, error) {
	// TODO
	return nil, nil
}

// Unregister handler
//
// Handle incoming UnregisterRequest messages
func (s *platformServer) Unregister(ctx context.Context, in *api.Empty) (*api.ErrorCode, error) {
	// TODO
	return nil, nil
}

// PostContract handler
//
// Handle incoming PostContractRequest messages
func (s *platformServer) PostContract(ctx context.Context, in *api.PostContractRequest) (*api.ErrorCode, error) {
	return contract.PostRoute(s.DB, in), nil
}

// JoinSignature handler
//
// Handle incoming JoinSignatureRequest messages
func (s *platformServer) JoinSignature(ctx context.Context, in *api.JoinSignatureRequest) (*api.ErrorCode, error) {
	// TODO
	return nil, nil
}

// ReadySign handler
//
// Handle incoming ReadySignRequest messages
func (s *platformServer) ReadySign(ctx context.Context, in *api.ReadySignRequest) (*api.ErrorCode, error) {
	// TODO
	return nil, nil
}

// GetServer returns the GRPC server associated with the platform
func GetServer(keyPath, db string, verbose bool) *grpc.Server {
	pid, err := authority.Start(keyPath)
	if err != nil {
		fmt.Println("An error occured during the private key and root certificate retrieval:", err)
		os.Exit(1)
	}

	dbManager, err := mgdb.NewManager(db)
	if err != nil {
		fmt.Println("An error occured during the connection to MongoDB:", err)
		os.Exit(1)
	}

	server := net.NewServer(pid.RootCA, pid.Pkey, pid.RootCA)
	api.RegisterPlatformServer(server, &platformServer{
		Pid:     pid,
		DB:      dbManager,
		Verbose: verbose,
	})
	return server
}
