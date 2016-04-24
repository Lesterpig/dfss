package userform

import (
	"io/ioutil"

	"dfss/dfssc/user"
	"dfss/gui/common"
	"dfss/gui/config"

	"github.com/spf13/viper"
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	*ui.QWidget
}

func NewWidget(onRegistered func(pw string)) *Widget {
	file := ui.NewFileWithName(":/userform/userform.ui")
	loader := ui.NewUiLoader()
	form := loader.Load(file)

	emailField := ui.NewLineEditFromDriver(form.FindChild("emailField"))
	hostField := ui.NewLineEditFromDriver(form.FindChild("hostField"))
	passwordField := ui.NewLineEditFromDriver(form.FindChild("passwordField"))
	passwordField.SetEchoMode(ui.QLineEdit_Password)
	registerButton := ui.NewPushButtonFromDriver(form.FindChild("registerButton"))

	home := viper.GetString("home_dir")

	// Events

	registerButton.OnClicked(func() {
		form.SetDisabled(true)
		filter := "Root Certificates (*.pem);;Any (*.*)"
		caFilename := ui.QFileDialogGetOpenFileNameWithParentCaptionDirFilterSelectedfilterOptions(form, "Select the CA file for the platform", home, filter, &filter, 0)
		_ = copyCA(caFilename)
		viper.Set("platform_addrport", hostField.Text())

		err := user.Register(
			passwordField.Text(),
			"", "", "", emailField.Text(), 2048,
		)
		if err != nil {
			common.ShowMsgBox(err.Error(), true)
		} else {
			viper.Set("email", emailField.Text())
			onRegistered(passwordField.Text())
			config.Save()
		}
		form.SetDisabled(false)
	})

	return &Widget{QWidget: form}
}

func copyCA(from string) error {
	if from == viper.GetString("file_ca") {
		return nil
	}

	file, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(viper.GetString("file_ca"), file, 0600)
}

func (w *Widget) Q() *ui.QWidget {
	return w.QWidget
}

func (w *Widget) Tick() {}
