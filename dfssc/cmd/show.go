package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"dfss/dfssc/common"
	"dfss/dfssp/contract"

	"github.com/spf13/cobra"
)

const contractShowTemplate = `UUID       : {{.UUID}}
Filename   : {{.File.Name}}
Filehash   : {{.File.Hash}}
Created on : {{.Date.Format "2006-01-02 15:04:05 MST"}}

Comment    :
  {{.Comment}}

Signers    :
{{range .Signers}}  - {{.Email}}
{{end}}`

var showCmd = &cobra.Command{
	Use:   "show <c>",
	Short: "print contract information from file c",
	Run:   showContract,
}

func showContract(cmd *cobra.Command, args []string) {
	filename := args[0]
	c := getContract(filename)
	if c == nil {
		return
	}

	b := new(bytes.Buffer)
	tmpl, err := template.New("contract").Parse(contractShowTemplate)
	if err != nil {
		fmt.Println("Internal error:", err)
		return
	}

	err = tmpl.Execute(b, c)
	if err != nil {
		fmt.Println("Cannot print contract:", err)
	}
	fmt.Print(b.String())
}

func getContract(filename string) *contract.JSON {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return nil
	}

	c, err := common.UnmarshalDFSSFile(data)
	if err != nil {
		fmt.Println("Corrupted file:", err)
		return nil
	}
	return c
}
