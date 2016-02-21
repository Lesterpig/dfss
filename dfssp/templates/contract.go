package templates

const contract = `Dear Sir or Madam,

Someone asked you to sign a contract on the DFSS platform.
Please download the attached file and open it with the DFSS client.

{{template "contractDetails" .}}
{{template "signature"}}
`

const contractDetails = `Signers :
{{range .Signers}}  - {{.Email}}
{{end}}
Contract name : {{.File.Name}}
SHA-512 hash  : {{printf "%x" .File.Hash}}
Comment       : {{.Comment}}
`
