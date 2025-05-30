package cmd

import (
	"docky/internal"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cleanCacheCmd = &cobra.Command{
	Use:   "clean-cache",
	Short: "Очищает кэш (директория _docker)",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.CleanCacheDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
		fmt.Println("✅ Кэш очищен!")
	},
}

func init() {
	rootCmd.AddCommand(cleanCacheCmd)
}
