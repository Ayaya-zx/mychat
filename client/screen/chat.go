package screen

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Chat представляет собой окно чата.
type Chat struct {
	// View - это предстваление окна чата.
	View        *tview.Grid
	users       *tview.TextView
	messages    *tview.TextView
	newMessages chan string
	Title       string
	Key         string
}

// UpdateUsers обновляет список пользователей.
func (c *Chat) UpdateUsers(users []string) {
	c.users.Clear()
	writer := c.users.BatchWriter()
	defer writer.Close()
	for _, user := range users {
		writer.Write([]byte(user + "\n"))
	}
}

// AddMessage добавляет новое сообщение к
// списку выводимых сообщений.
func (c *Chat) AddMessage(msg string) {
	writer := c.messages.BatchWriter()
	defer writer.Close()
	writer.Write([]byte(msg + "\n"))
	c.messages.ScrollToEnd()
}

func (c *Chat) SetMessageFieldFocus(app *tview.Application) {
	app.SetFocus(c.messages)
}

// NewMessages возвращает канал, по которому
// передаются сообщения, отправленный пользователем.
func (c *Chat) NewMessages() <-chan string {
	return c.newMessages
}

// Dispose очищает ресурсы чата.
func (c *Chat) Dispose() {
	close(c.newMessages)
}

// NewChat возвращает новый экземпляр чата.
func NewChat(name, title string) *Chat {
	newPrimitive := func(text string, align int) *tview.TextView {
		prim := tview.NewTextView().
			SetTextAlign(align).
			SetText(text)
		prim.SetBorderPadding(0, 0, 1, 0)
		return prim
	}

	newMessages := make(chan string)

	userLabel := newPrimitive("Users", tview.AlignCenter)
	chatLabel := newPrimitive("Chat", tview.AlignCenter)
	userList := newPrimitive("", tview.AlignLeft)
	chatField := newPrimitive("", tview.AlignLeft)
	messageLabel := newPrimitive("Your message", tview.AlignCenter)

	message := tview.NewTextArea()
	// message := tview.NewInputField()
	message.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			newMessages <- message.GetText()
			message.SetText("", true)
			// message.SetText("")
			return nil
		}
		return event
	})
	message.SetLabel(name + ": ")

	grid := tview.NewGrid().
		SetRows(1, 0, 0, 0, 1, 0).
		SetColumns(0, 25).
		SetBorders(true)

	grid.AddItem(userLabel, 0, 1, 1, 1, 0, 0, false).
		AddItem(userList, 1, 1, 5, 1, 0, 0, false).
		AddItem(chatLabel, 0, 0, 1, 1, 0, 0, false).
		AddItem(chatField, 1, 0, 3, 1, 0, 0, true).
		AddItem(messageLabel, 4, 0, 1, 1, 0, 0, false).
		AddItem(message, 5, 0, 1, 1, 0, 0, false)

	chat := &Chat{
		View:        grid,
		newMessages: newMessages,
		users:       userList,
		messages:    chatField,
		Title:       title,
		Key:         generateKey(),
	}

	return chat
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateKey() string {
	buff := make([]byte, 20)
	for i := 0; i < 20; i++ {
		j := rand.Intn(52)
		buff[i] = letters[j]
	}
	return string(buff)
}
