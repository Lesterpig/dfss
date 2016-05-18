// Package showcontract provides the show contract screen.
package showcontract

import (
	"io/ioutil"

	"dfss/dfssc/common"
	"dfss/dfssc/sign"
	"dfss/dfssp/contract"
	dialog "dfss/gui/common"
	"github.com/spf13/viper"
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	*ui.QWidget

	checked bool
}

func NewWidget(contract *contract.JSON, onSign func()) *Widget {
	file := ui.NewFileWithName(":/showcontract/showcontract.ui")
	loader := ui.NewUiLoader()
	form := loader.Load(file)

	fileField := ui.NewLabelFromDriver(form.FindChild("fileField"))
	commentField := ui.NewLabelFromDriver(form.FindChild("commentField"))
	signersField := ui.NewLabelFromDriver(form.FindChild("signersField"))
	informationField := ui.NewLabelFromDriver(form.FindChild("informationField"))
	signButton := ui.NewPushButtonFromDriver(form.FindChild("signButton"))
	fileButton := ui.NewPushButtonFromDriver(form.FindChild("fileButton"))

	fileField.SetText(contract.File.Name)
	commentField.SetText(contract.Comment)
	signersField.SetText(getSignersString(contract))
	informationField.SetText("Contract #" + contract.UUID + "\nCreated on " + contract.Date.Format("2006-01-02 15:04:05 MST") + ".")

	w := &Widget{
		QWidget: form,
		checked: false,
	}

	fileButton.OnClicked(func() {
		home := viper.GetString("home_dir")
		filter := "Any (*.*)"
		filename := ui.QFileDialogGetOpenFileNameWithParentCaptionDirFilterSelectedfilterOptions(nil, "Select the local contract to compare to "+contract.File.Name, home, filter, &filter, 0)
		if filename == "" {
			return
		}

		ok, err := sign.CheckContractHash(filename, contract.File.Hash)
		if err != nil {
			dialog.ShowMsgBox(err.Error(), true)
		} else if !ok {
			dialog.ShowMsgBox("The provided file is not the one hashed during contract creation, beware!", true)
		} else {
			dialog.ShowMsgBox("The provided file is correct for this contract!", false)
			w.checked = true
		}
	})

	signButton.OnClicked(func() {
		if !w.checked { // If the contract file has not been checked locally
			box := ui.NewMessageBox()
			box.SetWindowTitle("Warning")
			box.SetText("Do you want to check your local contract file before signing this DFSS contract?")
			box.SetIcon(ui.QMessageBox_Question)
			box.AddButton(ui.QMessageBox_Yes)
			box.AddButton(ui.QMessageBox_No)
			box.OnButtonClicked(func(b *ui.QAbstractButton) {
				if box.ButtonRole(b) == ui.QMessageBox_YesRole {
					fileButton.Click()
				} else {
					onSign()
				}
			})
			box.Show()
		} else {
			onSign()
		}
	})

	return w
}

func (w *Widget) Q() *ui.QWidget {
	return w.QWidget
}

func (w *Widget) Tick() {}

// Load loads a file and tries to unmarshall it into a DFSS contract
func Load(filename string) *contract.JSON {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	contract, err := common.UnmarshalDFSSFile(data)
	if err != nil {
		return nil
	}

	return contract
}

func getSignersString(contract *contract.JSON) string {
	var s string
	for i, signer := range contract.Signers {
		if i > 0 {
			s += ", "
		}
		s += signer.Email
	}
	return s
}
