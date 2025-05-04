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
		fmt.Println("✅ инициализация проекта завершена!")
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

func initNodeDir() error {
	path := filepath.Join(config.GetSiteDirPath(), yaml.NodePath)
	if utils.FileIsExists(path) {
		return nil
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("ошибка создания директории %s: %v", path, err)
	}
	return nil
}

func initNode() error {
	yaml.NodePath = os.Getenv(config.NodePathVarName)
	if yaml.NodePath == "" {
		yaml.NodePath = utils.ReadPath("Введите путь до директории с package.json относительно директории сайта. Например (local/js/vite или пустая строка): ")
	}
	yaml.NodePath = strings.TrimPrefix(yaml.NodePath, config.SitePathInContainer + "/")
	return initNodeDir()
}

func initEnvFile(recreate bool) error {
	envFileName := config.GetEnvFilePath()
	if recreate || !utils.FileIsExists(envFileName) {
		outFile, err := os.Create(envFileName)
		if err != nil {
			return err
		}
		defer outFile.Close()
		if yaml.PhpVersion == "" {
			yaml.PhpVersion = os.Getenv(config.PhpVersionVarName)
		}
		if yaml.MysqlVersion == "" {
			yaml.MysqlVersion = os.Getenv(config.MysqlVersionVarName)
		}
		if yaml.NodeVersion == "" {
			yaml.NodeVersion = os.Getenv(config.NodeVersionVarName)
		}
		if yaml.NodePath == "" {
			yaml.NodePath = os.Getenv(config.NodePathVarName)
		}
		if !strings.HasPrefix(yaml.NodePath, config.SitePathInContainer) {
			yaml.NodePath = filepath.Join(config.SitePathInContainer, yaml.NodePath)
		}
		data := []string{
			config.PhpVersionVarName + "=" + yaml.PhpVersion,
			config.MysqlVersionVarName + "=" + yaml.MysqlVersion,
			config.NodeVersionVarName + "=" + yaml.NodeVersion,
			config.NodePathVarName + "=" + yaml.NodePath,
		}
		if(yaml.SitePath == "") {
			yaml.SitePath = os.Getenv(config.SitePathVarName)
		}
		if yaml.SitePath != "" {
			data = append(data, config.SitePathVarName+"="+yaml.SitePath)
		}
		for _, line := range data {
			_, err := outFile.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func initDockerComposeFile() error {
	composeFilePath := config.GetDockerComposeFilePath()
	if utils.FileIsExists(composeFilePath) {
		fmt.Println()
		if utils.AskYesNo("Файл docker-compose.yml уже существует, создать новый?") {
			err := os.Rename(composeFilePath, composeFilePath+config.Timestamp)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
	}

	phpVersion := os.Getenv(config.PhpVersionVarName)
	mysqlVersion := os.Getenv(config.MysqlVersionVarName)
	isRecreate := false
	if phpVersion == "" {
		_, phpVersion = utils.ChooseFromList("Выберите версию php: ", yaml.AvailablePhpVersions[:])
		isRecreate = true
	}
	if mysqlVersion == "" {
		_, mysqlVersion = utils.ChooseFromList("Выберите версию mysql: ", yaml.AvailableMysqlVersions[:])
		isRecreate = true
	}
	yaml.PhpVersion = phpVersion
	yaml.MysqlVersion = mysqlVersion

	if utils.AskYesNo("Добавлять node js?") {
		yaml.CreateNode = true
		isRecreate = true
		initNode()
	}

	yaml.CreateSphinx = utils.AskYesNo("Добавлять sphinx?")
	err := initEnvFile(isRecreate)
	if err != nil {
		return err
	}
	return yaml.Create()
}
