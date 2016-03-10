package server

import (
	"fmt"
	"os"

	"dfss/dfssp/api"
	"dfss/dfssp/authority"
	"dfss/dfssp/common"
	"dfss/dfssp/contract"
	"dfss/dfssp/user"
	"dfss/mgdb"
	"dfss/net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type platformServer struct {
	Pid          *authority.PlatformID
	DB           *mgdb.MongoManager
	Rooms        *common.WaitingGroupMap
	CertDuration int
	Verbose      bool
}

// Register handler
//
// Handle incoming RegisterRequest messages
func (s *platformServer) Register(ctx context.Context, in *api.RegisterRequest) (*api.ErrorCode, error) {
	return user.Register(s.DB, in)
}

// Auth handler
//
// Handle incoming AuthRequest messages
func (s *platformServer) Auth(ctx context.Context, in *api.AuthRequest) (*api.RegisteredUser, error) {
	return user.Auth(s.Pid, s.DB, s.CertDuration, in)
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

	cn := net.GetCN(&ctx)
	if len(cn) == 0 {
		return &api.ErrorCode{Code: api.ErrorCode_BADAUTH}, nil
	}

	builder := contract.NewContractBuilder(s.DB, in)
	return builder.Execute(), nil
}

// JoinSignature handler
//
// Handle incoming JoinSignatureRequest messages
func (s *platformServer) JoinSignature(in *api.JoinSignatureRequest, stream api.Platform_JoinSignatureServer) error {

	ctx := stream.Context()
	cn := net.GetCN(&ctx)
	if len(cn) == 0 {
		_ = stream.Send(&api.UserConnected{
			ErrorCode: &api.ErrorCode{Code: api.ErrorCode_BADAUTH},
		})
		return nil
	}

	contract.JoinSignature(s.DB, s.Rooms, in, stream)
	return nil
}

// ReadySign handler
//
// Handle incoming ReadySignRequest messages
func (s *platformServer) ReadySign(ctx context.Context, in *api.ReadySignRequest) (*api.LaunchSignature, error) {
	cn := net.GetCN(&ctx)
	if len(cn) == 0 {
		return &api.LaunchSignature{ErrorCode: &api.ErrorCode{Code: api.ErrorCode_BADAUTH}}, nil
	}
	return contract.ReadySign(s.DB, s.Rooms, &ctx, in), nil
}

// GetServer returns the GRPC server associated with the platform
func GetServer(keyPath, db string, certValidity int, verbose bool) *grpc.Server {
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
		Pid:          pid,
		DB:           dbManager,
		Rooms:        common.NewWaitingGroupMap(),
		CertDuration: certValidity,
		Verbose:      verbose,
	})
	return server
}
