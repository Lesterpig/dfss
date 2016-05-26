package templates

const invitation = `Dear Sir or Madam,

Someone invited you to sign a multiparty contract using the DFSS platform.
Please download the latest version of DFSS and register a new account using
this mail address in order to sign the contract.

{{template "contractDetails" .}}
{{template "signature"}}
`
