package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Команда для создания (cайт, домен)",
}

var createSiteModuleCmd = &cobra.Command{
	Use:   "site",
	Short: "Создает новый сайт (директория, сертификаты, запись в hosts)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("✅ create site")
	},
}

func init() {
	createCmd.AddCommand(createSiteModuleCmd)
	rootCmd.AddCommand(createCmd)
}
