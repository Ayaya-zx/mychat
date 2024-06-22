package screen

import (
	"github.com/rivo/tview"
)

type Main struct {
	View  *tview.Flex
	pages *tview.Pages
	list  *tview.List
}

func NewMain(conn *Connection) *Main {
	flex := tview.NewFlex()
	list := tview.NewList()
	pages := tview.NewPages()

	flex.AddItem(list, 0, 1, true)
	flex.AddItem(pages, 0, 4, false)

	pages.AddPage("conn", conn.View, true, true)

	list.SetBorder(true)
	list.SetBorderPadding(0, 0, 1, 0)
	list.AddItem("Новый чат", "", 0, func() { pages.SwitchToPage("conn") })

	return &Main{View: flex, pages: pages, list: list}
}

func (m *Main) AddChat(chat *Chat) {
	m.list.AddItem(chat.title, "", 0, func() { m.pages.SwitchToPage(chat.title) })
	m.list.SetCurrentItem(m.list.GetItemCount() - 1)
	m.pages.AddPage(chat.GetTitle(), chat.View, true, false)
	m.pages.SwitchToPage(chat.GetTitle())
}

func (m *Main) RemoveChat(chat *Chat) {
	m.pages.RemovePage(chat.title)
	m.list.RemoveItem(
		m.list.FindItems(chat.title, "", false, false)[0],
	)
	next, _ := m.list.GetItemText(0)
	m.pages.SwitchToPage(next)
}
