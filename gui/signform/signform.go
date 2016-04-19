package signform

import (
	"dfss/gui/config"
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	*ui.QWidget
}

func NewWidget(conf *config.Config, onAuth func()) *Widget {
	file := ui.NewFileWithName(":/signform/signform.ui")
	loader := ui.NewUiLoader()
	form := loader.Load(file)

	w := &Widget{QWidget: form}

	return w
}
