package userform

import (
	"io/ioutil"

	"dfss/dfssc/user"
	"dfss/gui/config"
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	*ui.QWidget
}

func NewWidget(conf *config.Config, onRegistered func(pw string)) *Widget {
	file := ui.NewFileWithName(":/userform/userform.ui")
	loader := ui.NewUiLoader()
	form := loader.Load(file)

	emailField := ui.NewLineEditFromDriver(form.FindChild("emailField"))
	hostField := ui.NewLineEditFromDriver(form.FindChild("hostField"))
	passwordField := ui.NewLineEditFromDriver(form.FindChild("passwordField"))
	passwordField.SetEchoMode(ui.QLineEdit_Password)

	feedbackLabel := ui.NewLabelFromDriver(form.FindChild("feedbackLabel"))
	registerButton := ui.NewPushButtonFromDriver(form.FindChild("registerButton"))

	home := config.GetHomeDir()

	// Events

	registerButton.OnClicked(func() {
		form.SetDisabled(true)
		feedbackLabel.SetText("Registration in progress...")
		filter := "Root Certificates (*.pem);;Any (*.*)"
		caFilename := ui.QFileDialogGetOpenFileNameWithParentCaptionDirFilterSelectedfilterOptions(form, "Select the CA file for the platform", home, filter, &filter, 0)
		caDest := home + config.CAFile
		_ = copyCA(caFilename, caDest)

		err := user.Register(
			caDest,
			home+config.CertFile,
			home+config.KeyFile,
			hostField.Text(),
			passwordField.Text(),
			"", "", "", emailField.Text(), 2048,
		)
		if err != nil {
			feedbackLabel.SetText(err.Error())
		} else {
			conf.Email = emailField.Text()
			conf.Platform = hostField.Text()
			onRegistered(passwordField.Text())
			config.Save(*conf)
		}
		form.SetDisabled(false)
	})

	return &Widget{QWidget: form}
}

func copyCA(from string, to string) error {
	if from == to {
		return nil
	}

	file, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(to, file, 0600)
}

func (w *Widget) Q() *ui.QWidget {
	return w.QWidget
}

func (w *Widget) Tick() {}
