package main

import (
	// "docky/config"
	// "docky/models"
	// "fmt"
	"docky/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	// fmt.Println(config.GetCurrentDockerFileDirPath())
	// models.CreateYmlFile()
}
