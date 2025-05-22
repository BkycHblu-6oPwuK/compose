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
		text := "Файлы опубликованы!"
		if service != "" {
			err = publishService(service)
			text = "Сервис " + service + " опубликован!"
		} else {
			err = internal.PublishFiles()
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ " + text)
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
		return initEnvFile(yamlConfig)
	case yaml.Mysql:
		err := yaml.PublishMysqlService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		if yamlConfig.MysqlVersion == "" {
			yamlConfig.MysqlVersion = getOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, yaml.GetAvailableVersions(yaml.Mysql, yamlConfig))
		}
		return initEnvFile(yamlConfig)
	case yaml.Postgres:
		err := yaml.PublishPostgresService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		if yamlConfig.PostgresVersion == "" {
			yamlConfig.PostgresVersion = getOrChoose("Выберите версию postgres: ", yamlConfig.PostgresVersion, yaml.GetAvailableVersions(yaml.Postgres, yamlConfig))
		}
		return initEnvFile(yamlConfig)
	case yaml.Sphinx:
		return yaml.PublishSphinxService()
	case yaml.Redis:
		return yaml.PublishRedisService()
	case yaml.Memcached:
		return yaml.PublishMemcachedService()
	case yaml.Mailhog:
		return yaml.PublishMailhogService()
	case yaml.PhpMyAdmin:
		return yaml.PublishPhpMyAdminService()
	default:
		return fmt.Errorf("неизвестный сервис: %s", service)
	}
}
