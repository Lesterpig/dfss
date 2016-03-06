package sign

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

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
	nbReady   int
	cServer   *grpc.Server
	sequence  []uint32
	uuid      string
}

// NewSignatureManager populates a SignatureManager and connects to the platform.
func NewSignatureManager(fileCA, fileCert, fileKey, addrPort, passphrase string, port int, c *contract.JSON) (*SignatureManager, error) {
	m := &SignatureManager{
		auth:      security.NewAuthContainer(fileCA, fileCert, fileKey, addrPort, passphrase),
		localPort: port,
		contract:  c,
	}
	var err error
	_, _, _, err = m.auth.LoadFiles()
	if err != nil {
		return nil, err
	}

	m.cServer = m.GetServer()
	go func() { log.Fatalln(net.Listen("0.0.0.0:"+strconv.Itoa(port), m.cServer)) }()

	conn, err := net.Connect(m.auth.AddrPort, m.auth.Cert, m.auth.Key, m.auth.CA)
	if err != nil {
		return nil, err
	}

	m.platform = pAPI.NewPlatformClient(conn)

	m.peers = make(map[string]*cAPI.ClientClient)
	for _, u := range c.Signers {
		if u.Email != m.auth.Cert.Subject.CommonName {
			m.peers[u.Email] = nil
		}
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
			continue // Unable to connect to this user, ignore it for the moment
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
	if _, ok := m.peers[user.Email]; !ok {
		return // Ignore if unknown
	}

	addrPort := user.Ip + ":" + strconv.Itoa(int(user.Port))
	fmt.Println("- Trying to connect with", user.Email, "/", addrPort)

	conn, err := net.Connect(addrPort, m.auth.Cert, m.auth.Key, m.auth.CA)
	if err != nil {
		return false, err
	}

	// Sending Hello message
	client := cAPI.NewClientClient(conn)
	lastConnection := m.peers[user.Email]
	m.peers[user.Email] = &client

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	msg, err := client.Discover(ctx, &cAPI.Hello{Version: dfss.Version})
	if err != nil {
		return false, err
	}

	// Printing answer: application version
	// TODO check certificate
	fmt.Println("  Successfully connected!", "[", msg.Version, "]")

	// Check if we have any other peer to connect to
	if lastConnection == nil {
		m.nbReady++
		if m.nbReady == len(m.contract.Signers)-1 {
			return true, nil
		}
	}

	return false, nil
}

// SendReadySign sends the READY signal to the platform, and wait (potentially a long time) for START signal.
func (m *SignatureManager) SendReadySign() (signatureUUID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	launch, err := m.platform.ReadySign(ctx, &pAPI.ReadySignRequest{
		ContractUuid: m.contract.UUID,
	})
	if err != nil {
		return
	}

	errorCode := launch.GetErrorCode()
	if errorCode.Code != pAPI.ErrorCode_SUCCESS {
		err = errors.New(errorCode.Code.String() + " " + errorCode.Message)
		return
	}

	// Check signers from platform data
	if len(m.contract.Signers) != len(launch.KeyHash) {
		err = errors.New("Corrupted DFSS file: bad number of signers, unable to sign safely")
		return
	}

	for i, s := range m.contract.Signers {
		if s.Hash != fmt.Sprintf("%x", launch.KeyHash[i]) {
			err = errors.New("Corrupted DFSS file: signer " + s.Email + " has an invalid hash, unable to sign safely")
			return
		}
	}

	m.sequence = launch.Sequence
	m.uuid = launch.SignatureUuid
	signatureUUID = m.uuid
	return
}