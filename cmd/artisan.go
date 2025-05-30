package cmd

import (
	"docky/config"
	"docky/utils/globalHelper"
	"docky/yaml/helper"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var artisanCmd = &cobra.Command{
	Use:                "artisan",
	Short:              "Запускает artisan команду в контейнере " + helper.App,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globalHelper.ValidateWorkDir()
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
