package net

import (
	"fmt"
	"net"
	"testing"
	"time"

	"dfss/auth"
	pb "dfss/net/fixtures"
	"golang.org/x/net/context"
)

const caFixture = `-----BEGIN CERTIFICATE-----
MIIB5TCCAY+gAwIBAgIJAKId2y6Lo9T8MA0GCSqGSIb3DQEBCwUAME0xCzAJBgNV
BAYTAkZSMQ0wCwYDVQQKDARERlNTMRswGQYDVQQLDBJERlNTIFBsYXRmb3JtIHYw
LjExEjAQBgNVBAMMCWxvY2FsaG9zdDAgFw0xNjAxMjYxNTM2NTNaGA80NDgwMDMw
ODE1MzY1M1owTTELMAkGA1UEBhMCRlIxDTALBgNVBAoMBERGU1MxGzAZBgNVBAsM
EkRGU1MgUGxhdGZvcm0gdjAuMTESMBAGA1UEAwwJbG9jYWxob3N0MFwwDQYJKoZI
hvcNAQEBBQADSwAwSAJBAMGAgCtkRLePYFRTUN0V/0v/6phm0guHGS6f0TkSEas4
CGZTKFJVTBksMGIBtfyYw3XQx2bO8myeypDN5nV05DcCAwEAAaNQME4wHQYDVR0O
BBYEFO09nxx5/qeLK5Wig1+3kg66gn/mMB8GA1UdIwQYMBaAFO09nxx5/qeLK5Wi
g1+3kg66gn/mMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADQQCqNSH+rt/Z
ru2rkabLiHOGjI+AenSOvqWZ2dWAlLksYcyuQHKwjGWgpmqkiQCnkIDwIxZvu69Y
OBz0ASFn7eym
-----END CERTIFICATE-----
`

const serverKeyFixture = `-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBAMGAgCtkRLePYFRTUN0V/0v/6phm0guHGS6f0TkSEas4CGZTKFJV
TBksMGIBtfyYw3XQx2bO8myeypDN5nV05DcCAwEAAQJAHSdRKDh5KfbOGqZa3pR7
3GV4YPHM37PBFYc6rJCOXO9W8L4Q1kvEhjKXp7ke18Cge7bVmlKspvxvC62gxSQm
QQIhAPMYwpp29ZREdk8yU65Sp6w+EbZS9TjZkC+pk3syYjaxAiEAy8XWnnDMsUxb
6vp1SaaIfxI441AYzh3+8c56CAvt02cCIQDQ2jfvHz7zyDHg7rsILMkTaSwseW9n
DTwcRtOHZ40LsQIgDWEVAVwopG9+DYSaVNahWa6Jm6szpbzkc136NzMJT3sCIQDv
T2KSQQIYEvPYZmE+1b9f3rs/w7setrGtqVFkm/fTWQ==
-----END RSA PRIVATE KEY-----
`

const clientCertFixture = `-----BEGIN CERTIFICATE-----
MIIBkDCCAToCCQDSSWVk2vWTdjANBgkqhkiG9w0BAQsFADBNMQswCQYDVQQGEwJG
UjENMAsGA1UECgwEREZTUzEbMBkGA1UECwwSREZTUyBQbGF0Zm9ybSB2MC4xMRIw
EAYDVQQDDAlsb2NhbGhvc3QwIBcNMTYwMTI2MTUzNjUzWhgPNDQ4MDAzMDgxNTM2
NTNaME8xCzAJBgNVBAYTAkZSMQ0wCwYDVQQKDARERlNTMRkwFwYDVQQLDBBERlNT
IENsaWVudCB2MC4xMRYwFAYDVQQDDA10ZXN0QHRlc3QuY29tMFwwDQYJKoZIhvcN
AQEBBQADSwAwSAJBAMxHU0NP/elQbmM5HDZS5iWXr4wllaJ2bWWD0cZPI1p+jty0
wwkKwxEklPGZCDWq1+/C4EawaqMrtZW4HQVxdu8CAwEAATANBgkqhkiG9w0BAQsF
AANBACl2/KBGR8N4qzpNecr1yDdyfyE4nGYgr8aktAeHHNWFg53q3/VHokK0jEus
iM6sQlvDCoaE01s6gXrarE+APfU=
-----END CERTIFICATE-----
`

const clientKeyFixture = `-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBAMxHU0NP/elQbmM5HDZS5iWXr4wllaJ2bWWD0cZPI1p+jty0wwkK
wxEklPGZCDWq1+/C4EawaqMrtZW4HQVxdu8CAwEAAQJBAIBwYCO0idtGnQF6CQkG
+nmsc83UW874UzQ+u4jKfVoJlB+Llp7Y6iFhyS+Trw+ffpT32hYtoX+cRxWVh+g3
xOECIQD6S/gPG3dCJWGsZ7/V6WVOfzdz6VAHtBjfgYNMLZQhZwIhANDu7FpmhxWm
nVr3u6hLKY28Zhp03LqsljfEoWEkAmE5AiBPreV+8bBqZzoLx09jipRMg+UkSi7G
9QdCB5nDo3LXmwIgFGdArZNVncenljqbGNQ+OpkrX2oKJDC2eru5BsN9eAECIFXi
HBs0FRZQlpn1kXgfvakOtkifTcLaksnn3Y5PAEhP
-----END RSA PRIVATE KEY-----
`

