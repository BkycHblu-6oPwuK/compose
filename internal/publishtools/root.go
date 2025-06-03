package publishtools

import (
	"docky/internal/composefiletools"
	"docky/internal/config"
	"docky/internal/files"
	"docky/internal/globaltools"
	"docky/internal/symlinkstools"
	"docky/pkg/composefile/service"
	"docky/pkg/filetools"
	"docky/pkg/readertools"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func PublishFile(file string) error {
	switch file {
	case "php.ini", "xdebug.ini":
		yamlConfig := config.GetYamlConfig()
		pathToIni := filepath.Join(composefiletools.App, "php-"+yamlConfig.PhpVersion, file)
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework(), pathToIni)
		if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToIni), false); err != nil {
			return err
		}
		return composefiletools.PublishVolumes([]string{
			composefiletools.App,
		}, map[string][]string{
			composefiletools.App: {
				composefiletools.GetPhpConfVolumePath(file, true),
			},
		}, nil)
	case symlinkstools.FileName:
		pathToSymlinks := filepath.Join(config.GetConfFilesDirPath(), composefiletools.App, symlinkstools.FileName)
		if exists, _ := filetools.FileIsExists(pathToSymlinks); !exists {
			if err := filetools.InitDirs(filepath.Dir(pathToSymlinks)); err != nil {
				return err
			}

			file, err := os.Create(pathToSymlinks)
			if err != nil {
				return fmt.Errorf("ошибка при создании файла: %w", err)
			}
			defer file.Close()
		}
		return composefiletools.PublishVolumes([]string{
			composefiletools.App,
		}, map[string][]string{
			composefiletools.App: {
				composefiletools.GetSymlinksConfVolumePath(),
			},
		}, nil)
	case "cron_tasks":
		curFramework := config.GetCurFramework()
		switch curFramework {
		case config.Bitrix, config.Symfony, config.Vanilla:
			pathToCron := filepath.Join(composefiletools.App, "cron")
			filePath := filepath.Join(config.DockerFilesDirName, curFramework, pathToCron)
			if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToCron), true); err != nil {
				return err
			}
			return composefiletools.PublishVolumes([]string{
				composefiletools.App,
			}, map[string][]string{
				composefiletools.App: {
					composefiletools.GetCronConfVolumePath(),
				},
			}, nil)
		default:
			return fmt.Errorf("для вашего фреймворка cron не предустановлен: %s", curFramework)
		}
	case "nginx_conf":
		pathToConf := filepath.Join(composefiletools.Nginx, "conf.d")
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework(), pathToConf)
		if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToConf), true); err != nil {
			return err
		}
		return composefiletools.PublishVolumes([]string{
			composefiletools.Nginx,
		}, map[string][]string{
			composefiletools.Nginx: {
				composefiletools.GetNginxConfVolumePath(""),
			},
		}, func(s *service.Service) (isContinue bool, err error) {
			filtered := s.Volumes[:0]
			for _, volume := range s.Volumes {
				if !strings.Contains(volume, composefiletools.GetNginxConfPathInContainer()) {
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
	case composefiletools.Node:
		err := composefiletools.PublishNodeService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		err = globaltools.InitNode(yamlConfig)
		if err != nil {
			return err
		}
		return globaltools.InitEnvFile(yamlConfig)
	case composefiletools.Mysql:
		err := composefiletools.PublishMysqlService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		yamlConfig.MysqlVersion = readertools.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, composefiletools.GetAvailableVersions(composefiletools.Mysql, yamlConfig))
		return globaltools.InitEnvFile(yamlConfig)
	case composefiletools.Postgres:
		err := composefiletools.PublishPostgresService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		yamlConfig.PostgresVersion = readertools.GetOrChoose("Выберите версию postgres: ", yamlConfig.PostgresVersion, composefiletools.GetAvailableVersions(composefiletools.Postgres, yamlConfig))
		return globaltools.InitEnvFile(yamlConfig)
	case composefiletools.Sphinx:
		return composefiletools.PublishSphinxService()
	case composefiletools.Redis:
		return composefiletools.PublishRedisService()
	case composefiletools.Memcached:
		return composefiletools.PublishMemcachedService()
	case composefiletools.Mailhog:
		return composefiletools.PublishMailhogService()
	case composefiletools.PhpMyAdmin:
		return composefiletools.PublishPhpMyAdminService()
	default:
		return fmt.Errorf("неизвестный сервис: %s", service)
	}
}

func PublishDockerFile(service string) error {
	switch service {
	case composefiletools.App:
		yamlConfig := config.GetYamlConfig()
		pathToDockerfile := filepath.Join(composefiletools.App, "php-"+yamlConfig.PhpVersion, composefiletools.Dockerfile)
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework(), pathToDockerfile)
		if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToDockerfile), true); err != nil {
			return err
		}
		return composefiletools.PublishDockerfile(service, composefiletools.GetPhpConfComposePath(composefiletools.Dockerfile, true))
	case composefiletools.Node, composefiletools.Nginx:
		pathToDockerfile := filepath.Join(service, composefiletools.Dockerfile)
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework(), pathToDockerfile)
		if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToDockerfile), true); err != nil {
			return err
		}
		return composefiletools.PublishDockerfile(service, composefiletools.GetVarNameString(config.ConfPathVarName)+"/"+service+"/"+composefiletools.Dockerfile)
	default:
		return fmt.Errorf("неизвестный докерфайл для публикации: %s", service)
	}
}
