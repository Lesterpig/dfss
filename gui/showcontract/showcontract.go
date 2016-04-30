// Package showcontract provides the show contract screen.
package showcontract

import (
	"encoding/json"
	"io/ioutil"

	"dfss/dfssp/contract"
	"github.com/visualfc/goqt/ui"
)

type Widget struct {
	*ui.QWidget
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

	fileField.SetText(contract.File.Name)
	commentField.SetText(contract.Comment)
	signersField.SetText(getSignersString(contract))
	informationField.SetText("Contract #" + contract.UUID + "\nCreated on " + contract.Date.Format("2006-01-02 15:04:05 MST") + ".")

	signButton.OnClicked(onSign)

	return &Widget{
		QWidget: form,
	}
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

	contract := new(contract.JSON)
	err = json.Unmarshal(data, contract)
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
