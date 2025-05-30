package cmd

import (
	"docky/cmd/publish"
	"docky/config"
	"docky/internal"
	"docky/utils/globalHelper"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var service string
var file string

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Публикует файлы конфигурации",
	Run: func(cmd *cobra.Command, args []string) {
		globalHelper.ValidateWorkDir()
		var err error = nil
		text := "Файлы опубликованы!"

		if service != "" {
			err = publish.PublishService(service)
			text = "Сервис " + service + " опубликован!"
		} else if file != "" {
			err = publish.PublishFile(file)
			text = "Файл " + file + " опубликован!"
		} else {
			err = internal.PublishFiles()
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
		fmt.Println("✅ " + text)
	},
}

func init() {
	publishCmd.Flags().StringVar(&service, "service", "", "Опубликовать сервис в docker-compose")
	publishCmd.Flags().StringVar(&file, "file", "", "Опубликовать файл в директории"+config.ConfFilesDirName)
	rootCmd.AddCommand(publishCmd)
}
