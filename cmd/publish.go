package cmd

import (
	"docky/internal"
	"docky/yaml"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var service string

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Публикует файлы конфигурации",
	Run: func(cmd *cobra.Command, args []string) {
		validateWorkDir()
		var err error = nil
		if service != "" {
			err = publishService(service)
		} else {
			err = internal.PublishFiles()
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
		}
		fmt.Println("✅ Файлы опубликованы!")
	},
}

func init() {
	publishCmd.Flags().StringVar(&service, "service", "", "Опубликовать сервис в docker-compose")
	rootCmd.AddCommand(publishCmd)
}

func publishService(service string) error {
	if service == "node" {
		err := yaml.PublishNodeService()
		if err != nil {
			return err
		}
		err = initNode()
		if err != nil {
			return err
		}
		return initEnvFile(true)
	} else if service == "sphinx" {
		return yaml.PublishSphinxService()
	} else {
		return fmt.Errorf("неизвестный сервис: %s", service)
	}
}