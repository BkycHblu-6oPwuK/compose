package cmd

import (
	"docky/config"
	"docky/internal"
	"fmt"

	"github.com/spf13/cobra"
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