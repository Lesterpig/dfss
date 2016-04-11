package server

import (
	"sort"
	"sync"
	"time"

	"dfss/dfssd/api"
)

// NOTE: buffers are managed with slices since lists are pretty much not used in go
// @see https://github.com/golang/go/wiki/SliceTricks

var in []*api.Log // incoming msg buffer
var inMutex = &sync.Mutex{}

// addMessage to storage
func addMessage(msg *api.Log) {
	inMutex.Lock()
	in = append(in, msg)
	inMutex.Unlock()
}

// display logs that are more than since (ms) old
func display(since int64, fn func(*api.Log)) {
	var out []*api.Log      // sorted messages to display
	var recycled []*api.Log // messages too recent to be displayed

	present := time.Now().UnixNano()

	inMutex.Lock()

	for _, v := range in {
		if present-(*v).Timestamp > 1000000*since {
			out = append(out, v)
		} else {
			recycled = append(recycled, v)
		}
	}

	in = recycled
	inMutex.Unlock()

	sort.Sort(ByTimestamp(out))

	for _, v := range out {
		//lfn()
		fn(v)
	}
}

// refresh every second
func displayHandler(fn func(*api.Log)) {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		display(1000, fn)
	}
}

// ByTimestamp sorting interface
type ByTimestamp []*api.Log

func (l ByTimestamp) Len() int {
	return len(l)
}
func (l ByTimestamp) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l ByTimestamp) Less(i, j int) bool {
	return (*l[i]).Timestamp < (*l[j]).Timestamp
}
