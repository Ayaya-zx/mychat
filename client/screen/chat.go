package screen

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Chat struct {
	View        *tview.Grid
	users       *tview.TextView
	messages    *tview.TextView
	newMessages chan string
}

func (c *Chat) UpdateUsers(users []string) {
	c.users.Clear()
	writer := c.users.BatchWriter()
	defer writer.Close()
	for _, user := range users {
		writer.Write([]byte(user + "\n"))
	}
}

func (c *Chat) AddMessage(msg string) {
	writer := c.messages.BatchWriter()
	defer writer.Close()
	writer.Write([]byte(msg + "\n"))
}

func (c *Chat) NewMessages() <-chan string {
	return c.newMessages
}

func (c *Chat) Dispose() {
	close(c.newMessages)
}

func NewChat() *Chat {
	newPrimitive := func(text string, align int) *tview.TextView {
		prim := tview.NewTextView().
			SetTextAlign(align).
			SetText(text)
		prim.SetBorderPadding(0, 0, 1, 0)
		return prim
	}

	newMessages := make(chan string)

	users := newPrimitive("Users", tview.AlignCenter)
	chatLabel := newPrimitive("Chat", tview.AlignCenter)
	userList := newPrimitive("", tview.AlignLeft)
	chatField := newPrimitive("", tview.AlignLeft)
	messageLabel := newPrimitive("Your message", tview.AlignCenter)

	message := tview.NewTextArea()
	message.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			newMessages <- message.GetText()
			message.SetText("", true)
			return nil
		}
		return event
	})

	grid := tview.NewGrid().
		SetRows(1, 0, 0, 1, 0).
		SetColumns(30, 0).
		SetBorders(true)

	grid.AddItem(users, 0, 0, 1, 1, 0, 0, false).
		AddItem(userList, 1, 0, 4, 1, 0, 0, false).
		AddItem(chatLabel, 0, 1, 1, 1, 0, 0, false).
		AddItem(chatField, 1, 1, 2, 1, 0, 0, false).
		AddItem(messageLabel, 3, 1, 1, 1, 0, 0, false).
		AddItem(message, 4, 1, 1, 1, 0, 0, false)

	chat := &Chat{
		View:        grid,
		newMessages: newMessages,
		users:       userList,
		messages:    chatField,
	}

	return chat
}
