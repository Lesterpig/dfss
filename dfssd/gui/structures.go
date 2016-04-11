package gui

// This file stores strucutures used in GUI for fast documentation.

import (
	"time"

	"github.com/visualfc/goqt/ui"
)

// Window contains all information used to make the demonstrator works.
// It extends QMainWindow and cache several graphic informations.
// Do not attempt to instantiante it directly, use `NewWindow` function instead.
type Window struct {
	*ui.QMainWindow

	logField     *ui.QTextEdit
	graphics     *ui.QGraphicsView
	progress     *ui.QLabel
	playButton   *ui.QPushButton
	stopButton   *ui.QPushButton
	replayButton *ui.QPushButton
	quantumField *ui.QSpinBox
	speedSlider  *ui.QSlider
	scene        *Scene
	circleSize   float64
	pixmaps      map[string]*ui.QPixmap

	currentArrows []*ui.QGraphicsPathItem
	timer         *ui.QTimer
}

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
	Date     time.Time
}

// Scene holds the global scene for registered clients and signature events
type Scene struct {
	Clients []Client
	Events  []Event

	currentTime  time.Time
	currentEvent int
}
