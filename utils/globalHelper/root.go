package globalHelper

import (
	"docky/config"
	"docky/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ValidateWorkDir() {
	if fileExists, _ := utils.FileIsExists(config.GetDockerComposeFilePath()); !fileExists {
		fmt.Fprintf(os.Stderr, "❌ Ошибка: не инициализирован docker-compose.yml\n")
		os.Exit(1)
	}
}

func InitSiteDir() error {
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

func InitConfDir() error {
	path := config.GetConfFilesDirPath()
	if fileExists, _ := utils.FileIsExists(path); fileExists {
		return nil
	}

	err := os.Mkdir(path, 0755)
	if err != nil {
		return fmt.Errorf("ошибка создания директории для файлов конфигурации: %v", err)
	}
	return nil
}

func InitNodeDir(yamlConfig *config.YamlConfig) error {
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

func InitNode(yamlConfig *config.YamlConfig) error {
	switch yamlConfig.FrameworkName {
	case config.Bitrix:
		if yamlConfig.NodePath == "" {
			yamlConfig.NodePath = utils.ReadPath("Введите путь до директории с package.json относительно директории сайта. Например (local/js/vite или пустая строка): ")
		}
		yamlConfig.NodePath = strings.TrimPrefix(yamlConfig.NodePath, config.SitePathInContainer)
		return InitNodeDir(yamlConfig)
	case config.Laravel:
		yamlConfig.NodePath = config.SitePathInContainer
	}
	return nil
}

func InitEnvFile(yamlConfig *config.YamlConfig) error {
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