package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Определяем команду `publish`
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Публикует файлы",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("✅ Файлы опубликованы!")
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
}