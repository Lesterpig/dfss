package gui

import (
	"time"
)

// Client represents a DFSSC instance
type Client struct {
	Name string
}

// EventType is used as an enum for event types, to differenciate promises, signatures...
type EventType int

const (
	PROMISE EventType = iota
	SIGNATURE
	OTHER
)

// Event represents a single signature event
type Event struct {
	Type     EventType
	Sender   int
	Receiver int
	Date     *time.Time
	Duration *time.Duration
}

// Scene holds the global scene for registered clients and signature events
type Scene struct {
	Clients []Client
	Events  []Event

	currentEvent int
}
