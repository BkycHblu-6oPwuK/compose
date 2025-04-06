package main

import (
	"docky/cmd"
	//"fmt"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
