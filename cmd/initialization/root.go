package initialization

import (
	"docky/config"
	"docky/utils"
	"docky/utils/globalHelper"
	"docky/yaml/helper"
	"fmt"
	"os"
	"path/filepath"
)

func InitDockerComposeFile() error {
	if err := handleExistingComposeFile(); err != nil {
		return err
	}

	yamlConfig := config.GetYamlConfig()
	yamlConfig.FrameworkName = utils.GetOrChoose("Ваш фреймворк: ", yamlConfig.FrameworkName, helper.AvailableFramework[:])
	yamlConfig.PhpVersion = utils.GetOrChoose("Выберите версию php: ", "", helper.GetAvailableVersions(helper.App, yamlConfig))

	switch yamlConfig.FrameworkName {
	case config.Laravel:
		if err := initLaravelConfig(yamlConfig); err != nil {
			return err
		}
	case config.Vanilla:
		initVanillaConfig(yamlConfig)
	case config.Symfony:
		initSymfonyConfig(yamlConfig)
	default:
		initDefaultConfig(yamlConfig)
	}

	if err := globalHelper.InitEnvFile(yamlConfig); err != nil {
		return err
	}

	return helper.BuildYaml(yamlConfig).Save()
}

func InitLaravel() error {
	siteDir := config.GetSiteDirPath()

	if !utils.IsDirEmpty(siteDir) {
		if !utils.AskYesNo("Директория с сайтом не пуста. Удалить всё и установить Laravel?") {
			return nil
		}
		if err := recreateDir(siteDir); err != nil {
			return err
		}
	}

	if err := globalHelper.ExecDockerCompose([]string{"build", helper.App}); err != nil {
		return err
	}

	if err := installLaravelProject(); err != nil {
		return err
	}

	if err := setupNodePackages(siteDir); err != nil {
		return err
	}

	globalHelper.DownContainers()
	return nil
}

func InitSymfony() error {
	siteDir := config.GetSiteDirPath()

	if !utils.IsDirEmpty(siteDir) {
		if !utils.AskYesNo("Директория с сайтом не пуста. Удалить всё и установить Symfony?") {
			return nil
		}
		if err := recreateDir(siteDir); err != nil {
			return err
		}
	}

	if err := globalHelper.ExecDockerCompose([]string{"build", helper.App}); err != nil {
		return err
	}

	if err := installSymfonyProject(); err != nil {
		return err
	}

	// if err := setupNodePackages(siteDir); err != nil {
	// 	return err
	// }

	globalHelper.DownContainers()
	return nil
}

func handleExistingComposeFile() error {
	composeFilePath := config.GetDockerComposeFilePath()
	if exists, _ := utils.FileIsExists(composeFilePath); !exists {
		return nil
	}

	if !utils.AskYesNo("Файл docker-compose.yml уже существует, создать новый?") {
		return nil
	}
	return os.Rename(composeFilePath, composeFilePath+config.Timestamp)
}

func chooseDbAndCache(yamlConfig *config.YamlConfig) {
	yamlConfig.DbType = utils.GetOrChoose("Выберите базу данных: ", "", helper.AvailableDb[:])
	switch yamlConfig.DbType {
	case helper.Mysql:
		yamlConfig.MysqlVersion = utils.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, helper.GetAvailableVersions(helper.Mysql, yamlConfig))
	case helper.Postgres:
		yamlConfig.PostgresVersion = utils.GetOrChoose("Выберите версию postgres: ", yamlConfig.PostgresVersion, helper.GetAvailableVersions(helper.Postgres, yamlConfig))
	}

	cache := utils.GetOrChoose("Выберите сервер кеширования: ", "", append(helper.AvailableServerCache[:], "Пропуск"))
	if cache != "Пропуск" {
		yamlConfig.ServerCache = cache
	}
}

func chooseNode(yamlConfig *config.YamlConfig) {
	if utils.AskYesNo("Добавлять node js?") {
		yamlConfig.CreateNode = true
		globalHelper.InitNode(yamlConfig)
	}
}

func initLaravelConfig(yamlConfig *config.YamlConfig) error {
	if _, err := globalHelper.IsDockerComposeAvailable(); err != nil {
		return err
	}

	chooseDbAndCache(yamlConfig)
	yamlConfig.CreateNode = true
	globalHelper.InitNode(yamlConfig)
	return nil
}

func initVanillaConfig(yamlConfig *config.YamlConfig) {
	chooseDbAndCache(yamlConfig)
	chooseNode(yamlConfig)
}

func initSymfonyConfig(yamlConfig *config.YamlConfig) {
	chooseDbAndCache(yamlConfig)
}

func initDefaultConfig(yamlConfig *config.YamlConfig) {
	yamlConfig.DbType = helper.Mysql
	if yamlConfig.MysqlVersion == "" {
		yamlConfig.MysqlVersion = utils.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, helper.GetAvailableVersions(helper.Mysql, yamlConfig))
	}

	chooseNode(yamlConfig)
	yamlConfig.CreateSphinx = utils.AskYesNo("Добавлять sphinx?")
}

func recreateDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("не удалось очистить директорию: %w", err)
	}
	return os.MkdirAll(dir, 0755)
}

func installLaravelProject() error {
	siteDir := config.GetSiteDirPath()
	dir := "laravel"

	args := []string{
		"run", "--rm",
		"--user", "docky", "--entrypoint", "php",
		helper.App, "/home/docky/.config/composer/vendor/bin/laravel", "new", dir,
	}

	if err := globalHelper.ExecDockerCompose(args); err != nil {
		return err
	}

	newPath := filepath.Join(siteDir, dir)
	if exists, _ := utils.FileIsExists(newPath); exists {
		return utils.MoveDirContents(newPath, siteDir)
	}
	return nil
}

func installSymfonyProject() error {
	isCli := utils.AskYesNo("Вы создаете консольное приложение Symfony?")

	args := []string{
		"run", "--rm",
		"--user", "docky", "--entrypoint", "composer",
		helper.App, "create-project", "symfony/skeleton", ".",
	}

	if err := globalHelper.ExecDockerCompose(args); err != nil {
		return err
	}

	if !isCli {
		if err := globalHelper.ExecDockerCompose([]string{
			"run", "--rm",
			"--user", "docky", "--entrypoint", "composer",
			helper.App, "require", "webapp",
		}); err != nil {
			return err
		}
	}

	return nil
}

func setupNodePackages(siteDir string) error {
	if exists, _ := utils.FileIsExists(filepath.Join(siteDir, "package.json")); !exists {
		return nil
	}

	if err := globalHelper.ExecDockerCompose([]string{"build", helper.Node}); err != nil {
		return err
	}

	return globalHelper.ExecDockerCompose([]string{
		"run", "--rm",
		"--user", "docky", "--entrypoint", "npm",
		helper.Node, "install",
	})
}
