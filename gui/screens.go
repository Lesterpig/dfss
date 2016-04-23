package main

import (
	"dfss/dfssc/sign"
	"dfss/gui/authform"
	"dfss/gui/common"
	"dfss/gui/contractform"
	"dfss/gui/showcontract"
	"dfss/gui/signform"
	"dfss/gui/userform"
	"dfss/gui/welcome"
	"github.com/spf13/viper"
	"github.com/visualfc/goqt/ui"
)

func (w *window) showWelcome() {
	w.setScreen(welcome.NewWidget(
		w.showNewContractForm,
		func() { w.showShowContract("") },
		w.showFetchForm,
	))
}

func (w *window) showUserForm() {
	w.setScreen(userform.NewWidget(func(pwd string) {
		w.showAuthForm()
	}))
}

func (w *window) showAuthForm() {
	w.setScreen(authform.NewWidget(func() {
		w.showNewContractForm()
		w.addActions()
	}))
}

func (w *window) showNewContractForm() {
	w.setScreen(contractform.NewWidget())
}

func (w *window) showShowContract(filename string) {
	if filename == "" {
		home := viper.GetString("home_dir")
		filter := "Contract file (*.json);;Any (*.*)"
		filename = ui.QFileDialogGetOpenFileNameWithParentCaptionDirFilterSelectedfilterOptions(w, "Select the contract file", home, filter, &filter, 0)
		if filename == "" {
			return
		}
	}

	w.contract = showcontract.Load(filename)
	if w.contract == nil {
		common.ShowMsgBox("Unable to load file", true)
		return
	}
	w.setScreen(showcontract.NewWidget(w.contract, w.showSignForm))
}

func (w *window) showSignForm() {
	common.PasswordDialog(func(err error, pwd string) {
		widget := signform.NewWidget(w.contract, pwd)
		if widget == nil {
			common.ShowMsgBox("Unable to start the signing procedure", true)
			return
		}
		w.setScreen(widget)
	})
}

func (w *window) showFetchForm() {
	w.current.Q().SetDisabled(true)
	common.PasswordDialog(func(err error, pwd string) {
		if err != nil {
			w.current.Q().SetDisabled(false)
			return
		}

		dialog := ui.NewInputDialog()
		dialog.SetWindowTitle("Fetch a contract from the platform")
		dialog.SetLabelText("Please paste the contract identifier here:")
		dialog.Show()

		dialog.OnAccepted(func() {
			uuid := dialog.TextValue()
			path := viper.GetString("home_dir") + uuid + ".json"

			err := sign.FetchContract(pwd, uuid, path)

			if err != nil {
				common.ShowMsgBox(err.Error(), true)
				return
			}
			w.showShowContract(path)
			common.ShowMsgBox("Contract stored as "+path, false)
		})

		dialog.OnFinished(func(_ int32) {
			w.current.Q().SetDisabled(false)
		})
	})
}
