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
	ready        bool     // If true, this is the ready signal. If not, this is a new connection signal
	data         string   // Various data (CN or SignatureUUID)
	documentHash []byte   // Contract document SHA-512 hash
	chain        [][]byte // Only used to broadcast hash chain (signers hashes in order)
	sequence     []uint32 // Only used to broadcast signature sequence
}

// ReadySignTimeout is the delay users have to confirm the signature.
// A high value is not recommended, as there is no way to create any other signature on the same contract before the timeout
// if a client has connection issues.
var ReadySignTimeout = time.Minute

// ReadySign is the last job of the platform before the signature can occur.
// When a new client is ready, it joins a waitingGroup and waits for a master broadcast announcing that everybody is ready.
//
// Doing it this way is efficient in time, as only one goroutine deals with the database and do global checks.
func ReadySign(db *mgdb.MongoManager, rooms *common.WaitingGroupMap, ctx *context.Context, in *api.ReadySignRequest) *api.LaunchSignature {
	roomID := "ready_" + in.ContractUuid
	channel, _, first := rooms.Join(roomID)
	defer rooms.Unjoin(roomID, channel)

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
	timeout := time.After(ReadySignTimeout)
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
						DocumentHash:  s.documentHash,
						KeyHash:       s.chain,
						Sequence:      s.sequence,
					}
				} // data == "" means the contractUUID is bad
				return &api.LaunchSignature{ErrorCode: &api.ErrorCode{Code: api.ErrorCode_INVARG}}
			}
		case <-(*ctx).Done(): // Client's disconnection
			return &api.LaunchSignature{ErrorCode: &api.ErrorCode{Code: api.ErrorCode_INVARG}}
		case <-timeout: // Someone has not confirmed the signature within the delay
			return &api.LaunchSignature{ErrorCode: &api.ErrorCode{Code: api.ErrorCode_TIMEOUT, Message: "timeout for ready signal"}}
		}
	}

}

// masterReadyRoutine is a function to be started by the first signer ready as a goroutine.
// It will join the associated ready room and check ready status of each signer when a new signer signals its readiness.
func masterReadyRoutine(db *mgdb.MongoManager, rooms *common.WaitingGroupMap, contractUUID string) {
	roomID := "ready_" + contractUUID
	channel, oldMessages, _ := rooms.Join(roomID)
	defer rooms.Unjoin(roomID, channel)

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
		return
	}

	signersReady := make([]bool, len(contract.Signers))
	work := true
	timeout := time.After(ReadySignTimeout)
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
					ready:        true,
					data:         bson.NewObjectId().Hex(),
					documentHash: contract.File.Hash,
					chain:        contract.GetHashChain(),
					sequence:     GenerateSignSequence(len(contract.Signers)),
				})
				work = false
			}
		case <-timeout:
			work = false
		}
	}

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
		if !s {
			return
		}
	}

	ready = true
	return
}
