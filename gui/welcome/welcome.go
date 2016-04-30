// Package welcome provides the home screen for authenticated users.
package welcome

import (
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	*ui.QWidget
}

func NewWidget(onNew, onOpen, onFetch func()) *Widget {
	file := ui.NewFileWithName(":/welcome/welcome.ui")
	loader := ui.NewUiLoader()
	form := loader.Load(file)

	(ui.NewPushButtonFromDriver(form.FindChild("createButton"))).OnClicked(onNew)
	(ui.NewPushButtonFromDriver(form.FindChild("openButton"))).OnClicked(onOpen)
	(ui.NewPushButtonFromDriver(form.FindChild("fetchButton"))).OnClicked(onFetch)

	return &Widget{QWidget: form}
}

func (w *Widget) Q() *ui.QWidget {
	return w.QWidget
}

func (w *Widget) Tick() {}
