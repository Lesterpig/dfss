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
			layout.RemoveWidget(newauth.W)
			newauth.W.Hide()
		})

		newuser = userform.NewWidget(&conf, func(pwd string) {
			layout.RemoveWidget(newuser.W)
			newuser.W.Hide()
			layout.AddWidget(newauth.W)
		})

		if conf.Authenticated {
			// TODO
		} else if conf.Registered {
			layout.AddWidget(newauth.W)
		} else {
			layout.AddWidget(newuser.W)
		}

		w := ui.NewWidget()
		w.SetLayout(layout)
		w.SetWindowTitle("DFSS Client v" + dfss.Version)
		w.SetFixedSizeWithWidthHeight(WIDTH, HEIGHT)
		w.Show()

		ev := ui.NewCloseEvent()
		w.CloseEvent(ev)
	})
}
