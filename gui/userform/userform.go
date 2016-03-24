package userform

import (
	"io/ioutil"
	"os"
	osuser "os/user"
	"path/filepath"

	"dfss/dfssc/user"
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	W *ui.QWidget
}

func NewWidget() *Widget {
	file := ui.NewFileWithName(":/userform/userform.ui")
	loader := ui.NewUiLoader()
	form := loader.Load(file)

	emailField := ui.NewLineEditFromDriver(form.FindChild("emailField"))
	hostField := ui.NewLineEditFromDriver(form.FindChild("hostField"))
	passwordField := ui.NewLineEditFromDriver(form.FindChild("passwordField"))
	passwordField.SetEchoMode(ui.QLineEdit_Password)

	feedbackLabel := ui.NewLabelFromDriver(form.FindChild("feedbackLabel"))
	registerButton := ui.NewPushButtonFromDriver(form.FindChild("registerButton"))

	home := getHomeDir()
	fileDialog := ui.NewFileDialogWithParentCaptionDirectoryFilter(nil, "Select the CA file for the platform", home, "Root Certificates (*.pem);;Any (*.*)")

	// Events

	registerButton.OnClicked(func() {
		form.SetDisabled(true)
		feedbackLabel.SetText("Registration in progress...")
		fileDialog.Open()
	})

	fileDialog.OnFileSelected(func(ca string) {
		fileDialog.Hide()
		caDest := home + "ca.pem"
		_ = copyCA(ca, caDest)

		err := user.Register(
			home+"ca.pem",
			home+"cert.pem",
			home+"key.pem",
			hostField.Text(),
			passwordField.Text(),
			"", "", "", emailField.Text(), 2048,
		)
		if err != nil {
			feedbackLabel.SetText(err.Error())
		} else {
			feedbackLabel.SetText("Registration done! Please check your mails.")
		}
		form.SetDisabled(false)
	})

	fileDialog.OnRejected(func() {
		form.SetDisabled(false)
		feedbackLabel.SetText("Registration aborted.")
	})

	return &Widget{W: form}
}

func getHomeDir() string {
	u, err := osuser.Current()
	if err != nil {
		return ""
	}

	dfssPath := filepath.Join(u.HomeDir, ".dfss")
	if err := os.MkdirAll(dfssPath, os.ModeDir|0700); err != nil {
		return ""
	}

	return dfssPath + string(filepath.Separator)
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
