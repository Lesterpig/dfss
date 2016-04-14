package main

import (
	"dfss"
	"dfss/gui/authform"
	"dfss/gui/config"
	"dfss/gui/contractform"
	"dfss/gui/userform"
	"github.com/visualfc/goqt/ui"
)

func main() {
	// Load configuration
	conf := config.Load()

	// Start first window
	ui.Run(func() {
		window := ui.NewMainWindow()

		var newuser *userform.Widget
		var newauth *authform.Widget
		var newcontract *contractform.Widget

		newauth = authform.NewWidget(&conf, func() {
			window.SetCentralWidget(newcontract)
		})

		newuser = userform.NewWidget(&conf, func(pwd string) {
			window.SetCentralWidget(newauth)
		})

		newcontract = contractform.NewWidget(&conf)

		if conf.Authenticated {
			window.SetCentralWidget(newcontract)
		} else if conf.Registered {
			window.SetCentralWidget(newauth)
		} else {
			window.SetCentralWidget(newuser)
		}

		window.SetWindowTitle("DFSS Client v" + dfss.Version)
		window.Show()
	})
}
