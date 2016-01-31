package templates

const verificationMail = `Dear sir or Madam,

You asked to register to the DFSS platform.
Please send us your authentication request with
the following text as token:

{{.Token}}

If you did not asked for registration, we deeply excuse
for the error.

{{template "signature"}}
`

// VerificationMail contains the token to be sent in the verification mail
type VerificationMail struct {
	Token string
}
