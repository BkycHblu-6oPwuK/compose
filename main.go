package main

import (
	"docky/cmd"
	"fmt"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("❌ Ошибка:", err)
		os.Exit(1)
	}
}
