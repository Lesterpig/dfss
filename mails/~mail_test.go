package mails

import (
  "testing"
  "fmt"
  "os"
)

var c *CustomClient
var err error

func TestMain(m *testing.M) {
  // Setup
  c,err = Initiate("Jenkar","127.0.0.1","5000","")
  if err != nil{
    fmt.Println(err)
  }

  code := m.Run()

  err = c.Close()
  if err != nil{
     fmt.Println(err)
  }
  os.Exit(code)
}

func TestSingleMail(t *testing.T){
  err = c.Send([]string{"a@lesterpig.com"},"TestSingleMail","Gros espoirs!");
  if err != nil{
    t.Fatal(err)
  }
}

func TestDoubleMail(t *testing.T){
  err = c.Send([]string{"a@lesterpig.com","b@lesterpig.com"},"TestDoubleMail","Gros espoirs!");
  if err != nil{
    t.Fatal(err)
  }
}

func TestRuneMail(t *testing.T){
  err = c.Send([]string{"a@lesterpig.com"},"TestRuneMail","测试")
  if err!= nil{
    t.Fatal(err)
  }
}
