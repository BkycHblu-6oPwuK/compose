package commands

import (
	"fmt"
	"os"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"

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
	case framework.Symfony:
		rootCmd.AddCommand(symfonyCmd)
	}
}

func execSymfonyInContainer(args []string) error {
	execArgs := append([]string{
		"bin/console",
	}, args...)

	return execPhpInContainer(execArgs)
}
