package sign

// SignerStatus represents the current state of a signer.
type SignerStatus int

// These constants represent the different states of a signer.
const (
	StatusWaiting SignerStatus = iota
	StatusConnecting
	StatusConnected
	StatusError
)
