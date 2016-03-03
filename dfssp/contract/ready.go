package contract

import (
	"time"

	"dfss/dfssp/api"
	"dfss/dfssp/common"
	"dfss/dfssp/entities"
	"dfss/mgdb"
	"dfss/net"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

// readySignal is the structure that is transmitted accross goroutines
type readySignal struct {
	ready    bool     // If true, this is the ready signal. If not, this is a new connection signal
	data     string   // Various data (CN or SignatureUUID)
	chain    [][]byte // Only used to broadcast hash chain (signers hashes in order)
	sequence []uint32 // Only used to broadcast signature sequence
}

// ReadySign is the last job of the platform before the signature can occur.
// When a new client is ready, it joins a waitingGroup a waits for a master broadcast announcing that everybody is ready.
//
// Doing it this way is efficient in time, as only one goroutine deals with the database and do global checks.
func ReadySign(db *mgdb.MongoManager, rooms *common.WaitingGroupMap, ctx *context.Context, in *api.ReadySignRequest) *api.LaunchSignature {
	roomID := "ready_" + in.ContractUuid
	channel, _, first := rooms.Join(roomID)
	cn := net.GetCN(ctx)

	// Check UUID
	if !bson.IsObjectIdHex(in.ContractUuid) {
		return &api.LaunchSignature{ErrorCode: &api.ErrorCode{Code: api.ErrorCode_INVARG}}
	}

	// If first in the room, create a goroutine for ready check.
	// It is absolutely thread safe thanks to a mutex applied on the `first` variable.
	if first {
		go masterReadyRoutine(db, rooms, in.ContractUuid)
	}

	// Broadcast identity
	rooms.Broadcast(roomID, &readySignal{data: cn})

	// Wait for ready signal
	for {
		select {
		case signal, ok := <-channel:
			if !ok {
				return &api.LaunchSignature{ErrorCode: &api.ErrorCode{Code: api.ErrorCode_INTERR}}
			}
			s := signal.(*readySignal)
			if s.ready {
				if len(s.data) > 0 {
					return &api.LaunchSignature{
						ErrorCode:     &api.ErrorCode{Code: api.ErrorCode_SUCCESS},
						SignatureUuid: s.data,
						KeyHash:       s.chain,
					}
				} // data == "" means the contractUUID is bad
				return &api.LaunchSignature{ErrorCode: &api.ErrorCode{Code: api.ErrorCode_INVARG}}
			}
		case <-(*ctx).Done():
			rooms.Unjoin(roomID, channel)
			return nil
		case <-time.After(10 * time.Minute):
			rooms.Unjoin(roomID, channel)
			return &api.LaunchSignature{ErrorCode: &api.ErrorCode{Code: api.ErrorCode_INTERR, Message: "timeout"}}
		}
	}

}

// masterReadyRoutine is a function to be started by the first signer ready as a goroutine.
// It will join the associated ready room and check ready status of each signer when a new signer signals its readiness.
func masterReadyRoutine(db *mgdb.MongoManager, rooms *common.WaitingGroupMap, contractUUID string) {
	roomID := "ready_" + contractUUID
	channel, oldMessages, _ := rooms.Join(roomID)

	// Push oldMessages into the channel.
	// It is safe as this sould be a very small slice (the room is just created).
	for _, v := range oldMessages {
		channel <- v
	}

	// Get contract signers from database
	fetch := entities.Contract{ID: bson.ObjectIdHex(contractUUID)}
	contract := entities.Contract{}
	err := db.Get("contracts").FindByID(fetch, &contract)
	if err != nil {
		rooms.Broadcast(roomID, &readySignal{
			ready: true,
			data:  "",
		}) // This represents a "error" response
		rooms.Unjoin(roomID, channel)
		return
	}

	signersReady := make([]bool, len(contract.Signers))
	work := true
	for work {
		select {
		case signal, ok := <-channel:
			if !ok { // Channel closed, aborting everything
				return
			}
			cn := signal.(*readySignal).data
			ready := FindAndUpdatePendingSigner(cn, &signersReady, &contract.Signers)
			if ready {
				rooms.Broadcast(roomID, &readySignal{
					ready:    true,
					data:     bson.NewObjectId().Hex(),
					chain:    contract.GetHashChain(),
					sequence: GenerateSignSequence(len(contract.Signers)),
				})
				work = false
			}
		case <-time.After(10 * time.Minute):
			work = false
		}
	}

	rooms.Unjoin(roomID, channel)
}

// FindAndUpdatePendingSigner is a utility function to return the state of current signers readiness.
// It has absolutely no interaction with the database.
func FindAndUpdatePendingSigner(mail string, signersReady *[]bool, signers *[]entities.Signer) (ready bool) {
	// Find an update ready status
	for i, s := range *signers {
		if s.Email == mail {
			(*signersReady)[i] = true
			break
		}
	}

	// Check if everyone is ready
	for _, s := range *signersReady {
		if s == false {
			return
		}
	}

	ready = true
	return
}
