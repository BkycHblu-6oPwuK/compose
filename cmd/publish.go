package cmd

import (
	"docky/config"
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
	switch service {
	case yaml.Node:
		err := yaml.PublishNodeService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		err = initNode(yamlConfig)
		if err != nil {
			return err
		}
		return initEnvFile(yamlConfig, true)
	case yaml.Sphinx:
		return yaml.PublishSphinxService()
	case yaml.Redis:
		return yaml.PublisRedisService()
	case yaml.Memcached:
		return yaml.PublishMemcachedService()
	default:
		return fmt.Errorf("неизвестный сервис: %s", service)
	}
}
