package authform

import (
	"dfss/dfssc/user"
	"dfss/gui/common"
	"github.com/spf13/viper"
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	*ui.QWidget
}

func NewWidget(onAuth func()) *Widget {
	file := ui.NewFileWithName(":/authform/authform.ui")
	loader := ui.NewUiLoader()
	form := loader.Load(file)

	tokenField := ui.NewLineEditFromDriver(form.FindChild("tokenField"))
	authButton := ui.NewPushButtonFromDriver(form.FindChild("authButton"))

	authButton.OnClicked(func() {
		form.SetDisabled(true)
		err := user.Authenticate(
			viper.GetString("email"),
			tokenField.Text(),
		)
		form.SetDisabled(false)
		if err != nil {
			common.ShowMsgBox(err.Error(), true)
			tokenField.SetFocus()
			tokenField.SelectAll()
		} else {
			onAuth()
		}
	})

	return &Widget{QWidget: form}
}

func (w *Widget) Q() *ui.QWidget {
	return w.QWidget
}

func (w *Widget) Tick() {}
