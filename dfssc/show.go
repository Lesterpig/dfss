package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"dfss/dfssc/common"
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

func showContract(filename string) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return
	}

	c, err := common.UnmarshalDFSSFile(data)
	if err != nil {
		fmt.Println("Corrupted file:", err)
		return
	}

	tmpl, err := template.New("contract").Parse(contractShowTemplate)
	if err != nil {
		fmt.Println("Internal error:", err)
		return
	}

	b := new(bytes.Buffer)
	err = tmpl.Execute(b, c)
	if err != nil {
		fmt.Println("Cannot print contract:", err)
	}

	fmt.Print(b.String())

}
