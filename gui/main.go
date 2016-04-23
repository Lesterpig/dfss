package main

import (
	"dfss"
	"dfss/dfssp/contract"
	"dfss/gui/config"
	"github.com/spf13/viper"
	"github.com/visualfc/goqt/ui"
)

type window struct {
	*ui.QMainWindow

	current  widget
	contract *contract.JSON
}

type widget interface {
	Q() *ui.QWidget
	Tick()
}

func init() {
	viper.Set("filename_ca", "ca.pem")
	viper.Set("filename_cert", "cert.pem")
	viper.Set("filename_key", "key.pem")
	viper.Set("filename_config", "config")
}

func main() {
	// Load configuration
	config.Load()

	// Start first window
	ui.Run(func() {
		w := &window{
			QMainWindow: ui.NewMainWindow(),
		}

		if viper.GetBool("authenticated") {
			w.addActions()
			w.showNewContractForm()
		} else if viper.GetBool("registered") {
			w.showAuthForm()
		} else {
			w.showUserForm()
		}

		timer := ui.NewTimerWithParent(w)
		timer.OnTimeout(func() {
			w.current.Tick()
		})
		timer.StartWithMsec(1000)

		w.SetWindowTitle("DFSS Client v" + dfss.Version)
		w.SetWindowIcon(ui.NewIconWithFilename(":/images/digital_signature_pen.png"))
		w.Show()
	})
}

func (w *window) addActions() {
	newAct := ui.NewActionWithTextParent("&New", w)
	newAct.SetShortcuts(ui.QKeySequence_New)
	newAct.OnTriggered(w.showNewContractForm)

	openAct := ui.NewActionWithTextParent("&Open", w)
	openAct.SetShortcuts(ui.QKeySequence_Open)
	openAct.OnTriggered(func() {
		w.showShowContract("")
	})

	fetchAct := ui.NewActionWithTextParent("&Fetch", w)
	fetchAct.OnTriggered(w.showFetchForm)

	aboutAct := ui.NewActionWithTextParent("&About", w)
	aboutAct.OnTriggered(func() {
		ui.QMessageBoxAbout(w, "About DFSS Client", about)
	})

	aboutQtAct := ui.NewActionWithTextParent("About &Qt", w)
	aboutQtAct.OnTriggered(func() {
		ui.QApplicationAboutQt()
	})

	fileMenu := w.MenuBar().AddMenuWithTitle("&File")
	fileMenu.AddAction(newAct)
	fileMenu.AddAction(openAct)
	fileMenu.AddSeparator()
	fileMenu.AddAction(fetchAct)

	helpMenu := w.MenuBar().AddMenuWithTitle("&Help")
	helpMenu.AddAction(aboutAct)
	helpMenu.AddSeparator()
	helpMenu.AddAction(aboutQtAct)
}

func (w *window) setScreen(wi widget) {
	old := w.CentralWidget()
	w.SetCentralWidget(wi.Q())
	w.current = wi
	if old != nil {
		old.DeleteLater()
	}
}
