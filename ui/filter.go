package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewInputField returns a configured input field.
// onDone is a callback when Enter or Esc is pressed.
func NewFilter(onDone func(key tcell.Key)) *tview.InputField {
	filter := tview.NewInputField()
	filter.SetLabel("> ")
	filter.SetBorder(true)
	filter.SetBorderColor(0x00FF00)
	filter.SetBackgroundColor(0x000000)
	filter.SetBorder(false) // hidden by default
	return filter
}
