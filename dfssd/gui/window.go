// Package gui is the graphic part of the dfssd program.
package gui

// This file is the entry point of the gui package.
// It handles window instantiation and basic operations on it.

import (
	"math"
	"time"

	"dfss"
	"github.com/visualfc/goqt/ui"
)

// NewWindow creates and initialiaze a new dfssd main window.
func NewWindow() *Window {
	file := ui.NewFileWithName(":/widget.ui")
	loader := ui.NewUiLoader()
	widget := loader.Load(file)

	// Init main window
	window := ui.NewMainWindow()
	window.SetCentralWidget(widget)
	window.SetWindowTitle("DFSS Demonstrator v" + dfss.Version)

	w := &Window{
		QMainWindow: window,
		scene:       &Scene{},
	}
	w.InstallEventFilter(w)

	// Load dynamic elements from driver
	w.logField = ui.NewTextEditFromDriver(widget.FindChild("logField"))
	w.graphics = ui.NewGraphicsViewFromDriver(widget.FindChild("graphicsView"))
	w.progress = ui.NewLabelFromDriver(widget.FindChild("progressLabel"))

	w.playButton = ui.NewPushButtonFromDriver(widget.FindChild("playButton"))
	w.stopButton = ui.NewPushButtonFromDriver(widget.FindChild("stopButton"))
	w.replayButton = ui.NewPushButtonFromDriver(widget.FindChild("replayButton"))

	// Load pixmaps
	w.pixmaps = map[string]*ui.QPixmap{
		"ttp":      ui.NewPixmapWithFilenameFormatFlags(":/images/server_key.png", "", ui.Qt_AutoColor),
		"platform": ui.NewPixmapWithFilenameFormatFlags(":/images/server_connect.png", "", ui.Qt_AutoColor),
	}

	// Load icons
	w.addIcons()

	// Add actions
	w.addActions()
	w.initScene()
	w.initTimer()

	w.StatusBar().ShowMessage("Ready")
	w.PrintQuantumInformation()
	return w
}

// OnResizeEvent is called by Qt each time an user tries to resize the window.
// We have to redraw the whole scene to adapt.
func (w *Window) OnResizeEvent(ev *ui.QResizeEvent) bool {
	w.initScene()
	return true
}

// Log is used to print a new line in the log area of the window.
// It should be thread-safe.
func (w *Window) Log(str string) {
	str = time.Now().Format("[15:04:05.000] ") + str
	w.logField.Append(str)
	w.logField.EnsureCursorVisible()
}

// addIcons adds icons to control buttons, since we cannot add them directly in QtCreator.
func (w *Window) addIcons() {
	w.SetWindowIcon(ui.NewIconWithFilename(":/images/node_magnifier.png"))

	var i *ui.QIcon
	i = ui.NewIconWithFilename(":/images/control_play_blue.png")
	i.AddFileWithFilenameSizeModeState(":/images/control_play.png", ui.NewSizeWithWidthHeight(32, 32), ui.QIcon_Disabled, ui.QIcon_Off)
	w.playButton.SetIcon(i)

	i = ui.NewIconWithFilename(":/images/control_pause_blue.png")
	i.AddFileWithFilenameSizeModeState(":/images/control_pause.png", ui.NewSizeWithWidthHeight(32, 32), ui.QIcon_Disabled, ui.QIcon_Off)
	w.stopButton.SetIcon(i)

	i = ui.NewIconWithFilename(":/images/control_rewind_blue.png")
	w.replayButton.SetIcon(i)
}

// addActions adds action listenners to interactive parts of the window.
func (w *Window) addActions() {
	// MENU BAR
	openAct := ui.NewActionWithTextParent("&Open", w)
	openAct.SetShortcuts(ui.QKeySequence_Open)
	openAct.SetStatusTip("Open a demonstration file")
	openAct.OnTriggered(func() {
		filename := ui.QFileDialogGetOpenFileName()
		if filename != "" {
			w.Open(filename)
		}
	})

	saveAct := ui.NewActionWithTextParent("&Save", w)
	saveAct.SetShortcuts(ui.QKeySequence_Save)
	saveAct.SetStatusTip("Save a demonstration file")
	saveAct.OnTriggered(func() {
		filename := ui.QFileDialogGetSaveFileName()
		if filename != "" {
			w.Save(filename)
		}
	})

	w.MenuBar().AddAction(openAct)
	w.MenuBar().AddAction(saveAct)

	// SIMULATION CONTROL
	w.playButton.OnClicked(func() {
		w.playButton.SetDisabled(true)
		w.stopButton.SetDisabled(false)
		w.timer.StartWithMsec(speed)
	})

	w.stopButton.OnClicked(func() {
		w.playButton.SetDisabled(false)
		w.stopButton.SetDisabled(true)
		w.timer.Stop()
	})
	w.stopButton.SetDisabled(true)

	w.replayButton.OnClicked(func() {
		w.RemoveArrows()
		w.scene.currentEvent = 0
		w.PrintQuantumInformation()
	})
}

// initScene creates the Qt graphic scene associated to our custom scene.
// It draws the base circle, clients and servers, and do some memory management for us.
func (w *Window) initScene() {
	// Save old scene
	oldScene := w.graphics.Scene()

	scene := ui.NewGraphicsScene()
	w.graphics.SetScene(scene)

	// Draw base circle
	w.circleSize = math.Min(float64(w.graphics.Width()), float64(w.graphics.Height())) - 50
	r := w.circleSize / 2
	scene.AddEllipseFWithXYWidthHeightPenBrush(-r, -r, w.circleSize, w.circleSize, pen_gray, brush_none)

	// Draw clients
	w.DrawClients()
	w.DrawServers()

	// Purge
	if oldScene != nil {
		w.RemoveArrows()
		defer oldScene.Delete()
	}
}
