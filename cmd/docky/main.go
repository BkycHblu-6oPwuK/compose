package main

import (
	"os"

	"github.com/BkycHblu-6oPwuK/docky/cmd/docky/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
