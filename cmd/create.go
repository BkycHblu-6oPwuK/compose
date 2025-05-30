package cmd

import (
	"docky/cmd/create"
	"docky/config"
	"docky/utils/globalHelper"
	"fmt"
	"os"

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
		globalHelper.ValidateWorkDir()
		err := create.CreateSite()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
		fmt.Println("✅ сайт создан!")
	},
}

var createDomainModuleCmd = &cobra.Command{
	Use:   "domain",
	Short: "Создает новый домен (сертификаты, запись в hosts)",
	Run: func(cmd *cobra.Command, args []string) {
		err := create.CreateDomain()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
		}
		fmt.Println("✅ домен создан!")
	},
}

func init() {
	switch config.GetCurFramework() {
	case config.Bitrix:
		createCmd.AddCommand(createSiteModuleCmd)
	}
	createCmd.AddCommand(createDomainModuleCmd)
	rootCmd.AddCommand(createCmd)
}
