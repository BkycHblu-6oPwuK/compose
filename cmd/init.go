package cmd

import (
	"docky/config"
	"docky/utils"
	"docky/utils/globalHelper"
	"docky/yaml/helper"
	"fmt"
	"os"
	"path/filepath"

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
		err = globalHelper.InitSiteDir()
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

	yamlConfig.FrameworkName = utils.GetOrChoose("Ваш фреймворк: ", yamlConfig.FrameworkName, helper.AvailableFramework[:])
	yamlConfig.PhpVersion = utils.GetOrChoose("Выберите версию php: ", "", helper.GetAvailableVersions(helper.App, yamlConfig))

	switch yamlConfig.FrameworkName {
	case config.Laravel:
		_, err := isDockerComposeAvailable()
		if err != nil {
			return err
		}

		yamlConfig.DbType = utils.GetOrChoose("Выберите базу данных: ", "", helper.AvailableDb[:])

		switch yamlConfig.DbType {
		case helper.Mysql:
			yamlConfig.MysqlVersion = utils.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, helper.GetAvailableVersions(helper.Mysql, yamlConfig))
		case helper.Postgres:
			yamlConfig.PostgresVersion = utils.GetOrChoose("Выберите версию postgres: ", yamlConfig.PostgresVersion, helper.GetAvailableVersions(helper.Postgres, yamlConfig))
		}

		serverCache := utils.GetOrChoose("Выберите сервер кеширования: ", "", append(helper.AvailableServerCache[:], "Пропуск"))
		if serverCache != "Пропуск" {
			yamlConfig.ServerCache = serverCache
		}

		yamlConfig.CreateNode = true
		globalHelper.InitNode(yamlConfig)
	default:
		yamlConfig.DbType = helper.Mysql
		if yamlConfig.MysqlVersion == "" {
			yamlConfig.MysqlVersion = utils.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, helper.GetAvailableVersions(helper.Mysql, yamlConfig))
		}

		if utils.AskYesNo("Добавлять node js?") {
			yamlConfig.CreateNode = true
			globalHelper.InitNode(yamlConfig)
		}

		yamlConfig.CreateSphinx = utils.AskYesNo("Добавлять sphinx?")
	}
	if err := globalHelper.InitEnvFile(yamlConfig); err != nil {
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
