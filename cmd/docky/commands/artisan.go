package commands

import (
	"fmt"
	"os"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"

	"github.com/spf13/cobra"
)

var artisanCmd = &cobra.Command{
	Use:                "artisan",
	Short:              "Запускает artisan команду в контейнере " + composefiletools.App,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globaltools.ValidateWorkDir()
		if err := execArtisanInContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
	},
}

func init() {
	switch config.GetCurFramework() {
	case config.Laravel:
		rootCmd.AddCommand(artisanCmd)
	}
}

func execArtisanInContainer(args []string) error {
	execArgs := append([]string{
		"artisan",
	}, args...)

	return execPhpInContainer(execArgs)
}
