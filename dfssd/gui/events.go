package gui

import (
	"fmt"
	"math"
	"time"

	"dfss/dfssd/api"
	"github.com/visualfc/goqt/ui"
)

// TEMPORARY
const quantum = 100 // discretization argument for events (ns)
const speed = 1000  // duration of a quantum (ms)

// AddEvent interprets an incoming event into a graphic one.
// Expected format:
//
// Timestamp: unix nano timestamp
// Identifier: either "platform", "ttp" or "<email>"
// Log: one of the following
//			"sent promise to <email>"
//      "sent signature to <email>"
//
// Other messages are currently ignored.
func (w *Window) AddEvent(e *api.Log) {
	event := Event{
		Sender: w.scene.identifierToIndex(e.Identifier),
		Date:   time.Unix(0, e.Timestamp),
	}

	w.Log(fmt.Sprint(e.Identifier, " ", e.Log))

	var receiver string
	if n, _ := fmt.Sscanf(e.Log, "sent promise to %s", &receiver); n > 0 {
		event.Type = PROMISE
		event.Receiver = w.scene.identifierToIndex(receiver)
	} else if n, _ := fmt.Sscanf(e.Log, "sent signature to %s", &receiver); n > 0 {
		event.Type = SIGNATURE
		event.Receiver = w.scene.identifierToIndex(receiver)
	}

	if receiver != "" {
		w.scene.Events = append(w.scene.Events, event)
	}
}

func (w *Window) DrawEvent(e *Event) {
	xa, ya := w.GetClientPosition(e.Sender)
	xb, yb := w.GetClientPosition(e.Receiver)

	var color string
	switch e.Type {
	case PROMISE:
		color = "blue"
	case SIGNATURE:
		color = "green"
	default:
		color = "black"
	}

	w.DrawArrow(xa, ya, xb, yb, colors[color])
}

func (w *Window) PrintQuantumInformation() {
	if len(w.scene.Events) == 0 {
		w.progress.SetText("No event")
		return
	}

	beginning := w.scene.Events[0].Date.UnixNano()
	totalDuration := w.scene.Events[len(w.scene.Events)-1].Date.UnixNano() - beginning
	nbQuantum := math.Max(1, math.Ceil(float64(totalDuration)/quantum))
	durationFromBeginning := w.scene.currentTime.UnixNano() - beginning
	currentQuantum := math.Ceil(float64(durationFromBeginning)/quantum) + 1

	if w.scene.currentEvent == 0 {
		currentQuantum = 0
	}
	w.progress.SetText(fmt.Sprint(currentQuantum, " / ", nbQuantum))
}

func (w *Window) initTimer() {
	w.timer = ui.NewTimerWithParent(w)

	lastNbOfClients := len(w.scene.Clients)

	w.timer.OnTimeout(func() {
		nbEvents := len(w.scene.Events)
		if w.scene.currentEvent >= nbEvents {
			w.replayButton.Click()
			return
		}

		// Remove arrows from last tick
		w.RemoveArrows()

		// Check that we have a least one event to read
		if nbEvents == 0 {
			return
		}

		// Check if need to redraw everything
		if lastNbOfClients != len(w.scene.Clients) {
			w.initScene()
		}

		// Init first time
		if w.scene.currentEvent == 0 {
			w.scene.currentTime = w.scene.Events[0].Date
		}

		endOfQuantum := w.scene.currentTime.Add(quantum * time.Nanosecond)

		for i := w.scene.currentEvent; i < nbEvents; i++ {
			e := w.scene.Events[i]

			if e.Date.After(endOfQuantum) || e.Date.Equal(endOfQuantum) {
				break
			}

			w.DrawEvent(&e)
			w.scene.currentEvent++
		}

		w.PrintQuantumInformation()
		w.scene.currentTime = endOfQuantum
	})
}

func (s *Scene) identifierToIndex(identifier string) int {
	if identifier == "platform" {
		return -1
	}
	if identifier == "ttp" {
		return -2
	}

	for i, c := range s.Clients {
		if c.Name == identifier {
			return i
		}
	}

	s.Clients = append(s.Clients, Client{Name: identifier})
	return len(s.Clients) - 1
}
