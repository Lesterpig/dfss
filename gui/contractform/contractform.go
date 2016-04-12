package contractform

import (
	"strings"

	"dfss/dfssc/sign"
	"dfss/gui/config"
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	*ui.QWidget

	signers *ui.QPlainTextEdit
}

func NewWidget(conf *config.Config) *Widget {
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

	signersField.SetPlainText(conf.Email + "\n")

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

			home := config.GetHomeDir()
			err = sign.SendNewContract(
				home+config.CAFile,
				home+config.CertFile,
				home+config.KeyFile,
				conf.Platform,
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
