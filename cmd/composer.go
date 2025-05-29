package cmd

import (
	"docky/yaml/helper"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var composerCmd = &cobra.Command{
	Use:                "composer",
	Short:              "Запускает composer команду в контейнере " + helper.App,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		validateWorkDir()
		if err := execComposerInContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(composerCmd)
}

func execComposerInContainer(args []string) error {
	execArgs := append([]string{
		"exec", "-it",
		"--user", "docky",
		helper.App, "composer",
	}, args...)

	return execDockerCompose(execArgs)
}
