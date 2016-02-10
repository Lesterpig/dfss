package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"os"
	"testing"
)

var pkey *rsa.PrivateKey
var csrFixture = `-----BEGIN CERTIFICATE REQUEST-----
MIIBgjCB7AIBADBDMQswCQYDVQQGEwJGUjEMMAoGA1UECgwDT1JHMRAwDgYDVQQL
DAdPUkdVTklUMRQwEgYDVQQDDAt0ZXN0QGdvLnRsZDCBnzANBgkqhkiG9w0BAQEF
AAOBjQAwgYkCgYEAtbEFS3VyXHAcNzZ49XKgXzv9SaBszbHWAXmuQlgH4dyjL7OX
w6NOpjSrIW2MVN99/boW1CilMpJzyMRkfkYg2u/HQw1KRUqP62Tl9FIbFjO+rITC
JI4fHsMpOh6+oWw62wf9mbKfL+kmFTTTfAWZpcE/R8IM+vJK4R+6DE7qvXMCAwEA
AaAAMA0GCSqGSIb3DQEBCwUAA4GBAJKNtd8IsxkJyWnOoJjckX+djFxoCNgo7JS1
6evVTU3esDRQ0P4T6oqn4D+yGQlRNtO6/Ko1D9Vv8v14hG7ZJ23Xr6PNBCQEB1a4
vzcnqUbk1ftU8qbOoTEEElEEeGu/gaDYjHPt/P9apngZpV3KXVAepAyRRLdXPfKa
Shc+gMEf
-----END CERTIFICATE REQUEST-----
`
var crtFixture = `-----BEGIN CERTIFICATE-----
MIICRDCCAa2gAwIBAgIJAIf+q5v9t+rTMA0GCSqGSIb3DQEBCwUAMDoxCzAJBgNV
BAYTAkZSMQwwCgYDVQQKDANPUkcxEDAOBgNVBAsMB09SR1VOSVQxCzAJBgNVBAMM
AkNBMCAXDTE1MTEyMDE1NDMwOVoYDzQ3NTMxMDE2MTU0MzA5WjA6MQswCQYDVQQG
EwJGUjEMMAoGA1UECgwDT1JHMRAwDgYDVQQLDAdPUkdVTklUMQswCQYDVQQDDAJD
QTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAqNbQWl2UZiOmJcZDA5x2H2U5
m2qC0D3NNPv0jNOm6shGTKhcLH1W8DrDtv5NjRWN1XTpfy0VkrsoyPpsU6PFFzZC
GmkCoXKBD/dvDNrid2MgbURzx+0a+EmUfFh0+tVP2Dzy0zgb/FZWkM6HT0VQ8KAb
SmlRBctiujDV1RBUOm8CAwEAAaNQME4wHQYDVR0OBBYEFNSIxFzdlyGUGBnqsKpd
bS4te57xMB8GA1UdIwQYMBaAFNSIxFzdlyGUGBnqsKpdbS4te57xMAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQELBQADgYEAECSnGMqnlgyBdqTC02/Uo1jiPqJjLZV1
TRJFHxs4JPAsff+rdAQ1TQVfaNnvAkAoXVzM34xGPkJserMUBc7aQ61WrByGImai
RqEe6wUqHGuH2blNt+2LSSuFWuR02+LxsJARDVSLViAS3lNgXlgGnzOaRs31iwwU
czHnSiYoCog=
-----END CERTIFICATE-----
`
var keyFixture = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCo1tBaXZRmI6YlxkMDnHYfZTmbaoLQPc00+/SM06bqyEZMqFws
fVbwOsO2/k2NFY3VdOl/LRWSuyjI+mxTo8UXNkIaaQKhcoEP928M2uJ3YyBtRHPH
7Rr4SZR8WHT61U/YPPLTOBv8VlaQzodPRVDwoBtKaVEFy2K6MNXVEFQ6bwIDAQAB
AoGAWaKZwK/XvhYE+h70qvEgwPAzkjAMvNNio1Nz9GPVROYIdGAZd0Efq6/3Aaqm
r1UXFJDZ+buMrXaRY4mXgxv54MjkX4d4KRVfAfIRbZQlP1jjnT6eRFhLFGOZe6pT
FUxMWO8wJDPwuNJgHFZOm+Ja3v2Hgvt6wgqu+l1onx4RNtECQQDZZbPPtgqTIsrv
BO0PZvL2BLFpg55NUYvpm57JU0/wU+rqG27gVUx7IYGWQ7J61BRIwmMR9HwrDRtE
EwPDelXnAkEAxtHJqHdElDcJGH0Um1WVjd2U229Imo9FBTbvFjj8uaJwk1MYhwks
fUcvD6+s0uoKZAdHogQBM7nN5OrtrnjWOQJBAM+Fpf/BZpbNv6oqqaDqRUNTd4eh
fJuSHF0DkK/eN5DSioyvY0gCJN/lPC6UsOtPR42tAaVCHMV73Ws+O3l+bkECQGvT
pRGHtZrIildMpttjvBtXe/7SSMcCQoWEeIBN4cpvraxI2bmKoSVEcOKJ/SnaIk6D
oDbfAyPhdifbvZHtGQkCQQCDH0Jo3JY7TlOripsIWm8hyikOzw9Lfonhbvnaofjt
amR9w6/SM5D0y20NqMVCmJxHWYW9sRIfZOmRjprYbczH
-----END RSA PRIVATE KEY-----
`

func TestMain(m *testing.M) {
	pkey, _ = GeneratePrivateKey(1024)
	os.Exit(m.Run())
}

func TestGetCertificateRequest(t *testing.T) {

	res, err := GetCertificateRequest("FR", "ORG", "ORGUNIT", "test@go.tld", pkey)

	if res == nil || err != nil {
		t.Fatal(err)
	}

	if res[0] != '-' {
		t.Fatalf("Bad format\n%s", res)
	}
}

func TestPEMToCertificateRequest(t *testing.T) {

	res, err := PEMToCertificateRequest([]byte(csrFixture))

	if res == nil || err != nil {
		t.Fatal(err)
	}

	if res.Subject.Country[0] != "FR" {
		t.Fatal("Wrong country: ", res.Subject.Country)
	}

	if res.Subject.CommonName != "test@go.tld" {
		t.Fatal("Wrong CN: ", res.Subject.CommonName)
	}

	res, err = PEMToCertificateRequest([]byte("invalid"))

	if err == nil {
		t.Fatal("The request should not have been decoded as is was invalid format")
	}

}

func TestGetSelfSignedCertificate(t *testing.T) {
	res, err := GetSelfSignedCertificate(10, 20, "FR", "TEST", "TEST UNIT", "My Cn", pkey)

	if res == nil || err != nil {
		t.Fatal(err)
	}

	if res[0] != '-' {
		t.Fatalf("Bad format\n%s", res)
	}

}

func TestPEMToCertificate(t *testing.T) {

	res, err := PEMToCertificate([]byte(crtFixture))

	if res == nil || err != nil {
		t.Fatal(err)
	}

	if res.Subject.CommonName != "CA" {
		t.Fatal("Wrong CN: ", res.Subject.CommonName)
	}

	if res.Issuer.CommonName != "CA" {
		t.Fatal("Wrong issuer: ", res.Issuer.CommonName)
	}

}

func TestGetCertificate(t *testing.T) {

	req, _ := PEMToCertificateRequest([]byte(csrFixture))
	crt, _ := PEMToCertificate([]byte(crtFixture))
	key, _ := PEMToPrivateKey([]byte(keyFixture))

	res, err := GetCertificate(10, 21, req, crt, key)

	if res == nil || err != nil {
		t.Fatal(err)
	}

	if res[0] != '-' {
		t.Fatalf("Bad format\n%s", res)
	}

}

func TestGetCertificateHash(t *testing.T) {
	crt, _ := PEMToCertificate([]byte(crtFixture))
	res := GetCertificateHash(crt)
	expected := "5c68072d750aa6c24c96e0b984ed399b726f8664456ed74259841b207b43159c5eaddf436d8f28c48b3ef6c83e69d68d115f0f0e4cf71001c2ca599a6ff7a0c1"

	if fmt.Sprintf("%x", res) != expected {
		t.Fatalf("Bad hash")
	}
}

func ExampleGetCertificate() {

	// Load elements from PEM files
	certificateRequest, _ := PEMToCertificateRequest([]byte(csrFixture))
	signerCertificate, _ := PEMToCertificate([]byte(crtFixture))
	signerKey, _ := PEMToPrivateKey([]byte(keyFixture))

	// Generate the certificate for 365 days with a serial of 0x10 (16)
	cert, err := GetCertificate(365, uint64(0x10), certificateRequest, signerCertificate, signerKey)

	if cert == nil || err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Certificate generated")
	}

	// Check certificate validity

	roots := x509.NewCertPool()
	roots.AddCert(signerCertificate)

	signeeCertificate, _ := PEMToCertificate(cert)

	_, err = signeeCertificate.Verify(x509.VerifyOptions{
		Roots: roots,
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Certificate authenticated")
	}

	// Output:
	// Certificate generated
	// Certificate authenticated

}
