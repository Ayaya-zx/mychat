package main

import (
	"regexp"

	"github.com/rivo/tview"
)

var ipreg = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
var portreg = regexp.MustCompile(`^[1-9][0-9]{1,3}$`)

func main() {
	pages := tview.NewPages()
	pages.AddPage("connection", connectScreen(pages), true, true)
	pages.AddPage("bad ip", createErrorModal("Некорректный ip адрес", pages), true, false)
	pages.AddPage("bad port", createErrorModal("Некорректный порт", pages), true, false)

	app := tview.NewApplication()
	app.SetRoot(pages, true)
	app.EnableMouse(true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func createErrorModal(msg string, pages *tview.Pages) *tview.Modal {
	modal := tview.NewModal()
	modal.SetText(msg)
	modal.AddButtons([]string{"OK"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		pages.SwitchToPage("connection")
	})
	return modal
}

func connectScreen(pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	var ip, port, nickname string

	form.AddInputField("IP адрес", "", 20, nil,
		func(text string) {
			ip = text
		},
	)

	form.AddInputField("Порт", "", 20, nil,
		func(text string) {
			port = text
		},
	)

	form.AddInputField("Имя", "", 20, nil,
		func(text string) {
			nickname = text
		},
	)

	form.AddButton("Подключиться",
		func() {
			if !ipreg.Match([]byte(ip)) {
				pages.SwitchToPage("bad ip")
				return
			}
			if !portreg.Match([]byte(port)) {
				pages.SwitchToPage("bad port")
				return
			}
		},
	)

	_ = nickname
	return form
}
