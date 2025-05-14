package ui

import (
	"github.com/rivo/tview"
	"github.com/a0xAi/kubeve/kube"
)

func StartUI() {
	app := tview.NewApplication()
	textView := tview.NewTextView().SetDynamicColors(true).SetScrollable(true)
	textView.SetBorder(true).SetTitle("Kube Events")

	go kube.WatchEvents("default", func(event *v1.Event) {
		app.QueueUpdateDraw(func() {
			msg := "[" + event.Type + "] " + event.Message + "\n"
			textView.Write([]byte(msg))
		})
	})

	if err := app.SetRoot(textView, true).Run(); err != nil {
		panic(err)
	}
}