package main

import (
	"regexp"

	"github.com/rivo/tview"
)

var ipreg = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
var portreg = regexp.MustCompile(`^[1-9][0-9]{1,3}$`)

func main() {
	pages := tview.NewPages()
	// pages.AddPage("connection", connectScreen(pages), true, true)
	pages.AddPage("connection", connectScreen(pages), true, false)
	pages.AddPage("bad ip", createErrorModal("Некорректный ip адрес", pages), true, false)
	pages.AddPage("bad port", createErrorModal("Некорректный порт", pages), true, false)
	pages.AddPage("chat", chatScreen(), true, true)

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

func chatScreen() *tview.Grid {
	newPrimitive := func(text string, align int) tview.Primitive {
		prim := tview.NewTextView().
			SetTextAlign(align).
			SetText(text)
		prim.SetBorderPadding(0, 0, 1, 0)
		return prim
	}
	users := newPrimitive("Users", tview.AlignCenter)
	chatLabel := newPrimitive("Chat", tview.AlignCenter)
	usersList := newPrimitive("Aboba\nSnake123\nSoska228", tview.AlignLeft)
	chat := newPrimitive("Aboba: Hello, everyone!\nSnake123: Hi!", tview.AlignLeft)
	messageLabel := newPrimitive("Your message", tview.AlignCenter)
	message := newPrimitive("Hi, guys!", tview.AlignLeft)

	grid := tview.NewGrid().
		SetRows(1, 0, 0, 1, 0).
		SetColumns(30, 0).
		SetBorders(true)

	grid.AddItem(users, 0, 0, 1, 1, 0, 0, false).
		AddItem(chatLabel, 0, 1, 1, 1, 0, 0, false).
		AddItem(usersList, 1, 0, 4, 1, 0, 0, false).
		AddItem(chat, 1, 1, 2, 1, 0, 0, false).
		AddItem(messageLabel, 3, 1, 1, 1, 0, 0, false).
		AddItem(message, 4, 1, 1, 1, 0, 0, false)

	return grid
}
