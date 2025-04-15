package cmd

import (
	"docky/yaml"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var phpCmd = &cobra.Command{
	Use:                "php",
	Short:              "Запускает php команду в контейнере " + yaml.App,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if err := execPhpInContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(phpCmd)
}

func execPhpInContainer(args []string) error {
	execArgs := append([]string{
		"exec", "-it",
		"--user", "docky",
		yaml.App, "php",
	}, args...)

	return execDockerCompose(execArgs)
}
