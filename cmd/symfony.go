package cmd

import (
	"docky/config"
	"docky/utils/globalHelper"
	"docky/yaml/helper"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var symfonyCmd = &cobra.Command{
	Use:                "symfony",
	Short:              "Запускает bin/console команду в контейнере " + helper.App,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globalHelper.ValidateWorkDir()
		if err := execSymfonyInContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
	},
}

func init() {
	switch config.GetCurFramework() {
	case config.Symfony:
		rootCmd.AddCommand(symfonyCmd)
	}
}

func execSymfonyInContainer(args []string) error {
	execArgs := append([]string{
		"bin/console",
	}, args...)

	return execPhpInContainer(execArgs)
}
