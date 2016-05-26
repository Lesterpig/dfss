package user

import (
	"errors"

	pb "dfss/dfssp/api"
	"dfss/net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Unregister a user from the platform
func Unregister() error {
	client, err := connect()
	if err != nil {
		return err
	}

	// Stop the context if it takes too long for the platform to answer
	ctx, cancel := context.WithTimeout(context.Background(), net.DefaultTimeout)
	defer cancel()
	response, err := client.Unregister(ctx, &pb.Empty{})
	if err != nil {
		return errors.New(grpc.ErrorDesc(err))
	}
	if response.Code != pb.ErrorCode_SUCCESS {
		return errors.New(response.Message)
	}

	return nil
}
