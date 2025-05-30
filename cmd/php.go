package cmd

import (
	"docky/utils/globalHelper"
	"docky/yaml/helper"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var phpCmd = &cobra.Command{
	Use:                "php",
	Short:              "Запускает php команду в контейнере " + helper.App,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globalHelper.ValidateWorkDir()
		if err := execPhpInContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(phpCmd)
}

func execPhpInContainer(args []string) error {
	execArgs := append([]string{
		"exec", "-it",
		"--user", "docky",
		helper.App, "php",
	}, args...)

	return execDockerCompose(execArgs)
}
