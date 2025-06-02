package cmd

import (
	"docky/cmd/initialization"
	"docky/config"
	"docky/utils/globalHelper"
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Инициализация проекта",
	Run: func(cmd *cobra.Command, args []string) {
		err := initialization.InitDockerComposeFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
		err = globalHelper.InitSiteDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
		yamlConfig := config.GetYamlConfig()
		switch yamlConfig.FrameworkName {
		case config.Laravel:
			fmt.Println("Инициализация ларавел")
			initialization.InitLaravel()
		case config.Symfony:
			fmt.Println("Инициализация симфони")
			initialization.InitSymfony()
		}
		fmt.Println("✅ Инициализация проекта завершена!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}