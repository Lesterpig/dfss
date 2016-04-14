package common

import (
	"sync"
	"testing"
	"time"
)

func TestWaitingGroup(t *testing.T) {
	nb := 60
	w := NewWaitingGroupMap()

	waitGroup := &sync.WaitGroup{} // test only
	waitGroup.Add(nb)

	// Spawn nb emitters waiting for (nb-1) other emitters
	for i := 0; i < nb; i++ {
		go func(i int) {
			// Add some virtual latency
			time.Sleep(time.Duration(i) * time.Millisecond)
			// Join the waitingGroupMap
			myChan, nbs, _ := w.Join("A")
			w.Broadcast("A", i)
			// Wait for other msg
			for m := range myChan {
				nbs = append(nbs, m)
				if len(nbs) == nb {
					break
				}
			}
			// Free the waitingGroupMap
			w.Unjoin("A", myChan)
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait() // test only, wait for test to fully happen
}

func TestCloseWaitingGroup(t *testing.T) {
	w := NewWaitingGroupMap()

	waitGroup := &sync.WaitGroup{} // test only
	waitGroup.Add(1)

	go func() {
		myChan, _, _ := w.Join("A")
		for range myChan {
			t.Fatal("Should not be here")
		}
		// No need to call Unjoin here: if we do, we will try to unjoin a unknown room
		waitGroup.Done()
	}()

	time.Sleep(10 * time.Millisecond)
	w.CloseAll()
	waitGroup.Wait()
}
