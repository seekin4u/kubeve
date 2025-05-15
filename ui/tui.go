package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/a0xAi/kubeve/kube"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func StartUI(version string, overrideNamespace string) {
	var filterText string
	var allEvents []string
	var inputField *tview.InputField
	var namespaceList []string
	showTimestampColumn := true
	var recentNamespaces []string
	var header *Header

	autoScroll := true

	app := tview.NewApplication()
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})
	// load current context and namespace
	configLoadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	frame := tview.NewFrame(nil).
		SetBorders(1, 1, 1, 1, 1, 1)
	frame.SetBackgroundColor(0x000000)
	frame.SetPrimitive(flex)

	table := NewTable("[::b][green]Autoscroll ✓")

	namespace := overrideNamespace
	if namespace == "" {
		var err error
		namespace, _, err = clientConfig.Namespace()
		if err != nil {
			namespace = "default"
		}
	}
	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		// ignore RawConfig error
	}
	currentContext := rawConfig.CurrentContext

	// build Kubernetes REST config and clients
	restConfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}
	nsList, err := kubeClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		for _, ns := range nsList.Items {
			namespaceList = append(namespaceList, ns.Name)
		}
	}
	// server version
	versionInfo, _ := kubeClient.Discovery().ServerVersion()

	// derive cluster and user from kubeconfig contexts
	ctxConfig := rawConfig.Contexts[currentContext]
	clusterName := ctxConfig.Cluster
	userName := ctxConfig.AuthInfo

	currentNamespace := namespace
	showNamespaceColumn := false
	showStatusColumn := true
	showActionColumn := true
	showResourceColumn := true

	// header row

	header = NewHeader(
		currentContext,
		clusterName,
		namespace,
		userName,
		versionInfo.GitVersion,
		recentNamespaces,
	)

	var updateNamespace func(string)

	updateNamespace = func(newNS string) {
		if newNS == "" {
			namespace = metav1.NamespaceAll
		} else {
			namespace = newNS
		}
		// Update recent namespaces list (no duplicates, max 3)
		if newNS != "" {
			// remove if already present
			for i, ns := range recentNamespaces {
				if ns == newNS {
					recentNamespaces = append(recentNamespaces[:i], recentNamespaces[i+1:]...)
					break
				}
			}
			recentNamespaces = append([]string{newNS}, recentNamespaces...)
			if len(recentNamespaces) > 3 {
				recentNamespaces = recentNamespaces[:3]
			}
		}
		// Refresh RecentNSBox in header
		var recentLines []string
		recentLines = append(recentLines, "[blue]<0> [white]All Namespaces")
		for i, ns := range recentNamespaces {
			recentLines = append(recentLines, fmt.Sprintf("[blue]<%d> [white]%s", i+1, ns))
		}
		header.RecentNSBox.SetText(strings.Join(recentLines, "\n"))
		namespaceText := namespace
		if namespace == "" {
			namespaceText = "All namespaces"
		}
		header.InfoView.SetText(fmt.Sprintf(
			"[yellow]Cluster:[-] %s\n"+
				"[yellow]Namespace:[-] %s\n"+
				"[yellow]K8s Rev:[-] %s\n"+
				"[yellow]Kubeve Rev:[-] %s\n",
			clusterName, namespaceText, versionInfo.GitVersion, version,
		))
		allEvents = nil
		table.Clear()
		showNamespaceColumn = namespace == metav1.NamespaceAll
		renderTableHeader(table, showTimestampColumn, showNamespaceColumn, showStatusColumn, showActionColumn, showResourceColumn)

		// go kube.WatchEvents(namespace, false, func(event *corev1.Event) {
		go kube.WatchEvents(namespace, func(event *corev1.Event) {
			app.QueueUpdateDraw(func() {
				resource := fmt.Sprintf("%s/%s", event.InvolvedObject.Kind, event.InvolvedObject.Name)
				msg := fmt.Sprintf("%-25s │ %-60s │ %-10s │ %-20s │ %-10s │ %s\n",
					event.LastTimestamp.Time.Format(time.RFC3339),
					resource,
					event.Type,
					event.Reason,
					event.Namespace,
					event.Message,
				)
				if autoScroll {
					allEvents = append(allEvents, msg)
					if strings.Contains(msg, filterText) &&
						(namespace == metav1.NamespaceAll || event.Namespace == currentNamespace) {
						parts := strings.SplitN(msg, "│", 6)
						if len(parts) == 6 {
							row := table.GetRowCount()
							renderRow(table, row, parts, showTimestampColumn, showNamespaceColumn, showStatusColumn, showActionColumn, showResourceColumn)
							table.ScrollToEnd()
							table.Select(table.GetRowCount()-1, 0)
						}
					}
				}
			})
		})
	}

	inputField = tview.NewInputField()
	inputField.SetLabel("> ")
	inputField.SetBorder(true)
	inputField.SetBorderColor(0x00FF00)
	inputField.SetBackgroundColor(0x000000)
	inputField.SetFieldTextColor(0xFFFFFF)
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			filterText = inputField.GetText()
			// showNamespaceColumn := namespace == metav1.NamespaceAll
			table.Clear()
			// header row
			renderTableHeader(table, showTimestampColumn, showNamespaceColumn, showStatusColumn, showActionColumn, showResourceColumn)
			renderTableContent(table, allEvents, filterText, showTimestampColumn, showNamespaceColumn, showStatusColumn, showActionColumn, showResourceColumn)
			app.SetFocus(table)
		}
	})
	inputField.SetFieldBackgroundColor(0x000000)
	inputField.SetBackgroundColor(0x000000)
	inputField.SetBorder(false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyCtrlS:
			autoScroll = !autoScroll
			if autoScroll {
				table.SetTitle("[::b][green]Autoscroll ✓")
			} else {
				table.SetTitle("[::b][red]Autoscroll ✗")
			}
			return nil
		case event.Key() == tcell.KeyCtrlB:
			table.ScrollToEnd()
			table.Select(table.GetRowCount()-1, 0)
			return nil
		case event.Rune() == '/':
			inputField.SetText("")
			app.SetFocus(inputField)
			return nil
		case event.Key() == tcell.KeyCtrlN:
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
				AddItem(tview.NewBox(), 0, 1, false). // top spacer
				AddItem(
					tview.NewFlex().
						AddItem(tview.NewBox(), 0, 1, false). // left spacer
						AddItem(namespaceListView, 40, 0, true).
						AddItem(tview.NewBox(), 0, 1, false), // right spacer
									15, 0, true).
				AddItem(tview.NewBox(), 0, 1, false) // bottom spacer

			app.SetRoot(nsModal, true).SetFocus(namespaceListView)

			namespaceListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
					app.SetRoot(frame, true).SetFocus(table)
					return nil
				}
				return event
			})
			return nil
		case event.Rune() == 'T':
			showTimestampColumn = !showTimestampColumn
			table.Clear()
			renderTableHeader(table, showTimestampColumn, showNamespaceColumn, showStatusColumn, showActionColumn, showResourceColumn)
			renderTableContent(table, allEvents, filterText, showTimestampColumn, showNamespaceColumn, showStatusColumn, showActionColumn, showResourceColumn)
			return nil
		case event.Rune() == 'A':
			showActionColumn = !showActionColumn
			table.Clear()
			renderTableHeader(table, showTimestampColumn, showNamespaceColumn, showStatusColumn, showActionColumn, showResourceColumn)
			renderTableContent(table, allEvents, filterText, showTimestampColumn, showNamespaceColumn, showStatusColumn, showActionColumn, showResourceColumn)
			return nil
		case event.Rune() == 'R':
			showResourceColumn = !showResourceColumn
			table.Clear()
			renderTableHeader(table, showTimestampColumn, showNamespaceColumn, showStatusColumn, showActionColumn, showResourceColumn)
			renderTableContent(table, allEvents, filterText, showTimestampColumn, showNamespaceColumn, showStatusColumn, showActionColumn, showResourceColumn)
			return nil
		case event.Rune() == 'q', event.Key() == tcell.KeyCtrlC:
			app.Stop()
			return nil
		default:
			// Handle number key presses for namespace switching
			if event.Rune() >= '0' && event.Rune() <= '3' {
				switch event.Rune() {
				case '0':
					updateNamespace("") // all namespaces
				case '1', '2', '3':
					idx := int(event.Rune() - '1')
					if idx >= 0 && idx < len(namespaceList) {
						updateNamespace(namespaceList[idx])
					}
				}
				return nil
			}
			return event
		}
	})
	table.SetSelectedFunc(func(row int, column int) {
		if row > 0 && row-1 < len(allEvents) {
			parts := strings.SplitN(allEvents[row-1], "│", 5)
			if len(parts) == 5 {
				timeStr := strings.TrimSpace(parts[0])
				resource := strings.TrimSpace(parts[1])
				status := strings.TrimSpace(parts[2])
				action := strings.TrimSpace(parts[3])
				message := strings.TrimSpace(parts[4])

				detail := fmt.Sprintf(
					"[yellow]Event Detail[white]\n\n"+
						"[blue]Time:     [white]%s\n"+
						"[blue]Resource: [white]%s\n"+
						"[blue]Status:   [white]%s\n"+
						"[blue]Action:   [white]%s\n"+
						"[blue]Message:  [white]%s\n",
					timeStr, resource, status, action, message,
				)

				detailView := tview.NewTextView()
				detailView.SetDynamicColors(true)
				detailView.SetTextAlign(tview.AlignLeft)
				detailView.SetBorder(true)
				detailView.SetTitle(" Details ")
				detailView.SetBackgroundColor(0x000000)
				detailView.SetText(detail)

				modalFlex := tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(tview.NewBox(), 0, 1, false). // top spacer
					AddItem(
						tview.NewFlex().
							AddItem(tview.NewBox(), 0, 1, false). // left spacer
							AddItem(detailView, 80, 0, true).
							AddItem(tview.NewBox(), 0, 1, false), // right spacer
										15, 0, true).
					AddItem(tview.NewBox(), 0, 1, false) // bottom spacer

				app.SetRoot(modalFlex, true).SetFocus(detailView)

				detailView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
						app.SetRoot(frame, true).SetFocus(table)
						return nil
					}
					return event
				})
			}
		}
	})

	updateNamespace(namespace)

	flex.AddItem(header.Flex, 7, 0, false).
		AddItem(table, 0, 1, true).
		AddItem(inputField, 1, 0, false)
	if err := app.SetRoot(frame, true).Run(); err != nil {
		panic(err)
	}
}

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
		// ns := extractNamespace(event.Namespace)
		table.SetCell(row, col, tview.NewTableCell(strings.TrimSpace(parts[4])).SetExpansion(1))
		col++
	}
	table.SetCell(row, col, tview.NewTableCell(strings.TrimSpace(parts[2])).SetExpansion(1))
	col++
	if showActionColumn {
		actionText := strings.TrimSpace(parts[3])
		actionColor := "[white]"
		switch actionText {
		case "Created", "SuccessfulCreate":
			actionColor = "[green]"
		case "Started":
			actionColor = "[blue]"
		case "Pulled", "Pulling":
			actionColor = "[cyan]"
		case "Killing", "BackOff", "Unhealthy":
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

// extractNamespace extracts the namespace from a resource string of the form "namespace/resource" or "Kind/Name"
func extractNamespace(resource string) string {
	if parts := strings.Split(resource, "/"); len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// 507
