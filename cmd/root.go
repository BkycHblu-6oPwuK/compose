package cmd

import (
	"docky/config"
	"docky/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	curDirPath string // директория из которой запускается команда
	workDirPath string // директория с docker-compose.yml
)

var rootCmd = &cobra.Command{
	Use:   config.ScriptName,
	Short: "Программа для работы с docker-compose для битрикс проектов",
}

func init() {
	internal.ExtractFilesInCache()
	fmt.Println("🚀 Запуск docky...")
}

func Execute() error {
	return rootCmd.Execute()
}