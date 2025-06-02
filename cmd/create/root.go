package create

import (
	"docky/cmd/hosts"
	"docky/config"
	"docky/utils"
	"docky/yaml/helper"
	"fmt"
	"os"
	"path/filepath"
)

func CreateSite() error {
	domain, err := prepareDomain()
	if err != nil {
		return err
	}

	if err := pushToSimlinks(domain); err != nil {
		return err
	}

	dirPath := filepath.Join(config.GetSiteDirPath(), domain)
	if exists, _ := utils.FileIsExists(dirPath); !exists {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			return fmt.Errorf("ошибка создания директории сайта: %v", err)
		}
	}

	if err := createCerts(domain, filepath.Join(config.SitePathInContainer, domain), true); err != nil {
		return err
	}

	return hosts.PushToHosts()
}


func CreateDomain() error {
	domain, err := prepareDomain()
	if err != nil {
		return err
	}

	if err := createCerts(domain, config.SitePathInContainer, false); err != nil {
		return err
	}

	return hosts.PushToHosts()
}

func readDomain() string {
	return utils.ReadLine("Введите название сайта (доменное имя): ")
}

func pushToSimlinks(domain string) error {
	filePath := filepath.Join(config.GetConfFilesDirPath(), helper.App, "simlinks")
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	base := config.SitePathInContainer
	paths := []string{"bitrix", "local", "upload"}

	for _, p := range paths {
		line := fmt.Sprintf("%s/%s %s/%s/%s\n", base, p, base, domain, p)
		if _, err := file.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}


func prepareDomain() (string, error) {
	if err := initDir(); err != nil {
		return "", err
	}
	domain := readDomain()
	if err := hosts.PushToLocalHosts(domain); err != nil {
		return "", err
	}
	return domain, nil
}
