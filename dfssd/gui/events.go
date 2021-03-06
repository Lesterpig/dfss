package gui

// This file handles event timers and imports.

import (
	"fmt"
	"math"
	"time"

	"dfss/dfssd/api"
	"github.com/visualfc/goqt/ui"
)

// AddEvent interprets an incoming event into a graphic one.
// Expected format:
//
// Timestamp: unix nano timestamp
// Identifier: either "platform", "ttp" or "<email>"
// Log: one of the following
//      "sync with <entity>"
//			"sent promise to <email>"
//      "sent signature to <email>"
//      "contacting TTP with resolve index <index>"
//      "sent abort token to <email>"
//      "send signed contract to <email>"
//
// Other messages are currently ignored.
func (w *Window) AddEvent(e *api.Log) {
	event := Event{
		Sender: w.scene.identifierToIndex(e.Identifier),
		Date:   time.Unix(0, e.Timestamp),
	}

	w.Log(fmt.Sprint(e.Identifier, " ", e.Log))

	var receiver string
	var index int
	if n, _ := fmt.Sscanf(e.Log, "successfully sent promise to %s", &receiver); n > 0 {
		event.Type = PROMISE
		event.Receiver = w.scene.identifierToIndex(receiver)
	} else if n, _ := fmt.Sscanf(e.Log, "successfully sent signature to %s", &receiver); n > 0 {
		event.Type = SIGNATURE
		event.Receiver = w.scene.identifierToIndex(receiver)
	} else if n, _ := fmt.Sscanf(e.Log, "contacting TTP with resolve index %d", &index); n > 0 {
		event.Type = TTPCALL
		event.Receiver = w.scene.identifierToIndex("ttp")
	} else if n, _ := fmt.Sscanf(e.Log, "sent abort token to %s", &receiver); n > 0 {
		event.Type = TTPABORT
		event.Receiver = w.scene.identifierToIndex(receiver)
	} else if n, _ := fmt.Sscanf(e.Log, "sent signed contract to %s", &receiver); n > 0 {
		event.Type = TTPSIGNED
		event.Receiver = w.scene.identifierToIndex(receiver)
	} else if n, _ := fmt.Sscanf(e.Log, "sync with %s", &receiver); n > 0 {
		event.Type = SYNC
		event.Receiver = w.scene.identifierToIndex(receiver)
	} else {
		return
	}

	w.scene.Events = append(w.scene.Events, event)
}

// DrawEvent triggers the appropriate draw action for a spectific event.
func (w *Window) DrawEvent(e *Event) {
	xa, ya := w.GetClientPosition(e.Sender)
	xb, yb := w.GetClientPosition(e.Receiver)

	var color string
	switch e.Type {
	case PROMISE:
		color = "blue"
	case SIGNATURE:
		color = "green"
	case TTPCALL:
		color = "purple"
	case TTPABORT:
		color = "red"
	case TTPSIGNED:
		color = "lime"
	default:
		color = "grey"
	}

	w.DrawArrow(xa, ya, xb, yb, colors[color])
}

// PrintQuantumInformation triggers the update of the "x / y" quantum information.
func (w *Window) PrintQuantumInformation() {
	if len(w.scene.Events) == 0 {
		w.progress.SetText("No event")
		return
	}

	quantum := float64(w.quantumField.Value() * 1000)

	beginning := w.scene.Events[0].Date.UnixNano()
	totalDuration := w.scene.Events[len(w.scene.Events)-1].Date.UnixNano() - beginning
	nbQuantum := math.Floor(float64(totalDuration)/quantum) + 1
	durationFromBeginning := w.scene.currentTime.UnixNano() - beginning
	currentQuantum := math.Ceil(float64(durationFromBeginning)/quantum) + 1

	if w.scene.currentEvent == 0 {
		currentQuantum = 0
	}
	w.progress.SetText(fmt.Sprint(currentQuantum, " / ", nbQuantum))
}

// initTimer is called during window initialization. It initializes the timeout signal called for each refresh.
func (w *Window) initTimer() {
	w.timer = ui.NewTimerWithParent(w)

	lastNbOfClients := len(w.scene.Clients)

	w.timer.OnTimeout(func() {
		nbEvents := len(w.scene.Events)
		if w.scene.currentEvent >= nbEvents {
			w.stopButton.Click()
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

		quantum := time.Duration(w.quantumField.Value()) * time.Microsecond
		endOfQuantum := w.scene.currentTime.Add(quantum)

		drawnEvents := 0
		for i := w.scene.currentEvent; i < nbEvents; i++ {
			e := w.scene.Events[i]

			if e.Date.After(endOfQuantum) || e.Date.Equal(endOfQuantum) {
				break
			}

			w.DrawEvent(&e)
			w.scene.currentEvent++
			drawnEvents++
		}

		w.PrintQuantumInformation()
		w.scene.currentTime = endOfQuantum

		if w.speedSlider.Value() == 0 && drawnEvents > 0 {
			w.stopButton.Click() // step-by-step
		}
	})
}

// identifierToIndex is used to retrieve a client index from its name, inserting a new client if needed.
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
