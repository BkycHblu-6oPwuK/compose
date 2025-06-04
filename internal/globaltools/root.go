package globaltools

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/filetools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/readertools"
)

func ValidateWorkDir() {
	if fileExists, _ := filetools.FileIsExists(config.GetDockerComposeFilePath()); !fileExists {
		fmt.Fprintf(os.Stderr, "❌ Ошибка: не инициализирован docker-compose.yml\n")
		os.Exit(1)
	}
}

func InitSiteDir() error {
	siteDirPath := config.GetSiteDirPath()
	if fileExists, _ := filetools.FileIsExists(siteDirPath); fileExists {
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
	if fileExists, _ := filetools.FileIsExists(path); fileExists {
		return nil
	}

	err := os.Mkdir(path, 0755)
	if err != nil {
		return fmt.Errorf("ошибка создания директории для файлов конфигурации: %v", err)
	}
	return nil
}

func InitNodeDir(yamlConfig *config.YamlConfig) error {
	path := strings.TrimPrefix(yamlConfig.NodePath, composefiletools.SitePathInContainer)
	if path == "" {
		return nil
	}
	path = filepath.Join(config.GetSiteDirPath(), path)
	if err := filetools.InitDirs(path); err != nil {
		return err
	}
	return nil
}

func InitNode(yamlConfig *config.YamlConfig) error {
	switch yamlConfig.FrameworkName {
	case framework.Bitrix, framework.Symfony, framework.Vanilla:
		if yamlConfig.NodePath == "" {
			yamlConfig.NodePath = readertools.ReadPath("Введите путь до директории с package.json относительно директории сайта. Например (local/js/vite или пустая строка): ")
		}
		return InitNodeDir(yamlConfig)
	case framework.Laravel:
		yamlConfig.NodePath = composefiletools.SitePathInContainer
	}
	return nil
}

func InitEnvFile(yamlConfig *config.YamlConfig) error {
	outFile, err := os.Create(config.GetEnvFilePath())
	if err != nil {
		return err
	}
	defer outFile.Close()

	if !strings.HasPrefix(yamlConfig.NodePath, composefiletools.SitePathInContainer) {
		yamlConfig.NodePath = filepath.Join(composefiletools.SitePathInContainer, yamlConfig.NodePath)
	}

	data := []string{
		config.DockyFrameworkVarName + "=" + yamlConfig.FrameworkName.String(),
		config.PhpVersionVarName + "=" + yamlConfig.PhpVersion,
		config.MysqlVersionVarName + "=" + yamlConfig.MysqlVersion,
		config.MariadbVersionVarName + "=" + yamlConfig.MariadbVersion,
		config.PostgresVersionVarName + "=" + yamlConfig.PostgresVersion,
		config.NodeVersionVarName + "=" + yamlConfig.NodeVersion,
		config.NodePathVarName + "=" + yamlConfig.NodePath,
	}

	if yamlConfig.UserGroup != "" {
		data = append(data, config.UserGroupVarName+"="+yamlConfig.UserGroup)
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

func IsDockerComposeAvailable() ([]string, error) {
	if err := exec.Command("docker", "compose", "version").Run(); err == nil {
		return []string{"docker", "compose"}, nil
	}
	if err := exec.Command("docker-compose", "version").Run(); err == nil {
		return []string{"docker-compose"}, nil
	}
	return nil, errors.New("docker compose не установлен или не запущен")
}

func ExecDockerCompose(args []string) error {
	dockerCmd, err := IsDockerComposeAvailable()
	if err != nil {
		return err
	}
	os.Setenv(config.UserGroupVarName, config.GetUserGroup())
	os.Setenv(config.DockerPathVarName, config.GetCurrentDockerFileDirPath())
	os.Setenv(config.ConfPathVarName, config.GetConfFilesDirPath())
	os.Setenv(config.SitePathVarName, config.GetSiteDirPath())
	cmd := exec.Command(dockerCmd[0], append(dockerCmd[1:], args...)...)
	cmd.Dir = config.GetWorkDirPath()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func DownContainers() {
	_ = ExecDockerCompose([]string{"down"})
}
