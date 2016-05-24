package sign

import (
	"errors"
	"fmt"
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
	tAPI "dfss/dfsst/api"
	"dfss/dfsst/entities"
	"dfss/net"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Limit the buffer size of the channels
const chanBufferSize = 100

// SignatureManager handles the signature of a contract.
type SignatureManager struct {
	auth           *security.AuthContainer
	contract       *contract.JSON // contains the contractUUID, the list of the signers' hashes, the hash of the contract
	platform       pAPI.PlatformClient
	platformConn   *grpc.ClientConn
	ttp            tAPI.TTPClient
	ttpData        *pAPI.LaunchSignature_TTP
	peersConn      map[string]*grpc.ClientConn
	peers          map[string]*cAPI.ClientClient
	hashToID       map[string]uint32
	nbReady        int
	cServer        *grpc.Server
	cServerIface   clientServer
	sequence       []uint32
	lastValidIndex int // the last index at which we sent a promise
	currentIndex   int
	myID           uint32
	uuid           string
	keyHash        [][]byte
	mail           string
	archives       *Archives
	seal           []byte
	cancelled      bool
	finished       bool

	// Callbacks
	OnSignerStatusUpdate func(mail string, status SignerStatus, data string)
	OnProgressUpdate     func(current int, end int)
	Cancel               chan interface{}
}

// Archives stores the received and sent messages, as evidence if needed
type Archives struct {
	receivedPromises   []*cAPI.Promise // TODO: improve by using a map
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
			receivedPromises:   make([]*cAPI.Promise, 0),
			sentSignatures:     make([]*cAPI.Signature, 0),
			receivedSignatures: make([]*cAPI.Signature, 0),
		},
		Cancel: make(chan interface{}),
	}
	var err error
	_, _, _, err = m.auth.LoadFiles()
	if err != nil {
		return nil, err
	}

	m.mail = m.auth.Cert.Subject.CommonName
	dAPI.SetIdentifier(m.mail)

	net.DefaultTimeout = viper.GetDuration("timeout")

	m.cServer = m.GetServer()
	go func() { _ = net.Listen("0.0.0.0:"+strconv.Itoa(viper.GetInt("local_port")), m.cServer) }()

	connp, err := net.Connect(viper.GetString("platform_addrport"), m.auth.Cert, m.auth.Key, m.auth.CA, nil)
	if err != nil {
		return nil, err
	}

	m.platform = pAPI.NewPlatformClient(connp)
	m.platformConn = connp

	m.peersConn = make(map[string]*grpc.ClientConn)
	m.peers = make(map[string]*cAPI.ClientClient)
	for _, u := range c.Signers {
		if u.Email != m.auth.Cert.Subject.CommonName {
			m.peers[u.Email] = nil
		}
	}

	// Initialize TTP AuthContainer
	// This is needed to use platform seal verification client-side
	entities.AuthContainer = m.auth

	return m, nil
}

// ConnectToPeers tries to fetch the list of users for this contract, and tries to establish a connection to each peer.
func (m *SignatureManager) ConnectToPeers() error {
	localIps, err := net.ExternalInterfaceAddr()
	if err != nil {
		return err
	}

	stream, err := m.platform.JoinSignature(context.Background(), &pAPI.JoinSignatureRequest{
		ContractUuid: m.contract.UUID,
		Port:         uint32(viper.GetInt("local_port")),
		Ip:           localIps,
	})
	if err != nil {
		m.finished = true
		m.closeConnections()
		return err
	}

	c := make(chan error)
	go connectToPeersLoop(m, stream, c)

	select {
	case err = <-c:
		if err != nil {
			m.finished = true
			m.closeConnections()
		}
		return err
	case <-m.Cancel:
		m.cancelled = true
		m.closeConnections()
		return errors.New("Signature cancelled")
	}
}

