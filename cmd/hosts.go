package cmd

import (
	"docky/hosts"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Команда для работы с hosts",
}
var pushHostsModuleCmd = &cobra.Command{
	Use:   "push",
	Short: "Переносит записи в hosts из локального hosts в директории проекта",
	Run: func(cmd *cobra.Command, args []string) {
		validateWorkDir()
		err := hosts.PushToHosts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
		}
		fmt.Println("✅ Записи в hosts перенесены!")
	},
}

func init() {
	hostsCmd.AddCommand(pushHostsModuleCmd)
	rootCmd.AddCommand(hostsCmd)
}
