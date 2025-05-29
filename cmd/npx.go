package cmd

import (
	"docky/yaml/helper"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var npxCmd = &cobra.Command{
	Use:                "npx",
	Short:              "Запускает npx команду в контейнере " + helper.Node,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		validateWorkDir()
		if err := execNpxInContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			os.Exit(1)
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
		helper.Node, "npx",
	}, args...)

	return execDockerCompose(execArgs)
}
