package main

import (
	"flag"

	"github.com/a0xAi/kubeve/ui"
)

func main() {
	version := "0.3.10"
	namespace := flag.String("n", "", "Kubernetes namespace to use")
	flag.Parse()
	ui.StartUI(version, *namespace)
}
