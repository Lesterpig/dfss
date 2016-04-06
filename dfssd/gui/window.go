package gui

import (
	"math"

	"dfss"
	"github.com/visualfc/goqt/ui"
)

type Window struct {
	*ui.QMainWindow

	logField   *ui.QTextEdit
	graphics   *ui.QGraphicsView
	scene      *Scene
	circleSize float64
	pixmaps    map[string]*ui.QPixmap
}

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

	// Load pixmaps
	w.pixmaps = map[string]*ui.QPixmap{
		"ttp":      ui.NewPixmapWithFilenameFormatFlags(":/images/server_key.png", "", ui.Qt_AutoColor),
		"platform": ui.NewPixmapWithFilenameFormatFlags(":/images/server_connect.png", "", ui.Qt_AutoColor),
	}

	// Load icon
	w.SetWindowIcon(ui.NewIconWithFilename(":/images/node_magnifier.png"))

	// Add actions
	w.addActions()
	w.initScene()

	// TEST ONLY
	w.scene.Clients = []Client{
		Client{"signer1@lesterpig.com"},
		Client{"signer2@insa-rennes.fr"},
		Client{"signer3@dfss.com"},
	}

	w.StatusBar().ShowMessage("Ready")
	return w
}

func (w *Window) Log(str string) {
	w.logField.Append(str)
	w.logField.EnsureCursorVisible()
}

func (w *Window) addActions() {
	openAct := ui.NewActionWithTextParent("&Open", w)
	openAct.SetShortcuts(ui.QKeySequence_Open)
	openAct.SetStatusTip("Open a demonstration file")

	saveAct := ui.NewActionWithTextParent("&Save", w)
	saveAct.SetShortcuts(ui.QKeySequence_Save)
	saveAct.SetStatusTip("Save a demonstration file")

	w.MenuBar().AddAction(openAct)
	w.MenuBar().AddAction(saveAct)
}

func (w *Window) OnResizeEvent(ev *ui.QResizeEvent) bool {
	w.initScene()
	return true
}

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

	// TEST
	xa, ya := w.GetClientPosition(0)
	xb, yb := w.GetClientPosition(1)
	w.DrawArrow(xa, ya, xb, yb, colors["red"])
	w.DrawArrow(xb, yb, xa, ya, colors["red"])

	// Purge
	if oldScene != nil {
		oldScene.Delete()
	}
}
