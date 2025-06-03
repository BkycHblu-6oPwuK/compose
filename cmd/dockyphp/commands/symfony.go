package commands

import (
	"docky/internal/composefiletools"
	"docky/internal/config"
	"docky/internal/globaltools"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var symfonyCmd = &cobra.Command{
	Use:                "symfony",
	Short:              "Запускает bin/console команду в контейнере " + composefiletools.App,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globaltools.ValidateWorkDir()
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
