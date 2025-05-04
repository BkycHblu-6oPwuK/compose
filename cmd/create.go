package cmd

import (
	"docky/certs"
	"docky/config"
	"docky/hosts"
	"docky/internal"
	"docky/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Команда для создания (cайт, домен)",
}

var createSiteModuleCmd = &cobra.Command{
	Use:   "site",
	Short: "Создает новый сайт (директория, сертификаты, запись в hosts)",
	Run: func(cmd *cobra.Command, args []string) {
		err := createSite()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
		}
		fmt.Println("✅")
	},
}

var createDomainModuleCmd = &cobra.Command{
	Use:   "domain",
	Short: "Создает новый домен (сертификаты, запись в hosts)",
	Run: func(cmd *cobra.Command, args []string) {
		err := createDomain()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
		}
		fmt.Println("✅")
	},
}

func init() {
	createCmd.AddCommand(createSiteModuleCmd)
	createCmd.AddCommand(createDomainModuleCmd)
	rootCmd.AddCommand(createCmd)
}

func createSite() error {
	var err error = nil
	domain := readDomain()
	dirPath := filepath.Join(config.GetSiteDirPath(), domain)
	if !utils.FileIsExists(dirPath) {
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("ошибка создания директории сайта: %v", err)
		}
	}
	err = createCerts(domain, filepath.Join(config.SitePathInContainer, domain))
	if err != nil {
		return err
	}
	err = hosts.PushToLocalHosts(domain)
	if err != nil {
		return err
	}
	err = hosts.PushToHosts()
	return err
}

func createDomain() error {
	var err error = nil
	domain := readDomain()
	err = createCerts(domain, config.SitePathInContainer)
	if err != nil {
		return err
	}
	err = hosts.PushToLocalHosts(domain)
	if err != nil {
		return err
	}
	err = hosts.PushToHosts()
	return err
}

func readDomain() string {
	return utils.ReadLine("Введите название сайта (доменное имя): ")
}

func createCerts(domain string, rootPath string) error {
	var err error = nil
	if !utils.FileIsExists(config.GetDockerFilesDirPath()) {
		err = internal.PublishFiles()
		if err != nil {
			return err
		}
	}
	err = certs.Create(domain, rootPath)
	return err
}
