package gui

import (
	"fmt"
	"time"
	"math"
)

// TEMPORARY
const quantum = 100 // discretization argument for events (ns)
const speed = 500 // duration of a quantum (ms)

func (w *Window) StartSimulation() {
	w.ticker = time.NewTicker(speed * time.Millisecond)
	go subroutine(w)
}

func (w *Window) DrawEvent(e *Event) {
	xa, ya := w.GetClientPosition(e.Sender)
	xb, yb := w.GetClientPosition(e.Receiver)
	w.DrawArrow(xa, ya, xb, yb, colors["red"])
}

func (w *Window) PrintQuantumInformation() {
	beginning := w.scene.Events[0].Date.UnixNano()
	totalDuration := w.scene.Events[len(w.scene.Events) - 1].Date.UnixNano() - beginning
	nbQuantum := math.Ceil(float64(totalDuration) / quantum)
	durationFromBeginning := w.scene.currentTime.UnixNano() - beginning
	currentQuantum := math.Ceil(float64(durationFromBeginning) / quantum)+1

	w.progress.SetText(fmt.Sprint(currentQuantum, " / ", nbQuantum))
}

func subroutine(w *Window) {
	for _ = range w.ticker.C {
		nbEvents := len(w.scene.Events)
		if w.scene.currentEvent >= nbEvents {
			// TODO disable looping if needed
			w.RemoveArrows()
			w.scene.currentEvent = 0 // loop
			w.Log("Restarting simulation...")
			continue
		}

		// Remove arrows from last tick
		w.RemoveArrows()

		// Check that we have a least one event to read
		if nbEvents == 0 {
			continue
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
	}
}
