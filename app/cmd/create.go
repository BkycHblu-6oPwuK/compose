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
		validateWorkDir()
		err := createSite()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
		}
		fmt.Println("✅ сайт создан!")
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
		fmt.Println("✅ домен создан!")
	},
}

func init() {
	switch config.GetCurFramework() {
	case config.Bitrix:
		createCmd.AddCommand(createSiteModuleCmd)
	}
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
	err = pushToSimlinks(domain)
	if err != nil {
		return err
	}
	err = hosts.PushToHosts()
	return err
}

func pushToSimlinks(domain string) error {
	var file *os.File
	var err error
	filePath := filepath.Join(config.GetDockerFilesDirPath(), "app", "simlinks")

	if !utils.FileIsExists(filePath) {
		file, err = os.Create(filePath)
	} else {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	}
	if err != nil {
		return err
	}
	defer file.Close()

	base := config.SitePathInContainer

	lines := []string{
		fmt.Sprintf("%s/bitrix %s/%s/bitrix\n", base, base, domain),
		fmt.Sprintf("%s/local %s/%s/local\n", base, base, domain),
		fmt.Sprintf("%s/upload %s/%s/upload\n", base, base, domain),
	}

	for _, line := range lines {
		if _, err := file.WriteString(line); err != nil {
			return err
		}
	}

	return nil
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
