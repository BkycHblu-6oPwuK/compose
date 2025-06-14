package publishtools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/files"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/symlinkstools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/service"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/filetools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/readertools"
)

func PublishFile(file string) error {
	switch file {
	case "php.ini", "xdebug.ini":
		yamlConfig := config.GetYamlConfig()
		pathToIni := filepath.Join(composefiletools.App, "php-"+yamlConfig.PhpVersion, file)
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework().String(), pathToIni)
		if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToIni), false); err != nil {
			return err
		}
		return composefiletools.PublishVolumes(map[string][]string{
			composefiletools.App: {
				composefiletools.GetPhpConfVolumePath(file, true),
			},
		}, nil)
	case symlinkstools.FileName:
		pathTosymlinks := filepath.Join(config.GetConfFilesDirPath(), composefiletools.App, symlinkstools.FileName)
		if exists, _ := filetools.FileIsExists(pathTosymlinks); !exists {
			if err := filetools.InitDirs(filepath.Dir(pathTosymlinks)); err != nil {
				return err
			}

			file, err := os.Create(pathTosymlinks)
			if err != nil {
				return fmt.Errorf("ошибка при создании файла: %w", err)
			}
			defer file.Close()
		}
		return composefiletools.PublishVolumes(map[string][]string{
			composefiletools.App: {
				composefiletools.GetsymlinksConfVolumePath(),
			},
		}, nil)
	case "cron_tasks":
		curFramework := config.GetCurFramework()
		switch curFramework {
		case framework.Bitrix, framework.Symfony, framework.Vanilla:
			pathToCron := filepath.Join(composefiletools.App, "cron")
			filePath := filepath.Join(config.DockerFilesDirName, curFramework.String(), pathToCron)
			if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToCron), true); err != nil {
				return err
			}
			return composefiletools.PublishVolumes(map[string][]string{
				composefiletools.App: {
					composefiletools.GetCronConfVolumePath(),
				},
			}, nil)
		default:
			return fmt.Errorf("для вашего фреймворка cron не предустановлен: %s", curFramework)
		}
	case "nginx_conf":
		pathToConf := filepath.Join(composefiletools.Nginx, "conf.d")
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework().String(), pathToConf)
		if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToConf), true); err != nil {
			return err
		}
		return composefiletools.PublishVolumes(map[string][]string{
			composefiletools.Nginx: {
				composefiletools.GetNginxConfVolumePath(""),
			},
		}, func(b *service.ServiceBuilder) (isContinue bool, err error) {
			b.FilterVolumes(func(volume string) bool {
				return !strings.Contains(volume, composefiletools.GetNginxConfPathInContainer())
			})
			return true, nil
		})
	case "mysql_conf":
		pathToConf := filepath.Join(composefiletools.Mysql, "my.cnf")
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework().String(), pathToConf)
		if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToConf), true); err != nil {
			return err
		}
		return composefiletools.PublishVolumes(map[string][]string{
			composefiletools.Mysql: {
				composefiletools.GetMysqlCnfPath(true),
			},
			composefiletools.Mariadb: {
				composefiletools.GetMysqlCnfPath(true),
			},
		}, func(b *service.ServiceBuilder) (isContinue bool, err error) {
			b.RemoveVolume(composefiletools.GetMysqlCnfPath(false))
			return true, nil
		})
	case "postgres_conf":
		pathToConf := filepath.Join(composefiletools.Postgres, "postgresql.conf")
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework().String(), pathToConf)
		if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToConf), true); err != nil {
			return err
		}
		return composefiletools.PublishVolumes(map[string][]string{
			composefiletools.Postgres: {
				composefiletools.GetPostgresConfPath(true),
			},
		}, func(b *service.ServiceBuilder) (isContinue bool, err error) {
			b.RemoveVolume(composefiletools.GetPostgresConfPath(false))
			return true, nil
		})
	case "supervisord_conf":
		pathToConf := filepath.Join(composefiletools.App, "supervisord.conf")
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework().String(), pathToConf)
		if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToConf), true); err != nil {
			return err
		}
		return composefiletools.PublishVolumes(map[string][]string{
			composefiletools.App: {
				composefiletools.GetSupervisordConfPath(),
			},
		}, nil)
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
	case composefiletools.Mariadb:
		err := composefiletools.PublishMariadbService()
		if err != nil {
			return err
		}
		yamlConfig := config.GetYamlConfig()
		yamlConfig.MariadbVersion = readertools.GetOrChoose("Выберите версию mariadb: ", yamlConfig.MariadbVersion, composefiletools.GetAvailableVersions(composefiletools.Mariadb, yamlConfig))
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
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework().String(), pathToDockerfile)
		if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToDockerfile), true); err != nil {
			return err
		}
		return composefiletools.PublishDockerfile(service, composefiletools.GetPhpConfComposePath(composefiletools.Dockerfile, true))
	case composefiletools.Node, composefiletools.Nginx:
		pathToDockerfile := filepath.Join(service, composefiletools.Dockerfile)
		filePath := filepath.Join(config.DockerFilesDirName, config.GetCurFramework().String(), pathToDockerfile)
		if err := files.PublishFile(filePath, filepath.Join(config.GetConfFilesDirPath(), pathToDockerfile), true); err != nil {
			return err
		}
		return composefiletools.PublishDockerfile(service, composefiletools.GetVarNameString(config.ConfPathVarName)+"/"+service+"/"+composefiletools.Dockerfile)
	default:
		return fmt.Errorf("неизвестный докерфайл для публикации: %s", service)
	}
}
