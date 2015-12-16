package mails

import (
  "errors"
  "bytes"
  "crypto/tls"
  "net/smtp"
  "fmt"
  "encoding/base64"
)

// Modelizes the constants of a connection : an actual client, and a sender
type CustomClient struct {
  sender string
  client *smtp.Client
}

var errorNoTLS = errors.New("Connection failed : mail server doesn't support TLS");

var errorNoAuthentification = errors.New(`Connection failed : mail server doesn't
support authentication`);

// Start a connection. Returns errorNoTLS if server doesn't implement TLS.
// Returns errorNoAuthentification if there is no authentification possible
func initiate(sender string, host string, port string, password string) (*CustomClient,error){
  // Connect to the server
  c, err := smtp.Dial(host+port)
  if err != nil{
    return nil,err
  }

  // Authenticate to the server
  auth := smtp.PlainAuth("",sender,password,host)
  if ok,_ := c.Extension("AUTH");!ok{
    return nil,errorNoAuthentification
  }
  if err := c.Auth(auth); err != nil{
    return nil,err
  }

  // Check that the server does implement TLS
  if ok,_ := c.Extension("StartTLS");!ok{
    return nil,errorNoTLS;
  }

  // Start tls
  if err := c.StartTLS(&tls.Config{ServerName : host});err != nil{
    return nil, err
  }

  // Set the sender of the mail
  if err := c.Mail(sender);err != nil{
    return nil,err
  }

  return &CustomClient{sender,c},nil
}

// Send a mail with the custom client. Returns nil on success.
func (c *CustomClient) Send(dest []string, subject string, message string) error{
  // Keep the connection in a local variable for ease of access
  connection := c.client;

  // Set the header once & forget it (save for the receiver)
  // The following is the same for all connections (speficied by user)
  // - Content-Transfer-Encoding string : base64
  // - Content-Type string : "text/plain; charset = \"utf-8\""
  // - MIME-Version string : "1.0"
  var buffer bytes.Buffer
  buffer.WriteString("From: ")
  buffer.WriteString(c.sender)
  buffer.WriteString("\r\n")
  buffer.WriteString("Content-Transfer-Encoding: base64\r\n")
  buffer.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
  buffer.WriteString("MIME-Version: 1.0\r\n")
  buffer.WriteString("Subject: ")
  buffer.WriteString(subject)
  buffer.WriteString("\r\n")
  buffer.WriteString("To: ")

  header := buffer.String()

  // Encode the message in base64 ONCE
  base64Message := base64.StdEncoding.EncodeToString([]byte(message))

  for _,receiver := range dest{

    // Set the sender. If this is not the first call, this should also reset
    // the receivers
    if err := connection.Mail(c.sender);err != nil{
      return err
    }

    // Set the receiver. This modifies the header in the current instance
    if err := connection.Rcpt(receiver);err != nil{
      return err
    }
    localHeader := header + receiver + "\r\n"

    // Set the message : header, then message encoded in base64
    wc, err := connection.Data()
    if err != nil {
      return err
    }
    _, err = fmt.Fprintf(wc,localHeader+base64Message)
    if err != nil{
       return err;
    }
    err = wc.Close()
    if err != nil {
      return err
    }
  }
  return nil;
}

// Quits the connection of CustomClient
func (c *CustomClient) Close() error{
  return c.client.Close()
}
