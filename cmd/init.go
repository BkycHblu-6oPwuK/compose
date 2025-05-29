package cmd

import (
	"docky/config"
	"docky/utils"
	"docky/yaml/helper"
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
	if fileExists, _ := utils.FileIsExists(siteDirPath); fileExists {
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
	if fileExists, _ := utils.FileIsExists(path); fileExists {
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
	if fileExists, _ := utils.FileIsExists(composeFilePath); fileExists {
		if !utils.AskYesNo("Файл docker-compose.yml уже существует, создать новый?") {
			return nil
		}
		if err := os.Rename(composeFilePath, composeFilePath+config.Timestamp); err != nil {
			return err
		}
	}

	yamlConfig := config.GetYamlConfig()

	yamlConfig.FrameworkName = getOrChoose("Ваш фреймворк: ", yamlConfig.FrameworkName, helper.AvailableFramework[:])
	yamlConfig.PhpVersion = getOrChoose("Выберите версию php: ", "", helper.GetAvailableVersions(helper.App, yamlConfig))

	switch yamlConfig.FrameworkName {
	case config.Laravel:
		_, err := isDockerComposeAvailable()
		if err != nil {
			return err
		}

		yamlConfig.DbType = getOrChoose("Выберите базу данных: ", "", helper.AvailableDb[:])

		switch yamlConfig.DbType {
		case helper.Mysql:
			yamlConfig.MysqlVersion = getOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, helper.GetAvailableVersions(helper.Mysql, yamlConfig))
		case helper.Postgres:
			yamlConfig.PostgresVersion = getOrChoose("Выберите версию postgres: ", yamlConfig.PostgresVersion, helper.GetAvailableVersions(helper.Postgres, yamlConfig))
		}

		serverCache := getOrChoose("Выберите сервер кеширования: ", "", append(helper.AvailableServerCache[:], "Пропуск"))
		if serverCache != "Пропуск" {
			yamlConfig.ServerCache = serverCache
		}

		yamlConfig.CreateNode = true
		initNode(yamlConfig)
	default:
		yamlConfig.DbType = helper.Mysql
		if yamlConfig.MysqlVersion == "" {
			yamlConfig.MysqlVersion = getOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, helper.GetAvailableVersions(helper.Mysql, yamlConfig))
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

	return helper.BuildYaml(yamlConfig).Save()
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

	if err := execDockerCompose([]string{"build", helper.App}); err != nil {
		return err
	}

	dir := "laravel"
	execArgs := []string{
		"run", "--rm",
		"--user", "docky", "--entrypoint", "php",
		helper.App, "/home/docky/.config/composer/vendor/bin/laravel", "new", dir,
	}
	if err := execDockerCompose(execArgs); err != nil {
		return err
	}

	path := filepath.Join(siteDir, dir)
	if fileExists, _ := utils.FileIsExists(path); fileExists {
		if err := utils.MoveDirContents(path, siteDir); err != nil {
			return err
		}
	}

	if fileExists, _ := utils.FileIsExists(filepath.Join(siteDir, "package.json")); fileExists {
		if err := execDockerCompose([]string{"build", helper.Node}); err != nil {
			return err
		}
		execArgs := []string{
			"run", "--rm",
			"--user", "docky", "--entrypoint", "npm",
			helper.Node, "install",
		}
		if err := execDockerCompose(execArgs); err != nil {
			return err
		}
	}
	downContainers()
	return nil
}
