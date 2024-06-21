package screen

import (
	"log"

	"github.com/rivo/tview"
)

type ConnectionData struct {
	IP   string
	Port string
	Name string
}

type Connection struct {
	View           *tview.Form
	NewConnections <-chan ConnectionData
}

func NewConnection() *Connection {
	form := tview.NewForm()
	newConnections := make(chan ConnectionData)
	var data ConnectionData

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
			log.Println("Sending new connection data")
			newConnections <- data
		},
	)

	return &Connection{
		View:           form,
		NewConnections: newConnections,
	}
}

// func Connection(confirm func(ip, port, nickname string)) *tview.Form {
// 	form := tview.NewForm()
// 	var ip, port, nickname string

// 	form.AddInputField("IP адрес", "", 20, nil,
// 		func(text string) {
// 			ip = text
// 		},
// 	)

// 	form.AddInputField("Порт", "", 20, nil,
// 		func(text string) {
// 			port = text
// 		},
// 	)

// 	form.AddInputField("Имя", "", 20, nil,
// 		func(text string) {
// 			nickname = text
// 		},
// 	)

// 	form.AddButton("Подключиться",
// 		func() {
// 			confirm(ip, port, nickname)
// 		},
// 	)

// 	_ = nickname
// 	return form
// }
