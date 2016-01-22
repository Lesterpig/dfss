package api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"time"
)

var (
	// theses are the default parameters
	address    = "localhost:3000"
	identifier = "platform"
	demo       = false
	// lazy initializer
	dial       *grpc.ClientConn
	demoClient DemonstratorClient
)

// Switch for demo mode
//
// Should be used to pass the value of `-d` switch
func Switch(activationSwitch bool) {
	demo = activationSwitch
	return
}

// Lazy initialisation for demonstrator's connection to server
func dInit() error {
	var err error
	dial, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		grpclog.Printf("Fail to dial: %v", err)
	}

	demoClient = NewDemonstratorClient(dial)

	return err
}

// DClose close the connection with demonstrator server (if any)
//
// This should be called at the end of any program that import this library
func DClose() {
	if dial != nil {
		err := dial.Close()
		if err != nil {
			grpclog.Printf("Fail to close dialing: %v", err)
		}
	}
}

// DLog send a message to the demonstrator
//
// The client is dialed in a lazy way
// The default demonstrator server address is localhost:3000
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
