package config

// This file handles basic feedback messages.

import (
	"errors"
	"io/ioutil"

	"dfss/auth"

	"github.com/spf13/viper"
	"github.com/visualfc/goqt/ui"
)

// ShowMsgBox pops a new Qt Msg Box containing the appropriate message regarding the critical flag.
func ShowMsgBox(content string, critical bool) {
	m := ui.NewMessageBox()
	m.SetText(content)
	if critical {
		m.SetWindowTitle("Error")
		m.SetIcon(ui.QMessageBox_Critical)
	} else {
		m.SetWindowTitle("Information")
		m.SetIcon(ui.QMessageBox_Information)
	}
	m.Exec()
}

// PasswordDialog checks the current private key for any passphrase.
// If the key is protected, it spawns an inputDialog window to ask the user's passphrase,
// and then calls the callback function with the result.
//
// The callback is always called, even when an error occurs.
func PasswordDialog(callback func(err error, pwd string)) {
	// Try to get private key
	data, err := ioutil.ReadFile(viper.GetString("file_key"))
	if err != nil {
		callback(err, "")
		return
	}

	if !auth.IsPEMEncrypted(data) {
		callback(nil, "")
		return
	}

	dialog := ui.NewInputDialog()
	dialog.SetWindowTitle("Encrypted private key")
	dialog.SetLabelText("Please type your password to proceed:")
	dialog.SetTextEchoMode(ui.QLineEdit_Password)

	dialog.OnRejected(func() {
		callback(errors.New("user rejected"), "")
	})

	dialog.OnAccepted(func() {
		pwd := dialog.TextValue()
		if pwd == "" {
			pwd = " " // doing this to force the "wrong password" error msg
		}
		callback(nil, pwd)
	})

	dialog.Open()
}
