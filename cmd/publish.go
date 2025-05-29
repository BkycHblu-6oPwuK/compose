package cmd

import (
	"docky/config"
	"docky/internal"
	"docky/yaml/helper"
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
	case helper.Node:
		err := helper.PublishNodeService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		err = initNode(yamlConfig)
		if err != nil {
			return err
		}
		return initEnvFile(yamlConfig)
	case helper.Mysql:
		err := helper.PublishMysqlService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		if yamlConfig.MysqlVersion == "" {
			yamlConfig.MysqlVersion = getOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, helper.GetAvailableVersions(helper.Mysql, yamlConfig))
		}
		return initEnvFile(yamlConfig)
	case helper.Postgres:
		err := helper.PublishPostgresService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		if yamlConfig.PostgresVersion == "" {
			yamlConfig.PostgresVersion = getOrChoose("Выберите версию postgres: ", yamlConfig.PostgresVersion, helper.GetAvailableVersions(helper.Postgres, yamlConfig))
		}
		return initEnvFile(yamlConfig)
	case helper.Sphinx:
		return helper.PublishSphinxService()
	case helper.Redis:
		return helper.PublishRedisService()
	case helper.Memcached:
		return helper.PublishMemcachedService()
	case helper.Mailhog:
		return helper.PublishMailhogService()
	case helper.PhpMyAdmin:
		return helper.PublishPhpMyAdminService()
	default:
		return fmt.Errorf("неизвестный сервис: %s", service)
	}
}
