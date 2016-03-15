package main

import "github.com/visualfc/goqt/ui"

func main() {
	ui.Run(func() {
		widget := ui.NewWidget()

		label := ui.NewLabel()
		label.SetText("Welcome on dat fresh DFSS Client!")

		text := ui.NewTextEdit()
		text.SetText("Edit me!")

		button := ui.NewPushButton()
		button.SetText("Click me!")
		button.OnClicked(func() {
			button.SetText("Clicked!")
		})

		vbox := ui.NewVBoxLayout()
		vbox.SetAlignment(ui.Qt_AlignCenter)
		vbox.AddWidget(label)
		vbox.AddWidget(text)
		vbox.AddWidget(button)

		widget.SetWindowTitle("DFSS Client")
		widget.SetLayout(vbox)
		widget.SetMinimumWidth(400)
		widget.SetMinimumHeight(400)

		widget.Show()
	})
}
