package commands

import (
	"docky/internal/composefiletools"
	"docky/internal/globaltools"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var npxCmd = &cobra.Command{
	Use:                "npx",
	Short:              "Запускает npx команду в контейнере " + composefiletools.Node,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globaltools.ValidateWorkDir()
		if err := execNpxInContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(npxCmd)
}

func execNpxInContainer(args []string) error {
	execArgs := append([]string{
		"exec", "-it",
		"--user", "docky",
		composefiletools.Node, "npx",
	}, args...)

	return globaltools.ExecDockerCompose(execArgs)
}
