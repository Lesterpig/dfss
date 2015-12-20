package mails

import (
	"fmt"
	"os"
	"testing"
)

var client *CustomClient
var err error
var rcpt1 string
var rcpt2 string

func TestMain(m *testing.M) {
	// Setup is based on environment variables
	sender := os.Getenv("DFSS_TEST_MAIL_SENDER")
	host := os.Getenv("DFSS_TEST_MAIL_HOST")
	port := os.Getenv("DFSS_TEST_MAIL_PORT")
	username := os.Getenv("DFSS_TEST_MAIL_USER")
	password := os.Getenv("DFSS_TEST_MAIL_PASSWORD")
	rcpt1 = os.Getenv("DFSS_TEST_MAIL_RCPT1")
	rcpt2 = os.Getenv("DFSS_TEST_MAIL_RCPT2")
	client, err = NewCustomClient(sender, host, port, username, password)
	if err != nil {
		fmt.Println(err)
	}

	code := m.Run()

	err = client.Close()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func TestSingleMail(t *testing.T) {
	err = client.Send([]string{rcpt1}, "TestSingleMail", "Gros espoirs!", []string{}, []string{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoubleMail(t *testing.T) {
	err = client.Send([]string{rcpt1, rcpt2}, "TestDoubleMail", "Gros espoirs!", []string{}, []string{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRuneMail(t *testing.T) {
	err = client.Send([]string{rcpt1}, "TestRuneMail", "测试", []string{}, []string{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAttachmentMailText(t *testing.T) {
	err = client.Send([]string{rcpt1}, "TestAttachmentMailText", "What would make a good attachment?", []string{"text/plain"}, []string{"mail_test.go"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAttachmentMailImage(t *testing.T) {
	err = client.Send([]string{rcpt1}, "TestAttachmentMailImage", "What would make a good attachment?", []string{"image/gif"}, []string{"testImg.gif"})
	if err != nil {
		t.Fatal(err)
	}
}

func ExampleCustomClient() {

	// Create a connection
	client, err := NewCustomClient("qdauchy@insa-rennes.fr", "mailhost.insa-rennes.fr", "587", "qdauchy", "notreallymypass")
	if err != nil {
		fmt.Println(err)
	}

	// Some reused variables
	receiver1 := "lbonniot@insa-rennes.fr"
	receiver2 := "qdauchy@insa-rennes.fr"
	receiver3 := "tclaverie@insa-rennes.fr"
	subject := "Mail example"
	message := `Hello, this is a mail example. It's not like the cactus is going
  to be jealous or anything...`

	// Send a first mail, without attachments
	err = client.Send([]string{receiver1, receiver2}, subject, message, []string{}, []string{})
	if err != nil {
		fmt.Println(err)
	}

	// Send a second mail, with some attachments
	err = client.Send([]string{receiver1, receiver3}, subject, message, []string{"text/plain", "image/gif"}, []string{"email.go", "testImg.gif"})
	if err != nil {
		fmt.Println(err)
	}

	// Close the connection
	err = client.Close()
	if err != nil {
		fmt.Println(err)
	}

}
