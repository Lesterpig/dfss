package contract

import (
	n "net"

	"dfss/auth"
	"dfss/dfssp/api"
	"dfss/dfssp/common"
	"dfss/dfssp/entities"
	"dfss/mgdb"
	"dfss/net"
	"gopkg.in/mgo.v2/bson"
)

// JoinSignature allows a client to wait for other clients connections on a specific contract.
// Firstly, every client present BEFORE the call of this function is sent to the stream.
// Then, client information is sent to the stream as it's available.
//
// Please note that the current user will also receive it's own information.
// There is no timeout, this function will shut down on stream disconnection or on error.
func JoinSignature(db *mgdb.MongoManager, rooms *common.WaitingGroupMap, in *api.JoinSignatureRequest, stream api.Platform_JoinSignatureServer) {
	repository := entities.NewContractRepository(db.Get("contracts"))
	ctx := stream.Context()
	state, addr, _ := net.GetTLSState(&ctx)
	hash := auth.GetCertificateHash(state.VerifiedChains[0][0])

	if !bson.IsObjectIdHex(in.ContractUuid) {
		_ = stream.Send(&api.UserConnected{
			ErrorCode: &api.ErrorCode{
				Code:    api.ErrorCode_INVARG,
				Message: "invalid contract uuid",
			},
		})
		return
	}

	if !repository.CheckAuthorization(hash, bson.ObjectIdHex(in.ContractUuid)) {
		_ = stream.Send(&api.UserConnected{
			ErrorCode: &api.ErrorCode{
				Code:    api.ErrorCode_INVARG,
				Message: "unauthorized signature",
			},
		})
		return
	}

	// Join room
	roomID := "contract_" + in.ContractUuid
	channel, pendingSigners := rooms.Join(roomID)

	// Send pendingSigners
	for _, p := range pendingSigners {
		err := sendUserToStream(&stream, in.ContractUuid, p.(*api.User))
		if err != nil {
			rooms.Unjoin(roomID, channel)
			return
		}
	}

	// Broadcast self identity
	host, _, _ := n.SplitHostPort(addr.String())
	rooms.Broadcast(roomID, &api.User{
		KeyHash: hash,
		Email:   net.GetCN(&ctx),
		Ip:      host,
		Port:    in.Port,
	})

	// Listen for others
	for {
		select {
		case user, ok := <-channel:
			if !ok { // channel is closed, means that the room is closed
				return
			}
			err := sendUserToStream(&stream, in.ContractUuid, user.(*api.User))
			if err != nil {
				rooms.Unjoin(roomID, channel)
				return
			}
		case <-ctx.Done():
			rooms.Unjoin(roomID, channel)
			return
		}
	}
}

func sendUserToStream(stream *api.Platform_JoinSignatureServer, contractUUID string, user *api.User) error {
	return (*stream).Send(&api.UserConnected{
		ErrorCode:    &api.ErrorCode{Code: api.ErrorCode_SUCCESS},
		ContractUuid: contractUUID,
		User:         user,
	})
}
