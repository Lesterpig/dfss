package sign

import (
	"testing"

	cAPI "dfss/dfssc/api"
	"github.com/stretchr/testify/assert"
)

func TestContainsPromiseFrom(t *testing.T) {
	sender := []byte{0, 0, 7}

	var promises []*cAPI.Promise
	promises = append(promises, &cAPI.Promise{
		Context: &cAPI.Context{
			SenderKeyHash: []byte{0, 0, 8},
		},
		Index: uint32(1),
	})
	promises = append(promises, &cAPI.Promise{
		Context: &cAPI.Context{
			SenderKeyHash: sender,
		},
		Index: uint32(1),
	})

	m := SignatureManager{
		archives: &Archives{
			receivedPromises: promises,
		},
	}

	present, index := m.containsPromiseFrom([]byte{0, 0, 6}, 0)
	assert.Equal(t, present, false)
	assert.Equal(t, index, 0)

	present, index = m.containsPromiseFrom([]byte{0, 0, 7}, 0)
	assert.Equal(t, present, false)
	assert.Equal(t, index, 1)

	present, index = m.containsPromiseFrom([]byte{0, 0, 7}, 2)
	assert.Equal(t, present, true)
	assert.Equal(t, index, 1)
}

func TestRemoveReceivedPromise(t *testing.T) {
	sender := []byte{0, 0, 7}

	var promises []*cAPI.Promise
	promises = append(promises, &cAPI.Promise{
		Context: &cAPI.Context{
			SenderKeyHash: []byte{0, 0, 8},
		},
		Index: uint32(1),
	})
	promises = append(promises, &cAPI.Promise{
		Context: &cAPI.Context{
			SenderKeyHash: sender,
		},
		Index: uint32(1),
	})

	m := SignatureManager{
		archives: &Archives{
			receivedPromises: promises,
		},
	}

	assert.Equal(t, len(m.archives.receivedPromises), 2)

	err := m.removeReceivedPromise(-1)
	assert.Equal(t, err.Error(), "Index out of range")

	err = m.removeReceivedPromise(2)
	assert.Equal(t, err.Error(), "Index out of range")

	err = m.removeReceivedPromise(1)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(m.archives.receivedPromises), 1)

	present, index := m.containsPromiseFrom(sender, 2)
	assert.Equal(t, present, false)
	assert.Equal(t, index, 0)
	present, index = m.containsPromiseFrom([]byte{0, 0, 8}, 2)
	assert.Equal(t, present, true)
	assert.Equal(t, index, 0)
}

func TestUpdateRecievedPromises(t *testing.T) {
	sender := []byte{0, 0, 7}

	var promises []*cAPI.Promise
	promises = append(promises, &cAPI.Promise{
		Context: &cAPI.Context{
			SenderKeyHash: []byte{0, 0, 8},
		},
		Index: uint32(1),
	})
	promises = append(promises, &cAPI.Promise{
		Context: &cAPI.Context{
			SenderKeyHash: sender,
		},
		Index: uint32(1),
	})

	m := SignatureManager{
		archives: &Archives{
			receivedPromises: promises,
		},
	}

	newPromise0 := &cAPI.Promise{
		Context: &cAPI.Context{
			SenderKeyHash: []byte{0, 0, 0},
		},
		Index: uint32(2),
	}

	assert.Equal(t, len(m.archives.receivedPromises), 2)

	m.updateReceivedPromises([]*cAPI.Promise{newPromise0})
	assert.Equal(t, len(m.archives.receivedPromises), 3)

	present, index := m.containsPromiseFrom([]byte{0, 0, 0}, 3)
	assert.Equal(t, present, true)
	assert.Equal(t, index, 2)

	newPromise1 := &cAPI.Promise{
		Context: &cAPI.Context{
			SenderKeyHash: sender,
		},
		Index: uint32(2),
	}

	m.updateReceivedPromises([]*cAPI.Promise{newPromise1})
	assert.Equal(t, len(m.archives.receivedPromises), 3)

	present, index = m.containsPromiseFrom(sender, 3)
	assert.Equal(t, present, true)
	assert.Equal(t, index, 2)
}
