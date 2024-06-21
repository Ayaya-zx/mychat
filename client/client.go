package main

import (
	"fmt"
	"io"
	"log"
	"mychat/client/screen"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/rivo/tview"
)

var (
	ipreg   = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
	portreg = regexp.MustCompile(`^[1-9][0-9]{1,3}$`)
	app     *tview.Application
	pages   *tview.Pages
)

func main() {
	file, err := os.Create("/tmp/mychat.log")
	if err == nil {
		log.SetOutput(file)
	} else {
		fmt.Println("Can't set up log")
		log.SetOutput(io.Discard)
	}
	log.Println("App started")

	app = tview.NewApplication()

	connection := screen.NewConnection()
	go handleConnections(connection.NewConnections)

	pages = tview.NewPages()
	pages.AddPage("connection", connection.View, true, true)
	pages.AddPage("bad ip", screen.ErrorModal("Некорректный ip адрес", func() { pages.SwitchToPage("connection") }), true, false)
	pages.AddPage("bad port", screen.ErrorModal("Некорректный порт", func() { pages.SwitchToPage("connection") }), true, false)
	pages.AddPage("bad name", screen.ErrorModal("Некорректное имя", func() { pages.SwitchToPage("connection") }), true, false)
	pages.AddPage("bad connection", screen.ErrorModal("Не удалось подключиться", func() { pages.SwitchToPage("connection") }), true, false)

	connecting := tview.NewModal()
	connecting.SetText("Подключение...")
	pages.AddPage("connecting", connecting, true, false)

	app.SetRoot(pages, true)
	app.EnableMouse(true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func handleConnections(newConns <-chan screen.ConnectionData) {
	for data := range newConns {
		log.Println("Got new connection data!")
		if !ipreg.Match([]byte(data.IP)) {
			app.QueueUpdateDraw(func() {
				pages.SwitchToPage("bad ip")
			})
			continue
		}
		if !portreg.Match([]byte(data.Port)) {
			app.QueueUpdateDraw(func() {
				pages.SwitchToPage("bad port")
			})
			continue
		}
		if strings.TrimSpace(data.Name) == "" {
			app.QueueUpdateDraw(func() {
				pages.SwitchToPage("bad name")
			})
			continue
		}

		app.QueueUpdateDraw(func() {
			pages.SwitchToPage("connecting")
		})

		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s",
			data.IP, data.Port))
		if err != nil {
			log.Println(err)
			app.QueueUpdateDraw(func() {
				pages.SwitchToPage("bad connection")
			})
			continue
		}

		_, err = conn.Write([]byte("NAME" + data.Name))
		if err != nil {
			log.Println(err)
			app.QueueUpdateDraw(func() {
				pages.SwitchToPage("bad connection")
			})
			continue
		}

		chat := screen.NewChat()
		go handleChat(conn, chat)
		fmt.Println("Gonna open new chat")
		app.QueueUpdateDraw(
			func() {
				log.Println("Opening new chat")
				pages.AddPage("chat", chat.View, true, true)
				pages.SwitchToPage("chat")
			},
		)
	}
}

func handleChat(conn net.Conn, chat *screen.Chat) {
	msgEnd := "\xe2\x90\x9c"

	go func() {
		for msg := range chat.NewMessages() {
			log.Println("New message:", msg)
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Println(err)
			}
		}
	}()
	defer conn.Close()
	buff := make([]byte, 1024)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			break
		}
		input := string(buff[:n])
		msgs := strings.Split(input, msgEnd)
		fmt.Println("Got", len(msgs), "messages")
		for _, msg := range msgs {
			log.Println("Got message", msg)
			if len(msg) > 5 && msg[:5] == "USERS" {
				app.QueueUpdateDraw(func() {
					chat.UpdateUsers(strings.Split(msg[5:], ","))
				})
			} else if len(msg) > 7 && msg[:7] == "MESSAGE" {
				app.QueueUpdateDraw(func() {
					chat.AddMessage(msg[7:])
				})
			} else {
				log.Println("Strange message from server:" +
					msg)
			}
		}
	}
	chat.Dispose()
	app.QueueUpdateDraw(func() {
		pages.RemovePage("chat")
		pages.SwitchToPage("connection")
	})
}