// SERVER DEFINITION

type testServer struct{}

func (s *testServer) Ping(ctx context.Context, in *pb.Hop) (*pb.Hop, error) {
	(*in).Id++
	return in, nil
}
func (s *testServer) Auth(ctx context.Context, in *pb.Empty) (*pb.IsAuth, error) {
	return &pb.IsAuth{Auth: GetCN(&ctx) == "test@test.com"}, nil
}

func startTestServer(c chan bool) {

	ca, _ := auth.PEMToCertificate([]byte(caFixture))
	key, _ := auth.PEMToPrivateKey([]byte(serverKeyFixture))

	server := NewServer(ca, key, ca)
	pb.RegisterTestServer(server, &testServer{})
	go func() {
		_ = Listen("localhost:9000", server)
	}()
	<-c
	server.TestingCloseConns()
	server.Stop()
}

// TESTS

func TestServerOnly(t *testing.T) {
	c := make(chan bool)
	go startTestServer(c)
	time.Sleep(2 * time.Second) // Two seconds to start

	_, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		t.Fatal("Unable to bind to server:", err)
	}

	c <- true // Stop server
	time.Sleep(100 * time.Millisecond)
}

func TestServerClientAuth(t *testing.T) {

	// Start server
	c := make(chan bool)
	go startTestServer(c)
	time.Sleep(2 * time.Second)

	ca, _ := auth.PEMToCertificate([]byte(caFixture))
	cert, _ := auth.PEMToCertificate([]byte(clientCertFixture))
	key, _ := auth.PEMToPrivateKey([]byte(clientKeyFixture))

	conn, err := Connect("localhost:9000", cert, key, ca)

	if err != nil {
		t.Fatal("Unable to connect:", err)
	}

	client := pb.NewTestClient(conn)
	sharedServerClientTest(t, client, true)
	_ = conn.Close()

	c <- true
	time.Sleep(100 * time.Millisecond)
}

func TestServerClientNonAuth(t *testing.T) {

	// Start server
	c := make(chan bool)
	go startTestServer(c)
	time.Sleep(2 * time.Second)

	ca, _ := auth.PEMToCertificate([]byte(caFixture))
	conn, err := Connect("localhost:9000", nil, nil, ca)

	if err != nil {
		t.Fatal("Unable to connect:", err)
	}

	client := pb.NewTestClient(conn)
	sharedServerClientTest(t, client, false)
	_ = conn.Close()

	c <- true
	time.Sleep(100 * time.Millisecond)
}

func sharedServerClientTest(t *testing.T, client pb.TestClient, expectedAuth bool) {
	r, err := client.Ping(context.Background(), &pb.Hop{Id: 0})
	if err != nil {
		t.Fatal("Unable to ping:", err)
	}
	if (*r).Id != 1 {
		t.Fatal("Bad result, got", *r)
	}

	auth, err := client.Auth(context.Background(), &pb.Empty{})
	if err != nil {
		t.Fatal("Unable to auth:", err)
	}
	if expectedAuth != (*auth).Auth {
		t.Fatal("Bad result, got", *auth)
	}
}

// EXAMPLE

func Example() {

	// Load certs and private keys
	ca, _ := auth.PEMToCertificate([]byte(caFixture))
	cert, _ := auth.PEMToCertificate([]byte(clientCertFixture))
	ckey, _ := auth.PEMToPrivateKey([]byte(clientKeyFixture))
	skey, _ := auth.PEMToPrivateKey([]byte(serverKeyFixture))

	// Init server
	server := NewServer(ca, skey, ca)
	pb.RegisterTestServer(server, &testServer{})
	go func() {
		_ = Listen("localhost:9000", server)
	}()

	// Let the server enough time to start property
	time.Sleep(2 * time.Second)

	// Start a client
	// The second and third arguments can be empty for non-auth connection
	conn, err := Connect("localhost:9000", cert, ckey, ca)
	if err != nil {
		panic("Unable to connect")
	}

	client := pb.NewTestClient(conn)

	// During a ping, the server increments the Hop.Id field (test case only)
	r, err := client.Ping(context.Background(), &pb.Hop{Id: 41})
	if err != nil {
		panic("Unable to ping")
	}

	fmt.Println((*r).Id)

	// Close client
	_ = conn.Close()

	// Stop server
	server.Stop()

	// Output:
	// 42
}
