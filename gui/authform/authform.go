package authform

import (
	"dfss/dfssc/user"
	"dfss/gui/config"
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	*ui.QWidget
}

func NewWidget(conf *config.Config, onAuth func()) *Widget {
	file := ui.NewFileWithName(":/authform/authform.ui")
	loader := ui.NewUiLoader()
	form := loader.Load(file)

	tokenField := ui.NewLineEditFromDriver(form.FindChild("tokenField"))
	feedbackLabel := ui.NewLabelFromDriver(form.FindChild("feedbackLabel"))
	authButton := ui.NewPushButtonFromDriver(form.FindChild("authButton"))

	home := config.GetHomeDir()
	authButton.OnClicked(func() {
		form.SetDisabled(true)
		err := user.Authenticate(
			home+config.CAFile,
			home+config.CertFile,
			conf.Platform,
			conf.Email,
			tokenField.Text(),
		)
		form.SetDisabled(false)
		if err != nil {
			feedbackLabel.SetText(err.Error())
			tokenField.SetFocus()
			tokenField.SelectAll()
		} else {
			onAuth()
		}
	})

	return &Widget{QWidget: form}
}
