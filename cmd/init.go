package cmd

import (
	"docky/config"
	"docky/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Инициализация проекта",
	Run: func(cmd *cobra.Command, args []string) {
		err := initSiteDir()
		if err != nil {
			fmt.Println("❌ Ошибка инициализации проекта:", err)
			return
		}
		fmt.Println("✅ инициализация проекта завершена!")
	},
}

func init() {	
	rootCmd.AddCommand(initCmd)
}

func initSiteDir() error {
	siteDirPath := config.GetSiteDirPath()
	if utils.DirIsExists(siteDirPath) {
		fmt.Println("Директория сайта уже существует")
		return nil
	}

	err := os.Mkdir(siteDirPath, 0755)
	if err != nil {
		return fmt.Errorf("ошибка создания директории сайта: %v", err)
	}
	fmt.Println("Директория сайта успешно создана")
	return nil
}