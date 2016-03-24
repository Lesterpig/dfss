package main

import (
	"dfss"
	"dfss/gui/userform"
	"github.com/visualfc/goqt/ui"
)

const WIDTH = 650
const HEIGHT = 350

func main() {
	ui.Run(func() {
		form := userform.NewWidget()

		layout := ui.NewVBoxLayout()
		layout.AddWidget(form.W)

		w := ui.NewWidget()
		w.SetLayout(layout)
		w.SetWindowTitle("DFSS Client v" + dfss.Version)
		w.SetFixedSizeWithWidthHeight(WIDTH, HEIGHT)
		w.Show()
	})
}
