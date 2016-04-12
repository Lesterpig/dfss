package contractform

import (
	"dfss/gui/config"
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	*ui.QWidget
}

func NewWidget(conf *config.Config) *Widget {
	file := ui.NewFileWithName(":/contractform/contractform.ui")
	loader := ui.NewUiLoader()
	form := loader.Load(file)

	return &Widget{QWidget: form}
}
