// Package server provides the ttp server.
package server

import (
	"errors"
	"fmt"
	"os"
	"sync"

	cAPI "dfss/dfssc/api"
	"dfss/dfssc/security"
	dAPI "dfss/dfssd/api"
	tAPI "dfss/dfsst/api"
	"dfss/dfsst/entities"
	"dfss/dfsst/resolve"
	"dfss/mgdb"
	"dfss/net"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2/bson"
)

// InternalError : constant string used to return a generic error message through gRPC in case of an internal error.
const InternalError string = "Internal server error"

type ttpServer struct {
	DB        *mgdb.MongoManager
	globalMut *sync.Mutex
	mutMap    map[bson.ObjectId]*sync.Mutex
}

// Alert route for the TTP.
func (server *ttpServer) Alert(ctx context.Context, in *tAPI.AlertRequest) (*tAPI.TTPResponse, error) {
	valid, signatureUUID, signers, senderIndex := entities.IsRequestValid(ctx, in.Promises)
	if !valid {
		dAPI.DLog("invalid request from " + net.GetCN(&ctx))
		return nil, errors.New(InternalError)
	}

	server.globalMut.Lock()

	_, ok := server.mutMap[signatureUUID]
	if !ok {
		server.mutMap[signatureUUID] = &sync.Mutex{}
	}
	server.mutMap[signatureUUID].Lock()

	server.globalMut.Unlock()
	defer server.mutMap[signatureUUID].Unlock()

	dAPI.DLog("resolve index is: " + fmt.Sprint(in.Index))
	valid = int(in.Index) < len(in.Promises[0].Context.Sequence)
	if !valid {
		dAPI.DLog("invalid sequence index from " + net.GetCN(&ctx))
		return nil, errors.New(InternalError)
	}
	// Now we know that the request contains information correctly signed by the platform,
	// with the same signatureUUID (thus signed information) for all promises, and sent by a valid signer
	// wrt to the signed signers' hashes

	dAPI.DLog("Resolve request from " + net.GetCN(&ctx) + " is valid")

	manager := entities.NewArchivesManager(server.DB)
	err := manager.InitializeArchives(in.Promises[0], signatureUUID, &signers)
	if err != nil {
		dAPI.DLog("error occured during the initialization of the signature archives")
		return nil, err
	}
	// Now archives contains the new or already present SignatureArchives

	// We check if we have already sent an abort token to the sender of the request
	stop, message, err := server.handleAbortedSender(manager, senderIndex)
	if stop {
		dAPI.DLog("already sent an abort token to " + net.GetCN(&ctx))
		return message, err
	}

	// We check that the sender of the request sent valid and complete information
	stop, message, tmpPromises, err := server.handleInvalidPromises(manager, in.Promises, senderIndex, in.Index)
	if stop {
		dAPI.DLog("sent abort token to " + net.GetCN(&ctx))
		return message, err
	}
	// Now we are sure that the sender of the AlertRequest is not dishonest

	// We try to use the already generated contract if it exists
	generated, contract := manager.WasContractSigned()
	if generated {
		dAPI.DLog("sent signed contract to " + net.GetCN(&ctx))
		return &tAPI.TTPResponse{
			Abort:    false,
			Contract: contract,
		}, nil
	}

	// If we didn't already generate the signed contract, we take into account the new promises
	// Computing the dishonest signers wrt to the new evidence
	server.updateArchiveWithEvidence(manager, tmpPromises)
	// Try to generate the contract now
	message, err = server.handleContractGenerationTry(manager)
	// We manually update the database
	ok, err = server.DB.Get("signatures").UpdateByID(*(manager.Archives))
	if !ok {
		fmt.Fprintln(os.Stderr, err)
		return nil, errors.New(InternalError)
	}

	if message.Abort {
		dAPI.DLog("sent abort token to " + net.GetCN(&ctx))
	} else {
		dAPI.DLog("sent signed contract to " + net.GetCN(&ctx))
	}

	return message, err
}

// handleAbortedSender : if the specified signer has already recieved an abort token, adds him to the dishonest signers of the specified
// signatureArchives, and returns a boolean that states if we should stop the execution of the resolve protocol, and the response that should be sent back to him.
//
// Updates the database with the new aborted signers.
// If an error occurs during this process, it is returned.
func (server *ttpServer) handleAbortedSender(manager *entities.ArchivesManager, senderIndex uint32) (bool, *tAPI.TTPResponse, error) {
	if manager.HasReceivedAbortToken(senderIndex) {
		dAPI.DLog("Sender has already contacted the ttp. He is dishonnest.")
		manager.AddToDishonest(senderIndex)

		ok, err := manager.DB.Get("signatures").UpdateByID(*(manager.Archives))
		if !ok {
			fmt.Fprintln(os.Stderr, err)
			return true, nil, errors.New(InternalError)
		}

		return true, &tAPI.TTPResponse{
			Abort:    true,
			Contract: nil,
		}, nil
	}
	dAPI.DLog("sender has never contacted the ttp before")
	return false, nil, nil
}

