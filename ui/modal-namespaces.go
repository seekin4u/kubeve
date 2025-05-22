package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NamespacesModal(app *tview.Application, frame tview.Primitive, table *tview.Table, namespaceList []string, updateNamespace func(string)) {
	names := append([]string{}, namespaceList...)
	filtered := append([]string{}, names...)
	selection := 0
	filterText := ""

	input := tview.NewInputField().
		SetLabelStyle(tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.Color16)).
		SetFieldStyle(tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.Color16)).
		SetLabel("> ").
		SetFieldWidth(0)
	input.SetBorder(false)
	input.SetChangedFunc(func(text string) {
		filterText = text
		filtered = filtered[:0]
		for _, ns := range names {
			if strings.Contains(strings.ToLower(ns), strings.ToLower(filterText)) {
				filtered = append(filtered, ns)
			}
		}
		selection = 0
	})

	overlay := tview.NewBox().SetBackgroundColor(tcell.Color16).SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		listH := height - 1
		if selection < 0 {
			selection = 0
		} else if selection >= len(filtered) {
			selection = len(filtered) - 1
		}
		start := 0
		if len(filtered) > listH {
			if selection < listH {
				start = 0
			} else {
				start = selection - (listH - 1)
			}
		}
		visibleCount := len(filtered) - start
		if visibleCount > listH {
			visibleCount = listH
		}
		ofs := 0
		if visibleCount < listH {
			ofs = listH - visibleCount
		}

		borderColor := tcell.ColorRed
		inactiveBorder := tcell.ColorBlack

		for i := 0; i < visibleCount; i++ {
			row := start + i
			var borderBg tcell.Color
			if row == selection {
				borderBg = borderColor
			} else {
				borderBg = inactiveBorder
			}
			screen.SetContent(x, y+ofs+i, ' ', nil, tcell.StyleDefault.Background(borderBg))

			fg := tcell.ColorWhite
			if row == selection {
				fg = tcell.ColorYellow
			}
			tview.Print(screen, filtered[row], x+1, y+ofs+i, width-1, tview.AlignLeft, fg)
		}
		// draw filter input at bottom
		input.SetRect(x, y+listH, width, 1)
		input.Draw(screen)
		return x, y, width, height
	})

	prev := app.GetInputCapture()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			selection--
		case tcell.KeyDown:
			selection++
		case tcell.KeyEnter:
			if len(filtered) > 0 {
				updateNamespace(filtered[selection])
			}
			app.SetInputCapture(prev)
			app.SetRoot(frame, true).SetFocus(table)
		case tcell.KeyEsc:
			app.SetInputCapture(prev)
			app.SetRoot(frame, true).SetFocus(table)
		default:
			handler := input.InputHandler()
			if handler != nil {
				handler(event, nil)
			}
		}
		return nil
	})

	app.SetRoot(overlay, true).SetFocus(input)
}
