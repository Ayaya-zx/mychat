package screen

import "github.com/rivo/tview"

func ErrorModal(msg string, done func()) *tview.Modal {
	modal := tview.NewModal()
	modal.SetText(msg)
	modal.AddButtons([]string{"OK"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		done()
	})
	return modal
}
