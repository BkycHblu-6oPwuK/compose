package commands

import (
	"fmt"
	"os"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/files"

	"github.com/spf13/cobra"
)

var cleanCacheCmd = &cobra.Command{
	Use:   "clean-cache",
	Short: "Очищает кэш (директория _docker)",
	Run: func(cmd *cobra.Command, args []string) {
		err := files.CleanCacheDir()
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
