package commands

import (
	"fmt"
	"os"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/files"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/publishtools"

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
