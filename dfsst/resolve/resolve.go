package resolve

import (
	cAPI "dfss/dfssc/api"
	"dfss/dfsst/entities"
)

// ArePromisesComplete : determines if the set of promises present in the AlertRequest is EQUAL (not just included) to the one expected from the TTP
// for this signer, at this particular step of the signing sequence.
func ArePromisesComplete(promiseEntities []*entities.Promise, promise *cAPI.Promise) bool {
	// TODO
	// This requires to dig into the mathematical specifications, to determine the set of promises
	// expected to be persent in the AlertRequest from this signer at this position
	return true
}

// Solve : tries to generate the signed contract from present evidence.
func Solve(manager *entities.ArchivesManager) (bool, []byte) {
	// Test if we can generate the contract
	for i := range manager.Archives.Signers {
		ok := manager.HasSignerPromised(uint32(i))
		if !ok {
			return false, nil
		}
	}

	return true, GenerateSignedContract(manager.Archives)
}

// GenerateSignedContract : generates the signed contract.
// Does not take into account if we have the evidence to do it (see function 'CanGenerate').
//
// XXX : Implementation needs cryptographic promises
func GenerateSignedContract(archives *entities.SignatureArchives) []byte {
	// TODO
	return []byte{0, 0, 7}
}

// ComputeDishonestSigners : computes the dishonest signers from the provided evidence.
// ie: if an already aborted signer sent a message after his call to the ttp, then he is dishonest
func ComputeDishonestSigners(archives *entities.SignatureArchives, evidence []*entities.Promise) []uint32 {
	var res []uint32
	// Used to add only once the signers if several promises point to them as being dishonest
	added := false

	for _, signer := range archives.AbortedSigners {
		added = false
		sIndex := signer.SignerIndex
		for _, p := range evidence {
			if !added && (p.SenderKeyIndex == sIndex) && (p.SequenceIndex > signer.AbortIndex) {
				res = append(res, sIndex)
				added = true
			}
		}
	}

	return res
}