// handleInvalidPromises : if the specified signer has sent us a valid request, but with invalid promises, ie:
// - he sent us impossible information wrt the information signed by the platform
// OR
// - he sent not enough information wrt the information signed by the platform and the signing protocol
// then he is added to the dishonest signers.
//
// Returns a boolean that states if we should stop the execution of the resolve protocol, and the response that should be sent back to him.
// If the promises are valid, return them in the simplified form of an array of *entities.Promise
//
// Updates the database with the new aborted signer.
// If an error occurs during this process, it is returned.
func (server *ttpServer) handleInvalidPromises(manager *entities.ArchivesManager, promises []*cAPI.Promise, senderIndex, stepIndex uint32) (bool, *tAPI.TTPResponse, []*entities.Promise, error) {
	valid, tmpPromises := entities.ArePromisesValid(promises)
	if valid {
		dAPI.DLog("received promises are valid")
	}
	complete := resolve.ArePromisesComplete(tmpPromises, promises[0], stepIndex)
	if complete {
		dAPI.DLog("received promises are complete")
	}
	if !valid || !complete {
		if !valid {
			dAPI.DLog("received promises are not valid")
		}
		if !complete {
			dAPI.DLog("received promises are not complete")
		}
		manager.AddToAbort(senderIndex)
		manager.AddToDishonest(senderIndex)

		ok, err := manager.DB.Get("signatures").UpdateByID(*(manager.Archives))
		if !ok {
			fmt.Fprintln(os.Stderr, err)
			return true, nil, nil, errors.New(InternalError)
		}

		return true, &tAPI.TTPResponse{
			Abort:    true,
			Contract: nil,
		}, nil, nil
	}

	return false, nil, tmpPromises, nil
}

// updateArchiveWithEvidence : computes the dishonest signers from the new provided evidence, and updates the specified signatureArchives accordingly.
//
// DOES NOT UPDATE THE DATABASE (should be handled manually)
func (server *ttpServer) updateArchiveWithEvidence(manager *entities.ArchivesManager, tmpPromises []*entities.Promise) {
	computedDishonest := resolve.ComputeDishonestSigners(manager.Archives, tmpPromises)

	for _, di := range computedDishonest {
		manager.AddToDishonest(di)
	}

	for _, p := range tmpPromises {
		manager.AddPromise(p)
	}
}

// handleContractGenerationTry : tries to generate the signed contract from the specified signatureArchives.
// Returns the response to send back to the sender of the request.
//
// Does not take into account if the sender is dishonest.
// If the contract has been successfully generated, returns it. Otherwise, returns an abort token.
//
// DOES NOT UPDATE THE DATABASE (should be handled manually)
func (server *ttpServer) handleContractGenerationTry(manager *entities.ArchivesManager) (*tAPI.TTPResponse, error) {
	generated, contract := resolve.Solve(manager)
	if !generated {
		return &tAPI.TTPResponse{
			Abort:    true,
			Contract: nil,
		}, nil
	}

	// We add the generated contract to the signatureArchives
	manager.Archives.SignedContract = contract
	return &tAPI.TTPResponse{
		Abort:    false,
		Contract: contract,
	}, nil
}

// Recover route for the TTP.
func (server *ttpServer) Recover(ctx context.Context, in *tAPI.RecoverRequest) (*tAPI.TTPResponse, error) {
	if !bson.IsObjectIdHex(in.SignatureUUID) {
		return nil, errors.New("Invalid signature uuid.")
	}
	bsonUUID := bson.ObjectIdHex(in.SignatureUUID)

	manager := entities.NewArchivesManager(server.DB)
	present, archives := manager.ContainsSignature(bsonUUID)
	if !present {
		return nil, errors.New("Unknown signature uuid.")
	}
	manager.Archives = archives

	contract, err := handleRecover(ctx, manager)
	if err != nil {
		return nil, err
	}

	return &tAPI.TTPResponse{Contract: contract}, nil
}

func handleRecover(ctx context.Context, manager *entities.ArchivesManager) ([]byte, error) {
	senderHash := net.GetClientHash(&ctx)
	if senderHash == nil {
		return []byte{}, errors.New("Bad authentication.")
	}

	present, senderID := manager.Archives.ContainsSigner(senderHash)
	if !present {
		return []byte{}, errors.New("Signer was not part of the signature.")
	}

	aborted := manager.HasReceivedAbortToken(senderID)
	if aborted {
		return []byte{}, errors.New("Signer was aborted.")
	}

	_, contract := manager.WasContractSigned()
	return contract, nil
}

// GetServer returns the gRPC server.
func GetServer() *grpc.Server {
	// We can do that because NewAuthContainer is looking for "file_ca", "file_cert", and "file_key" in viper, which are set by the TTP
	entities.AuthContainer = security.NewAuthContainer(viper.GetString("password"))
	ca, cert, key, err := entities.AuthContainer.LoadFiles()
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occured during the private key and certificates retrieval:", err)
		os.Exit(1)
	}

	dbManager, err := mgdb.NewManager(viper.GetString("dbURI"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occured during the connection to MongoDB:", err)
		os.Exit(1)
	}

	mutmap := make(map[bson.ObjectId]*sync.Mutex)

	server := &ttpServer{
		DB:        dbManager,
		globalMut: &sync.Mutex{},
		mutMap:    mutmap,
	}

	netServer := net.NewServer(cert, key, ca)
	tAPI.RegisterTTPServer(netServer, server)
	return netServer
}
