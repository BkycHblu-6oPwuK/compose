package commands

/**
* @todo отрефакторить сброс, с флагом --all вообще сбрасывать всё, что можно, по умолчанию оставлять кастомизацию. Перенести сброс в пакет docky/yaml
 */

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/files"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/service"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/filetools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/orderedmap"

	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "reset",
	Short: "Сбрасывает docker-compose.yml под актуальный формат",
	Run: func(cmd *cobra.Command, args []string) {
		globaltools.ValidateWorkDir()
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
	if err := files.CleanCacheDir(); err != nil {
		return err
	}
	if err := files.ExtractFilesInCache(); err != nil {
		return err
	}
	dockerComposeFilePath := config.GetDockerComposeFilePath()
	if fileExists, _ := filetools.FileIsExists(dockerComposeFilePath); !fileExists {
		return fmt.Errorf("файл %s не найден", dockerComposeFilePath)
	}
	composeFile, err := composefile.Load(dockerComposeFilePath)
	if err != nil {
		return fmt.Errorf("ошибка при загрузке docker-compose.yml: %v", err)
	}
	err = resetServices(composeFile.Services)
	if err != nil {
		return fmt.Errorf("ошибка при сбросе сервисов: %v", err)
	}
	err = os.Rename(dockerComposeFilePath, dockerComposeFilePath+config.GetTimeStamp())
	if err != nil {
		return err
	}
	hostsRename()
	if err := composeFile.Save(dockerComposeFilePath); err != nil {
		return err
	}
	return globaltools.ExecDockerCompose([]string{"build"})
}

func resetServices(services *orderedmap.OrderedMap[string, service.Service]) error {
	nodeVersion := ""
	nodePath := ""
	sitePath := ""
	phpVersion := ""
	mysqlVersion := ""
	services.ForEach(func(name string, curService service.Service) {
		if curService.Build.Context != "" {
			curService.Build.Context = "${" + config.DockerPathVarName + "}"
		}
		if curService.Build.Dockerfile != "" {
			curService.Build.Dockerfile = replaceDockerPath(curService.Build.Dockerfile, &phpVersion)
		}
		if curService.Build.Args != nil {
			if curService.Build.Args["NODE_VERSION"] != "" {
				nodeVersion = curService.Build.Args["NODE_VERSION"]
				curService.Build.Args["NODE_VERSION"] = "${" + config.NodeVersionVarName + "}"
			}
			if curService.Build.Args["DOCKER_PATH"] != "" {
				delete(curService.Build.Args, "DOCKER_PATH")
			}
			if curService.Build.Args["NODE_PATH"] != "" {
				delete(curService.Build.Args, "NODE_PATH")
			}
		}
		if curService.Image != "" && strings.HasPrefix(curService.Image, "mysql:") {
			parts := strings.Split(curService.Image, ":")
			if len(parts) == 2 {
				mysqlVersion = parts[1]
			}
			curService.Image = "mysql:${" + config.MysqlVersionVarName + "}"
		}
		if curService.Volumes != nil {
			filtered := curService.Volumes[:0]
			for i, volume := range curService.Volumes {
				if strings.HasPrefix(volume, "./_docker") || strings.HasPrefix(volume, "./vendor/beeralex/compose/src/_docker") {
					curService.Volumes[i] = replaceDockerPath(volume, &phpVersion)
				} else {
					parts := strings.Split(volume, ":")
					if len(parts) == 2 && parts[1] == "/var/www" {
						if !(strings.HasPrefix(volume, "./site:") || strings.HasPrefix(volume, "./site/:")) {
							sitePath = parts[0]
						}
						curService.Volumes[i] = "${SITE_PATH}:" + parts[1]
					}
				}
				switch curService.Volumes[i] {
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
					filtered = append(filtered, curService.Volumes[i])
				}
			}
			curService.Volumes = filtered
		}
		if name == composefiletools.App {
			if curService.Environment == nil {
				curService.Environment = make(map[string]string)
			}
			curService.ExtraHosts = []string{"host.docker.internal:host-gateway"}
		} else if name == composefiletools.Node {
			if curService.Command == nil {
				curService.Command = "tail -f /dev/null"
			}
			if curService.WorkingDir == "" {
				curService.WorkingDir = "${" + config.NodePathVarName + "}"
			}
		}
		if curService.Environment != nil {
			newEnv := make(map[string]string)
			for key, value := range curService.Environment {
				if key != "DOCKER_PATH" {
					newEnv[key] = value
				}
			}
			if name == composefiletools.App {
				newEnv["PHP_IDE_CONFIG"] = "serverName=xdebugServer"
				newEnv["XDEBUG_TRIGGER"] = "testTrig"
			}
			curService.Environment = newEnv
		}
		services.Set(name, curService)
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
		if fileExists, _ := filetools.FileIsExists(sitepath); fileExists {
			err := os.Rename(sitepath, config.GetSiteDirPath())
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Ошибка при переименовании директории сайта: %v\n", err)
			}
		}
	}

	return globaltools.InitEnvFile(yamlConfig)
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
	if fileExists, isDir := filetools.FileIsExists(hostsPath); fileExists && !isDir {
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
