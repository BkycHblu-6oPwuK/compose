package cmd

import (
	"docky/yaml/helper"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var pm2Cmd = &cobra.Command{
	Use:                "pm2",
	Short:              "Запускает pm2 команду в контейнере " + helper.Node,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		validateWorkDir()
		if err := execPm2InContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			os.Exit(1)
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
		helper.Node, "pm2",
	}, args...)

	return execDockerCompose(execArgs)
}
