package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"net"

	api "dfss/dfssd/api"
)

type server struct{}

// Sendlog Handler
//
// Handle incoming log messages
func (s *server) SendLog(ctx context.Context, in *api.Log) (*api.Ack, error) {
	addMessage(in)
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

	// log display manager
	go displayHandler()

	// bootstrap gRPC service !
	grpcServer := grpc.NewServer()
	api.RegisterDemonstratorServer(grpcServer, &server{})
	err = grpcServer.Serve(lis)

	return err
}
