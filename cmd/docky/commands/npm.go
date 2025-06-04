package commands

import (
	"fmt"
	"os"

	"github.com/BkycHblu-6oPwuK/docky/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/internal/globaltools"

	"github.com/spf13/cobra"
)

var npmCmd = &cobra.Command{
	Use:                "npm",
	Short:              "Запускает npm команду в контейнере " + composefiletools.Node,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globaltools.ValidateWorkDir()
		if err := execNpmInContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(npmCmd)
}

func execNpmInContainer(args []string) error {
	execArgs := append([]string{
		"exec", "-it",
		"--user", "docky",
		composefiletools.Node, "npm",
	}, args...)

	return globaltools.ExecDockerCompose(execArgs)
}
