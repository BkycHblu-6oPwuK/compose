package publish

import (
	"docky/config"
	"docky/internal"
	"docky/utils"
	"docky/utils/globalHelper"
	"docky/yaml/helper"
	"fmt"
	"path/filepath"
)

func PublishFile(file string) error {
	switch file {
	case "php.ini", "xdebug.ini":
		yamlConfig := config.GetYamlConfig()
		pathToIni := filepath.Join(helper.App, "php-" + yamlConfig.PhpVersion, file)
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework(), pathToIni)
		if err := internal.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToIni)); err != nil {
			return err
		}
		return helper.PublishVolumes([]string{
			helper.App,
		}, map[string][]string{
			helper.App: {
				"${" + config.ConfPathVarName + "}" + "/app/php-${" + config.PhpVersionVarName + "}/"+file+":/usr/local/etc/php/conf.d/" + file,
			},
		})
	default:
		return fmt.Errorf("неизвестный файл: %s", file)
	}
}

func PublishService(service string) error {
	switch service {
	case helper.Node:
		err := helper.PublishNodeService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		err = globalHelper.InitNode(yamlConfig)
		if err != nil {
			return err
		}
		return globalHelper.InitEnvFile(yamlConfig)
	case helper.Mysql:
		err := helper.PublishMysqlService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		if yamlConfig.MysqlVersion == "" {
			yamlConfig.MysqlVersion = utils.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, helper.GetAvailableVersions(helper.Mysql, yamlConfig))
		}
		return globalHelper.InitEnvFile(yamlConfig)
	case helper.Postgres:
		err := helper.PublishPostgresService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		if yamlConfig.PostgresVersion == "" {
			yamlConfig.PostgresVersion = utils.GetOrChoose("Выберите версию postgres: ", yamlConfig.PostgresVersion, helper.GetAvailableVersions(helper.Postgres, yamlConfig))
		}
		return globalHelper.InitEnvFile(yamlConfig)
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
