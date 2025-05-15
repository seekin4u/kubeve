package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

// Header exposes the header flex and infoView for live updates.
type Header struct {
	Flex        *tview.Flex
	InfoView    *tview.TextView
	RecentNSBox *tview.TextView
}

// NewHeader builds the top-bar with context info, shortcuts and ASCII logo.
func NewHeader(
	currentContext, clusterName, namespace, userName, kubeRev string,
	recentNamespaces []string,
) *Header {
	// Context/info pane
	infoView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	infoView.SetBackgroundColor(0x000000)
	namespaceText := namespace
	if namespace == "" {
		namespaceText = "All namespaces"
	}
	infoView.SetText(fmt.Sprintf(
		"[yellow]Cluster:[-] %s\n"+
			"[yellow]Namespace:[-] %s\n"+
			"[yellow]K8s Rev:[-] %s\n"+
			"[yellow]Kubeve Rev:[-] %s\n",
		clusterName, namespaceText, userName, kubeRev, "0.3.0",
	))

	// Recent namespace shortcuts pane
	recentNs := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	recentNs.SetBackgroundColor(0x000000)
	var recentLines []string
	recentLines = append(recentLines, "[blue]<0> [white]All Namespaces")
	for i, ns := range recentNamespaces {
		if i >= 3 {
			break
		}
		recentLines = append(recentLines, fmt.Sprintf("[blue]<%d> [white]%s", i+1, ns))
	}
	recentNs.SetText(strings.Join(recentLines, "\n"))

	// Shortcut keys pane
	shortcuts := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	shortcuts.SetBackgroundColor(0x000000)
	shortcuts.SetText(`[blue]</>       [white]Toggle filter
[blue]<ctrl+s>  [white]Toggle autoscroll
[blue]<ctrl+b>  [white]Go to last event
[blue]<ctrl+n>  [white]Change namespace
[blue]<↑↓>      [white]Scroll`)

	shortcuts2 := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	shortcuts2.SetBackgroundColor(0x000000)
	shortcuts2.SetText(`[blue]<shift+t>  [white]Toggle timestamp
[blue]<shift+r>  [white]Toggle resource
[blue]<shift+a>  [white]Toggle action`)
	// ASCII logo pane
	logo := `__        ___.                      
|  | ____ _\_ |__   [red]_______  __ ____ 
[white]|  |/ /  |  \ __ \[red]_/ __ \  \/ // __ \
[white]|    <|  |  / \_\ \  [red]___/\   /\  ___/
[white]|__|_ \____/|___  /[red]\___  >\_/  \___ >
     [white]\/         \/     [red]\/          \/ `
	logoView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)
	logoView.SetBackgroundColor(0x000000)
	logoView.SetText(logo)

	// Put them side by side
	headerFlex := tview.NewFlex().
		AddItem(infoView, 0, 2, false).
		AddItem(recentNs, 0, 1, false).
		AddItem(shortcuts, 0, 2, false).
		AddItem(shortcuts2, 0, 2, false).
		AddItem(logoView, 0, 2, false)

	return &Header{
		Flex:        headerFlex,
		InfoView:    infoView,
		RecentNSBox: recentNs,
	}
}
