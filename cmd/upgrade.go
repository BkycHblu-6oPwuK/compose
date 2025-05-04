package cmd

import (
	"docky/config"
	"docky/utils"
	myYaml "docky/yaml"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Очищает кэш (директория _docker)",
	Run: func(cmd *cobra.Command, args []string) {
		err := upgrade()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
		}
		fmt.Println("✅ Кэш очищен!")
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func upgrade() error {
	filePath := config.GetDockerComposeFilePath()
	if !utils.FileIsExists(filePath) {
		return fmt.Errorf("файл %s не найден", filePath)
	}
	oldData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(oldData, &root); err != nil {
		return err
	}

	var (
		phpVersion   string
		mysqlVersion string
		nodeVersion  string
		nodePath     string
		sitePath     string
	)

	// Ищем корень "services"
	for i := 0; i < len(root.Content[0].Content); i += 2 {
		keyNode := root.Content[0].Content[i]
		valNode := root.Content[0].Content[i+1]

		if keyNode.Value == "services" {
			for j := 0; j < len(valNode.Content); j += 2 {
				serviceName := valNode.Content[j].Value
				serviceMap := valNode.Content[j+1]

				//var buildNode *yaml.Node
				for k := 0; k < len(serviceMap.Content); k += 2 {
					key := serviceMap.Content[k]
					val := serviceMap.Content[k+1]

					// Обработка build
					if key.Value == "build" {
						//buildNode = val
						for b := 0; b < len(val.Content); b += 2 {
							bk := val.Content[b]
							bv := val.Content[b+1]

							replaceDockerPath := func(value string) string {
								value = strings.Replace(value, "./_docker", "${"+config.DockerPathVarName+"}", 1)
								value = strings.Replace(value, "./vendor/beeralex/compose/src/_docker", "${"+config.DockerPathVarName+"}", 1)
								return value
							}

							if bk.Value == "context" {
								bv.Value = "${" + config.DockerPathVarName + "}"
							}

							if bk.Value == "dockerfile" {
								if serviceName == "app" {
									re := regexp.MustCompile(`php-(\d+\.\d+)`)
									if match := re.FindStringSubmatch(bv.Value); len(match) > 1 {
										phpVersion = match[1]
										bv.Value = strings.Replace(bv.Value, match[0], "php-${"+config.PhpVersionVarName+"}", 1)
									}
								}

								bv.Value = replaceDockerPath(bv.Value)
							}

							if bk.Value == "args" && bv.Kind == yaml.MappingNode {
								for a := 0; a < len(bv.Content); a += 2 {
									argKey := bv.Content[a]
									argVal := bv.Content[a+1]

									if argKey.Value == "NODE_VERSION" {
										nodeVersion = argVal.Value
										argVal.Value = "${" + config.NodeVersionVarName + "}"
									}
									if argKey.Value == "NODE_PATH" {
										nodePath = argVal.Value
										argVal.Value = "${" + config.NodePathVarName + "}"
									}

									if argKey.Value == "DOCKER_PATH" {
										// удалить
										bv.Content = append(bv.Content[:a], bv.Content[a+2:]...)
										a -= 2
									}
								}
							}
						}
					}

					if key.Value == "volumes" && val.Kind == yaml.SequenceNode {
						for idx, item := range val.Content {
							if item.Kind == yaml.ScalarNode {
								path := item.Value
								parts := strings.Split(path, ":")
								if len(parts) == 2 && parts[1] == "/var/www" {
									if !(strings.HasPrefix(path, "./site:") || strings.HasPrefix(path, "./site/:")) {
										sitePath = parts[0]
									}
									path = "${SITE_PATH}:" + parts[1]
								} else {
									switch {
									case strings.HasPrefix(path, "./_docker"):
										path = strings.Replace(path, "./_docker", "${"+config.DockerPathVarName+"}", 1)
									case strings.HasPrefix(path, "./vendor/beeralex/compose/src/_docker"):
										path = strings.Replace(path, "./vendor/beeralex/compose/src/_docker", "${"+config.DockerPathVarName+"}", 1)
									}
								}

								// Подмена php-версии в пути
								re := regexp.MustCompile(`php-(\d+\.\d+)`)
								if match := re.FindStringSubmatch(path); len(match) > 1 {
									path = strings.Replace(path, match[0], "php-${PHP_VERSION}", 1)
								}

								// Назначаем обратно
								val.Content[idx].Value = path
							}
						}
					}

					if key.Value == "environment" {
						switch val.Kind {
						case yaml.MappingNode:
							newContent := []*yaml.Node{}
							for i := 0; i < len(val.Content); i += 2 {
								k := val.Content[i]
								v := val.Content[i+1]
								if k.Value != "DOCKER_PATH" {
									newContent = append(newContent, k, v)
								}
							}
							val.Content = newContent
						}
					}

					if key.Value == "image" && serviceName == "mysql" {
						parts := strings.Split(val.Value, ":")
						if len(parts) == 2 {
							mysqlVersion = parts[1]
						}
						val.Value = "mysql:${" + config.MysqlVersionVarName + "}"
					}
				}
			}
		}

		// volumes
		if keyNode.Value == "volumes" {
			for v := 1; v < len(valNode.Content); v += 2 {
				valNode.Content[v].Kind = yaml.MappingNode
				valNode.Content[v].Tag = "!!map"
				valNode.Content[v].Content = nil // очищаем
			}
		}
	}

	// Сохраняем результат с сохранением порядка
	err = os.Rename(filePath, filePath+config.Timestamp)
	if err != nil {
		return err
	}
	newYaml, err := yaml.Marshal(&root)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, newYaml, 0644); err != nil {
		return err
	}

	fmt.Println("Файл успешно преобразован: new-docker-compose.yml")
	fmt.Printf("PHP_VERSION=%s\n", phpVersion)
	fmt.Printf("MYSQL_VERSION=%s\n", mysqlVersion)
	fmt.Printf("NODE_VERSION=%s\n", nodeVersion)
	fmt.Printf("NODE_PATH=%s\n", nodePath)
	fmt.Printf("SITE_PATH=%s\n", sitePath)

	if phpVersion != "" && phpVersion != "${"+config.PhpVersionVarName+"}" {
		myYaml.PhpVersion = phpVersion
	}
	if mysqlVersion != "" && mysqlVersion != "${"+config.MysqlVersionVarName+"}" {
		myYaml.MysqlVersion = mysqlVersion
	}
	if nodeVersion != "" && nodeVersion != "${"+config.NodeVersionVarName+"}" {
		myYaml.NodeVersion = nodeVersion
	}
	if nodePath != "" && nodePath != "${"+config.NodePathVarName+"}" {
		myYaml.NodePath = nodePath
	}
	if sitePath != "" && sitePath != "${"+config.SitePathVarName+"}" {
		myYaml.SitePath = sitePath
	}

	initEnvFile(true)

	return nil
}
