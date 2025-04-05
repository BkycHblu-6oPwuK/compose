package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "docky",
	Short: "Программа для работы с docker-compose для битрикс проектов",
}

func Execute() error {
	return rootCmd.Execute()
}