func connectToPeersLoop(m *SignatureManager, stream pAPI.Platform_JoinSignatureClient, c chan error) {
	for !m.cancelled {
		userConnected, err := stream.Recv()
		if err != nil {
			c <- err
			return
		}
		errorCode := userConnected.GetErrorCode()
		if errorCode.Code != pAPI.ErrorCode_SUCCESS {
			c <- errors.New(errorCode.Message)
			return
		}
		ready, err := m.addPeer(userConnected.User)
		if err != nil {
			continue // Unable to connect to this user, ignore it for the moment
		}
		if ready {
			c <- nil
			return
		}
	}
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

	var conn *grpc.ClientConn
	for _, ip := range user.Ip {
		addrPort := ip + ":" + strconv.Itoa(int(user.Port))
		m.OnSignerStatusUpdate(user.Email, StatusConnecting, addrPort)

		// This is an certificate authentificated TLS connection
		conn, err = net.Connect(addrPort, m.auth.Cert, m.auth.Key, m.auth.CA, user.KeyHash)
		if err == nil {
			break
		}

		if m.cancelled {
			err = errors.New("Signature cancelled")
			break
		}
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
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

	c := make(chan *pAPI.LaunchSignature)
	go func() {
		launch, _ := m.platform.ReadySign(ctx, &pAPI.ReadySignRequest{
			ContractUuid: m.contract.UUID,
		})
		c <- launch
	}()

	var launch *pAPI.LaunchSignature
	select {
	case launch = <-c: // OK
	case <-m.Cancel:
		m.cancelled = true
		m.closeConnections()
		err = errors.New("Signature cancelled")
		return
	}

	errorCode := launch.GetErrorCode()
	if errorCode.Code != pAPI.ErrorCode_SUCCESS {
		err = errors.New(errorCode.Code.String() + " " + errorCode.Message)
		m.finished = true
		m.closeConnections()
		return
	}

	// Check signers from platform data
	if len(m.contract.Signers) != len(launch.KeyHash) {
		err = errors.New("Corrupted DFSS file: bad number of signers, unable to sign safely")
		m.finished = true
		m.closeConnections()
		return
	}

	for i, s := range m.contract.Signers {
		if s.Hash != fmt.Sprintf("%x", launch.KeyHash[i]) {
			err = errors.New("Corrupted DFSS file: signer " + s.Email + " has an invalid hash, unable to sign safely")
			m.finished = true
			m.closeConnections()
			return
		}
	}

	// Connect to TTP, if any
	err = m.connectToTTP(launch.Ttp)

	m.sequence = launch.Sequence
	m.uuid = launch.SignatureUuid
	m.keyHash = launch.KeyHash
	m.seal = launch.Seal
	signatureUUID = m.uuid
	return
}

// connectToTTP : tries to open a connection with the ttp specified in the contract.
func (m *SignatureManager) connectToTTP(ttp *pAPI.LaunchSignature_TTP) error {
	if ttp == nil {
		m.ttpData = &pAPI.LaunchSignature_TTP{
			Addrport: "",
			Hash:     []byte{},
		}
		return nil
	}

	// TODO check that the connection spots missing TTP and returns an error quickly enough
	conn, err := net.Connect(ttp.Addrport, m.auth.Cert, m.auth.Key, m.auth.CA, ttp.Hash)
	if err != nil {
		return err
	}

	m.ttpData = ttp
	m.ttp = tAPI.NewTTPClient(conn)
	return nil
}

// Initialize computes the values needed for the start of the signing
func (m *SignatureManager) Initialize() (int, error) {
	myID, err := m.FindID()
	if err != nil {
		return 0, err
	}
	m.myID = myID

	m.currentIndex, err = common.FindNextIndex(m.sequence, myID, -1)
	if err != nil {
		return 0, err
	}

	nextIndex, err := common.FindNextIndex(m.sequence, myID, m.currentIndex)
	if err != nil {
		return 0, err
	}

	m.lastValidIndex = 0

	return nextIndex, nil
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

// IsTerminated returns true if the signature is cancelled or finished
func (m *SignatureManager) IsTerminated() bool {
	return m.cancelled || m.finished
}
