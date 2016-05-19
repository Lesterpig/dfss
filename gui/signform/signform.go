// Package signform provides the signature screen.
package signform

import (
	"dfss/dfssc/sign"
	"dfss/dfssp/contract"
	"github.com/spf13/viper"
	"github.com/visualfc/goqt/ui"
)

type line struct {
	status              sign.SignerStatus
	info, mail          string
	cellA, cellB, cellC *ui.QTableWidgetItem
}

type Widget struct {
	*ui.QWidget

	manager                  *sign.SignatureManager
	contract                 *contract.JSON
	table                    *ui.QTableWidget
	progressBar              *ui.QProgressBar
	feedbackLabel            *ui.QLabel
	cancelButton             *ui.QPushButton
	lines                    []line
	statusMax, statusCurrent int32
	feedback                 string
	running                  bool
}

func NewWidget(contract *contract.JSON, pwd string) *Widget {
	loadIcons()
	file := ui.NewFileWithName(":/signform/signform.ui")
	loader := ui.NewUiLoader()
	form := loader.Load(file)
	w := &Widget{
		QWidget:  form,
		contract: contract,
	}

	w.feedbackLabel = ui.NewLabelFromDriver(w.FindChild("mainLabel"))
	w.table = ui.NewTableWidgetFromDriver(w.FindChild("signersTable"))
	w.progressBar = ui.NewProgressBarFromDriver(w.FindChild("progressBar"))
	w.cancelButton = ui.NewPushButtonFromDriver(w.FindChild("cancelButton"))

	m, err := sign.NewSignatureManager(
		pwd,
		w.contract,
	)
	if err != nil {
		return nil
	}

	w.manager = m
	w.manager.OnSignerStatusUpdate = w.signerUpdated
	w.manager.OnProgressUpdate = func(current int, max int) {
		w.statusCurrent = int32(current)
		w.statusMax = int32(max)
	}

	w.cancelButton.OnClicked(func() {
		// Render an immediate feedback to user
		f := "Cancelling signature process..."
		w.feedback = f
		w.feedbackLabel.SetText(f)
		w.cancelButton.SetDisabled(true)
		w.statusMax = 1
		w.statusCurrent = 0
		// Ask for cancellation in a separate goroutine to avoid blocking Qt
		go func() { w.manager.Cancel <- true }()
	})

	w.initLines()
	w.signerUpdated(viper.GetString("email"), sign.StatusConnected, "It's you!")
	go func() {
		err = w.execute()
		if err != nil {
			w.feedback = err.Error()
		} else {
			w.feedback = "Contract signed successfully!"
		}
	}()

	return w
}

// execute() is called in a goroutine OUTSIDE of Qt loop.
// WE SHOULD NOT CALL ANY QT FUNCTION FROM IT.
func (w *Widget) execute() error {
	w.feedback = "Connecting to peers..."
	err := w.manager.ConnectToPeers()
	if err != nil {
		return err
	}

	w.feedback = "Waiting for peers..."
	_, err = w.manager.SendReadySign()
	if err != nil {
		return err
	}

	w.feedback = "Signature in progress..."
	w.running = true
	err = w.manager.Sign()
	if err != nil {
		return err
	}

	w.feedback = "Storing file..."
	return w.manager.PersistSignaturesToFile() // TODO choose destination
}

func (w *Widget) signerUpdated(mail string, status sign.SignerStatus, data string) {
	for i, s := range w.contract.Signers {
		if s.Email == mail {
			w.lines[i].status = status
			w.lines[i].info = data
			break
		}
	}
}

// Tick updates the whole screen, as we cannot directly update the screen on each callback.
// It is called by the screen coordinator via a timer.
func (w *Widget) Tick() {
	w.feedbackLabel.SetText(w.feedback)
	w.progressBar.SetMaximum(w.statusMax)
	w.progressBar.SetValue(w.statusCurrent)
	w.cancelButton.SetDisabled(w.running || w.manager.IsTerminated())
	for _, l := range w.lines {
		l.cellA.SetIcon(icons[l.status])
		l.cellA.SetText(icons_labels[l.status])
		l.cellB.SetText(l.info)
	}
}

func (w *Widget) initLines() {
	w.table.SetRowCount(int32(len(w.contract.Signers)))
	for i, s := range w.contract.Signers {
		status := ui.NewTableWidgetItemWithIconTextType(icons[sign.StatusWaiting], icons_labels[sign.StatusWaiting], int32(ui.QTableWidgetItem_UserType))
		w.table.SetItem(int32(i), 0, status)

		info := ui.NewTableWidgetItemWithTextType("", int32(ui.QTableWidgetItem_UserType))
		w.table.SetItem(int32(i), 1, info)

		mail := ui.NewTableWidgetItemWithTextType(s.Email, int32(ui.QTableWidgetItem_UserType))
		w.table.SetItem(int32(i), 2, mail)

		// We must store items, otherwise GC will remove it and cause the application to crash.
		// Trust me, hard to debug...
		w.lines = append(w.lines, line{sign.StatusWaiting, "", s.Email, status, info, mail})
	}
}

func (w *Widget) Q() *ui.QWidget {
	return w.QWidget
}
