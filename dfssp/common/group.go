package common

import (
	"container/list"
	"sync"
)

type waitingGroup struct {
	channels    *list.List
	oldMessages []interface{}
	mutex       sync.Mutex
}

// WaitingGroupMap is a synchronisation tool for goroutines.
// It enables several goroutines to wait for other goroutines in a specific "room".
//
// After joining a room, a goroutine can broadcast or wait for events.
// To avoid memory leaks, always call Unjoin when leaving a goroutine.
//
// See group_test.go for some examples.
type WaitingGroupMap struct {
	data  map[string]*waitingGroup
	mutex sync.Mutex
}

// NewWaitingGroupMap returns a ready to use WaitingGroupMap.
func NewWaitingGroupMap() *WaitingGroupMap {
	return &WaitingGroupMap{
		data: make(map[string]*waitingGroup),
	}
}

// Join permits the current goroutine to join a room.
// It returns the listenning channel and a slice containing messages already sent by other members of the room.
// The room is automatically created if unknown.
func (g *WaitingGroupMap) Join(room string) (listen chan interface{}, oldMessages []interface{}, newRoom bool) {
	// Check if the current waiting group knows this room
	g.mutex.Lock()
	_, present := g.data[room]
	if !present {
		g.data[room] = &waitingGroup{
			channels:    list.New(),
			oldMessages: make([]interface{}, 0),
		}
		newRoom = true
	}
	g.mutex.Unlock()

	g.data[room].mutex.Lock()
	listen = make(chan interface{}, 100)
	g.data[room].channels.PushBack(listen)
	oldMessages = g.data[room].oldMessages
	g.data[room].mutex.Unlock()

	return
}

// Unjoin remove the given chan from the current room, freeing memory if needed.
// If there is nobody remaining in the room, it is destroyed by calling Close.
func (g *WaitingGroupMap) Unjoin(room string, i chan interface{}) {
	g.data[room].mutex.Lock()
	for e := g.data[room].channels.Front(); e != nil; e = e.Next() {
		if e.Value == i {
			g.data[room].channels.Remove(e) // Remove element from list
			break
		}
	}
	if g.data[room].channels.Len() == 0 {
		g.Close(room)
		return
	}
	g.data[room].mutex.Unlock()
}

// Broadcast emits a message to every member of the room, including the sender.
func (g *WaitingGroupMap) Broadcast(room string, value interface{}) {
	g.data[room].mutex.Lock()
	g.data[room].oldMessages = append(g.data[room].oldMessages, value)
	g.data[room].mutex.Unlock()

	for e := g.data[room].channels.Front(); e != nil; e = e.Next() {
		e.Value.(chan interface{}) <- value
	}
}

// Close removes the room from the current WaitingGroupMap, closing all opened channels,
// and clearing oldMessages.
func (g *WaitingGroupMap) Close(room string) {
	for e := g.data[room].channels.Front(); e != nil; e = e.Next() {
		close(e.Value.(chan interface{}))
	}
	delete(g.data, room)
}

// CloseAll clears every available room in the current WaitingGroupMap.
func (g *WaitingGroupMap) CloseAll() {
	for k := range g.data {
		g.Close(k)
	}
}
