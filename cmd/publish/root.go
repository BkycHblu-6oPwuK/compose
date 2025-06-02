package publish

import (
	"docky/config"
	"docky/internal"
	"docky/utils"
	"docky/utils/globalHelper"
	"docky/yaml/helper"
	"docky/yaml/service"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func PublishFile(file string) error {
	switch file {
	case "php.ini", "xdebug.ini":
		yamlConfig := config.GetYamlConfig()
		pathToIni := filepath.Join(helper.App, "php-"+yamlConfig.PhpVersion, file)
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework(), pathToIni)
		if err := internal.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToIni), false); err != nil {
			return err
		}
		return helper.PublishVolumes([]string{
			helper.App,
		}, map[string][]string{
			helper.App: {
				"${" + config.ConfPathVarName + "}" + "/" + helper.App + "/php-${" + config.PhpVersionVarName + "}/" + file + ":/usr/local/etc/php/conf.d/" + file,
			},
		}, nil)
	case "simlinks":
		dirPath := filepath.Join(config.GetConfFilesDirPath(), helper.App)
		pathToSimlinks := filepath.Join(dirPath, "simlinks")
		if exists, _ := utils.FileIsExists(pathToSimlinks); !exists {
			if exists, _ := utils.FileIsExists(dirPath); !exists {
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					return fmt.Errorf("ошибка при создании директории: %w", err)
				}
			}
			file, err := os.Create(pathToSimlinks)
			if err != nil {
				return fmt.Errorf("ошибка при создании файла: %w", err)
			}
			defer file.Close()
		}
		return helper.PublishVolumes([]string{
			helper.App,
		}, map[string][]string{
			helper.App: {
				"${" + config.ConfPathVarName + "}" + "/" + helper.App + "/" + "simlinks:/usr/simlinks_extra",
			},
		}, nil)
	case "cron_tasks":
		curFramework := config.GetCurFramework()
		switch curFramework {
		case config.Bitrix, config.Symfony, config.Vanilla:
			pathToCron := filepath.Join(helper.App, "cron")
			filePath := filepath.Join(config.DockerFilesDirName, curFramework, pathToCron)
			if err := internal.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToCron), true); err != nil {
				return err
			}
			return helper.PublishVolumes([]string{
				helper.App,
			}, map[string][]string{
				helper.App: {
					"${" + config.ConfPathVarName + "}" + "/" + helper.App + "/cron:/var/spool/cron/crontabs",
				},
			}, nil)
		default:
			return fmt.Errorf("для вашего фреймворка cron не предустановлен: %s", curFramework)
		}
	case "nginx_conf":
		pathToConf := filepath.Join(helper.Nginx, "conf.d")
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework(), pathToConf)
		if err := internal.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToConf), true); err != nil {
			return err
		}
		return helper.PublishVolumes([]string{
			helper.Nginx,
		}, map[string][]string{
			helper.Nginx: {
				"${" + config.ConfPathVarName + "}" + "/" + helper.Nginx + "/conf.d:/etc/nginx/conf.d",
			},
		}, func(s *service.Service) (isContinue bool, err error) {
			filtered := s.Volumes[:0]
			for _, volume := range s.Volumes {
				if !strings.Contains(volume, "/etc/nginx/conf.d") {
					filtered = append(filtered, volume)
				}
			}
			s.Volumes = filtered
			return true, nil
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
		yamlConfig.MysqlVersion = utils.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, helper.GetAvailableVersions(helper.Mysql, yamlConfig))
		return globalHelper.InitEnvFile(yamlConfig)
	case helper.Postgres:
		err := helper.PublishPostgresService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		yamlConfig.PostgresVersion = utils.GetOrChoose("Выберите версию postgres: ", yamlConfig.PostgresVersion, helper.GetAvailableVersions(helper.Postgres, yamlConfig))
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

func PublishDockerFile(service string) error {
	switch service {
	case helper.App:
		yamlConfig := config.GetYamlConfig()
		pathToDockerfile := filepath.Join(helper.App, "php-"+yamlConfig.PhpVersion, helper.Dockerfile)
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework(), pathToDockerfile)
		if err := internal.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToDockerfile), true); err != nil {
			return err
		}
		return helper.PublishDockerfile(service, "${"+config.ConfPathVarName+"}/"+helper.App+"/php-${"+config.PhpVersionVarName+"}/"+helper.Dockerfile)
	case helper.Node, helper.Nginx:
		pathToDockerfile := filepath.Join(service, helper.Dockerfile)
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework(), pathToDockerfile)
		if err := internal.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToDockerfile), true); err != nil {
			return err
		}
		return helper.PublishDockerfile(service, "${"+config.ConfPathVarName+"}/"+service+"/"+helper.Dockerfile)
	default:
		return fmt.Errorf("неизвестный докерфайл для публикации: %s", service)
	}
}
