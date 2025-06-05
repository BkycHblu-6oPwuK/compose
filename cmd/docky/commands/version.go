package commands

import (
	"fmt"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Aliases: []string{"v"},
	Short: "Версия " + config.ScriptName,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("✅ " + config.ScriptVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
