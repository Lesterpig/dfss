package sign

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"strconv"

	"dfss"
	cAPI "dfss/dfssc/api"
	"dfss/dfssc/security"
	pAPI "dfss/dfssp/api"
	"dfss/dfssp/contract"
	"dfss/net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// SignatureManager handles the signature of a contract.
type SignatureManager struct {
	auth      *security.AuthContainer
	localPort int
	contract  *contract.JSON
	platform  pAPI.PlatformClient
	peers     map[string]*cAPI.ClientClient
	cServer   *grpc.Server
	cert, ca  *x509.Certificate
	key       *rsa.PrivateKey
}

// NewSignatureManager populates a SignatureManager and connects to the platform.
func NewSignatureManager(fileCA, fileCert, fileKey, addrPort, passphrase string, port int, c *contract.JSON) (*SignatureManager, error) {
	m := &SignatureManager{
		auth:      security.NewAuthContainer(fileCA, fileCert, fileKey, addrPort, passphrase),
		localPort: port,
		contract:  c,
	}
	var err error
	m.ca, m.cert, m.key, err = m.auth.LoadFiles()
	if err != nil {
		return nil, err
	}

	m.cServer = m.GetServer()
	go func() { _ = net.Listen("0.0.0.0:"+strconv.Itoa(port), m.cServer) }()

	conn, err := net.Connect(m.auth.AddrPort, m.cert, m.key, m.ca)
	if err != nil {
		return nil, err
	}

	m.platform = pAPI.NewPlatformClient(conn)

	m.peers = make(map[string]*cAPI.ClientClient)
	for _, u := range c.Signers {
		m.peers[u.Email] = nil
	}

	return m, nil
}

// ConnectToPeers tries to fetch the list of users for this contract, and tries to establish a connection to each peer.
func (m *SignatureManager) ConnectToPeers() error {
	stream, err := m.platform.JoinSignature(context.Background(), &pAPI.JoinSignatureRequest{
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
		if errorCode.Code != pAPI.ErrorCode_SUCCESS {
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
func (m *SignatureManager) addPeer(user *pAPI.User) (ready bool, err error) {
	if user == nil {
		err = errors.New("unexpected user format")
		return
	}

	addrPort := user.Ip + ":" + strconv.Itoa(int(user.Port))
	fmt.Println("Trying to connect with", user.Email, "/", addrPort)

	conn, err := net.Connect(addrPort, m.cert, m.key, m.ca)
	if err != nil {
		return false, err
	}

	// Sending Hello message
	client := cAPI.NewClientClient(conn)
	m.peers[user.Email] = &client
	msg, err := client.Discover(context.Background(), &cAPI.Hello{Version: dfss.Version})
	if err != nil {
		return false, err
	}
	// Printing answer: application version
	fmt.Println("Recieved:", msg.Version)

	// Check if we have any other peer to connect to
	for _, u := range m.peers {
		if u == nil {
			return false, nil
		}
	}

	return true, nil
}
