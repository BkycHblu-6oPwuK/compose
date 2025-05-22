package cmd

import (
	"docky/config"
	"docky/utils"
	"docky/yaml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Инициализация проекта",
	Run: func(cmd *cobra.Command, args []string) {
		err := initDockerComposeFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
		err = initSiteDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
		yamlConfig := config.GetYamlConfig()
		switch yamlConfig.FrameworkName {
		case config.Laravel:
			fmt.Println("Инициализация ларавел")
			initLaravel()
		}
		fmt.Println("✅ Инициализация проекта завершена!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initSiteDir() error {
	siteDirPath := config.GetSiteDirPath()
	if utils.FileIsExists(siteDirPath) {
		return nil
	}

	err := os.Mkdir(siteDirPath, 0755)
	if err != nil {
		return fmt.Errorf("ошибка создания директории сайта: %v", err)
	}
	return nil
}

func initNodeDir(yamlConfig *config.YamlConfig) error {
	path := filepath.Join(config.GetSiteDirPath(), yamlConfig.NodePath)
	if utils.FileIsExists(path) {
		return nil
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("ошибка создания директории %s: %v", path, err)
	}
	return nil
}

func initNode(yamlConfig *config.YamlConfig) error {
	switch yamlConfig.FrameworkName {
	case config.Bitrix:
		if yamlConfig.NodePath == "" {
			yamlConfig.NodePath = utils.ReadPath("Введите путь до директории с package.json относительно директории сайта. Например (local/js/vite или пустая строка): ")
		}
		yamlConfig.NodePath = strings.TrimPrefix(yamlConfig.NodePath, config.SitePathInContainer)
		return initNodeDir(yamlConfig)
	case config.Laravel:
		yamlConfig.NodePath = config.SitePathInContainer
	}
	return nil
}

func initEnvFile(yamlConfig *config.YamlConfig) error {
	outFile, err := os.Create(config.GetEnvFilePath())
	if err != nil {
		return err
	}
	defer outFile.Close()

	if !strings.HasPrefix(yamlConfig.NodePath, config.SitePathInContainer) {
		yamlConfig.NodePath = filepath.Join(config.SitePathInContainer, yamlConfig.NodePath)
	}
	
	data := []string{
		config.DockyFrameworkVarName + "=" + yamlConfig.FrameworkName,
		config.PhpVersionVarName + "=" + yamlConfig.PhpVersion,
		config.MysqlVersionVarName + "=" + yamlConfig.MysqlVersion,
		config.PostgresVersionVarName + "=" + yamlConfig.PostgresVersion,
		config.NodeVersionVarName + "=" + yamlConfig.NodeVersion,
		config.NodePathVarName + "=" + yamlConfig.NodePath,
	}

	if yamlConfig.SitePath != "" {
		data = append(data, config.SitePathVarName+"="+yamlConfig.SitePath)
	}

	for _, line := range data {
		if _, err := outFile.WriteString(line + "\n"); err != nil {
			return err
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			if err := os.Setenv(parts[0], parts[1]); err != nil {
				return err
			}
		}
	}

	return nil
}

func getOrChoose(prompt, value string, options []string) string {
	if value == "" {
		_, value = utils.ChooseFromList(prompt, options)
	}
	return value
}

func initDockerComposeFile() error {
	composeFilePath := config.GetDockerComposeFilePath()
	if utils.FileIsExists(composeFilePath) {
		if !utils.AskYesNo("Файл docker-compose.yml уже существует, создать новый?") {
			return nil
		}
		if err := os.Rename(composeFilePath, composeFilePath+config.Timestamp); err != nil {
			return err
		}
	}

	yamlConfig := config.GetYamlConfig()
	yamlFile := yaml.NewYamlFile(yamlConfig)

	yamlConfig.FrameworkName = getOrChoose("Ваш фреймворк: ", yamlConfig.FrameworkName, yaml.AvailableFramework[:])
	yamlConfig.PhpVersion = getOrChoose("Выберите версию php: ", "", yaml.GetAvailableVersions(yaml.App, yamlFile.Config))

	switch yamlConfig.FrameworkName {
	case config.Laravel:
		_, err := isDockerComposeAvailable()
		if err != nil {
			return err
		}

		yamlConfig.DbType = getOrChoose("Выберите базу данных: ", "", yaml.AvailableDb[:])

		switch yamlConfig.DbType {
		case yaml.Mysql:
			yamlConfig.MysqlVersion = getOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, yaml.GetAvailableVersions(yaml.Mysql, yamlFile.Config))
		case yaml.Postgres:
			yamlConfig.PostgresVersion = getOrChoose("Выберите версию postgres: ", yamlConfig.PostgresVersion, yaml.GetAvailableVersions(yaml.Postgres, yamlFile.Config))
		}

		serverCache := getOrChoose("Выберите сервер кеширования: ", "", append(yaml.AvailableServerCache[:], "Пропуск"))
		if serverCache != "Пропуск" {
			yamlConfig.ServerCache = serverCache
		}

		yamlConfig.CreateNode = true
		initNode(yamlConfig)
	default:
		yamlConfig.DbType = yaml.Mysql
		if yamlConfig.MysqlVersion == "" {
			yamlConfig.MysqlVersion = getOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, yaml.GetAvailableVersions(yaml.Mysql, yamlFile.Config))
		}

		if utils.AskYesNo("Добавлять node js?") {
			yamlConfig.CreateNode = true
			initNode(yamlConfig)
		}

		yamlConfig.CreateSphinx = utils.AskYesNo("Добавлять sphinx?")
	}
	if err := initEnvFile(yamlConfig); err != nil {
		return err
	}

	return yamlFile.Create()
}

func initLaravel() error {
	siteDir := config.GetSiteDirPath()

	siteIsEmpty := utils.IsDirEmpty(siteDir)
	if !siteIsEmpty && !utils.AskYesNo("Директория с сайтом не пуста. Удалить всё и установить Laravel?") {
		return nil
	}

	if !siteIsEmpty {
		if err := os.RemoveAll(siteDir); err != nil {
			return fmt.Errorf("не удалось очистить директорию: %w", err)
		}
		if err := os.MkdirAll(siteDir, 0755); err != nil {
			return fmt.Errorf("не удалось создать директорию: %w", err)
		}
	}

	if err := execDockerCompose([]string{"build", yaml.App}); err != nil {
		return err
	}

	dir := "laravel"
	execArgs := []string{
		"run", "--rm",
		"--user", "docky", "--entrypoint", "php",
		yaml.App, "/home/docky/.config/composer/vendor/bin/laravel", "new", dir,
	}
	if err := execDockerCompose(execArgs); err != nil {
		return err
	}

	path := filepath.Join(siteDir, dir)
	if utils.FileIsExists(path) {
		if err := utils.MoveDirContents(path, siteDir); err != nil {
			return err
		}
	}

	if utils.FileIsExists(filepath.Join(siteDir, "package.json")) {
		if err := execDockerCompose([]string{"build", yaml.Node}); err != nil {
			return err
		}
		execArgs := []string{
			"run", "--rm",
			"--user", "docky", "--entrypoint", "npm",
			yaml.Node, "install",
		}
		if err := execDockerCompose(execArgs); err != nil {
			return err
		}
	}
	downContainers()
	return nil
}
