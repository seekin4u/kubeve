package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewTable(status string) *tview.Table {
	table := tview.NewTable().SetBorders(false).SetFixed(1, 0)
	table.SetSelectable(true, false)
	table.SetBorder(true).SetTitle(status)
	table.SetBackgroundColor(0x000000)
	return table
}

// renderTableHeader sets up the table header cells based on visible columns.
func renderTableHeader(table *tview.Table, showTimestampColumn bool, showNamespaceColumn bool, showStatusColumn bool, showActionColumn bool, showResourceColumn bool) {
	col := 0
	if showTimestampColumn {
		table.SetCell(0, col, tview.NewTableCell("TIME").
			SetSelectable(false).SetAttributes(tcell.AttrBold).SetExpansion(1))
		col++
	}
	if showNamespaceColumn {
		table.SetCell(0, col, tview.NewTableCell("NAMESPACE").
			SetSelectable(false).SetAttributes(tcell.AttrBold).SetExpansion(1))
		col++
	}
	if showStatusColumn {
		table.SetCell(0, col, tview.NewTableCell("STATUS").
			SetSelectable(false).SetAttributes(tcell.AttrBold).SetExpansion(1))
		col++
	}
	if showActionColumn {
		table.SetCell(0, col, tview.NewTableCell("ACTION").
			SetSelectable(false).SetAttributes(tcell.AttrBold).SetExpansion(1))
		col++
	}
	if showResourceColumn {
		table.SetCell(0, col, tview.NewTableCell("RESOURCE").
			SetSelectable(false).SetAttributes(tcell.AttrBold).SetExpansion(2))
		col++
	}
	table.SetCell(0, col, tview.NewTableCell("MESSAGE").
		SetSelectable(false).SetAttributes(tcell.AttrBold).SetExpansion(5))
}

func renderRow(table *tview.Table, row int, parts []string, showTimestampColumn bool, showNamespaceColumn bool, showStatusColumn bool, showActionColumn bool, showResourceColumn bool) {
	col := 0
	if showTimestampColumn {
		table.SetCell(row, col, tview.NewTableCell(strings.TrimSpace(parts[0])).SetExpansion(1))
		col++
	}
	if showNamespaceColumn {
		table.SetCell(row, col, tview.NewTableCell(strings.TrimSpace(parts[4])).SetExpansion(1))
		col++
	}
	statusText := strings.TrimSpace(parts[2])
	statusColor := "[white]"
	switch statusText {
	case "Warning":
		statusColor = "[yellow]"
	}
	table.SetCell(row, col, tview.NewTableCell(fmt.Sprintf("%s%s", statusColor, statusText)).SetExpansion(1))
	col++
	if showActionColumn {
		actionText := strings.TrimSpace(parts[3])
		actionColor := "[white]"
		switch actionText {
		case "Created", "SuccessfulCreate", "Completed":
			actionColor = "[green]"
		case "Started", "Pulled", "Pulling":
			actionColor = "[blue]"
		case "Killing", "BackOff", "Unhealthy", "FailedToRetrieveImagePullSecret":
			actionColor = "[red]"
		}
		table.SetCell(row, col, tview.NewTableCell(fmt.Sprintf("%s%s", actionColor, actionText)).
			SetExpansion(1).SetTextColor(tcell.ColorWhite))
		col++
	}
	if showResourceColumn {
		table.SetCell(row, col, tview.NewTableCell(strings.TrimSpace(parts[1])).SetExpansion(2))
		col++
	}
	table.SetCell(row, col, tview.NewTableCell(strings.TrimSpace(parts[5])).SetExpansion(5))
}

func renderTableContent(table *tview.Table, events []string, filterText string, showTimestampColumn bool, showNamespaceColumn bool, showStatusColumn bool, showActionColumn bool, showResourceColumn bool) {
	row := 1
	for _, line := range events {
		if strings.Contains(line, filterText) {
			parts := strings.SplitN(line, "│", 6)
			if len(parts) == 6 {
				renderRow(table, row, parts, showTimestampColumn, showNamespaceColumn, showStatusColumn, showActionColumn, showResourceColumn)
				row++
			}
		}
	}
}
