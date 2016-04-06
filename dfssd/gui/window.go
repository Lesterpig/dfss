package gui

import (
	"dfss"
	"github.com/visualfc/goqt/ui"
)

type Window struct {
	*ui.QMainWindow

	logField *ui.QTextEdit
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
	}

	// Load dynamic elements from driver
	w.logField = ui.NewTextEditFromDriver(widget.FindChild("logField"))

	// Add actions
	w.addActions()

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
