package sign

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"dfss"
	cAPI "dfss/dfssc/api"
	"dfss/dfssc/common"
	"dfss/dfssc/security"
	dAPI "dfss/dfssd/api"
	pAPI "dfss/dfssp/api"
	"dfss/dfssp/contract"
	"dfss/net"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Limit the buffer size of the channels
const chanBufferSize = 100

// SignatureManager handles the signature of a contract.
type SignatureManager struct {
	auth         *security.AuthContainer
	contract     *contract.JSON // contains the contractUUID, the list of the signers' hashes, the hash of the contract
	platform     pAPI.PlatformClient
	platformConn *grpc.ClientConn
	peersConn    map[string]*grpc.ClientConn
	peers        map[string]*cAPI.ClientClient
	hashToID     map[string]uint32
	nbReady      int
	cServer      *grpc.Server
	cServerIface clientServer
	sequence     []uint32
	currentIndex int
	uuid         string
	keyHash      [][]byte
	mail         string
	archives     *Archives

	// Callbacks
	OnSignerStatusUpdate func(mail string, status SignerStatus, data string)
	OnProgressUpdate     func(current int, end int)
}

// Archives stores the received and sent messages, as evidence if needed
type Archives struct {
	sentPromises       []*cAPI.Promise
	receivedPromises   []*cAPI.Promise
	sentSignatures     []*cAPI.Signature
	receivedSignatures []*cAPI.Signature
	mutex              sync.Mutex
}

// NewSignatureManager populates a SignatureManager and connects to the platform.
func NewSignatureManager(passphrase string, c *contract.JSON) (*SignatureManager, error) {
	m := &SignatureManager{
		auth:     security.NewAuthContainer(passphrase),
		contract: c,
		archives: &Archives{
			sentPromises:       make([]*cAPI.Promise, 0),
			receivedPromises:   make([]*cAPI.Promise, 0),
			sentSignatures:     make([]*cAPI.Signature, 0),
			receivedSignatures: make([]*cAPI.Signature, 0),
		},
	}
	var err error
	_, _, _, err = m.auth.LoadFiles()
	if err != nil {
		return nil, err
	}

	m.mail = m.auth.Cert.Subject.CommonName
	dAPI.SetIdentifier(m.mail)

	m.cServer = m.GetServer()
	go func() { log.Fatalln(net.Listen("0.0.0.0:"+strconv.Itoa(viper.GetInt("local_port")), m.cServer)) }()

	conn, err := net.Connect(viper.GetString("platform_addrport"), m.auth.Cert, m.auth.Key, m.auth.CA, nil)
	if err != nil {
		return nil, err
	}

	m.platform = pAPI.NewPlatformClient(conn)
	m.platformConn = conn

	m.peersConn = make(map[string]*grpc.ClientConn)
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
		Port:         uint32(viper.GetInt("local_port")),
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
	m.OnSignerStatusUpdate(user.Email, StatusConnecting, addrPort)

	// This is an certificate authentificated TLS connection
	conn, err := net.Connect(addrPort, m.auth.Cert, m.auth.Key, m.auth.CA, user.KeyHash)
	if err != nil {
		m.OnSignerStatusUpdate(user.Email, StatusError, err.Error())
		return false, err
	}

	// Sending Hello message
	client := cAPI.NewClientClient(conn)
	lastConnection := m.peers[user.Email]
	m.peers[user.Email] = &client
	// The connection is encapsulated into the interface, so we
	// need to create another way to access it
	m.peersConn[user.Email] = conn

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	msg, err := client.Discover(ctx, &cAPI.Hello{Version: dfss.Version})
	if err != nil {
		m.OnSignerStatusUpdate(user.Email, StatusError, err.Error())
		return false, err
	}

	// Printing answer: application version
	// TODO check certificate
	m.OnSignerStatusUpdate(user.Email, StatusConnected, msg.Version)

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
	m.keyHash = launch.KeyHash
	signatureUUID = m.uuid
	return
}

// Initialize computes the values needed for the start of the signing
func (m *SignatureManager) Initialize() (uint32, int, error) {
	myID, err := m.FindID()
	if err != nil {
		return 0, 0, err
	}

	m.currentIndex, err = common.FindNextIndex(m.sequence, myID, -1)
	if err != nil {
		return 0, 0, err
	}

	nextIndex, err := common.FindNextIndex(m.sequence, myID, m.currentIndex)
	if err != nil {
		return 0, 0, err
	}

	return myID, nextIndex, nil
}

// FindID finds the sequence id for the user's email and the contract to sign
func (m *SignatureManager) FindID() (uint32, error) {
	signers := m.contract.Signers
	for id, signer := range signers {
		if signer.Email == m.mail {
			return uint32(id), nil
		}
	}
	return 0, errors.New("Mail couldn't be found amongst signers")
}
