package main

import (
	"docky/cmd/dockyphp/commands"
	"os"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
