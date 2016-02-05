// Package user handles all the user-related logic
package user

// Register a user using the provided parameters
func Register(fileCA, fileCert, fileKey, addrPort, passphrase, country, organization, unit, mail string, bits int) error {
	manager, err := NewRegisterManager(fileCA, fileCert, fileKey, addrPort, passphrase, country, organization, unit, mail, bits)
	if err != nil {
		return err
	}

	err = manager.GetCertificate()
	if err != nil {
		return err
	}

	return nil
}
