package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:                "update",
	Short:              "Запускает обновление скрипта",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if err := execUpdate(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func execUpdate() error {
	cmd := exec.Command("bash", "-c", "curl -sSL https://raw.githubusercontent.com/BkycHblu-6oPwuK/docky/main/scripts/install.sh | sudo sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
