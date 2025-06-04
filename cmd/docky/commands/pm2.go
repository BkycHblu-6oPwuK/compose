package commands

import (
	"docky/internal/composefiletools"
	"docky/internal/globaltools"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var pm2Cmd = &cobra.Command{
	Use:                "pm2",
	Short:              "Запускает pm2 команду в контейнере " + composefiletools.Node,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globaltools.ValidateWorkDir()
		if err := execPm2InContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pm2Cmd)
}

func execPm2InContainer(args []string) error {
	execArgs := append([]string{
		"exec", "-it",
		"--user", "docky",
		composefiletools.Node, "pm2",
	}, args...)

	return globaltools.ExecDockerCompose(execArgs)
}
