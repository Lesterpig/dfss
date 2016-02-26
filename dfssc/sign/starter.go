package sign

import (
	"errors"
	"fmt"
	"strconv"

	"dfss/dfssc/security"
	"dfss/dfssp/api"
	"dfss/dfssp/contract"
	"dfss/net"
	"golang.org/x/net/context"
)

// SignatureManager handles the signature of a contract.
type SignatureManager struct {
	auth      *security.AuthContainer
	localPort int
	contract  *contract.JSON
	platform  api.PlatformClient
	peers     map[string]*api.User
}

// NewSignatureManager populates a SignatureManager and connects to the platform.
func NewSignatureManager(fileCA, fileCert, fileKey, addrPort, passphrase string, port int, c *contract.JSON) (*SignatureManager, error) {
	m := &SignatureManager{
		auth:      security.NewAuthContainer(fileCA, fileCert, fileKey, addrPort, passphrase),
		localPort: port,
		contract:  c,
	}

	ca, cert, key, err := m.auth.LoadFiles()
	if err != nil {
		return nil, err
	}

	conn, err := net.Connect(m.auth.AddrPort, cert, key, ca)
	if err != nil {
		return nil, err
	}

	m.platform = api.NewPlatformClient(conn)

	m.peers = make(map[string]*api.User)
	for _, u := range c.Signers {
		m.peers[u.Email] = nil
	}

	return m, nil
}

// ConnectToPeers tries to fetch the list of users for this contract, and tries to establish a connection to each peer.
func (m *SignatureManager) ConnectToPeers() error {
	stream, err := m.platform.JoinSignature(context.Background(), &api.JoinSignatureRequest{
		ContractUuid: m.contract.UUID,
		Port:         uint32(m.localPort),
	})
	if err != nil {
		return err
	}

	for {
		userConnected, err := stream.Recv()
		if err != nil {
			return err
		}
		errorCode := userConnected.GetErrorCode()
		if errorCode.Code != api.ErrorCode_SUCCESS {
			return errors.New(errorCode.Message)
		}
		ready, err := m.addPeer(userConnected.User)
		if err != nil {
			return err
		}
		if ready {
			break
		}
	}

	return nil
}

// addPeer stores a peer from the platform and tries to establish a connection to this peer.
func (m *SignatureManager) addPeer(user *api.User) (ready bool, err error) {
	if user == nil {
		err = errors.New("unexpected user format")
		return
	}

	m.peers[user.Email] = user
	fmt.Println("Trying to connect with", user.Email, "/", user.Ip+":"+strconv.Itoa(int(user.Port)))

	// TODO do the connection

	// Check if we have any other peer to connect to
	for _, u := range m.peers {
		if u == nil {
			return
		}
	}

	ready = true
	return
}
