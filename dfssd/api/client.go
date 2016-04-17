package api

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	address, identifier string
	demo                bool
	// lazy initializer
	dial       *grpc.ClientConn
	demoClient DemonstratorClient
)

// Configure is used to update current parameters.
// Call it at least one time before the first DLog call.
func Configure(activated bool, addrport, id string) {
	address = addrport
	identifier = id
	demo = activated
}

// SetIdentifier updates the current client identifier.
func SetIdentifier(id string) {
	identifier = id
}

// Lazy initialisation for demonstrator's connection to server
func dInit() error {
	var err error
	dial, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		grpclog.Printf("Fail to dial: %v", err)
		return err
	}

	demoClient = NewDemonstratorClient(dial)
	return nil
}

// DClose close the connection with demonstrator server (if any)
//
// This should be called at the end of any program that import this library
func DClose() {
	if dial != nil {
		err := dial.Close()
		if err != nil {
			grpclog.Printf("Failed to close dialing with demonstrator: %v", err)
		}
	}
}

// DLog send a message to the demonstrator
//
// The client is dialed in a lazy way
func DLog(log string) {
	// check demo switch
	if !demo {
		return
	}

	// lazy initialisation
	if dial == nil {
		err := dInit()
		if err != nil {
			return // fail silently
		}
	}
	_, err := demoClient.SendLog(
		context.Background(),
		&Log{Timestamp: time.Now().UnixNano(), Identifier: identifier, Log: log})

	if err != nil {
		grpclog.Printf("Fail to send message: %v", err)
	}
}
