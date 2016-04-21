package main

import (
	"dfss"
	"dfss/gui/authform"
	"dfss/gui/config"
	"dfss/gui/contractform"
	"dfss/gui/signform"
	"dfss/gui/userform"
	"github.com/spf13/viper"
	"github.com/visualfc/goqt/ui"
)

type window struct {
	*ui.QMainWindow

	current widget
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
	openAct := ui.NewActionWithTextParent("&Open", w)
	openAct.SetShortcuts(ui.QKeySequence_Open)
	openAct.OnTriggered(func() {
		w.showSignForm()
	})
	w.MenuBar().AddAction(openAct)
}

func (w *window) setScreen(wi widget) {
	old := w.CentralWidget()
	w.SetCentralWidget(wi.Q())
	w.current = wi
	if old != nil {
		old.DeleteLater()
	}
}

func (w *window) showUserForm() {
	w.setScreen(userform.NewWidget(func(pwd string) {
		w.showAuthForm()
	}))
}

func (w *window) showAuthForm() {
	w.setScreen(authform.NewWidget(func() {
		w.showNewContractForm()
		w.addActions()
	}))
}

func (w *window) showNewContractForm() {
	w.setScreen(contractform.NewWidget())
}

func (w *window) showSignForm() {
	home := viper.GetString("home_dir")
	filter := "Contract file (*.json);;Any (*.*)"
	filename := ui.QFileDialogGetOpenFileNameWithParentCaptionDirFilterSelectedfilterOptions(w, "Select the contract file", home, filter, &filter, 0)
	if filename != "" {
		config.PasswordDialog(func(err error, pwd string) {
			widget := signform.NewWidget(filename, pwd)
			if widget != nil {
				w.setScreen(widget)
			}
		})
	}
}
