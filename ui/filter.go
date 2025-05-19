package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewFilter() *tview.InputField {
	filter := tview.NewInputField()
	filter.SetLabelStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.Color16))
	filter.SetFieldStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.Color16))
	filter.SetLabel("> ")
	return filter
}
