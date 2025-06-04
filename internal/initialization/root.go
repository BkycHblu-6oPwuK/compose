package initialization

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/filetools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/readertools"
)

func InitDockerComposeFile() error {
	if err := handleExistingComposeFile(); err != nil {
		return err
	}

	yamlConfig := config.GetYamlConfig()
	yamlConfig.FrameworkName = framework.ParseFramework(readertools.GetOrChoose("Ваш фреймворк: ", yamlConfig.FrameworkName.String(), framework.GetAllStrings()))
	yamlConfig.PhpVersion = readertools.GetOrChoose("Выберите версию php: ", "", composefiletools.GetAvailableVersions(composefiletools.App, yamlConfig))

	switch yamlConfig.FrameworkName {
	case framework.Laravel:
		if err := initLaravelConfig(yamlConfig); err != nil {
			return err
		}
	case framework.Vanilla:
		initVanillaConfig(yamlConfig)
	case framework.Symfony:
		initSymfonyConfig(yamlConfig)
	default:
		initDefaultConfig(yamlConfig)
	}

	if err := globaltools.InitEnvFile(yamlConfig); err != nil {
		return err
	}

	return composefiletools.BuildYaml(yamlConfig).Save(config.GetDockerComposeFilePath())
}

func InitLaravel() error {
	siteDir := config.GetSiteDirPath()

	if !filetools.IsDirEmpty(siteDir) {
		if !readertools.AskYesNo("Директория с сайтом не пуста. Удалить всё и установить Laravel?") {
			return nil
		}
		if err := recreateDir(siteDir); err != nil {
			return err
		}
	}

	if err := globaltools.ExecDockerCompose([]string{"build", composefiletools.App}); err != nil {
		return err
	}

	if err := installLaravelProject(); err != nil {
		return err
	}

	if err := setupNodePackages(siteDir); err != nil {
		return err
	}

	globaltools.DownContainers()
	return nil
}

func InitSymfony() error {
	siteDir := config.GetSiteDirPath()

	if !filetools.IsDirEmpty(siteDir) {
		if !readertools.AskYesNo("Директория с сайтом не пуста. Удалить всё и установить Symfony?") {
			return nil
		}
		if err := recreateDir(siteDir); err != nil {
			return err
		}
	}

	if err := globaltools.ExecDockerCompose([]string{"build", composefiletools.App}); err != nil {
		return err
	}

	if err := installSymfonyProject(); err != nil {
		return err
	}

	globaltools.DownContainers()
	return nil
}

func handleExistingComposeFile() error {
	composeFilePath := config.GetDockerComposeFilePath()
	if exists, _ := filetools.FileIsExists(composeFilePath); !exists {
		return nil
	}

	if !readertools.AskYesNo("Файл docker-compose.yml уже существует, создать новый?") {
		return nil
	}
	return os.Rename(composeFilePath, composeFilePath+config.GetTimeStamp())
}

func chooseDbAndCache(yamlConfig *config.YamlConfig) {
	yamlConfig.DbType = readertools.GetOrChoose("Выберите базу данных: ", "", composefiletools.AvailableDb[:])
	switch yamlConfig.DbType {
	case composefiletools.Mysql:
		yamlConfig.MysqlVersion = readertools.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, composefiletools.GetAvailableVersions(composefiletools.Mysql, yamlConfig))
	case composefiletools.Postgres:
		yamlConfig.PostgresVersion = readertools.GetOrChoose("Выберите версию postgres: ", yamlConfig.PostgresVersion, composefiletools.GetAvailableVersions(composefiletools.Postgres, yamlConfig))
	}

	cache := readertools.GetOrChoose("Выберите сервер кеширования: ", "", append(composefiletools.AvailableServerCache[:], "Пропуск"))
	if cache != "Пропуск" {
		yamlConfig.ServerCache = cache
	}
}

func chooseNode(yamlConfig *config.YamlConfig) {
	if readertools.AskYesNo("Добавлять node js?") {
		yamlConfig.CreateNode = true
		globaltools.InitNode(yamlConfig)
	}
}

func initLaravelConfig(yamlConfig *config.YamlConfig) error {
	if _, err := globaltools.IsDockerComposeAvailable(); err != nil {
		return err
	}

	chooseDbAndCache(yamlConfig)
	yamlConfig.CreateNode = true
	globaltools.InitNode(yamlConfig)
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
	yamlConfig.DbType = composefiletools.Mysql
	if yamlConfig.MysqlVersion == "" {
		yamlConfig.MysqlVersion = readertools.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, composefiletools.GetAvailableVersions(composefiletools.Mysql, yamlConfig))
	}

	chooseNode(yamlConfig)
	yamlConfig.CreateSphinx = readertools.AskYesNo("Добавлять sphinx?")
}

func recreateDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("не удалось очистить директорию: %w", err)
	}
	return filetools.InitDirs(dir)
}

func installLaravelProject() error {
	siteDir := config.GetSiteDirPath()
	dir := "laravel"

	args := []string{
		"run", "--rm",
		"--user", "docky", "--entrypoint", "php",
		composefiletools.App, "/home/docky/.config/composer/vendor/bin/laravel", "new", dir,
	}

	if err := globaltools.ExecDockerCompose(args); err != nil {
		return err
	}

	newPath := filepath.Join(siteDir, dir)
	if exists, _ := filetools.FileIsExists(newPath); exists {
		return filetools.MoveDirContents(newPath, siteDir)
	}
	return nil
}

func installSymfonyProject() error {
	isCli := readertools.AskYesNo("Вы создаете консольное приложение Symfony?")

	args := []string{
		"run", "--rm",
		"--user", "docky", "--entrypoint", "composer",
		composefiletools.App, "create-project", "symfony/skeleton", ".",
	}

	if err := globaltools.ExecDockerCompose(args); err != nil {
		return err
	}

	if !isCli {
		if err := globaltools.ExecDockerCompose([]string{
			"run", "--rm",
			"--user", "docky", "--entrypoint", "composer",
			composefiletools.App, "require", "webapp",
		}); err != nil {
			return err
		}
	}

	return nil
}

func setupNodePackages(siteDir string) error {
	if exists, _ := filetools.FileIsExists(filepath.Join(siteDir, "package.json")); !exists {
		return nil
	}

	if err := globaltools.ExecDockerCompose([]string{"build", composefiletools.Node}); err != nil {
		return err
	}

	return globaltools.ExecDockerCompose([]string{
		"run", "--rm",
		"--user", "docky", "--entrypoint", "npm",
		composefiletools.Node, "install",
	})
}
