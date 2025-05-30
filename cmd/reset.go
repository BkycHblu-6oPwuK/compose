package cmd

/**
* @todo отрефакторить сброс, с флагом --all вообще сбрасывать всё, что можно, по умолчанию оставлять кастомизацию. Перенести сброс в пакет docky/yaml
 */

import (
	"docky/config"
	"docky/internal"
	"docky/utils"
	"docky/utils/globalHelper"
	myYaml "docky/yaml"
	"docky/yaml/helper"
	myService "docky/yaml/service"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "reset",
	Short: "Сбрасывает docker-compose.yml под актуальный формат",
	Run: func(cmd *cobra.Command, args []string) {
		globalHelper.ValidateWorkDir()
		err := reset()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
		fmt.Println("✅ docker-compose.yml обновлён, проверьте его на наличие ошибок. Проверьте файл .env на наличие новых переменных окружения. Старый файл docker-compose.yml переименован.")
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func reset() error {
	if err := internal.CleanCacheDir(); err != nil {
		return err
	}
	if err := internal.ExtractFilesInCache(); err != nil {
		return err
	}
	dockerComposeFilePath := config.GetDockerComposeFilePath()
	if fileExists, _ := utils.FileIsExists(dockerComposeFilePath); !fileExists {
		return fmt.Errorf("файл %s не найден", dockerComposeFilePath)
	}
	composeFile, err := myYaml.Load()
	if err != nil {
		return fmt.Errorf("ошибка при загрузке docker-compose.yml: %v", err)
	}
	err = resetServices(composeFile.Services)
	if err != nil {
		return fmt.Errorf("ошибка при сбросе сервисов: %v", err)
	}
	err = os.Rename(dockerComposeFilePath, dockerComposeFilePath+config.Timestamp)
	if err != nil {
		return err
	}
	hostsRename()
	if err := composeFile.Save(); err != nil {
		return err
	}
	return execDockerCompose([]string{"build"})
}

func resetServices(services *utils.OrderedMap[string, myService.Service]) error {
	nodeVersion := ""
	nodePath := ""
	sitePath := ""
	phpVersion := ""
	mysqlVersion := ""
	services.ForEach(func(name string, service myService.Service) {
		if service.Build.Context != "" {
			service.Build.Context = "${" + config.DockerPathVarName + "}"
		}
		if service.Build.Dockerfile != "" {
			service.Build.Dockerfile = replaceDockerPath(service.Build.Dockerfile, &phpVersion)
		}
		if service.Build.Args != nil {
			if service.Build.Args["NODE_VERSION"] != "" {
				nodeVersion = service.Build.Args["NODE_VERSION"]
				service.Build.Args["NODE_VERSION"] = "${" + config.NodeVersionVarName + "}"
			}
			if service.Build.Args["DOCKER_PATH"] != "" {
				delete(service.Build.Args, "DOCKER_PATH")
			}
			if service.Build.Args["NODE_PATH"] != "" {
				delete(service.Build.Args, "NODE_PATH")
			}
		}
		if service.Image != "" && strings.HasPrefix(service.Image, "mysql:") {
			parts := strings.Split(service.Image, ":")
			if len(parts) == 2 {
				mysqlVersion = parts[1]
			}
			service.Image = "mysql:${" + config.MysqlVersionVarName + "}"
		}
		if service.Volumes != nil {
			filtered := service.Volumes[:0]
			for i, volume := range service.Volumes {
				if strings.HasPrefix(volume, "./_docker") || strings.HasPrefix(volume, "./vendor/beeralex/compose/src/_docker") {
					service.Volumes[i] = replaceDockerPath(volume, &phpVersion)
				} else {
					parts := strings.Split(volume, ":")
					if len(parts) == 2 && parts[1] == "/var/www" {
						if !(strings.HasPrefix(volume, "./site:") || strings.HasPrefix(volume, "./site/:")) {
							sitePath = parts[0]
						}
						service.Volumes[i] = "${SITE_PATH}:" + parts[1]
					}
				}
				switch service.Volumes[i] {
				case "${DOCKER_PATH}/app/php-${PHP_VERSION}/php.ini:/usr/local/etc/php/conf.d/php.ini":
					continue
				case "${DOCKER_PATH}/app/php-fpm.conf:/usr/local/etc/php-fpm.d/zzzzwww.conf":
					continue
				case "${DOCKER_PATH}/app/nginx.conf:/etc/nginx/conf.d/nginx.conf":
					continue
				case "${DOCKER_PATH}/sphinx/sphinx.conf:/usr/local/etc/sphinx.conf":
					continue
				case "${DOCKER_PATH}/nginx/conf.d/:/etc/nginx/conf.d":
					continue
				case "${DOCKER_PATH}/app/php-${PHP_VERSION}/xdebug.ini:/usr/local/etc/php/conf.d/xdebug.ini":
					continue
				case "${DOCKER_PATH}/app/nginx:/etc/nginx/conf.d":
					continue
				case "${DOCKER_PATH}/nginx/conf.d:/etc/nginx/conf.d":
					continue
				case "${DOCKER_PATH}/app/nginx/:/etc/nginx/conf.d":
					continue
				default:
					filtered = append(filtered, service.Volumes[i])
				}
			}
			service.Volumes = filtered
		}
		if name == helper.App {
			if service.Environment == nil {
				service.Environment = make(map[string]string)
			}
			service.ExtraHosts = []string{"host.docker.internal:host-gateway"}
		} else if name == helper.Node {
			if service.Command == nil {
				service.Command = "tail -f /dev/null"
			}
			if service.WorkingDir == "" {
				service.WorkingDir = "${" + config.NodePathVarName + "}"
			}
		}
		if service.Environment != nil {
			newEnv := make(map[string]string)
			for key, value := range service.Environment {
				if key != "DOCKER_PATH" {
					newEnv[key] = value
				}
			}
			if name == helper.App {
				newEnv["PHP_IDE_CONFIG"] = "serverName=xdebugServer"
				newEnv["XDEBUG_TRIGGER"] = "testTrig"
			}
			service.Environment = newEnv
		}
		service.Networks = []string{helper.Compose}
		services.Set(name, service)
	})
	yamlConfig := config.GetYamlConfig()
	if phpVersion != "" && phpVersion != "${"+config.PhpVersionVarName+"}" {
		yamlConfig.PhpVersion = phpVersion
	}
	if mysqlVersion != "" && mysqlVersion != "${"+config.MysqlVersionVarName+"}" {
		yamlConfig.MysqlVersion = mysqlVersion
	}
	if nodeVersion != "" && nodeVersion != "${"+config.NodeVersionVarName+"}" {
		yamlConfig.NodeVersion = nodeVersion
	}
	if nodePath != "" && nodePath != "${"+config.NodePathVarName+"}" {
		yamlConfig.NodePath = nodePath
	}
	if sitePath != "" && sitePath != "${"+config.SitePathVarName+"}" {
		sitepath := filepath.Join(config.GetWorkDirPath(), sitePath)
		if fileExists, _ := utils.FileIsExists(sitepath); fileExists {
			err := os.Rename(sitepath, config.GetSiteDirPath())
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Ошибка при переименовании директории сайта: %v\n", err)
			}
		}
	}

	return globalHelper.InitEnvFile(yamlConfig)
}

func replaceDockerPath(value string, phpVersion *string) string {
	value = strings.Replace(value, "./_docker", "${"+config.DockerPathVarName+"}", 1)
	value = strings.Replace(value, "./vendor/beeralex/compose/src/_docker", "${"+config.DockerPathVarName+"}", 1)
	re := regexp.MustCompile(`php-(\d+\.\d+)`)
	if match := re.FindStringSubmatch(value); len(match) > 1 {
		*phpVersion = match[1]
		value = strings.Replace(value, match[0], "php-${"+config.PhpVersionVarName+"}", 1)
	}
	return value
}

func hostsRename() {
	hostsPath := filepath.Join(config.GetWorkDirPath(), "hosts")
	targetPath := config.GetLocalHostsFilePath()
	if fileExists, isDir := utils.FileIsExists(hostsPath); fileExists && !isDir {
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			fmt.Printf("ошибка при создании директории: %v\n", err)
			return
		}
		err := os.Rename(hostsPath, targetPath)
		if err != nil {
			fmt.Printf("ошибка при переименовывании файла hosts: %v\n", err)
		}
	}
}