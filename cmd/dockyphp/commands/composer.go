package commands

import (
	"docky/internal/composefiletools"
	"docky/internal/globaltools"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var composerCmd = &cobra.Command{
	Use:                "composer",
	Short:              "Запускает composer команду в контейнере " + composefiletools.App,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globaltools.ValidateWorkDir()
		if err := execComposerInContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
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
		composefiletools.App, "composer",
	}, args...)

	return globaltools.ExecDockerCompose(execArgs)
}
