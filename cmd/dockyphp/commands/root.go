package commands

import (
	"docky/internal/config"
	"docky/internal/files"
	"docky/internal/globaltools"
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
		if err := globaltools.ExecDockerCompose(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
	},
}

func init() {
	err := files.ExtractFilesInCache()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}
