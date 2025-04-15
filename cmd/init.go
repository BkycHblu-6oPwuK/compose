package cmd

import (
	"docky/config"
	"docky/utils"
	"docky/yaml"
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

func initNodeDir(path string) error {
	if utils.FileIsExists(path) {
		return nil
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("ошибка создания директории %s: %v", path, err)
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

	_, yaml.PhpVersion = utils.ChooseFromList("Выберите версию php: ", yaml.AvailablePhpVersions[:])
	_, yaml.MysqlVersion = utils.ChooseFromList("Выберите версию php: ", yaml.AvailableMysqlVersions[:])

	if utils.AskYesNo("Добавлять node js?") {
		yaml.CreateNode = true
		yaml.NodePath = utils.ReadPath("Введите путь до директории с package.json относительно директории сайта. Например (local/js/vite или пустая строка): ")
		initNodeDir(filepath.Join(config.GetSiteDirPath(), yaml.NodePath))
	}

	yaml.CreateSphinx = utils.AskYesNo("Добавлять sphinx?")

	return yaml.Create()
}
