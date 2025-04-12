package cmd

import (
	"docky/config"
	"docky/internal"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                config.ScriptName,
	Short:              "Программа для работы с docker-compose для битрикс проектов",
	DisableFlagParsing: true,                // важно!
	Args:               cobra.ArbitraryArgs, // принимает любые аргументы
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		if err := execDockerCompose(args); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	os.Setenv(config.UserGroupVarName, strconv.Itoa(os.Getegid()))
	internal.ExtractFilesInCache()
	fmt.Println("🚀 Запуск docky...")
}

func Execute() error {
	return rootCmd.Execute()
}

func execDockerCompose(args []string) error {
	cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
	cmd.Dir = config.GetWorkDirPath()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
