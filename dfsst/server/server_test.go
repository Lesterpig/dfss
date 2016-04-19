package server

import (
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	"dfss/auth"
	"dfss/dfsst/entities"
	"dfss/mgdb"

	"gopkg.in/mgo.v2/bson"
)

var (
	fca, fcert, fkey string
	db               string
	platformCA, cert *x509.Certificate
	pkey             *rsa.PrivateKey

	sequence             []uint32
	signers              [][]byte
	contractDocumentHash []byte
	signatureUUID        string
	signatureUUIDBson    bson.ObjectId

	signedHash []byte

	ttp            *ttpServer
	ttpAddressPort string

	err error

	signersEntities []entities.Signer
)

func init() {
	fca = filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "dfsst", "testdata", "dfssp_rootCA.pem")
	CAData, _ := ioutil.ReadFile(fca)

	platformCA, _ = auth.PEMToCertificate(CAData)

	fcert = filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "dfsst", "testdata", "cert.pem")
	CertData, _ := ioutil.ReadFile(fcert)

	cert, _ = auth.PEMToCertificate(CertData)

	fkey = filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "dfsst", "testdata", "key.pem")
	KeyData, _ := ioutil.ReadFile(fkey)

	pkey, _ = auth.PEMToPrivateKey(KeyData)

	db = os.Getenv("DFSS_MONGO_URI")
	if db == "" {
		db = "mongodb://localhost/dfss-test"
	}
	ttpAddressPort = "localhost:9090"

	sequence = []uint32{0, 1, 2, 0, 1, 2, 0, 1, 2}

	for i := 0; i < 3; i++ {
		h := sha512.Sum512([]byte{byte(i)})
		signer := h[:]
		signers = append(signers, signer)
	}

	contractDocumentHash = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	signatureUUIDBson = bson.NewObjectId()
	signatureUUID = signatureUUIDBson.Hex()

	signedHash = []byte{}

	signersEntities = make([]entities.Signer, 0)
	for _, s := range signers {
		signerEntity := entities.NewSigner(s)
		signersEntities = append(signersEntities, *signerEntity)
	}
}

func TestMain(m *testing.M) {
	/*ttp = GetServer(fca, fcert, fkey, "", db, true)
	go func() { _ = net.Listen(ttpAddressPort, ttp) }()
	*/
	dbManager, err := mgdb.NewManager(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occured during the connection to MongoDB:", err)
		os.Exit(2)
	}

	viper.Set("verbose", true)
	ttp = &ttpServer{
		DB: dbManager,
	}
	code := m.Run()

	ttp.DB.Close()

	os.Exit(code)
}

func TestAlert(t *testing.T) {
	// TODO
	// This requires the user of a real Alert message
}
