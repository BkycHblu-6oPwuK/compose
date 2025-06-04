package main

import (
	"docky/cmd/docky/commands"
	"os"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
