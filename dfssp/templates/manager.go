package templates

import (
	"bytes"
	"log"
	"os"
	"text/template"

	"dfss/mails"
)

var tpl *template.Template
var ready bool

const signature = `Yours faithfully,

The DFSS Platform`

// Init compiles templates and panics if encountered an error.
// This method is called automatically one time by `Get`.
func Init() {

	tpl = template.Must(template.New("main").Parse("{{define `contract`}}" + contract + "{{end}}"))
	_ = template.Must(tpl.Parse("{{define `signature`}}" + signature + "{{end}}"))
	_ = template.Must(tpl.Parse("{{define `invitation`}}" + invitation + "{{end}}"))
	_ = template.Must(tpl.Parse("{{define `contractDetails`}}" + contractDetails + "{{end}}"))
	_ = template.Must(tpl.Parse("{{define `verificationMail`}}" + verificationMail + "{{end}}"))
	ready = true

}

// Get computes the asked template with the data provided.
func Get(name string, data interface{}) (string, error) {
	if !ready {
		Init()
	}

	b := new(bytes.Buffer)
	err := tpl.ExecuteTemplate(b, name, data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// MailConn is a helper to return a new mail connection if available.
// The result is null if the mail server is not ready.
// Do not forget to defer the close operation!
func MailConn() *mails.CustomClient {

	sender := os.Getenv("DFSS_MAIL_SENDER")
	host := os.Getenv("DFSS_MAIL_HOST")
	port := os.Getenv("DFSS_MAIL_PORT")
	username := os.Getenv("DFSS_MAIL_USERNAME")
	password := os.Getenv("DFSS_MAIL_PASSWORD")

	if len(sender) == 0 {
		return nil
	}

	c, err := mails.NewCustomClient(sender, host, port, username, password)
	if err != nil {
		log.Println(err)
		return nil
	}
	return c
}
