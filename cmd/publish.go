package cmd

import (
	"docky/internal"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Определяем команду `publish`
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Публикует файлы",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.PublishFiles()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
		}
		fmt.Println("✅ Файлы опубликованы!")
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
}
