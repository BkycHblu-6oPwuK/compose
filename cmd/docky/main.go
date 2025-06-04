package main

import (
	"os"

	"github.com/BkycHblu-6oPwuK/docky/v2/cmd/docky/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
