package cmd

import (
	"docky/utils/globalHelper"
	"docky/yaml/helper"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var npmCmd = &cobra.Command{
	Use:                "npm",
	Short:              "Запускает npm команду в контейнере " + helper.Node,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globalHelper.ValidateWorkDir()
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
		helper.Node, "npm",
	}, args...)

	return execDockerCompose(execArgs)
}
