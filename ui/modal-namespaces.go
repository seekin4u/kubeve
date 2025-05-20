package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NamespacesModal(app *tview.Application, frame *tview.Frame, table *tview.Table, namespaceList []string, updateNamespace func(name string)) {
	namespaceListView := tview.NewList()
	for _, ns := range namespaceList {
		namespaceListView.AddItem(ns, "", 0, nil)
	}
	namespaceListView.SetSelectedFunc(func(index int, name string, secondary string, shortcut rune) {
		updateNamespace(name)
		app.SetRoot(frame, true).SetFocus(table)
	})
	namespaceListView.SetBorder(true).SetTitle(" Select Namespace ")
	namespaceListView.SetBackgroundColor(0x000000)

	nsModal := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(
			tview.NewFlex().
				AddItem(tview.NewBox(), 0, 1, false).
				AddItem(namespaceListView, 40, 0, true).
				AddItem(tview.NewBox(), 0, 1, false),
			15, 0, true).
		AddItem(tview.NewBox(), 0, 1, false)

	app.SetRoot(nsModal, true).SetFocus(namespaceListView)

	namespaceListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
			app.SetRoot(frame, true).SetFocus(table)
			return nil
		}
		return event
	})
}
