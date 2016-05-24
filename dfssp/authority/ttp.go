package authority

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"dfss/dfssp/api"
)

// TTPHolder stores available TTPs (trusted third parties)
type TTPHolder struct {
	ttps  []*api.LaunchSignature_TTP
	next  int
	mutex *sync.Mutex
}

// NewTTPHolder loads available TTPs from the specified file.
// The format of this file should be as-is:
//
// <addr ttp 1>[:<port ttp 1] <SHA-512 hash of the ttp certificate (hex format)>\n
// ...
//
// Example: see testdata/ttps.
// If an error occurs during the retrieval of the file, an empty TTPHolder will be provided.
// If the file is corrupted (wrong format), and error will be thrown.
func NewTTPHolder(filename string) (*TTPHolder, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		data = []byte{}
	}

	lines := strings.Split(string(data), "\n")
	ttps := make([]*api.LaunchSignature_TTP, len(lines)-1) // -1 to ignore last string (empty)
	for i := 0; i < len(lines)-1; i++ {
		line := lines[i]
		words := strings.Split(line, " ")
		if len(words) < 2 {
			return nil, errors.New("corrupted ttp file: not enough words at line " + strconv.Itoa(i))
		}
		hash, err := hex.DecodeString(words[1])
		if err != nil {
			return nil, errors.New("corrupted ttp file: invalid hash at line " + strconv.Itoa(i))
		}
		ttps[i] = &api.LaunchSignature_TTP{
			Addrport: words[0],
			Hash:     hash,
		}
	}

	holder := &TTPHolder{
		ttps:  ttps,
		next:  0,
		mutex: &sync.Mutex{},
	}

	return holder, nil
}

// Nb returns the number of loaded TTP in this holder.
func (h *TTPHolder) Nb() int {
	return len(h.ttps)
}

// Get returns a TTP from the TTP holder.
// It is thread-safe, and base on a round-robin system.
//
// If the TTPHolder is empty, returns nil.
func (h *TTPHolder) Get() *api.LaunchSignature_TTP {
	if h.Nb() == 0 {
		return nil
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if len(h.ttps) == h.next {
		h.next = 0
	}

	ttp := h.ttps[h.next]
	h.next++
	return ttp
}

// Add adds the provided TTP to the TTP holder.
// It is thread-safe.
func (h *TTPHolder) Add(addrport string, hash []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.ttps = append(h.ttps, &api.LaunchSignature_TTP{
		Addrport: addrport,
		Hash:     hash,
	})
}

// Save saves the TTP holder in a file, respecting the same format as presented in the loader.
func (h *TTPHolder) Save(filename string) error {
	data := ""
	for _, ttp := range h.ttps {
		data += ttp.Addrport + " " + hex.EncodeToString(ttp.Hash) + "\n"
	}
	return ioutil.WriteFile(filename, []byte(data), 0600)
}
