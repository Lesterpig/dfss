package main

import (
	"dfss"
	"dfss/gui/authform"
	"dfss/gui/config"
	"dfss/gui/userform"
	"github.com/visualfc/goqt/ui"
)

const WIDTH = 650
const HEIGHT = 350

func main() {
	// Load configuration
	conf := config.Load()

	// Start first window
	ui.Run(func() {
		layout := ui.NewVBoxLayout()

		var newuser *userform.Widget
		var newauth *authform.Widget

		newauth = authform.NewWidget(&conf, func() {
			layout.RemoveWidget(newauth)
			newauth.Hide()
		})

		newuser = userform.NewWidget(&conf, func(pwd string) {
			layout.RemoveWidget(newuser)
			newuser.Hide()
			layout.AddWidget(newauth)
		})

		if conf.Authenticated {
			// TODO
		} else if conf.Registered {
			layout.AddWidget(newauth)
		} else {
			layout.AddWidget(newuser)
		}

		w := ui.NewWidget()
		w.SetLayout(layout)
		w.SetWindowTitle("DFSS Client v" + dfss.Version)
		w.SetFixedSizeWithWidthHeight(WIDTH, HEIGHT)
		w.Show()
	})
}
