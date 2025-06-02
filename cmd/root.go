package cmd

import (
	"docky/config"
	"docky/internal"
	"docky/utils/globalHelper"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                config.ScriptName + " [docker compose commands]",
	Short:              "Утилита для работы с docker compose в Bitrix-проектах",
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		if err := globalHelper.ExecDockerCompose(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
	},
}

func init() {
	err := internal.ExtractFilesInCache()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}
