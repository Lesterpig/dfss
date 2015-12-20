// Package mails provides a simple interface with the smtp library
package mails

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"
)

// CustomClient :
// Modelizes the constants of a connection : an actual client, and a sender
type CustomClient struct {
	sender string
	client *smtp.Client
}

const rfc2822 = "Fri 18 Dec 2015 10:01:17 -0606" // used to format time to rfc2822. Not accurate but fmt can't see a ,

// NewCustomClient starts up a custom client.
func NewCustomClient(sender, host, port, username, password string) (*CustomClient, error) {

	// Connect to the server. Type of connection is smtp.Client
	connection, err := smtp.Dial(host + ":" + port)
	if err != nil {
		return nil, err
	}

	// Check that the server does implement TLS
	if ok, _ := connection.Extension("StartTLS"); !ok {
		return nil, errors.New("Connection failed : mail server doesn't support TLS")
	}

	// Start tls
	if err := connection.StartTLS(&tls.Config{InsecureSkipVerify: true, ServerName: host}); err != nil {
		return nil, err
	}

	// Authenticate to the server if it is supported
	if ok, _ := connection.Extension("AUTH"); ok {
		auth := smtp.PlainAuth(username, username, password, host)
		if err := connection.Auth(auth); err != nil {
			return nil, err
		}
	}

	return &CustomClient{sender, connection}, nil
}

// Send a mail with the custom client. Returns nil on success.
func (c *CustomClient) Send(receivers []string, subject, message string, extensions, filenames []string) error {
	// Keep the connection in a local variable for ease of access
	connection := c.client

	boundary := randomBoundary()
	header := createHeader(c.sender, subject, boundary)

	// Encode the message in base64 ONCE
	base64Message := base64.StdEncoding.EncodeToString([]byte(message))

	for _, receiver := range receivers {

		// Set the sender
		if err := connection.Mail(c.sender); err != nil {
			return err
		}

		// Set the receiver. This modifies the header in the current instance
		if err := connection.Rcpt(receiver); err != nil {
			return err
		}

		// Set the message : header, then message encoded in base64, then attachments
		var localBuffer bytes.Buffer
		if err := createFullMessage(&localBuffer, receiver, c.sender, header, base64Message, extensions, filenames, boundary); err != nil {
			return err
		}

		// Send it. Data returns a writer to which one can write to write the message itself
		emailWriter, err := connection.Data()
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(emailWriter, localBuffer.String())
		if err != nil {
			return err
		}
		err = emailWriter.Close()
		if err != nil {
			return err
		}

		// Reset the envellope
		err = connection.Reset()
		if err != nil {
			return err
		}
	}

	return nil
}

// Close the connection of CustomClient
func (c *CustomClient) Close() error {
	return c.client.Close()
}

// Creates the header for all messages
func createHeader(sender, subject, boundary string) string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "From: %s\r\n", sender)
	fmt.Fprintf(&buffer, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buffer, "Subject: %s\r\n", subject)
	// Replace the first space with a comma and a space to conform to rfc2822
	fmt.Fprintf(&buffer, "Date: %s%s", strings.Replace(time.Now().UTC().Format(rfc2822), " ", ", ", 1), "\r\n")
	fmt.Fprintf(&buffer, "Content-Type: multipart/mixed; boundary=\"%s\"; charset=\"UTF-8\"\r\n", boundary)
	fmt.Fprintf(&buffer, "To: ")
	return buffer.String()
}

// Create the full message for a single receiver
func createFullMessage(b *bytes.Buffer, receiver, sender, globalHeader, base64Message string, extensions, filenames []string, boundary string) error {
	fmt.Fprintf(b, "%s%s\r\n", globalHeader, receiver)

	writer := multipart.NewWriter(b)
	if err := writer.SetBoundary(boundary); err != nil {
		return err
	}
	// Set the message
	if err := createText(writer, base64Message); err != nil {
		return err
	}

	// Set attachments. Here for now because the boundaries are wanted unique
	for index, value := range filenames {
		if err := createAttachment(writer, extensions[index], value); err != nil {
			return err
		}
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}

// Create an attachment with a certain extension
func createAttachment(writer *multipart.Writer, extension, path string) error {
	// Create a header
	newHeader := make(textproto.MIMEHeader)
	newHeader.Add("Content-Type", extension)
	newHeader.Add("Content-Transfer-Encoding", "base64")
	newHeader.Add("Content-Disposition", "attachment; filename="+path+";")

	// Create a writer for the file
	output, err := writer.CreatePart(newHeader)
	if err != nil {
		return err
	}

	// Write the file to the message
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	fmt.Fprintf(output, base64.StdEncoding.EncodeToString([]byte(data)))

	return nil
}

// Creates the equivalent of the message wrapped in a boundary. The message is expected to have been encoded via base64
func createText(writer *multipart.Writer, message string) error {
	// Create the mime header for the message
	mimeHeaderMessage := make(textproto.MIMEHeader)
	mimeHeaderMessage.Add("Content-Transfer-Encoding", "base64")
	mimeHeaderMessage.Add("Content-Type", "text/plain; charset=\"UTF-8\"")

	// Set the message
	output, err := writer.CreatePart(mimeHeaderMessage)
	if err != nil {
		return err
	}

	fmt.Fprintf(output, "%s", message)
	return nil
}

// Totally copied from go stl
func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}
