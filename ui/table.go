package ui

import (
	"github.com/rivo/tview"
)

func NewTable(status string) *tview.Table {
	table := tview.NewTable().SetBorders(false).SetFixed(1, 0)
	table.SetSelectable(true, false)
	table.SetBorder(true).SetTitle(status)
	table.SetBackgroundColor(0x000000)
	return table
}
