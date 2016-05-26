package sign

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	cAPI "dfss/dfssc/api"
	"dfss/dfssc/common"
	"dfss/dfssp/contract"
)

// SignedContractJSON is an union of contract and related signatures
type SignedContractJSON struct {
	Contract   contract.JSON
	Signatures []cAPI.Signature
}

// PersistSignaturesToFile save contract informations and signatures to disk
func (m *SignatureManager) PersistSignaturesToFile() error {

	// Check content, don't write an empty file
	if len(m.archives.receivedSignatures) == 0 {
		return fmt.Errorf("No stored signatures, cannot create an empty file (yes I'm a coward)")
	}

	// Fill JSON struct
	signedContract := SignedContractJSON{
		Contract: *m.contract,
		Signatures: make(
			[]cAPI.Signature,
			len(m.archives.sentSignatures)+len(m.archives.receivedSignatures),
		),
	}

	for i, s := range m.archives.sentSignatures {
		signedContract.Signatures[i] = *s
	}

	for i, s := range m.archives.receivedSignatures {
		signedContract.Signatures[len(m.archives.sentSignatures)+i] = *s
	}

	proof, err := json.MarshalIndent(signedContract, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(m.mail+"-"+m.uuid+".proof", proof, 0600)
}

// PersistRecoverDataToFile : save recover informations to disk.
// returns the file name and an error if any occured
func (m *SignatureManager) PersistRecoverDataToFile() (string, error) {
	// Check content, don't write an empty file
	if len(m.uuid) == 0 || len(m.ttpData.Addrport) == 0 || len(m.ttpData.Hash) == 0 {
		return "", fmt.Errorf("Invalid recover data. Cannot persist file.")
	}

	// Fill JSON struct
	recData := common.RecoverDataJSON{
		SignatureUUID: m.uuid,
		TTPAddrport:   m.ttpData.Addrport,
		TTPHash:       m.ttpData.Hash,
	}

	file, err := json.MarshalIndent(recData, "", "  ")
	if err != nil {
		return "", err
	}

	filename := m.mail + "-" + m.uuid + ".run"
	err = ioutil.WriteFile(filename, file, 0600)
	if err != nil {
		return "", err
	}

	return filename, nil
}
