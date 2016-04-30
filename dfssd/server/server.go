// Package server provides the demonstration unsecure server.
package server

import (
	"net"

	api "dfss/dfssd/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type Server struct{}

// Sendlog Handler
//
// Handle incoming log messages
func (s *Server) SendLog(ctx context.Context, in *api.Log) (*api.Ack, error) {
	addMessage(in)
	return &api.Ack{}, nil
}

// Listen with gRPG service
func Listen(addrPort string, fn func(*api.Log)) error {
	// open tcp socket
	lis, err := net.Listen("tcp", addrPort)
	if err != nil {
		grpclog.Fatalf("Failed to open tcp socket: %v", err)
		return err
	}

	// log display manager
	go displayHandler(fn)

	// bootstrap gRPC service !
	grpcServer := grpc.NewServer()
	api.RegisterDemonstratorServer(grpcServer, &Server{})
	err = grpcServer.Serve(lis)

	return err
}
