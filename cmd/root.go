package cmd

import (
	"docky/config"
	"docky/internal"
	"docky/yaml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                config.ScriptName + " [docker compose commands]",
	Short:              "Утилита для работы с docker compose в Bitrix-проектах",
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		if err := execDockerCompose(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	err := internal.ExtractFilesInCache()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func isDockerComposeAvailable() ([]string, error) {
	if err := exec.Command("docker", "compose", "version").Run(); err == nil {
		return []string{"docker", "compose"}, nil
	}
	if err := exec.Command("docker-compose", "version").Run(); err == nil {
		return []string{"docker-compose"}, nil
	}
	return nil, errors.New("docker compose не установлен или не запущен")
}

func execDockerCompose(args []string) error {
	dockerCmd, err := isDockerComposeAvailable()
	if err != nil {
		return err
	}
	sitePath := config.GetSiteDirPath()
	os.Setenv(config.UserGroupVarName, strconv.Itoa(os.Getegid()))
	os.Setenv(config.DockerPathVarName, config.GetCurrentDockerFileDirPath())
	os.Setenv(config.SitePathVarName, sitePath)
	yaml.NodePath = filepath.Join(sitePath, strings.TrimPrefix(os.Getenv(config.NodePathVarName), config.SitePathInContainer + "/"))
	initSiteDir()
	if yaml.NodePath != "" {
		initNodeDir()
	}
 	cmd := exec.Command(dockerCmd[0], append(dockerCmd[1:], args...)...)
	cmd.Dir = config.GetWorkDirPath()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
