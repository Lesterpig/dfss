// Package server provides the platform server.
package server

import (
	"fmt"
	"os"

	"dfss/auth"
	"dfss/dfssp/api"
	"dfss/dfssp/authority"
	"dfss/dfssp/common"
	"dfss/dfssp/contract"
	"dfss/dfssp/user"
	"dfss/mgdb"
	"dfss/net"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type platformServer struct {
	Pid   *authority.PlatformID
	DB    *mgdb.MongoManager
	Rooms *common.WaitingGroupMap
	TTPs  *authority.TTPHolder
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
	return user.Auth(s.Pid, s.DB, in)
}

// Unregister handler
//
// Handle incoming UnregisterRequest messages
func (s *platformServer) Unregister(ctx context.Context, in *api.Empty) (*api.ErrorCode, error) {
	return user.Unregister(s.DB, net.GetClientHash(&ctx)), nil
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

// GetContract handler
//
// Handle incoming GetContractRequest messages
func (s *platformServer) GetContract(ctx context.Context, in *api.GetContractRequest) (*api.Contract, error) {
	hash := net.GetClientHash(&ctx)
	if hash == nil {
		return &api.Contract{ErrorCode: &api.ErrorCode{Code: api.ErrorCode_BADAUTH}}, nil
	}
	return contract.Fetch(s.DB, in.Uuid, hash), nil
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

	signal := contract.ReadySign(s.DB, s.Rooms, &ctx, in)
	if signal.ErrorCode.Code == api.ErrorCode_SUCCESS {
		signal.Ttp = s.TTPs.Get() // Assign a ttp to this signature, if any available
		sealedSignal := *signal
		sealedSignal.ErrorCode = nil
		sealedSignal.Seal = nil
		var err error
		signal.Seal, err = auth.SignStructure(s.Pid.Pkey, sealedSignal)
		if err != nil {
			return &api.LaunchSignature{ErrorCode: &api.ErrorCode{Code: api.ErrorCode_INTERR}}, nil
		}
	}

	return signal, nil
}

// GetServer returns the GRPC server associated with the platform
func GetServer() *grpc.Server {
	pid, err := authority.Start(viper.GetString("path"))
	if err != nil {
		fmt.Println("An error occured during the private key and root certificate retrieval:", err)
		os.Exit(1)
	}

	dbManager, err := mgdb.NewManager(viper.GetString("dbURI"))
	if err != nil {
		fmt.Println("An error occured during the connection to MongoDB:", err)
		os.Exit(1)
	}

	ttpholder, err := authority.NewTTPHolder(viper.GetString("ttps"))
	if err != nil {
		fmt.Println("An error occured during the ttp file load:", err)
	}

	if ttpholder.Nb() == 0 {
		fmt.Println("Warning: no TTP loaded. See `dfssp ttp --help`.")
	}

	server := net.NewServer(pid.RootCA, pid.Pkey, pid.RootCA)
	api.RegisterPlatformServer(server, &platformServer{
		Pid:   pid,
		DB:    dbManager,
		Rooms: common.NewWaitingGroupMap(),
		TTPs:  ttpholder,
	})
	return server
}
