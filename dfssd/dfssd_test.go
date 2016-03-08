package main

import (
	dapi "dfss/dfssd/api"
	"testing"
)

// Test for a client-server dialing
//
// Cannot add output statement because output includes timestamp
func TestServerAndClient(t *testing.T) {

	// Start a server
	go func() {
		err := listen("localhost:3000")
		if err != nil {
			t.Error("Unable to start server")
		}
	}()

	// Start a client
	go func() {
		defer dapi.DClose()
		// this one fails silently, so you can't really test it
		dapi.DLog("This is a log message from a client")
	}()

	// Start another client
	go func() {
		defer dapi.DClose()
		// this one fails silently, so you can't really test it
		dapi.DLog("This is a log message from another client")
	}()
}
