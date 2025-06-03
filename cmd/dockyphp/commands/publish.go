package commands

import (
	"docky/internal/config"
	"docky/internal/files"
	"docky/internal/globaltools"
	"docky/internal/publishtools"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var serviceFlag string
var dockerfileFlag string
var fileFlag string

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Публикует файлы конфигурации",
	Run: func(cmd *cobra.Command, args []string) {
		globaltools.ValidateWorkDir()
		var err error = nil
		text := "Файлы опубликованы!"

		if serviceFlag != "" {
			err = publishtools.PublishService(serviceFlag)
			text = "Сервис " + serviceFlag + " опубликован!"
		} else if fileFlag != "" {
			err = publishtools.PublishFile(fileFlag)
			text = "Файл " + fileFlag + " опубликован!"
		} else if dockerfileFlag != "" {
			err = publishtools.PublishDockerFile(dockerfileFlag)
			text = "Докерфайл " + dockerfileFlag + " опубликован!"
		} else {
			err = files.PublishFiles()
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
		fmt.Println("✅ " + text)
	},
}

func init() {
	publishCmd.Flags().StringVar(&serviceFlag, "service", "", "Опубликовать сервис в docker-compose")
	publishCmd.Flags().StringVar(&fileFlag, "file", "", "Опубликовать файл в директории"+config.ConfFilesDirName)
	publishCmd.Flags().StringVar(&dockerfileFlag, "dockerfile", "", "Опубликовать докерфайл в директории"+config.ConfFilesDirName)
	rootCmd.AddCommand(publishCmd)
}
