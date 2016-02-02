package templates

const invitation = `Dear Sir or Madam,

Someone invited you to sign a multiparty contract using the DFSS platform.
Please download the latest version of DFSS an register a new account using
this adress mail in order to sign the contract.

{{template "contractDetails" .}}
{{template "signature"}}
`
