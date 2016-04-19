package contractform

import (
	"strings"

	"dfss/dfssc/sign"
	"dfss/gui/config"

	"github.com/spf13/viper"
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	*ui.QWidget

	signers *ui.QPlainTextEdit
}

func NewWidget() *Widget {
	file := ui.NewFileWithName(":/contractform/contractform.ui")
	loader := ui.NewUiLoader()
	form := loader.Load(file)

	fileField := ui.NewLineEditFromDriver(form.FindChild("fileField"))
	commentField := ui.NewPlainTextEditFromDriver(form.FindChild("commentField"))
	signersField := ui.NewPlainTextEditFromDriver(form.FindChild("signersField"))
	fileButton := ui.NewPushButtonFromDriver(form.FindChild("fileButton"))
	createButton := ui.NewPushButtonFromDriver(form.FindChild("createButton"))
	feedbackLabel := ui.NewLabelFromDriver(form.FindChild("feedbackLabel"))

	w := &Widget{
		QWidget: form,
		signers: signersField,
	}

	signersField.SetPlainText(viper.GetString("email") + "\n")

	fileButton.OnClicked(func() {
		filter := "Any (*.*)"
		filename := ui.QFileDialogGetOpenFileNameWithParentCaptionDirFilterSelectedfilterOptions(form, "Select contract file", config.GetHomeDir(), filter, &filter, 0)
		fileField.SetText(filename)
	})

	createButton.OnClicked(func() {
		form.SetDisabled(true)
		feedbackLabel.SetText("Please wait...")
		config.PasswordDialog(func(err error, pwd string) {
			if err != nil {
				form.SetDisabled(false)
				feedbackLabel.SetText("Aborted.")
				return // wrong key or rejection, aborting
			}

			err = sign.SendNewContract(
				pwd,
				fileField.Text(),
				commentField.ToPlainText(),
				w.SignersList(),
			)

			if err != nil {
				feedbackLabel.SetText(err.Error())
			} else {
				feedbackLabel.SetText("Contract successfully sent to signers!")
				fileField.SetText("")
			}
			form.SetDisabled(false)
		})
	})

	return w
}

func (w *Widget) SignersList() (list []string) {
	rawList := strings.Split(w.signers.ToPlainText(), "\n")

	for _, e := range rawList {
		clean := strings.TrimSpace(e)
		if clean != "" {
			list = append(list, clean)
		}
	}

	return
}

func (w *Widget) Q() *ui.QWidget {
	return w.QWidget
}

func (w *Widget) Tick() {}
