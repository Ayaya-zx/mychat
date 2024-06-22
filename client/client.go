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
	ipreg      = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
	portreg    = regexp.MustCompile(`^[1-9][0-9]{1,3}$`)
	app        *tview.Application
	mainScreen *screen.Main
	conn       *screen.Connection
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

func init() {
	app = tview.NewApplication()
	conn = screen.NewConnection(app)
	mainScreen = screen.NewMain(conn)
}

func main() {
	file, err := os.Create("/tmp/mychat.log")
	if err == nil {
		log.SetOutput(file)
	} else {
		log.SetOutput(io.Discard)
	}
	log.Println("App started")

	// app.SetRoot(pages, true)
	app.SetRoot(mainScreen.View, true)
	app.EnableMouse(true)

	go handleConnections(conn.NewConnections)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func handleConnections(newConns <-chan screen.ConnectionData) {
	for data := range newConns {
		log.Println("Got new connection data!")
		if !ipreg.Match([]byte(data.IP)) {
			log.Println("Sending ERR_BAD_IP")
			app.QueueUpdateDraw(func() {
				conn.SendResult(screen.ERR_BAD_IP)
			})
			continue
		}
		if !portreg.Match([]byte(data.Port)) {
			log.Println("Sending ERR_BAD_PORT")
			app.QueueUpdateDraw(func() {
				conn.SendResult(screen.ERR_BAD_PORT)
			})
			continue
		}
		if strings.TrimSpace(data.Name) == "" {
			log.Println("Sending ERR_BAD_NAME")
			app.QueueUpdateDraw(func() {
				conn.SendResult(screen.ERR_BAD_NAME)
			})
			continue
		}

		c, err := net.Dial("tcp", fmt.Sprintf("%s:%s",
			data.IP, data.Port))
		if err != nil {
			log.Println(err)
			log.Println("Sending ERR_CONN")
			app.QueueUpdateDraw(func() {
				conn.SendResult(screen.ERR_CONN)
			})
			continue
		}

		_, err = c.Write([]byte("NAME" + data.Name))
		if err != nil {
			log.Println(err)
			log.Println("Sending ERR_CONN")
			app.QueueUpdateDraw(func() {
				conn.SendResult(screen.ERR_CONN)
			})
			continue
		}

		chat := screen.NewChat(data.Name, c.RemoteAddr().String())
		go handleChat(c, chat)
		app.QueueUpdateDraw(
			func() {
				log.Println("Opening new chat")
				mainScreen.AddChat(chat)
				chat.SetMessageFieldFocus(app)
				conn.SendResult(screen.OK)
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
		mainScreen.RemoveChat(chat)
	})
}
