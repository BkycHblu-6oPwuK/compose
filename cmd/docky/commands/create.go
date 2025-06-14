package commands

import (
	"fmt"
	"os"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/hoststools"

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
		globaltools.ValidateWorkDir()
		err := hoststools.CreateSite(createDomainNameFlag)
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
		globaltools.ValidateWorkDir()
		err := hoststools.CreateDomain(createDomainNameFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
		}
		fmt.Println("✅ домен создан!")
	},
}
var createDomainNameFlag string

func init() {
	createDomainModuleCmd.Flags().StringVar(&createDomainNameFlag, "name", "", "Задать имя домена")
	switch config.GetCurFramework() {
	case framework.Bitrix:
		createSiteModuleCmd.Flags().StringVar(&createDomainNameFlag, "name", "", "Задать имя домена")
		createCmd.AddCommand(createSiteModuleCmd)
	}
	createCmd.AddCommand(createDomainModuleCmd)
	rootCmd.AddCommand(createCmd)
}
