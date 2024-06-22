package screen

import (
	"log"

	"github.com/rivo/tview"
)

const (
	CONN_SCREEN     = "connection"
	CHAT_SCREEN     = "chat"
	CONN_WAIT_MODAL = "connecting"
	BAD_IP_MODAL    = "bad ip"
	BAD_PORT_MODAL  = "bad port"
	BAD_NAME_MODAL  = "bad name"
	CONN_ERR_MODAL  = "bad connection"
)

type ConnectionData struct {
	IP   string
	Port string
	Name string
}

type Connection struct {
	View           *tview.Pages
	NewConnections <-chan ConnectionData
	form           *tview.Form
	app            *tview.Application
}

type ConnectionResult int

const (
	OK ConnectionResult = iota
	ERR_BAD_IP
	ERR_BAD_PORT
	ERR_BAD_NAME
	ERR_CONN
)

type Validator func(string) bool

func (c *Connection) SendResult(result ConnectionResult) {
	log.Println("Got conn result")
	switch result {
	case OK:
		c.Clear()
		c.View.SwitchToPage(CONN_SCREEN)
	case ERR_BAD_IP:
		c.View.SwitchToPage(BAD_IP_MODAL)
	case ERR_BAD_PORT:
		c.View.SwitchToPage(BAD_PORT_MODAL)
	case ERR_BAD_NAME:
		c.View.SwitchToPage(BAD_NAME_MODAL)
	case ERR_CONN:
		c.View.SwitchToPage(CONN_ERR_MODAL)
	}
	_, p := c.View.GetFrontPage()
	c.app.SetFocus(p)
}

func (c *Connection) Clear() {
	for i := 0; i < c.form.GetFormItemCount(); i++ {
		input, ok := c.form.GetFormItem(i).(*tview.InputField)
		if ok {
			input.SetText("")
		}
	}
}

func NewConnection(app *tview.Application) *Connection {
	pages := tview.NewPages()
	form := tview.NewForm()
	newConnections := make(chan ConnectionData)
	var data ConnectionData

	pages.AddPage(CONN_SCREEN, form, true, true)
	pages.AddPage(CONN_WAIT_MODAL, tview.NewModal().SetText("Подключение..."), true, false)
	pages.AddPage(BAD_IP_MODAL, ErrorModal("Некорректный ip адрес", func() { pages.SwitchToPage(CONN_SCREEN) }), true, false)
	pages.AddPage(BAD_PORT_MODAL, ErrorModal("Некорректный порт", func() { pages.SwitchToPage(CONN_SCREEN) }), true, false)
	pages.AddPage(BAD_NAME_MODAL, ErrorModal("Некорректное имя", func() { pages.SwitchToPage(CONN_SCREEN) }), true, false)
	pages.AddPage(CONN_ERR_MODAL, ErrorModal("Не удалось подключиться", func() { pages.SwitchToPage(CONN_SCREEN) }), true, false)

	form.AddInputField("IP адрес", "", 20, nil,
		func(text string) {
			data.IP = text
		},
	)

	form.AddInputField("Порт", "", 20, nil,
		func(text string) {
			data.Port = text
		},
	)

	form.AddInputField("Имя", "", 20, nil,
		func(text string) {
			data.Name = text
		},
	)

	form.AddButton("Подключиться",
		func() {
			newConnections <- data
			pages.SwitchToPage(CONN_WAIT_MODAL)
			_, p := pages.GetFrontPage()
			app.SetFocus(p)
		},
	)

	form.SetBorder(true)

	return &Connection{
		View:           pages,
		NewConnections: newConnections,
		form:           form,
		app:            app,
	}
}
