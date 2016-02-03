package main

import (
	api "dfss/dfssd/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"net"
)

type server struct{}

// Sendlog Handler
//
// Handle incoming log messages
func (s *server) SendLog(ctx context.Context, in *api.Log) (*api.Ack, error) {
	// TODO send message to log management
	log.Printf("[%d] %s:: %s", in.Timestamp, in.Identifier, in.Log)
	return &api.Ack{}, nil
}

// Listen with gRPG service
func listen(addrPort string) error {
	// open tcp socket
	lis, err := net.Listen("tcp", addrPort)
	if err != nil {
		grpclog.Fatalf("Failed to open tcp socket: %v", err)
		return err
	}
	log.Printf("Server listening on %s", addrPort)

	// bootstrap gRPC service !
	grpcServer := grpc.NewServer()
	api.RegisterDemonstratorServer(grpcServer, &server{})
	err = grpcServer.Serve(lis)

	return err
}
